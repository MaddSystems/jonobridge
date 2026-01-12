# Lobosoftware 

## Testing Port
Test port 8501 jonobridge.dwim.mx
Client: Lobosoftware

### Plates URL
```
export PLATES_URL="https://pluto.dudewhereismy.com.mx/imei/search?appId=3059"
export ELASTIC_DOC_NAME="lobosoftware"
export LOBOSOFTWARE_TOKEN_URL="https://pruebascampanaws.lobos.com.mx/campana/authorization/usergpsauthorization"
export LOBOSOFTWARE_URL="https://pruebascampanaws.lobos.com.mx/campana/gps/datagps"
export LOBOSOFTWARE_USER="GPSCONTROL"
export LOBOSOFTWARE_USER_KEY="GPSC0NTR0L.2024"
```

Values for Chep - Almanza Business
```
https://bridge.dudewhereismy.mx/fields/autenticValues?port=8692

{
  "autentic_values": [
    {
      "value": "GPSCONTROL",
      "key": "user"
    },
    {
      "value": "GPSC0NTR0L.2024",
      "key": "user_key"
    }
  ]
}
```

Config port:
```
https://bridge.dudewhereismy.mx/fields/rutineConfig?port=8692

{
  "localaddress": [
    "l",
    "138.68.54.137",
    ":8692",
    "Local address",
    "meitrack"
  ],
  "protocol": "meitrack",
  "transport": "udp",
  "integratorRemote": [
    {
      "name": "LoboSoftware",
      "configuration": [
        "109",
        "https://pruebascampanaws.lobos.com.mx/campana/authorization/usergpsauthorization",
        "",
        ""
      ]
    }
  ]
}
```

Test Equipment
```
IMEI=867869061350229
PLACAS=MT4426E
```

