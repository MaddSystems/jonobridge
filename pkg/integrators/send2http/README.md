# Send2http 

## Testing Port
Test port 8501 jonobridge.dwim.mx
Client: Almanza Business - Chep

### Plates URL
```
export PLATES_URL="https://pluto.dudewhereismy.com.mx/imei/search?appId=97"
```

Values for Chep - Almanza Business
```
https://bridge.dudewhereismy.mx/fields/autenticValues?port=8611
{"autentic_values": [{"value": "GPScontrol", "key": "proveedor"}]}
```

Config port:
```
https://bridge.dudewhereismy.mx/fields/rutineConfig?port=8611
{"localaddress": ["l", "138.68.54.137", ":8611", "Local address", "meitrack"], "protocol": "meitrack", "transport": "udp", "integratorRemote": [{"name": "AltoTrack", "configuration": ["33", "http://ws4.send2http.com/WSPosiciones_WalmartMX/WSPosiciones_WalmartMX.svc?wsdl/IServicePositions", "", "meitrack"]}]}
```

Test Equipment
```
eco=31AK7K
IMEI=864507038426002
PLACAS=LB88279
```

find . -type f -exec sed -i -e 's/SEND2HTTP/AVOCADOCLOUD/g' -e 's/send2http/avocadocloud/g' -e 's/Send2http/Avocadocloud/g' {} \;
