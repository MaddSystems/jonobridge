# Controlt 

## Testing Port
Test port 8501 jonobridge.dwim.mx
Client: Transportes Rodira

### Plates URL
```
export PLATES_URL="https://pluto.dudewhereismy.com.mx/imei/search?appId=2812"
export CONTROLT_USER="gpscontrolmx" 
export CONTROLT_USER_KE="G$276Tv3$08" 
export CONTROLT_URL="http://controlt.net/APP/HUB/service.asmx?WSDL/InsertEventAndLogin"

	controlt_user = os.Getenv("CONTROLT_USER")
	if controlt_user == "" {
		controlt_user = "gpscontrolmx" // Default fallback
	}
	controlt_user_key = os.Getenv("CONTROLT_USER_KEY")
	if controlt_user_key == "" {
		controlt_user_key = "G$276Tv3$08" // Default fallback
	}
	controlt_url = os.Getenv("CONTROLT_URL")
	if controlt_url == "" {
		controlt_url = "http://controlt.net/APP/HUB/service.asmx?WSDL/InsertEventAndLogin" // Default fallback
	}
```


Test Equipment
```
eco=99AJ/Z
IMEI=864606048151541
PLACAS=
```

