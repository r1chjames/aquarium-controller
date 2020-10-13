package input

import (
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"gitlab.com/r1chjames/aquarium-controller/mqttBackend"
	"gitlab.com/r1chjames/aquarium-controller/utils"
	"log"
	"strconv"
	"time"
)

var moistureStateTopic string
var gpioChip string
var sensorPin string

func InitMoisture(stateTopic string, readInterval time.Duration, gpioChp string, gpioPin string) {
	log.Printf("Initialising moisture sensor module. State topic: %s, read interval: %s, GPIO chip: %s,  GPIO pin: %s", stateTopic, readInterval, gpioChip, gpioPin)
	moistureStateTopic = stateTopic
	gpioChip = gpioChp
	sensorPin = gpioPin
	utils.DoEvery(readInterval, processMoisture)
}

func processMoisture() {
	moistureValue := readRawMoisture()
	mqttBackend.Publish(moistureStateTopic, strconv.Itoa(moistureValue))
}

func readRawMoisture() int {
	sensorGpioPin, _ := rpi.Pin(sensorPin)

	chip, err := gpiod.NewChip(gpioChip)
	defer chip.Close()
	if err != nil {
		log.Fatalf("Unable to connect to GPIO chip: %s - %s", gpioChip, err)
	}

	line, err := chip.RequestLine(sensorGpioPin, gpiod.AsInput)
	defer line.Close()
	if err != nil {
		log.Fatalf("Unable to request line, GPIO pin: %d - %s", sensorGpioPin, err)
	}

	value, err := line.Value()
	if err != nil {
		log.Fatalf("Unable to read moisture - %s", err)
	}

	log.Printf("Moisture: %d%%", value)
	return value
}