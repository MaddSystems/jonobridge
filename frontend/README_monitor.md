# JonoBridge Monitoring System

The JonoBridge monitoring system consists of two Python scripts that monitor GPS tracking services deployed in Kubernetes. These monitors ensure that GPS data is being properly ingested into Elasticsearch and can automatically restart services or send alerts when issues are detected.

## Components

### 1. `monitor.py` - Elasticsearch Service Monitor

A comprehensive monitoring system that checks GPS tracking services for data ingestion health.

### 2. `monitor_whatsapp.py` - Monitor with WhatsApp Alerts

Extended version of the monitor that includes WhatsApp alert functionality using MessageBird API.

## What It Does

### Core Monitoring Features

- **Elasticsearch Health Checks**: Monitors Elasticsearch indices to ensure GPS tracking data is being received within configured time thresholds
- **Service Auto-Restart**: Automatically restarts Kubernetes deployments when services fail to receive data
- **Webhook Alerts**: Sends alerts to configured webhook URLs when issues are detected
- **Configurable Monitoring**: Each service type can have different check intervals, data thresholds, and retry settings
- **Database Integration**: Stores monitoring configuration and status in MySQL database
- **Comprehensive Logging**: Detailed logging with multiple debug levels for troubleshooting

### WhatsApp Alert Features (monitor_whatsapp.py only)

- **WhatsApp Notifications**: Sends formatted alerts to configured phone numbers using WhatsApp Business API
- **Message Templates**: Uses predefined WhatsApp templates for consistent messaging
- **Multi-Phone Support**: Can send alerts to multiple phone numbers simultaneously
- **Alert Throttling**: Prevents alert spam with configurable timing intervals
- **Failure Count Tracking**: Only sends alerts after a minimum number of consecutive failures

## Architecture

### Service Discovery
- Automatically discovers services that have Elasticsearch configuration
- Scans database tables for `elastic_doc_name` columns
- Supports multiple service types (listeners, interpreters, integrators)

### Monitoring Cycle
1. **Discovery**: Find all services with Elasticsearch configuration
2. **Check Timing**: Verify if it's time to check each service (based on intervals)
3. **Data Validation**: Query Elasticsearch for recent data within time thresholds
4. **Status Update**: Update database with check results and failure counts
5. **Action Taking**: Restart services or send alerts based on configuration
6. **WhatsApp Alerts**: Send formatted alerts if conditions are met

### Database Tables

#### monitor_config
```sql
- id: Primary key
- service_type: Service identifier (unique)
- check_interval_minutes: How often to check this service (default: 5)
- data_threshold_minutes: How recent data must be to be considered healthy (default: 5)
- restart_enabled: Whether to auto-restart failed services (default: true)
- alert_enabled: Whether to send webhook alerts (default: false)
- alert_webhook_url: URL for webhook notifications
- request_timeout_seconds: Elasticsearch query timeout (default: 30)
- retry_attempts: Number of retry attempts for failed queries (default: 3)
- retry_delay_seconds: Delay between retry attempts (default: 5)
- last_check: Timestamp of last check
- last_success: Timestamp of last successful check
- last_failure: Timestamp of last failed check
- failure_count: Consecutive failure count
- whatsapp_alerts_enabled: Whether WhatsApp alerts are enabled for this service
- last_whatsapp_alert: Timestamp of last WhatsApp alert sent
- created_at: Record creation timestamp
- updated_at: Record update timestamp
```

#### whatsapp_config (monitor_whatsapp.py only)
```sql
- id: Primary key
- enabled: Whether WhatsApp alerts are globally enabled (default: false)
- api_key: MessageBird API key
- channel_id: WhatsApp Business channel ID
- namespace: WhatsApp template namespace
- template_name: Template name for alerts (default: 'server_alerts')
- language_code: Language code for templates (default: 'es_MX')
- alert_interval_minutes: Minimum time between alerts (default: 30)
- min_failure_count: Minimum failures before alerting (default: 1)
- created_at: Record creation timestamp
- updated_at: Record update timestamp
```

#### whatsapp_phones (monitor_whatsapp.py only)
```sql
- id: Primary key
- phone_number: Phone number for alerts (unique)
- contact_name: Contact identifier
- enabled: Whether this number receives alerts (default: true)
- created_at: Record creation timestamp
- updated_at: Record update timestamp
```

## Configuration

### Environment Setup

Both monitors require:
- MySQL database connection (configured in code)
- Access to Elasticsearch/OpenSearch endpoints
- kubectl access for service restarts (optional)
- MessageBird API credentials (for WhatsApp alerts)

### Default Settings

- **Check Interval**: 5 minutes
- **Data Threshold**: 5 minutes (data must be newer than 5 minutes)
- **Request Timeout**: 30 seconds
- **Retry Attempts**: 3
- **Retry Delay**: 5 seconds
- **Auto-restart**: Enabled by default

## Usage

### Basic Monitoring

```bash
# Run continuous monitoring
python monitor.py

# Run single cycle for testing
python monitor.py --once

# Debug mode (detailed logging)
python monitor.py -d 1

# Elasticsearch-only debug (shows queries)
python monitor.py -d 2

# Minimal output (problems only)
python monitor.py -d 3
```

### WhatsApp Monitoring

```bash
# Run continuous monitoring with WhatsApp alerts
python monitor_whatsapp.py

# Force alerts (ignore timing restrictions)
python monitor_whatsapp.py --force

# Single cycle with WhatsApp alerts
python monitor_whatsapp.py --once --force

# Debug modes work the same as basic monitor
python monitor_whatsapp.py -d 2 --force
```

### Command Line Options

- `-d, --debug [LEVEL]`: Debug mode (1=full, 2=ES queries only, 3=problems only)
- `--once`: Run single monitoring cycle and exit
- `--force`: Force WhatsApp alerts (monitor_whatsapp.py only)

## Alert Types

### Webhook Alerts
Sent when `alert_enabled=true` and `alert_webhook_url` is configured. Payload:
```json
{
  "text": "Alert message with service details",
  "timestamp": "2025-11-07T12:00:00.000Z"
}
```

### WhatsApp Alerts
Sent using MessageBird API with configured templates. Messages include:
- Service type and document name
- Client namespace
- Error description
- Timestamp

## Troubleshooting

### Common Issues

1. **Database Connection Errors**
   - Verify MySQL credentials in code
   - Check database server connectivity

2. **Elasticsearch Connection Issues**
   - Verify Elasticsearch URL and credentials
   - Check network connectivity and SSL certificates
   - Review index names and permissions

3. **Kubernetes Access Problems**
   - Ensure kubectl is configured and authenticated
   - Verify namespace permissions for service restarts

4. **WhatsApp Alert Failures**
   - Verify MessageBird API credentials
   - Check WhatsApp Business account status
   - Confirm template approval and namespace

### Debug Modes

- **Level 1**: Full debug logging with all details
- **Level 2**: Elasticsearch queries only (good for API debugging)
- **Level 3**: Minimal output showing only problems (good for production)

### Log Files

Logs are written to `logs/elasticsearch_monitor.log` with rotation and both file and console output.

## Dependencies

- `mysql-connector-python`: MySQL database access
- `requests`: HTTP client for Elasticsearch and webhooks
- `urllib3`: HTTP library with SSL handling
- `messagebird`: WhatsApp API client (monitor_whatsapp.py only)

## Integration with JonoBridge

The monitoring system integrates with the main JonoBridge application:

- **Web Configuration Interface**: Monitor settings can be configured through the JonoBridge web interface at `/admin/monitor`
- **Database Sharing**: Uses the same MySQL database as the main application
- **Service Discovery**: Automatically monitors services configured through the JonoBridge web interface
- **Kubernetes Integration**: Works with services deployed by JonoBridge
- **Alert Coordination**: Can trigger alerts based on service health from the web UI

### Web Interface Configuration

The JonoBridge admin panel provides a user-friendly interface to configure:

- **Global Monitor Settings**: Check intervals, data thresholds, timeouts
- **WhatsApp Alert Configuration**: API keys, templates, phone numbers
- **Per-Service Settings**: Individual service monitoring parameters
- **Alert Management**: Enable/disable alerts and configure webhook URLs

## Best Practices

1. **Test Configurations**: Use `--once` and `-d 1` for initial testing
2. **Monitor Logs**: Regularly check log files for issues
3. **Alert Testing**: Use `--force` to test WhatsApp alerts without waiting
4. **Resource Usage**: Monitor system resources as the monitor runs continuously
5. **Backup Configurations**: Regularly backup monitoring configurations

## Production Deployment

### Running as a Service

For production environments, run the monitors as background services:

```bash
# Basic monitor
nohup python monitor.py > monitor.log 2>&1 &

# WhatsApp monitor
nohup python monitor_whatsapp.py > monitor_whatsapp.log 2>&1 &
```

### Process Management

Use process managers like `systemd` or `supervisor` for production:

**systemd service example (/etc/systemd/system/jonobridge-monitor.service):**
```ini
[Unit]
Description=JonoBridge Monitor Service
After=network.target

[Service]
Type=simple
User=jonobridge
WorkingDirectory=/home/ubuntu/jonobridge/frontend
ExecStart=/usr/bin/python3 monitor_whatsapp.py
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Management commands:**
```bash
sudo systemctl enable jonobridge-monitor
sudo systemctl start jonobridge-monitor
sudo systemctl status jonobridge-monitor
sudo systemctl stop jonobridge-monitor
```

### Log Rotation

Configure logrotate for monitor logs:

**/etc/logrotate.d/jonobridge-monitor:**
```
/home/ubuntu/jonobridge/frontend/logs/*.log {
    daily
    rotate 7
    compress
    missingok
    notifempty
    create 644 jonobridge jonobridge
}
```</content>
<parameter name="filePath">/home/ubuntu/jonobridge/frontend/README_monitor.md