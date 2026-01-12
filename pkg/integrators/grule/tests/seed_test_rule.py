import mysql.connector
import os

# Database configuration from environment or defaults
DB_HOST = os.getenv("MYSQL_HOST", "127.0.0.1")
DB_USER = os.getenv("MYSQL_USER", "gpscontrol")
DB_PASS = os.getenv("MYSQL_PASS", "qazwsxedc")
DB_NAME = os.getenv("MYSQL_DB", "grule")

def seed():
    print(f"🌱 Seeding test rule to {DB_NAME} at {DB_HOST}...")
    try:
        conn = mysql.connector.connect(
            host=DB_HOST,
            user=DB_USER,
            password=DB_PASS,
            database=DB_NAME
        )
        cursor = conn.cursor()
        
        # Clear existing rules for clean test
        cursor.execute("DELETE FROM fleet_rules WHERE name = 'AuditIntegrationTest'")
        
        # Insert test rule
        sql = """
        INSERT INTO fleet_rules (name, description, grl_content, audit_manifest, priority, active)
        VALUES (%s, %s, %s, %s, %s, %s)
        """
        rule_name = "AuditIntegrationTest"
        description = "Test rule for declarative audit"
        grl = 'rule AuditIntegrationTest salience 100 { when IncomingPacket.Speed >= 0 then actions.Log("Audit test fired"); IncomingPacket.Speed = -1; }'
        manifest = """
stages:
  - rule: AuditIntegrationTest
    order: 1
    audit:
      enabled: true
      description: "Integration Test Step"
      level: info
      is_alert: false
"""
        cursor.execute(sql, (rule_name, description, grl, manifest, 100, 1))
        
        conn.commit()
        print("✅ Test rule seeded successfully")
        cursor.close()
        conn.close()
    except Exception as e:
        print(f"❌ Failed to seed rule: {e}")

if __name__ == "__main__":
    seed()
