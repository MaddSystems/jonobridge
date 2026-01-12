# Avocadocloud 

## Testing Port
Test port 8501 jonobridge.dwim.mx
Client: Almanza Business - Chep

### Plates URL
```
export PLATES_URL="https://pluto.dudewhereismy.com.mx/imei/search?appId=2834"
export ELASTIC_DOC_NAME="DafneGarrido_avocadocloud"
export AVOCADO_USER="phoenix"
export AVOCADO_PASSWORD="Ph03niX-2018_"
export AVOCADO_USER_ADM="2135"
export AVOCADO_URL="https://cerberusenlinea.com/WEB_SERVICE_PHOENIX_CLOUD_PRD-1.0/PH_PHOENIX_CLOUD_PRD_v01?wsdl/recibirEventosGPS"
```

Values for Dafne Garrido
```
https://bridge.dudewhereismy.mx/fields/autenticValues?port=8665

{"autentic_values": [{"value": "phoenix", "key": "user"}, {"value": "Ph03niX-2018_", "key": "user_key"}, {"value": "2135", "key": "adm"}, {"value": "https://cerberusenlinea.com/WEB_SERVICE_PHOENIX_CLOUD_PRD-1.0/PH_PHOENIX_CLOUD_PRD_v01?wsdl/recibirEventosGPS", "key": "url"}]}
```

Config port:
```
Configuración
https://bridge.dudewhereismy.mx/fields/rutineConfig?port=8665

{"localaddress": ["l", "138.68.54.137", ":8665", "Local address", "meitrack"], "protocol": "meitrack", "transport": "udp", "integratorRemote": [{"name": "IntegratorAvocadoCloud", "configuration": ["88", "https://cerberusenlinea.com/WEB_SERVICE_PHOENIX_CLOUD_PRD-1.0/PH_PHOENIX_CLOUD_PRD_v01?WSDL", "", ""]}, {"name": "Motumcloud", "configuration": ["88", "https://positions.apis.motumcloud.com/v2/positions", "", ""]}]}

```

Test Equipment
```
eco=84AGY
IMEI=864292049109851
PLACAS=84AGY
```
