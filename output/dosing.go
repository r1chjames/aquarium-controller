package output

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"gitlab.com/r1chjames/aquarium-controller/mqttBackend"
	"log"
	"strconv"
	"time"
)

var stateAckTopic string
var gpioDirectory string

func InitDosing(commandTopic string, stateTopic string, gpioDir string) {
	log.Print(fmt.Sprintf("Initialising dosing pump module. Command topic: %s, state topic: %s", commandTopic, stateTopic))
	mqttSub(commandTopic)
	stateAckTopic = stateTopic
	gpioDirectory = gpioDir
}

func parseIncomingMessage(_ mqtt.Client, msg mqtt.Message) {
	dosingMessage := parseJsonMessage(msg.Payload())
	log.Print(fmt.Sprintf("Processing incoming dosing message: %s", string(msg.Payload())))
	actuatePump(dosingMessage)
	message := time.Now().String()
	mqttBackend.Publish(fmt.Sprintf("%s%d", stateAckTopic, dosingMessage.Pump), message)
}

func actuatePump(message dosingMessage) {
	pumpPin, _ := rpi.Pin(strconv.Itoa(message.Pump))
	durationToActuate, _ := time.ParseDuration(fmt.Sprintf("%ds", message.Seconds))
	log.Print(fmt.Sprintf("Starting pump: %d, seconds: %d", pumpPin, durationToActuate))

	chip, err := gpiod.NewChip(gpioDirectory)
	defer chip.Close()
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect to GPIO chip: %s - %s", gpioDirectory, err))
	}

	line, err := chip.RequestLine(pumpPin, gpiod.AsOutput(1))
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to send message to pump, GPIO pin: %d - %s", pumpPin, err))
	}

	time.AfterFunc(durationToActuate, func() {
		err := line.SetValue(0)
		if err != nil {
			log.Fatal(fmt.Sprintf("Unable to send message to pump, GPIO pin: %d - %s", pumpPin, err))
		}
	})
	line.Close()
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