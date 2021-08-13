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

var (
	stateAckTopic string
	gpioChip string
	location *time.Location
)

func InitDosing(commandTopic string, stateTopic string, gpioChp string, tz string) {
	log.Printf("Initialising dosing pump module. Command topic: %s, state topic: %s, GPIO chip: %s", commandTopic, stateTopic, gpioChp)
	stateAckTopic = stateTopic
	gpioChip = gpioChp
	parseTimezone(tz)
	mqttSub(commandTopic)
}

func parseTimezone(tz string) {
	timezone, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatalf("Invalid Timezone: %s", timezone)
	}
	location = timezone
}

func parseIncomingMessage(_ mqtt.Client, msg mqtt.Message) {
	log.Printf("Processing incoming dosing message: %s", string(msg.Payload()))
	go handleMessageAndAck(msg)
}

func handleMessageAndAck(msg mqtt.Message) {
	token := make(chan bool)
	dosingMessage := parseJsonMessage(msg.Payload())
	actuatePump(dosingMessage, token)
	<-token
	t := time.Now().In(location)
	mqttBackend.Publish(fmt.Sprintf("%s%d", stateAckTopic, dosingMessage.Pump), t.Format("2006-01-02 15:04:05"), true)
}

func actuatePump(message dosingMessage, token chan bool) {
	pumpPin, _ := rpi.Pin(strconv.Itoa(message.Pump))
	log.Printf("Starting pump: %d, seconds: %d", pumpPin, message.Seconds)

	chip, err := gpiod.NewChip(gpioChip)
	defer chip.Close()
	if err != nil {
		log.Fatalf("Unable to connect to GPIO chip: %s - %s", gpioChip, err)
	}

	line, err := chip.RequestLine(pumpPin, gpiod.AsOutput(1))
	if err != nil {
		log.Fatalf("Unable to request line, GPIO pin: %d - %s", pumpPin, err)
	}

	err = line.SetValue(1)
	if err != nil {
		log.Fatalf("Unable to send message to pump, GPIO pin: %d - %s", pumpPin, err)
	}

	durationToActuate, _ := time.ParseDuration(fmt.Sprintf("%ds", message.Seconds))
	time.AfterFunc(durationToActuate, func() {
		err := line.SetValue(0)
		if err != nil {
			log.Fatalf("Unable to send message to pump, GPIO pin: %d - %s", pumpPin, err)
		}
		log.Printf("Finished dosing pump: %d", pumpPin)
		_ = line.Close()
		token <- true
	})
}

func mqttSub(topic string) {
	mqttBackend.Subscribe(topic, parseIncomingMessage)
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