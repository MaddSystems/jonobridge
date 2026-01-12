USE automation_app;

-- Drop existing tables if they exist to ensure clean recreation
DROP TABLE IF EXISTS `service_connections`;
DROP TABLE IF EXISTS `listener`;
DROP TABLE IF EXISTS `forwarder`;
DROP TABLE IF EXISTS `automation`;
DROP TABLE IF EXISTS `mqtt`;
DROP TABLE IF EXISTS `proxy`;
DROP TABLE IF EXISTS `meitrackprotocol`;
DROP TABLE IF EXISTS `pinoprotocol`;
DROP TABLE IF EXISTS `suntechprotocol`;
DROP TABLE IF EXISTS `hubaiprotocol`;
DROP TABLE IF EXISTS `queclinkprotocol`;
DROP TABLE IF EXISTS `ruptelaprotocol`;
DROP TABLE IF EXISTS `clients`;
DROP TABLE IF EXISTS `meitrack`;
DROP TABLE IF EXISTS `recursoconfiable`;
DROP TABLE IF EXISTS `send2elastic`;
DROP TABLE IF EXISTS `unigis`;
DROP TABLE IF EXISTS `vecfleet`;
DROP TABLE IF EXISTS `unigis`;
DROP TABLE IF EXISTS `activetrack`;
DROP TABLE IF EXISTS `altotrack`;
DROP TABLE IF EXISTS `xpot`;
DROP TABLE IF EXISTS `clients`;
DROP TABLE IF EXISTS `avocadocloud`;
DROP TABLE IF EXISTS `motumcloud`;
DROP TABLE IF EXISTS `skyangel`;
DROP TABLE IF EXISTS `lobosoftware`;
DROP TABLE IF EXISTS `gt06_serials`;

-- Create clients table
CREATE TABLE IF NOT EXISTS `clients` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `port` INT NOT NULL DEFAULT 0,
    `api_port` INT NOT NULL DEFAULT 0,
    `web_port` INT NOT NULL DEFAULT 0,
    `namespace` VARCHAR(255) NOT NULL DEFAULT 'default_namespace',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create table for storing GT06 serial numbers by IMEI
CREATE TABLE IF NOT EXISTS `gt06_serials` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `imei` VARCHAR(50) NOT NULL,
    `serial_number` INT UNSIGNED NOT NULL DEFAULT 0,
    `last_updated` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `unique_imei` (`imei`)
);

-- Create table for storing client automation settings
CREATE TABLE IF NOT EXISTS `automation` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `enabled` BOOLEAN DEFAULT TRUE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing MQTT broker configuration
CREATE TABLE IF NOT EXISTS `mqtt` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `broker_host` VARCHAR(255) NOT NULL,
    `broker_port` INT NOT NULL DEFAULT 1883,
    `topic` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing recursoconfiable configuration
CREATE TABLE IF NOT EXISTS `recursoconfiable` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `user_name` VARCHAR(255),
    `password` VARCHAR(255),
    `url` VARCHAR(255),
    `customer_id` VARCHAR(255),
    `customer_name` VARCHAR(255),
    `plates_url` VARCHAR(255),

    `replicas` INT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing forwarder configuration
CREATE TABLE IF NOT EXISTS `forwarder` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `forwarder_host` VARCHAR(255) NOT NULL,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing meitrack proxy configuration
CREATE TABLE IF NOT EXISTS `proxy` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `platform_host` VARCHAR(255) NOT NULL,
    `port` INT NOT NULL,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing httprequest configuration
CREATE TABLE IF NOT EXISTS `httprequest` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `http_url` TEXT NOT NULL,
    `http_polling_time` INT NOT NULL,
    `skywave_access_id` VARCHAR(255),
    `skywave_password` VARCHAR(255),
    `skywave_from_id` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing listener configuration
CREATE TABLE IF NOT EXISTS `listener` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `port` INT NOT NULL,
    `api_port` INT NOT NULL,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing meitrackprotocol configuration
CREATE TABLE IF NOT EXISTS `meitrackprotocol` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing pinoprotocol configuration
CREATE TABLE IF NOT EXISTS `pinoprotocol` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing suntechprotocol configuration
CREATE TABLE IF NOT EXISTS `suntechprotocol` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing huabaoprotocol configuration
CREATE TABLE IF NOT EXISTS `huabaoprotocol` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing queclinkprotocol configuration
CREATE TABLE IF NOT EXISTS `queclinkprotocol` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing skywaveprotocol configuration
CREATE TABLE IF NOT EXISTS `skywaveprotocol` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);


-- Create table for storing ruptelaprotocol configuration
CREATE TABLE IF NOT EXISTS `ruptelaprotocol` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing xpot configuration
CREATE TABLE IF NOT EXISTS `xpot` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `xpot_forward_host` VARCHAR(255) NOT NULL,
    `mysql_user` VARCHAR(255) NOT NULL,
    `mysql_password` VARCHAR(255) NOT NULL,
    `mysql_database` VARCHAR(255) NOT NULL,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing service connections
CREATE TABLE IF NOT EXISTS `service_connections` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `source_service` VARCHAR(255) NOT NULL,
    `source_id` INT NOT NULL,
    `target_service` VARCHAR(255) NOT NULL,
    `target_id` INT NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE,
    INDEX `idx_client_source` (`client_id`, `source_service`, `source_id`),
    INDEX `idx_client_target` (`client_id`, `target_service`, `target_id`)
);


-- Create table for storing meitrack configuration
CREATE TABLE IF NOT EXISTS `meitrack` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `meitrack_host` VARCHAR(255) NOT NULL,
    `meitrack_fwd_only_adas` VARCHAR(1) NOT NULL DEFAULT 'Y',   
    `meitrack_mock_imei` VARCHAR(1) NOT NULL DEFAULT 'N',
    `meitrack_mock_value` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `send_to_elastic` VARCHAR(1) NOT NULL DEFAULT 'N',
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing send2elastic configuration
CREATE TABLE IF NOT EXISTS `send2elastic` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `plates_url` VARCHAR(255),
    `replicas` INT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);
-- Create table for storing unigis configuration
CREATE TABLE IF NOT EXISTS `unigis` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    
    `client_id` INT NOT NULL,`user_name` VARCHAR(255),
    `user` VARCHAR(255),
    `user_key` VARCHAR(255),
    `url` VARCHAR(255),
    `service` VARCHAR(255),
    `plates_url` VARCHAR(255),
    `replicas` INT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);
-- Create table for storing vecfleet configuration
CREATE TABLE IF NOT EXISTS `vecfleet` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    
    `client_id` INT NOT NULL,
    `vecfleet_url` VARCHAR(255),
    `vecfleet_post_url` VARCHAR(255),
    `vecfleet_email` VARCHAR(255),
    `vecfleet_name` VARCHAR(255),
    `vecfleet_password` VARCHAR(255),
    `service` VARCHAR(255),
    `plates_url` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing unigis configuration
CREATE TABLE IF NOT EXISTS `unigis` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `unigis_user` VARCHAR(255),
    `unigis_key` VARCHAR(255),
    `unigis_url` VARCHAR(255),
    `plates_url` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing activetrack configuration
CREATE TABLE IF NOT EXISTS `activetrack` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `activetrack_token` VARCHAR(255),
    `activetrack_url` VARCHAR(255),
    `activetrack_virtual_imei_url` VARCHAR(255),
    `plates_url` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing altotrack configuration
CREATE TABLE IF NOT EXISTS `altotrack` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `altotrack_proveedor` VARCHAR(255),
    `altotrack_url` VARCHAR(255),
    `plates_url` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing avocadocloud configuration
CREATE TABLE IF NOT EXISTS `avocadocloud` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `plates_url` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `avocado_user` VARCHAR(255),
    `avocado_password` VARCHAR(255),
    `avocado_user_adm` VARCHAR(255),
    `avocado_url` VARCHAR(255) NOT NULL,
    `mysql_user` VARCHAR(255) NOT NULL,
    `mysql_password` VARCHAR(255) NOT NULL,
    `mysql_database` VARCHAR(255) NOT NULL,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing motumcloud configuration
CREATE TABLE IF NOT EXISTS `motumcloud` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `plates_url` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `motum_user` VARCHAR(255),        
    `motum_password` VARCHAR(255),        
    `motum_referer` VARCHAR(255),
    `motum_apikey` VARCHAR(255), 
    `motum_carrier` VARCHAR(255),     
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing motumcloud configuration
CREATE TABLE IF NOT EXISTS `skyangel` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `skyangel_user` VARCHAR(255),
    `skyangel_key` VARCHAR(255),
    `skyangel_url` VARCHAR(255),
    `plates_url` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing lobosoftware
CREATE TABLE IF NOT EXISTS `lobosoftware` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `plates_url` VARCHAR(255),
    `lobosoftware_url` VARCHAR(500),
    `lobosoftware_token_url` VARCHAR(500),
    `lobosoftware_user` VARCHAR(255),
    `lobosoftware_user_key` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing movup configuration
CREATE TABLE IF NOT EXISTS `movup` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `plates_url` VARCHAR(255),
    `movup_user` VARCHAR(255),
    `movup_userkey` VARCHAR(255),
    `movup_url_path` VARCHAR(255),
    `movup_token` VARCHAR(255),
    `movup_provider` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing rfl3
CREATE TABLE IF NOT EXISTS `rfl3` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `plates_url` VARCHAR(255),
    `rfl3_url` VARCHAR(500),
    `rfl3_xapikey` VARCHAR(500),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing gpsinfo
CREATE TABLE IF NOT EXISTS `gpsinfo` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `plates_url` VARCHAR(255),
    `gpsinfo_user` VARCHAR(500),
    `gpsinfo_password` VARCHAR(500),
    `gpsinfo_host` VARCHAR(500),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing logitrack
CREATE TABLE IF NOT EXISTS `logitrack` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `logitrack_user` VARCHAR(255),
    `logitrack_user_key` VARCHAR(255),
    `logitrack_urlway` VARCHAR(500),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing controlnavigation
CREATE TABLE IF NOT EXISTS `controlnavigation` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `controlnavigation_user` VARCHAR(255),
    `controlnavigation_user_key` VARCHAR(255),
    `controlnavigation_url` VARCHAR(500),
    `controlnavigation_token_url` VARCHAR(500),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing controlt
CREATE TABLE IF NOT EXISTS `controlt` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `controlt_user` VARCHAR(255),
    `controlt_user_key` VARCHAR(255),
    `controlt_url` VARCHAR(500),
    `plates_url` VARCHAR(500),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing telemetry configuration
CREATE TABLE IF NOT EXISTS `telemetry` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `plates_url` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `telemetry_url` VARCHAR(255),        
    `telemetry_owner_id` VARCHAR(255),          
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing uscentral
CREATE TABLE IF NOT EXISTS `uscentral` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `uscentral_user` VARCHAR(255),
    `uscentral_user_key` VARCHAR(255),
    `uscentral_url` VARCHAR(500),
    `uscentral_token` VARCHAR(500),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing portal configuration
CREATE TABLE IF NOT EXISTS `portal` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT,
    `mysql_user` VARCHAR(255) DEFAULT '',
    `mysql_password` VARCHAR(255) DEFAULT '',
    `mysql_database` VARCHAR(255) DEFAULT '',
    `plates_url` VARCHAR(255) DEFAULT '',
    `portal_endpoint` VARCHAR(255) DEFAULT '',
    `portal_port` INT DEFAULT 0,
    `portal_user` VARCHAR(255) DEFAULT '',
    `portal_password` VARCHAR(255) DEFAULT '',
    `portal_script` TEXT,
    `replicas` INT DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing device information
CREATE TABLE IF NOT EXISTS `devices` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `imei` VARCHAR(16),
    `sn` VARCHAR(255),
    `password` VARCHAR(16),
    `creation_date` DATETIME,
    `ff0` INT,
    `log` TEXT,
    `mygroup` VARCHAR(40),
    `email` VARCHAR(64),
    `plates` VARCHAR(15),
    `protocol` INT,
    `lastupdate` DATETIME,
    `eco` VARCHAR(30),
    `latitude` DECIMAL(10,6),
    `longitude` DECIMAL(10,6),
    `altitude` INT,
    `speed` FLOAT,
    `angle` FLOAT,
    `url` VARCHAR(255),
    `vin` VARCHAR(40),
    `enterprise` VARCHAR(255),
    `event_code` VARCHAR(5),
    `telephone` VARCHAR(255),
    `device_key` VARCHAR(255),
    `name` VARCHAR(255),
    `last_alarm` DATETIME,
    `alarm_status` VARCHAR(25),
    `last_followme` DATETIME,
    `followme_status` VARCHAR(25),
    `last_name` VARCHAR(30),
    `maiden_name` VARCHAR(30),
    `street` VARCHAR(255),
    `delegacion` VARCHAR(255),
    `number` VARCHAR(15),
    `zip` VARCHAR(10),
    `colonia` VARCHAR(255),
    `panic` VARCHAR(5),
    `dvr` VARCHAR(1),
    `alarmcount` INT,
    `alt_lat` DECIMAL(10,6),
    `alt_lon` DECIMAL(10,6),
    `tipodeunidad` VARCHAR(255),
    `marca` VARCHAR(255),
    `submarca` VARCHAR(255),
    `fechamodelo` INT,
    `zona` VARCHAR(255),
    `municipio` VARCHAR(255),
    `numconsesion` VARCHAR(255)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create table for storing alephri
CREATE TABLE IF NOT EXISTS `alephri` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `plates_url` VARCHAR(255),
    `alephri_company` VARCHAR(255),
    `alephri_url` VARCHAR(255),
    `alephri_token` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing webrelay
CREATE TABLE IF NOT EXISTS `webrelay` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `ggs_user` VARCHAR(255),
    `ggs_password` VARCHAR(255),
    `app_id` INT,
    `webrelay_token` VARCHAR(255),
    `portal_endpoint` VARCHAR(255),
    `webrelay_port` INT,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for monitor
CREATE TABLE `monitor_config` (
  `id` int NOT NULL AUTO_INCREMENT,
  `service_type` varchar(50) NOT NULL,
  `check_interval_minutes` int DEFAULT '5',
  `data_threshold_minutes` int DEFAULT '5',
  `restart_enabled` tinyint(1) DEFAULT '1',
  `alert_enabled` tinyint(1) DEFAULT '0',
  `alert_webhook_url` varchar(500) DEFAULT NULL,
  `request_timeout_seconds` int DEFAULT '30',
  `retry_attempts` int DEFAULT '3',
  `retry_delay_seconds` int DEFAULT '5',
  `last_check` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `last_success` timestamp NULL DEFAULT NULL,
  `last_failure` timestamp NULL DEFAULT NULL,
  `failure_count` int DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `whatsapp_alerts_enabled` tinyint(1) DEFAULT '0',
  `last_whatsapp_alert` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `service_type` (`service_type`)
);

-- Create table for whatsapp
CREATE TABLE `whatsapp_config` (
  `id` int NOT NULL AUTO_INCREMENT,
  `enabled` tinyint(1) DEFAULT '0',
  `api_key` varchar(255) NOT NULL,
  `channel_id` varchar(255) NOT NULL,
  `namespace` varchar(255) NOT NULL,
  `template_name` varchar(100) DEFAULT 'server_alerts',
  `language_code` varchar(10) DEFAULT 'es_MX',
  `alert_interval_minutes` int DEFAULT '30',
  `min_failure_count` int DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)

-- Create table for whatsapp phones
CREATE TABLE `whatsapp_phones` (
  `id` int NOT NULL AUTO_INCREMENT,
  `phone_number` varchar(20) NOT NULL,
  `contact_name` varchar(100) DEFAULT NULL,
  `enabled` tinyint(1) DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_phone` (`phone_number`)
)

-- Create table for storing dicsa
CREATE TABLE IF NOT EXISTS `dicsa` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `plates_url` VARCHAR(500),
    `dicsa_url` VARCHAR(500),
    `dicsa_token_url` VARCHAR(500),
    `dicsa_user` VARCHAR(255),
    `dicsa_user_key` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(500),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing nstech configuration
CREATE TABLE IF NOT EXISTS `nstech` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `plates_url` VARCHAR(255),
    `nstech_client_id` VARCHAR(255),
    `nstech_client_secret` VARCHAR(255),
    `nstech_technology_id` VARCHAR(255),
    `nstech_account_id` VARCHAR(255),
    `nstech_url` VARCHAR(255),
    `nstech_token_url` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing httpinput
CREATE TABLE IF NOT EXISTS `httpinput` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `portal_endpoint` VARCHAR(255),
    `webrelay_port` INT,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

-- Create table for storing gpsgatehttp
CREATE TABLE IF NOT EXISTS `gpsgatehttp` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `gpsgate_test_telegram` VARCHAR(255),
    `telegram_api_url` VARCHAR(255),
    `telegram_bot_token` VARCHAR(255),
    `gpsgate_test_chat_id` VARCHAR(255),
    `telegram_message_header` VARCHAR(255),
    `telegram_additional_fields` VARCHAR(255),
    `ego_api_url` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `gruasaya` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `gruasaya_url` VARCHAR(255),
    `gruasaya_token` VARCHAR(255),
    `elastic_doc_name` VARCHAR(255),
    `elastic_url` VARCHAR(255),
    `elastic_user` VARCHAR(255),
    `elastic_password` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `websqldump` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `mysql_user` VARCHAR(255),
    `mysql_password` VARCHAR(255),
    `mysql_database` VARCHAR(255),
    `web_user` VARCHAR(255),
    `web_password` VARCHAR(255),
    `portal_endpoint` VARCHAR(255),
    `websqldump_port` INT,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `send2mysql` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `mysql_user` VARCHAR(255),
    `mysql_password` VARCHAR(255),
    `mysql_database` VARCHAR(255),
    `mysql_update` VARCHAR(255),
    `plates_url` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `send2http` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `send2http_url` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);
    `telegram_chat_id` VARCHAR(255),
    `websqldump_port` INT,
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `send2http` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `send2http_url` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `grule` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `mysql_user` VARCHAR(255),
    `mysql_password` VARCHAR(255),
    `mysql_database` VARCHAR(255),
    `mysql_update` VARCHAR(255),
    `portal_endpoint` VARCHAR(255),
    `telegram_bot_token` VARCHAR(255),
    `telegram_chat_id` VARCHAR(255),
    `grule_web_port` INT,
    `grule_audit_enabled` VARCHAR(1),
    `grule_audit_level` VARCHAR(255),
    `replicas` INT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE
);