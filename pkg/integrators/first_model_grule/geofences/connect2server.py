import mysql.connector
import os

def get_db_connection(suppress_message=False, database="geofences"):
    """Establece y retorna una conexión a la base de datos MySQL."""
    try:
        conn = mysql.connector.connect(
            host=os.getenv("MYSQL_HOST"),
            port=os.getenv("MYSQL_PORT", 3306),
            user=os.getenv("MYSQL_USER"),
            password=os.getenv("MYSQL_PASS"),
            database=database
        )
        if not suppress_message:
            print("Conexión a la base de datos exitosa.")
        return conn
    except mysql.connector.Error as err:
        print(f"Error al conectar a MySQL: {err}")
        return None

def create_database_and_tables():
    """Crea la base de datos 'geofences' y las tablas necesarias si no existen."""
    # Connect without specifying database first
    try:
        conn = mysql.connector.connect(
            host=os.getenv("MYSQL_HOST"),
            port=os.getenv("MYSQL_PORT", 3306),
            user=os.getenv("MYSQL_USER"),
            password=os.getenv("MYSQL_PASS")
        )
        cursor = conn.cursor()
        
        # Create database if not exists
        print("Creando base de datos 'geofences' si no existe...")
        cursor.execute("CREATE DATABASE IF NOT EXISTS geofences")
        cursor.execute("USE geofences")
        
        # Create geofences table
        print("Creando tabla 'geofences' si no existe...")
        cursor.execute("""
            CREATE TABLE IF NOT EXISTS geofences (
                id int NOT NULL,
                name varchar(255) NOT NULL,
                description text,
                shapeType varchar(50) DEFAULT NULL,
                centerLat double DEFAULT NULL,
                centerLon double DEFAULT NULL,
                radius double DEFAULT NULL,
                boundingBoxMinX double DEFAULT NULL,
                boundingBoxMaxX double DEFAULT NULL,
                boundingBoxMinY double DEFAULT NULL,
                boundingBoxMaxY double DEFAULT NULL,
                created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                PRIMARY KEY (id)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
        """)
        
        # Create geofence_groups table
        print("Creando tabla 'geofence_groups' si no existe...")
        cursor.execute("""
            CREATE TABLE IF NOT EXISTS geofence_groups (
                id int NOT NULL,
                name varchar(255) NOT NULL,
                description text,
                colour varchar(7) DEFAULT NULL,
                priority int DEFAULT NULL,
                pinned tinyint(1) DEFAULT NULL,
                useInGeocoding tinyint(1) DEFAULT NULL,
                PRIMARY KEY (id),
                UNIQUE KEY name (name)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
        """)
        
        # Create geofence_group_mapping table
        print("Creando tabla 'geofence_group_mapping' si no existe...")
        cursor.execute("""
            CREATE TABLE IF NOT EXISTS geofence_group_mapping (
                group_id int NOT NULL,
                geofence_id int NOT NULL,
                PRIMARY KEY (group_id, geofence_id),
                KEY idx_group (group_id),
                KEY idx_geofence (geofence_id),
                CONSTRAINT geofence_group_mapping_ibfk_1 FOREIGN KEY (group_id) REFERENCES geofence_groups (id) ON DELETE CASCADE,
                CONSTRAINT geofence_group_mapping_ibfk_2 FOREIGN KEY (geofence_id) REFERENCES geofences (id) ON DELETE CASCADE
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
        """)
        
        conn.commit()
        print("✅ Base de datos y tablas creadas exitosamente.\n")
        cursor.close()
        conn.close()
        return True
        
    except mysql.connector.Error as err:
        print(f"Error al crear base de datos y tablas: {err}")
        return False

if __name__ == "__main__":
    # Example usage:
    create_database_and_tables()
    conn = get_db_connection()
    if conn:
        conn.close()
        print("Conexión cerrada.")
