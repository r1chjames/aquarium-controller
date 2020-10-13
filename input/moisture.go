package input

import (
	"fmt"
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
	log.Print(fmt.Sprintf("Initialising moisture sensor module. State topic: %s, read interval: %s", stateTopic, readInterval))
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
		log.Fatal(fmt.Sprintf("Unable to connect to GPIO chip: %s - %s", gpioChip, err))
	}

	line, err := chip.RequestLine(sensorGpioPin, gpiod.AsInput)
	defer line.Close()
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to request line, GPIO pin: %d - %s", sensorGpioPin, err))
	}

	val, err := line.Value()
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to read moisture - %s", err))
	}

	return val
}