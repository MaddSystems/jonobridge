# Altotrack 

## Testing Port
Test port 8501 jonobridge.dwim.mx
Client: Almanza Business - Chep

### Plates URL
```
export PLATES_URL="https://pluto.dudewhereismy.com.mx/imei/search?appId=97"
export ELASTIC_DOC_NAME="Almanza_Chep_AltroTrack"
export ALTOTRACK_URL="http://ws4.altotrack.com/WSPosiciones_Chep/WSPosiciones_Chep.svc?wsdl/IServicePositions"
export ALTOTRACK_PROVEEDOR="GPScontrol"
```

### Values from plates URL
```
      {
         "imei":"864507038426002",
         "status":" ",
         "ccid":"89520703100006313752",
         "model":"1",
         "color":"Blanco",
         "vin":"1",
         "year":"1",
         "eco":"31AK7K",
         "brand":"Dina",
         "device":"Meitrack T366",
         "tel":"528126494890",
         "phone":"+528126494890",
         "last_report":"2025-03-19 14:16:24",
         "motor":"",
         "plates":"LB88279"
      },
```

### Autentic Values for Chep - Almanza Business
```
https://bridge.dudewhereismy.mx/fields/autenticValues?port=8611
```

### Result
```
{
   "autentic_values":[
      {
         "value":"GPScontrol",
         "key":"proveedor"
      }
   ]
}
```

### Ports Config for port 8611:
```
https://bridge.dudewhereismy.mx/fields/rutineConfig?port=8611
{
   "localaddress":[
      "l",
      "138.68.54.137",
      ":8611",
      "Local address",
      "meitrack"
   ],
   "protocol":"meitrack",
   "transport":"udp",
   "integratorRemote":[
      {
         "name":"AltoTrack",
         "configuration":[
            "33",
            "http://ws4.altotrack.com/WSPosiciones_WalmartMX/WSPosiciones_WalmartMX.svc?wsdl/IServicePositions",
            "",
            "meitrack"
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

### Debug

kubectl get pods -n almanzachepaltrotrack
kubectl exec -it altotrack-6c4c9bbcb5-69d7w --namespace=almanzachepaltrotrack -- /bin/sh
kubectl logs altotrack-6c4c9bbcb5-nnkcr -n almanzachepaltrotrack --previous
