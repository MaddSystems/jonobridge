#!/usr/bin/env python3
import requests
import json
import os
from urllib3.exceptions import InsecureRequestWarning

# Suppress SSL warnings
requests.urllib3.disable_warnings(InsecureRequestWarning)

# Configuration
SOURCE_URL = "http://elasticserver.dwim.mx:9200"
TARGET_URL = "https://opensearch.madd.com.mx:9200"
TARGET_USER = "admin"
TARGET_PASS = "GPSc0ntr0l1"

# List of all indices
indices = [
    "scania-mx_vecfleet"
]

def create_index_with_zero_replicas(index_name):
    """Create index with 0 replicas in OpenSearch"""
    url = f"{TARGET_URL}/{index_name}"
    headers = {"Content-Type": "application/json"}
    auth = (TARGET_USER, TARGET_PASS)
    data = {
        "settings": {
            "number_of_shards": 1,
            "number_of_replicas": 0
        }
    }
    
    try:
        response = requests.put(url, headers=headers, auth=auth, json=data, verify=False)
        print(f"    Create index response: {response.status_code}")
        if response.status_code == 400 and "already exists" in response.text:
            print(f"    Index {index_name} already exists - continuing with migration")
            return True
        elif response.status_code not in [200, 201]:
            print(f"    Create index error: {response.text}")
        return response.status_code in [200, 201, 400]  # 400 is OK if index exists
    except Exception as e:
        print(f"    Create index exception: {e}")
        return False

def export_all_from_elasticsearch(index_name):
    """Export ALL data from Elasticsearch using scroll API"""
    # Initial search request
    url = f"{SOURCE_URL}/{index_name}/_search"
    headers = {"Content-Type": "application/json"}
    params = {"scroll": "5m", "size": 1000}
    data = {"query": {"match_all": {}}}
    
    all_hits = []
    
    try:
        # Initial request
        response = requests.get(url, headers=headers, params=params, json=data)
        response.raise_for_status()
        result = response.json()
        
        print(f"    Total documents available: {result['hits']['total']['value'] if isinstance(result['hits']['total'], dict) else result['hits']['total']}")
        
        scroll_id = result.get('_scroll_id')
        hits = result['hits']['hits']
        all_hits.extend(hits)
        
        print(f"    Fetched batch 1: {len(hits)} documents")
        
        # Continue scrolling while there are more documents
        batch_num = 2
        while hits:
            scroll_url = f"{SOURCE_URL}/_search/scroll"
            scroll_data = {
                "scroll": "5m",
                "scroll_id": scroll_id
            }
            
            response = requests.post(scroll_url, headers=headers, json=scroll_data)
            response.raise_for_status()
            result = response.json()
            
            scroll_id = result.get('_scroll_id')
            hits = result['hits']['hits']
            
            if hits:
                all_hits.extend(hits)
                print(f"    Fetched batch {batch_num}: {len(hits)} documents (total so far: {len(all_hits)})")
                batch_num += 1
        
        print(f"    Total documents fetched: {len(all_hits)}")
        return {"hits": {"hits": all_hits}}
        
    except requests.exceptions.RequestException as e:
        print(f"    Error exporting {index_name}: {e}")
        return None

def convert_to_bulk_format(data, index_name):
    """Convert Elasticsearch response to bulk format"""
    bulk_data = []
    
    if data and 'hits' in data and 'hits' in data['hits']:
        for hit in data['hits']['hits']:
            # Add index action
            bulk_data.append(json.dumps({
                "index": {
                    "_index": index_name,
                    "_id": hit["_id"]
                }
            }))
            # Add document source
            bulk_data.append(json.dumps(hit["_source"]))
    
    result = "\n".join(bulk_data) + "\n" if bulk_data else ""
    print(f"    Bulk data size: {len(result)} characters")
    return result

def import_to_opensearch(index_name, bulk_data):
    """Import data to OpenSearch using bulk API"""
    if not bulk_data.strip():
        print("    No bulk data to import")
        return False
        
    url = f"{TARGET_URL}/{index_name}/_bulk"
    headers = {"Content-Type": "application/json"}
    auth = (TARGET_USER, TARGET_PASS)
    
    try:
        response = requests.post(url, headers=headers, auth=auth, data=bulk_data, verify=False)
        print(f"    Bulk import response: {response.status_code}")
        
        if response.status_code == 200:
            result = response.json()
            if result.get('errors'):
                print(f"    Bulk import had errors: {json.dumps(result, indent=2)}")
                return False
            else:
                print(f"    Bulk import successful")
                return True
        else:
            print(f"    Bulk import failed: {response.text}")
            return False
            
    except requests.exceptions.RequestException as e:
        print(f"    Error importing to {index_name}: {e}")
        return False

def import_to_opensearch_batch(index_name, documents, batch_size=500):
    """Import documents to OpenSearch in batches"""
    total_docs = len(documents)
    print(f"    Importing {total_docs} documents in batches of {batch_size}")
    
    success_count = 0
    
    for i in range(0, total_docs, batch_size):
        batch_docs = documents[i:i + batch_size]
        
        # Create bulk data for this batch
        bulk_data = []
        for doc in batch_docs:
            bulk_data.append(json.dumps({
                "index": {
                    "_index": index_name,
                    "_id": doc["_id"]
                }
            }))
            bulk_data.append(json.dumps(doc["_source"]))
        
        bulk_string = "\n".join(bulk_data) + "\n"
        
        # Import this batch
        url = f"{TARGET_URL}/{index_name}/_bulk"
        headers = {"Content-Type": "application/json"}
        auth = (TARGET_USER, TARGET_PASS)
        
        try:
            response = requests.post(url, headers=headers, auth=auth, data=bulk_string, verify=False)
            
            if response.status_code == 200:
                result = response.json()
                if result.get('errors'):
                    print(f"    Batch {i//batch_size + 1} had errors")
                    # Count successful items in this batch
                    for item in result.get('items', []):
                        if 'index' in item and item['index'].get('status') in [200, 201]:
                            success_count += 1
                else:
                    success_count += len(batch_docs)
                    print(f"    Batch {i//batch_size + 1}/{(total_docs-1)//batch_size + 1} imported successfully ({len(batch_docs)} docs)")
            else:
                print(f"    Batch {i//batch_size + 1} failed with status {response.status_code}: {response.text[:200]}")
                
        except requests.exceptions.RequestException as e:
            print(f"    Error importing batch {i//batch_size + 1}: {e}")
    
    print(f"    Successfully imported {success_count}/{total_docs} documents")
    return success_count > 0

def get_document_count(index_name, is_target=False):
    """Get document count from an index"""
    base_url = TARGET_URL if is_target else SOURCE_URL
    url = f"{base_url}/{index_name}/_count"
    headers = {"Content-Type": "application/json"}
    
    try:
        if is_target:
            auth = (TARGET_USER, TARGET_PASS)
            response = requests.get(url, headers=headers, auth=auth, verify=False)
        else:
            response = requests.get(url, headers=headers)
        
        response.raise_for_status()
        result = response.json()
        return result.get('count', 0)
    except Exception as e:
        print(f"    Error getting count for {index_name}: {e}")
        return -1

def migrate_index(index_name):
    """Migrate a single index"""
    print(f"Migrating index: {index_name}")
    
    # Check source document count
    source_count = get_document_count(index_name, is_target=False)
    print(f"  Source documents: {source_count}")
    
    # Create index with 0 replicas
    print("  Creating index with 0 replicas...")
    if not create_index_with_zero_replicas(index_name):
        print(f"  Error: Could not create or access index {index_name}")
        return
    
    # Export from Elasticsearch
    print("  Exporting from source...")
    exported_data = export_all_from_elasticsearch(index_name)
    
    if exported_data and 'hits' in exported_data and 'hits' in exported_data['hits']:
        documents = exported_data['hits']['hits']
        print(f"  Converting {len(documents)} documents...")
        
        if documents:
            print("  Importing to target in batches...")
            if import_to_opensearch_batch(index_name, documents):
                # Verify target document count
                target_count = get_document_count(index_name, is_target=True)
                print(f"  Target documents after import: {target_count}")
                if target_count == source_count:
                    print(f"  ✓ Migration completed successfully for {index_name}")
                else:
                    print(f"  ⚠ Migration completed with discrepancy: {target_count}/{source_count}")
            else:
                print(f"  Error importing {index_name}")
        else:
            print(f"  No documents to import for {index_name}")
    else:
        print(f"  Error exporting {index_name}")
    
    print("----------------------------------------")

def main():
    """Main migration function"""
    for index in indices:
        migrate_index(index)
    
    print("All migrations completed!")

if __name__ == "__main__":
    main()

