package main

// Message represents the structure of the JSON data received from the MQTT topic.
// Using a map for dynamic keys.
type Message map[string]interface{}
