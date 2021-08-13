package mqttBackend

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/jasonlvhit/gocron"
	"log"
	"net/url"
	"os"
)

var mqttClient mqtt.Client

func Connect(clientId string, uri *url.URL) {
	setUpLogging()
	opts := createClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	mqttClient = client
	go startConnectionMonitor()
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

func setUpLogging() {
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
}

func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	opts.ConnectRetry = true
	opts.AutoReconnect = true
	opts.OnConnectionLost = func(cl mqtt.Client, err error) {
		log.Fatalf("Connection Lost: %s\n", err.Error())
	}
	opts.OnConnect = func(cl mqtt.Client) {
		log.Println("Connected")
	}
	opts.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) {
		log.Println("Attempting to reconnect")
	}
	return opts
}

func Subscribe(topic string, callback func(client mqtt.Client, msg mqtt.Message)) {
	token := mqttClient.Subscribe(topic, 2, callback)
	go func() {
		<-token.Done()
		if token.Error() != nil {
			log.Fatalf("Unable to subscribe. Topic: %s, error: %s", topic, token.Error())
		} else {
			log.Printf("Subscribed to: %s", topic)
		}
	}()
}

func IsConnected() bool {
	connected := mqttClient.IsConnected()
	if !connected {
		log.Fatal("MQTT disconnected")
	}
	return connected
}

func startConnectionMonitor() {
	_ = gocron.Every(5).Minutes().Do(IsConnected)
	<- gocron.Start()
}

func Publish(topic string, message string, retain bool) {
	token := mqttClient.Publish(topic, 0, retain, message)
	go func() {
		 <-token.Done()
		if token.Error() != nil {
			log.Fatalf("Unable to publish ack message. Topic: %s, message: %s, error: %s", topic, message, token.Error())
		}
	}()
}