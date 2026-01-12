# Telemetry 

## Testing Port
Test port 8501 jonobridge.dwim.mx
Client:

### Plates URL
```
export PLATES_URL="https://pluto.dudewhereismy.com.mx/imei/search?appId=97"
export ELASTIC_DOC_NAME="telemetry"
export TELEMETRY_URL="https://telemetry.europe.ghtrack.com:8392/data/tracking_data/v1/devices?user=gpscontrolmadd&password=qQJMCjXYxgLQPHvK" 
export TELEMETRY_OWNER_ID="CORPORATIVO_HALCONES_CONTINENTAL"
```


Test Equipment
```

```

find . -type f -exec sed -i -e 's/TELEMETRY/AVOCADOCLOUD/g' -e 's/telemetry/avocadocloud/g' -e 's/Telemetry/Avocadocloud/g' {} \;
