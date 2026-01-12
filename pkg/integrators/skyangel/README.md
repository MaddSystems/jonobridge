# Skyyangel 

## Testing Port
Test port 8501 jonobridge.dwim.mx
Client: Transportes peninsular

### Plates URL
```
export SKYANGEL_USER="gpscontrol"
export SKYANGEL_KEY="skygerenciador"
export ELASTIC_DOC_NAME="transportes_peninsular_skyangel"
export PLATES_URL="https://pluto.dudewhereismy.com.mx/imei/search?appId=2911"
```

Values for Chep - Almanza Business
```
https://bridge.dudewhereismy.mx/fields/autenticValues?port=8638
{
   "autentic_values":[
      {
         "value":"gpscontrol",
         "key":"user"
      },
      {
         "value":"skygerenciador",
         "key":"user_key"
      }
   ]
}
```

Config port:
```
https://bridge.dudewhereismy.mx/fields/rutineConfig?port=8638
{
   "localaddress":[
      "l",
      "138.68.54.137",
      ":8638",
      "Local address",
      "meitrack"
   ],
   "protocol":"meitrack",
   "transport":"udp",
   "integratorRemote":[
      {
         "name":"SkyAngel",
         "configuration":[
            "60",
            "http://ws.skyangel.com.mx/Gerenciador.php",
            "",
            ""
         ]
      }
   ]
}
```

Test Equipment
```

```

find . -type f -exec sed -i -e 's/SKYANGEL/AVOCADOCLOUD/g' -e 's/skyangel/avocadocloud/g' -e 's/Skyyangel/Avocadocloud/g' {} \;
