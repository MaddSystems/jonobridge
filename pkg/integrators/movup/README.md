# Movup 

## Testing Port
Test port 8501 jonobridge.dwim.mx
Client: Logistic ARG

### Plates URL
```
export PLATES_URL="https://pluto.dudewhereismy.com.mx/imei/search?appId=2826"
export ELASTIC_DOC_NAME="logisticsarg_movup"
export MOVUP_USER="Arguelles"
export MOVUP_USERKEY="ERhcowRY5ofl"
export MOVUP_URL_PATH="roshfrans-mx"
export MOVUP_TOKEN="Basic QXJndWVsbGVzOkVSaGNvd1JZNW9mbA=="
export MOVUP_PROVIDER="GPScontrol"
```

Values for Chep - Almanza Business
```
https://bridge.dudewhereismy.mx/fields/autenticValues?port=8667
{
   "autentic_values":[
      {
         "value":"Arguelles",
         "key":"user"
      },
      {
         "value":"ERhcowRY5ofl",
         "key":"user_key"
      },
      {
         "value":"roshfrans-mx",
         "key":"url_path"
      },
      {
         "value":"Basic QXJndWVsbGVzOkVSaGNvd1JZNW9mbA==",
         "key":"token"
      },
      {
         "value":"GPScontrol",
         "key":"provider"
      }
   ]
}
```

Config port:
```
https://bridge.dudewhereismy.mx/fields/rutineConfig?port=8667

{
   "localaddress":[
      "l",
      "138.68.54.137",
      ":8667",
      "Local address",
      "meitrack"
   ],
   "protocol":"meitrack",
   "transport":"udp",
   "integratorRemote":[
      {
         "name":"Unigis",
         "configuration":[
            "90",
            "http://unigis2.unisolutions.com.ar",
            "Unigis",
            ""
         ]
      },
      {
         "name":"RecursoConfiable",
         "configuration":[
            "90",
            "http://gps.rcontrol.com.mx/Tracking/wcf/RCService.svc?wsdl/IRCService",
            "RecursoConfiable",
            ""
         ]
      },
      {
         "name":"MovUp",
         "configuration":[
            "90",
            "http://gps.rcontrol.com.mx/Tracking/wcf/RCService.svc?wsdl/IRCService",
            "MovUp",
            ""
         ]
      }
   ]
}

```

Test Equipment
```
eco=31AK7K
IMEI=864507038426002
PLACAS=LB88279
```

find . -type f -exec sed -i -e 's/MOVUP/AVOCADOCLOUD/g' -e 's/movup/avocadocloud/g' -e 's/Movup/Avocadocloud/g' {} \;
