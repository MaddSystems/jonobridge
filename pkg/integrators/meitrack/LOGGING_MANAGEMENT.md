# Meitrack Logging Management Guide

## Logging Levels

### Default (Production) Mode
By default, the application runs in **quiet mode** with minimal logging:
- Only errors and critical status messages
- Summary reports every 5 minutes
- Connection failures and important events

### Verbose Mode
Set `MEITRACK_VERBOSE_LOGGING=Y` to enable detailed logging:
- Individual message processing logs
- Connection details for each server
- Message forwarding details
- All debug information

## Environment Variable Configuration

```bash
# Disable verbose logging (DEFAULT - recommended for production)
MEITRACK_VERBOSE_LOGGING=N

# Enable verbose logging (for debugging only)
MEITRACK_VERBOSE_LOGGING=Y
```

## Log Volume Comparison

### Production Mode (MEITRACK_VERBOSE_LOGGING=N)
```
2025/07/08 18:30:00 Successfully connected to the MQTT broker
2025/07/08 18:30:00 Subscribed to topic: tracker/jonoprotocol
2025/07/08 18:35:00 Messages processed: 245 total, 245 in last 5 minutes
2025/07/08 18:40:00 Messages processed: 523 total, 278 in last 5 minutes
```
**Log Rate**: ~4 lines per 5 minutes = **~1 MB per day**

### Verbose Mode (MEITRACK_VERBOSE_LOGGING=Y)
```
2025/07/08 18:26:31 Processing message #4562 on topic tracker/from-tcp
2025/07/08 18:26:31 Posting ADAS event
2025/07/08 18:26:31 Posting to server: server1.gpscontrol.com.mx:8500
2025/07/08 18:26:31 Successfully connected to server1.gpscontrol.com.mx:8500
2025/07/08 18:26:31 Sent 156 bytes to server1.gpscontrol.com.mx:8500
2025/07/08 18:26:31 Posting:$$Q128,866811062034643,CCE,...
```
**Log Rate**: ~6 lines per message = **~50-100 MB per day** (depending on message volume)

## Message Counter Behavior

### Counter Growth
```go
var messageCounter int64  // 64-bit signed integer
```

- **Maximum Value**: 9,223,372,036,854,775,807
- **At 1 msg/sec**: 292 billion years to overflow
- **At 1000 msg/sec**: 292 million years to overflow
- **Practical Impact**: Counter will never overflow in real-world usage

### Counter Display
- Shows cumulative count since application start
- Resets to 0 on application restart
- Used for monitoring message processing rates

## Log Management Strategies

### 1. Log Rotation (Recommended)

**Using systemd service with journald:**
```bash
# /etc/systemd/system/meitrack.service
[Unit]
Description=Meitrack Integrator
After=network.target

[Service]
Type=simple
User=meitrack
WorkingDirectory=/opt/meitrack
ExecStart=/opt/meitrack/meitrack
Restart=always
RestartSec=5

# Log management
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

**Configure journald rotation:**
```bash
# /etc/systemd/journald.conf
[Journal]
SystemMaxUse=1G          # Maximum disk space for all logs
SystemMaxFileSize=100M   # Maximum size per log file
MaxRetentionSec=7d       # Keep logs for 7 days
ForwardToSyslog=no       # Don't duplicate to syslog
```

### 2. Docker Logging

**For Docker deployments:**
```bash
docker run -d \
  --name meitrack \
  --log-driver=json-file \
  --log-opt max-size=100m \
  --log-opt max-file=3 \
  -e MEITRACK_VERBOSE_LOGGING=N \
  maddsystems/meitrack:1.1.0
```

### 3. External Log Aggregation

**Ship logs to external system:**
```bash
# Use fluentd, logstash, or similar to ship logs
# Configure with size limits and filtering
```

### 4. File-based Logging with Rotation

**Using logrotate:**
```bash
# /etc/logrotate.d/meitrack
/var/log/meitrack/*.log {
    daily
    missingok
    rotate 7
    compress
    delaycompress
    notifempty
    create 644 meitrack meitrack
    postrotate
        systemctl reload meitrack || true
    endscript
}
```

## Monitoring Commands

### Check Current Log Size
```bash
# For systemd/journald
journalctl -u meitrack --disk-usage

# For file-based logs
du -h /var/log/meitrack/

# For Docker
docker logs meitrack 2>/dev/null | wc -l
```

### Monitor Log Growth Rate
```bash
# Watch log growth in real-time
journalctl -u meitrack -f

# Check message processing rate
journalctl -u meitrack | grep "Messages processed" | tail -5
```

### Clear Old Logs
```bash
# Clear journald logs older than 2 days
journalctl --vacuum-time=2d

# Clear Docker logs
docker logs meitrack > /dev/null 2>&1 && docker restart meitrack
```

## Troubleshooting High Log Volume

### Symptoms
- Disk space filling rapidly
- Performance degradation due to I/O
- Log files growing > 100MB/day

### Solutions

1. **Disable Verbose Logging**
```bash
export MEITRACK_VERBOSE_LOGGING=N
# or remove the environment variable entirely
```

2. **Implement Log Filtering**
```bash
# Only log errors and status (if using custom logging)
export LOG_LEVEL=ERROR
```

3. **Increase Log Rotation Frequency**
```bash
# Rotate logs hourly instead of daily
# Update logrotate configuration
```

4. **Monitor Semaphore Issues**
The "context deadline exceeded" error suggests connection pool exhaustion. Consider:
```go
// Increase semaphore size if needed
var connectionSemaphore = semaphore.NewWeighted(20) // Instead of 10
```

## Recommended Production Settings

```bash
# Minimal logging for production
MEITRACK_VERBOSE_LOGGING=N

# Standard environment variables
MEITRACK_HOST=server1.gpscontrol.com.mx:8500
MQTT_BROKER_HOST=your-mqtt-broker
MEITRACK_FWD_ONLY_ADAS=Y  # If you only want ADAS events
SEND_TO_ELASTIC=N         # Unless specifically needed

# System-level log management
# - Enable log rotation (journald or logrotate)
# - Set maximum log retention (7 days)
# - Monitor disk usage alerts
```

## Performance Impact

### Low Verbosity (Production)
- **CPU Impact**: Minimal (~1% overhead)
- **Disk I/O**: ~1 MB/day
- **Memory**: Negligible

### High Verbosity (Debug)
- **CPU Impact**: 5-10% overhead (string formatting)
- **Disk I/O**: 50-100 MB/day
- **Memory**: Minimal increase

**Recommendation**: Always use `MEITRACK_VERBOSE_LOGGING=N` in production unless actively debugging issues.
