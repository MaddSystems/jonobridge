# Motumcloud 

export PLATES_URL="https://pluto.dudewhereismy.com.mx/imei/search?appId=2834"
export ELASTIC_DOC_NAME="DafneGarrido_Motumcloud"
export MOTUMCLOUD_USER="mcloud-dafgarri-gpsctrl-cfort@dafgarri-gpsctrl.com"
export MOTUMCLOUD_PASSWORD="Ap1CeCl0udAfG4rGp5Ctrl"
export MOTUMCLOUD_REFERER="mcloud-dafgarri-gpsctrl-cfort"
export MOTUMCLOUD_APIKEY="AIzaSyAgO6dk0ZKDO15_M6dJ7fClmcJ4_cHVm8c"
export MOTUMCLOUD_CARRIER="Dafne Garrido"

### Credentials:

https://bridge.dudewhereismy.mx/fields/autenticValues?port=8665&owner=motumcloud

{
    "autentic_values": [
        {
            "value": "mcloud-dafgarri-gpsctrl-cfort@dafgarri-gpsctrl.com",
            "key": "user"
        },
        {
            "value": "Ap1CeCl0udAfG4rGp5Ctrl",
            "key": "password"
        },
        {
            "value": "mcloud-dafgarri-gpsctrl-cfort",
            "key": "referer"
        },
        {
            "value": "AIzaSyAgO6dk0ZKDO15_M6dJ7fClmcJ4_cHVm8c",
            "key": "apikey"
        },
        {
            "value": "Dafne Garrido",
            "key": "carrier"
        }
    ]
}

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


https://bridge.dudewhereismy.mx/fields/autenticValues?port=8665&owner=Motumcloud


Config port:
```
https://bridge.dudewhereismy.mx/fields/rutineConfig?port=8611
{"localaddress": ["l", "138.68.54.137", ":8611", "Local address", "meitrack"], "protocol": "meitrack", "transport": "udp", "integratorRemote": [{"name": "AltoTrack", "configuration": ["33", "http://ws4.motumcloud.com/WSPosiciones_WalmartMX/WSPosiciones_WalmartMX.svc?wsdl/IServicePositions", "", "meitrack"]}]}
```

Test Equipment
```
eco=31AK7K
IMEI=864507038426002
PLACAS=LB88279
```

find . -type f -exec sed -i -e 's/MOTUMCLOUD/AVOCADOCLOUD/g' -e 's/motumcloud/avocadocloud/g' -e 's/Motumcloud/Avocadocloud/g' {} \;
