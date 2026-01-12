import requests
import json
import time

# Configuración
BASE_URL = "https://jonobridge.madd.com.mx/api/v1/diagnostic"
HEADERS = {
    "accept": "application/json"
}

def get_namespaces():
    """Obtiene la lista de namespaces"""
    url = f"{BASE_URL}/namespaces"
    response = requests.get(url, headers=HEADERS)
    
    if response.status_code != 200:
        print(f"Error al obtener namespaces: {response.status_code}")
        print(response.text)
        return []
    
    data = response.json()
    return data.get("namespaces", [])

def restart_namespace(namespace_name):
    """Reinicia todos los deployments en un namespace (lo que reinicia todos los pods)"""
    url = f"{BASE_URL}/restart-namespace/{namespace_name}"
    response = requests.post(url, headers=HEADERS)
    
    if response.status_code == 200:
        result = response.json()
        print(f"Reiniciado '{namespace_name}': {result['deployments_restarted']} deployments")
        if result.get("message"):
            print(f"   → {result['message']}")
        return True
    else:
        print(f"Error al reiniciar '{namespace_name}': {response.status_code}")
        print(response.text)
        return False

def main():
    print("Obteniendo lista de namespaces...")
    namespaces = get_namespaces()
    
    if not namespaces:
        print("No se encontraron namespaces o hubo un error.")
        return
    
    print(f"Se encontraron {len(namespaces)} namespaces.\n")
    
    # Filtrar solo los que están healthy (opcional)
    healthy_namespaces = [ns for ns in namespaces if ns["status"] == "healthy"]
    
    print(f"Reiniciando {len(healthy_namespaces)} namespaces saludables...\n")
    
    for ns in healthy_namespaces:
        name = ns["name"]
        print(f"Reiniciando namespace: {name}...")
        success = restart_namespace(name)
        
        if success:
            print("Éxito.\n")
        else:
            print("Falló.\n")
        
        # Opcional: esperar entre reinicios para no saturar
        time.sleep(2)
    
    print("Proceso completado.")

if __name__ == "__main__":
    main()
