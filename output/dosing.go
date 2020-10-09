package output

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/mgutz/logxi/v1"
	"gitlab.com/r1chjames/go-aquarium-sensors/mqttBackend"
	"time"
)

var stateAckTopic string

func InitDosing(commandTopic string, stateTopic string) {
	mqttSub(commandTopic)
	stateAckTopic = stateTopic
}

func parseIncomingMessage(client mqtt.Client, msg mqtt.Message) {
	log.Info("Starting pump x")
	actuatePump(1, 2)
	message := time.Now().String()
	mqttBackend.Publish(stateAckTopic, message)
}

func actuatePump(pin int, lengthOfTime int) {
	// actuate
}

func mqttSub(topic string) {
	token := mqttBackend.Subscribe(topic, parseIncomingMessage)
	if token.Error() != nil {
		log.Fatal(fmt.Sprintf("Unable to subscribe. Topic: %s, error: %s", topic, token.Error()))
	}
}