package mqttBackend

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"net/url"
	"time"
)

var mqttClient mqtt.Client

func Connect(clientId string, uri *url.URL) {
	opts := createClientOptions(clientId, uri)
	log.Print(opts)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	mqttClient = client
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
	return mqttClient.Subscribe(topic, 0, callback)
}

func Publish(topic string, message string) {
	token := mqttClient.Publish(topic, 0, false, message)
	if token.Error() != nil {
		log.Fatal(fmt.Sprintf("Unable to publish ack message. Topic: %s, message: %s, error: %s", topic, message, token.Error()))
	}
}