module activetrack

go 1.23.2

require (
	github.com/MaddSystems/jonobridge/common v0.0.0-00010101000000-000000000000
	github.com/eclipse/paho.mqtt.golang v1.5.0
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
)

// Replace directive pointing to the local common module
replace github.com/MaddSystems/jonobridge/common => ../../common
