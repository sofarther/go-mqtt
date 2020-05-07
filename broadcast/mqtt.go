package broadcast

import mqtt "github.com/eclipse/paho.mqtt.golang"

const mqttHost = "tcp://127.0.0.1:1883"
const mqttClientId = "pusher"
const mqttUsername = "guest"
const mqttPassword = "guest"

var client mqtt.Client

func newClient() mqtt.Client {
	opts := mqtt.NewClientOptions().
		AddBroker(mqttHost).
		SetUsername(mqttUsername).
		SetPassword(mqttPassword).
		SetClientID(mqttClientId)

	if client == nil {
		client = mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

	}

	return client

}
