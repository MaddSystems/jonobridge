# Logitrack 

## Testing Port
Test port 8505 jonobridge.dwim.mx
Client: Siempre a un click

### Environment variables from the URL
```
export PLATFORM_HOST="ms03.trackermexico.com.mx:10003"
export ELASTIC_DOC_NAME="logitrack"
export LOGITRACK_USER="86lLrlD6Ka2/GPo0MiAK1Tw.A/Z0ffPNl2RrUI9D"
export LOGITRACK_USER_KEY="\$2y\$10\$IZM28dkBOYwkIqptdD77OOgbhTBec5kpv"
export LOGITRACK_URLWAY="https://gps-homologations.logitrack.mx/integrations/api/v1"
```

### Account
ms03.trackermexico.com.mx
Soportetracker
0000
IMEI 867869061142014


### Values for Chep - Almanza Business
```
https://bridge.dudewhereismy.mx/fields/autenticValues?port=8611
{"autentic_values": [{"value": "GPScontrol", "key": "proveedor"}]}
```

### Config port:
```
https://bridge.dudewhereismy.mx/fields/rutineConfig?port=8611
{"localaddress": ["l", "138.68.54.137", ":8611", "Local address", "meitrack"], "protocol": "meitrack", "transport": "udp", "integratorRemote": [{"name": "AltoTrack", "configuration": ["33", "http://ws4.logitrack.com/WSPosiciones_WalmartMX/WSPosiciones_WalmartMX.svc?wsdl/IServicePositions", "", "meitrack"]}]}
```

### Test Equipment
```
eco=31AK7K
IMEI=867869061142014
PLACAS=
```

find . -type f -exec sed -i -e 's/LOGITRACK/AVOCADOCLOUD/g' -e 's/logitrack/avocadocloud/g' -e 's/Logitrack/Avocadocloud/g' {} \;
