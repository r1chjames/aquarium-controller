package output

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/warthog618/gpiod"
	"gitlab.com/r1chjames/aquarium-controller/mqttBackend"
	"log"
	"time"
)

var stateAckTopic string

func InitDosing(commandTopic string, stateTopic string) {
	log.Print(fmt.Sprintf("Initialising dosing pump module. Command topic: %s, state topic: %s", commandTopic, stateTopic))
	mqttSub(commandTopic)
	stateAckTopic = stateTopic
}

func parseIncomingMessage(_ mqtt.Client, msg mqtt.Message) {
	dosingMessage := parseJsonMessage(msg.Payload())
	log.Print(fmt.Sprintf("Processing incoming dosing message: %s", string(msg.Payload())))
	actuatePump(dosingMessage)
	message := time.Now().String()
	mqttBackend.Publish(fmt.Sprintf("%s%d", stateAckTopic, dosingMessage.Pump), message)
}

func actuatePump(message dosingMessage) {
	log.Print(fmt.Sprintf("Starting pump: %d, seconds: %d", message.Pump, message.Seconds))

	chip, err := gpiod.NewChip("gpiomem")
	defer chip.Close()
	if err != nil {
		log.Fatal("Unable to connect to GPIO")
	}

	line, err := chip.RequestLine(message.Pump, gpiod.AsOutput(1))
	defer line.Close()
	if err != nil {
		log.Fatal("Unable to send message to pump")
	}

	durationToActuate, _ := time.ParseDuration(fmt.Sprintf("%ds", message.Seconds))
	time.AfterFunc(durationToActuate, func() {
		err := line.SetValue(0)
		if err != nil {
			log.Fatal("Unable to send message to pump")
		}
	})
}

func mqttSub(topic string) {
	token := mqttBackend.Subscribe(topic, parseIncomingMessage)
	if token.Error() != nil {
		log.Fatal(fmt.Sprintf("Unable to subscribe. Topic: %s, error: %s", topic, token.Error()))
	}
}

func parseJsonMessage (message []byte) dosingMessage {
	var dosingMessage dosingMessage
	err := json.Unmarshal(message, &dosingMessage)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to parse dosing message from MQTT: %s", string(message)))
	}
	return dosingMessage
}

type dosingMessage struct {
	Pump int `json:"pump"`
	Seconds int `json:"seconds"`
}