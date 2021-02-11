package mqttBackend

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/jasonlvhit/gocron"
	"log"
	"net/url"
)

var mqttClient mqtt.Client

func Connect(clientId string, uri *url.URL) {
	opts := createClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	mqttClient = client
	go startConnectionMonitor()
}

func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}

func Subscribe(topic string, callback func(client mqtt.Client, msg mqtt.Message)) mqtt.Token {
	return mqttClient.Subscribe(topic, 2, callback)
}

func IsConnected() bool {
	connected := mqttClient.IsConnected()
	if !connected {
		log.Fatal("MQTT disconnected")
	}

	return connected
}

func startConnectionMonitor() {
	gocron.Every(5).Minutes().Do(IsConnected)
	<- gocron.Start()
}

func Publish(topic string, message string, retain bool) {
	token := mqttClient.Publish(topic, 0, retain, message)
	if token.Error() != nil {
		log.Fatalf("Unable to publish ack message. Topic: %s, message: %s, error: %s", topic, message, token.Error())
	}
}