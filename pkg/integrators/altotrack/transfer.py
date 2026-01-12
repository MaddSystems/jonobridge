#!/usr/bin/env python3
import requests
import json
import os
from urllib3.exceptions import InsecureRequestWarning
import gc
from datetime import datetime, timedelta
from dateutil.relativedelta import relativedelta

# Suppress SSL warnings
requests.urllib3.disable_warnings(InsecureRequestWarning)

# Configuration
SOURCE_URL = "http://elasticserver.dwim.mx:9200"
TARGET_URL = "https://opensearch.madd.com.mx:9200"
TARGET_USER = "admin"
TARGET_PASS = "GPSc0ntr0l1"

def get_available_indices():
    """Get list of available indices from source Elasticsearch"""
    url = f"{SOURCE_URL}/_cat/indices?format=json"
    headers = {"Content-Type": "application/json"}
    
    try:
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        indices_data = response.json()
        return [idx['index'] for idx in indices_data if not idx['index'].startswith('.')]
    except Exception as e:
        print(f"Error getting indices: {e}")
        return []

def select_index_interactively():
    """Prompt user to select an index"""
    available_indices = get_available_indices()
    
    if not available_indices:
        print("No indices found!")
        return None
    
    print("\nAvailable indices:")
    for i, index in enumerate(available_indices, 1):
        print(f"{i}. {index}")
    
    while True:
        try:
            choice = input("\nEnter the number of the index you want to transfer: ").strip()
            idx = int(choice) - 1
            if 0 <= idx < len(available_indices):
                return available_indices[idx]
            else:
                print("Invalid choice. Please try again.")
        except (ValueError, KeyboardInterrupt):
            print("Invalid input. Please enter a number.")
            return None

def get_date_range():
    """Get date range for last month and current month"""
    now = datetime.now()
    
    # Start of last month (beginning of day)
    last_month_start = (now.replace(day=1) - relativedelta(months=1)).strftime('%Y-%m-%dT00:00:00Z')
    
    # Start of next month (to include all of current month)
    next_month_start = (now.replace(day=1) + relativedelta(months=1)).strftime('%Y-%m-%dT00:00:00Z')
    
    print(f"Date range: {last_month_start} to {next_month_start}")
    return last_month_start, next_month_start

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

def import_batch_to_opensearch(index_name, documents):
    """Import a batch of documents to OpenSearch"""
    if not documents:
        return False
        
    # Create bulk data for this batch
    bulk_data = []
    for doc in documents:
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
                print(f"    Batch had errors")
                # Count successful items in this batch
                success_count = 0
                for item in result.get('items', []):
                    if 'index' in item and item['index'].get('status') in [200, 201]:
                        success_count += 1
                return success_count
            else:
                return len(documents)
        else:
            print(f"    Batch failed with status {response.status_code}: {response.text[:200]}")
            return 0
            
    except requests.exceptions.RequestException as e:
        print(f"    Error importing batch: {e}")
        return 0

def migrate_index_streaming(index_name, batch_size=500):
    """Migrate index using streaming approach - process batches without storing all in memory"""
    print(f"Migrating index: {index_name}")
    
    # Get date range
    start_date, end_date = get_date_range()
    
    # Check source document count with date filter
    source_count = get_document_count_with_date_filter(index_name, start_date, end_date, is_target=False)
    print(f"  Source documents (filtered by date): {source_count}")
    
    if source_count == 0:
        print(f"  No documents found for the specified date range")
        return
    
    # Create index with 0 replicas
    print("  Creating index with 0 replicas...")
    if not create_index_with_zero_replicas(index_name):
        print(f"  Error: Could not create or access index {index_name}")
        return
    
    # Initial search request with date filter
    url = f"{SOURCE_URL}/{index_name}/_search"
    headers = {"Content-Type": "application/json"}
    params = {"scroll": "5m", "size": 1000}
    data = {
        "query": {
            "range": {
                "time": {
                    "gte": start_date,
                    "lt": end_date
                }
            }
        }
    }
    
    total_imported = 0
    batch_buffer = []
    
    try:
        # Initial request
        print("  Starting streaming migration...")
        response = requests.get(url, headers=headers, params=params, json=data)
        response.raise_for_status()
        result = response.json()
        
        scroll_id = result.get('_scroll_id')
        hits = result['hits']['hits']
        
        print(f"  Processing in batches of {batch_size} documents...")
        batch_num = 1
        
        while hits:
            # Add hits to batch buffer
            batch_buffer.extend(hits)
            
            # Process buffer when it reaches batch_size or when we're done
            while len(batch_buffer) >= batch_size:
                batch_to_process = batch_buffer[:batch_size]
                batch_buffer = batch_buffer[batch_size:]
                
                # Import this batch
                imported_count = import_batch_to_opensearch(index_name, batch_to_process)
                total_imported += imported_count
                
                # Calculate and display progress percentage
                progress_percentage = (total_imported / source_count) * 100 if source_count > 0 else 0
                
                if imported_count == len(batch_to_process):
                    print(f"    Batch {(total_imported-1)//batch_size + 1} imported successfully ({imported_count} docs) - Total: {total_imported} ({progress_percentage:.1f}%)")
                else:
                    print(f"    Batch {(total_imported-1)//batch_size + 1} partially imported ({imported_count}/{len(batch_to_process)} docs) - Total: {total_imported} ({progress_percentage:.1f}%)")
                
                # Force garbage collection to free memory
                gc.collect()
            
            # Get next scroll batch
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
                print(f"    Fetched batch {batch_num}: {len(hits)} documents from source")
                batch_num += 1
        
        # Process any remaining documents in buffer
        if batch_buffer:
            imported_count = import_batch_to_opensearch(index_name, batch_buffer)
            total_imported += imported_count
            progress_percentage = (total_imported / source_count) * 100 if source_count > 0 else 0
            print(f"    Final batch imported: {imported_count} docs - Total: {total_imported} ({progress_percentage:.1f}%)")
        
        # Clean up scroll
        if scroll_id:
            try:
                requests.delete(f"{SOURCE_URL}/_search/scroll", 
                              headers=headers, 
                              json={"scroll_id": scroll_id})
            except:
                pass  # Ignore cleanup errors
        
        # Verify target document count with date filter
        start_date, end_date = get_date_range()
        target_count = get_document_count_with_date_filter(index_name, start_date, end_date, is_target=True)
        print(f"  Target documents after import: {target_count}")
        if target_count == source_count:
            print(f"  ✓ Migration completed successfully for {index_name} (100.0%)")
        else:
            print(f"  ⚠ Migration completed with discrepancy: {target_count}/{source_count}")
        
    except requests.exceptions.RequestException as e:
        print(f"  Error during streaming migration: {e}")
    
    print("----------------------------------------")

def get_document_count_with_date_filter(index_name, start_date, end_date, is_target=False):
    """Get document count from an index with date filter"""
    base_url = TARGET_URL if is_target else SOURCE_URL
    url = f"{base_url}/{index_name}/_count"
    headers = {"Content-Type": "application/json"}
    
    data = {
        "query": {
            "range": {
                "time": {
                    "gte": start_date,
                    "lt": end_date
                }
            }
        }
    }
    
    try:
        if is_target:
            auth = (TARGET_USER, TARGET_PASS)
            response = requests.post(url, headers=headers, auth=auth, json=data, verify=False)
        else:
            response = requests.post(url, headers=headers, json=data)
        
        response.raise_for_status()
        result = response.json()
        return result.get('count', 0)
    except Exception as e:
        print(f"    Error getting count for {index_name}: {e}")
        return -1

def get_document_count(index_name, is_target=False):
    """Get total document count from an index"""
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

def main():
    """Main migration function"""
    selected_index = select_index_interactively()
    
    if selected_index:
        print(f"\nSelected index: {selected_index}")
        migrate_index_streaming(selected_index, batch_size=500)
        print("Migration completed!")
    else:
        print("No index selected. Exiting.")

if __name__ == "__main__":
    main()