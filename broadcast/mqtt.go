package broadcast

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

const mqttHost = "tcp://127.0.0.1:1883"
const mqttClientId = "pusher"
const mqttUsername = "guest"
const mqttPassword = "guest"

const maxPublishTryTimes = 3

var client mqtt.Client

func init()  {
	newClient()
}

func newClient() mqtt.Client {
	opts := mqtt.NewClientOptions().
		AddBroker(mqttHost).
		SetUsername(mqttUsername).
		SetPassword(mqttPassword).
		SetClientID(mqttClientId)
   // 需 加锁 判断, 保证 client 只 初始化 一次
	if client == nil {

		client = mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

	}

	return client

}

func publish(topic string, m interface{}, qos byte) bool {
	if token := client.Publish(topic, qos, false, m); token.Wait() && token.Error() != nil {

		if qos == 0 {
			log.Printf("publis error: %s\n", token.Error())
			return false
		} else {
			return reTryPublish(topic, m, qos)
		}
	} else {
		return true
	}
}

func reTryPublish(topic string, m interface{}, qos byte) bool {

	for i := 0; i < maxPublishTryTimes; i++ {
		select {
		case <-time.After(1 * time.Second):
			if token := client.Publish(topic, qos, false, m); token.Wait() && token.Error() == nil {
				return true
			}else {
				log.Printf("publish try %d times error: %s \n", i+1, token.Error())
			}
		}
	}

	return false

}
