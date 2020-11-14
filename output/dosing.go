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
var gpioChip string
var timezone string

func InitDosing(commandTopic string, stateTopic string, gpioChp string, tz string) {
	log.Printf("Initialising dosing pump module. Command topic: %s, state topic: %s, GPIO chip: %s", commandTopic, stateTopic, gpioChp)
	mqttSub(commandTopic)
	stateAckTopic = stateTopic
	gpioChip = gpioChp
	timezone = tz
}

func parseIncomingMessage(_ mqtt.Client, msg mqtt.Message) {
	dosingMessage := parseJsonMessage(msg.Payload())
	log.Printf("Processing incoming dosing message: %s", string(msg.Payload()))
	actuatePump(dosingMessage)
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatalf("Invalid Timezone: %s", timezone)
	}
	t := time.Now().In(loc)
	mqttBackend.Publish(fmt.Sprintf("%s%d", stateAckTopic, dosingMessage.Pump), t.Format("2006-01-02 15:04:05"))
}

func actuatePump(message dosingMessage) {
	pumpPin, _ := rpi.Pin(strconv.Itoa(message.Pump))
	log.Printf("Starting pump: %d, seconds: %d", pumpPin, message.Seconds)

	chip, err := gpiod.NewChip(gpioChip)
	defer chip.Close()
	if err != nil {
		log.Fatalf("Unable to connect to GPIO chip: %s - %s", gpioChip, err)
	}

	line, err := chip.RequestLine(pumpPin, gpiod.AsOutput(0))
	if err != nil {
		log.Fatalf("Unable to request line, GPIO pin: %d - %s", pumpPin, err)
	}

	err = line.SetValue(0)
	if err != nil {
		log.Fatalf("Unable to send message to pump, GPIO pin: %d - %s", pumpPin, err)
	}

	durationToActuate, _ := time.ParseDuration(fmt.Sprintf("%ds", message.Seconds))
	time.AfterFunc(durationToActuate, func() {
		err := line.SetValue(1)
		if err != nil {
			log.Fatalf("Unable to send message to pump, GPIO pin: %d - %s", pumpPin, err)
		}
		log.Printf("Finished dosing pump: %d", pumpPin)
		line.Close()
	})
}

func mqttSub(topic string) {
	token := mqttBackend.Subscribe(topic, parseIncomingMessage)
	if token.Wait() && token.Error() != nil {
		log.Fatalf("Unable to subscribe. Topic: %s, error: %s", topic, token.Error())
	}
}

func parseJsonMessage (message []byte) dosingMessage {
	var dosingMessage dosingMessage
	err := json.Unmarshal(message, &dosingMessage)
	if err != nil {
		log.Fatalf("Unable to parse dosing message from MQTT: %s", string(message))
	}
	return dosingMessage
}

type dosingMessage struct {
	Pump int `json:"pump"`
	Seconds int `json:"seconds"`
}