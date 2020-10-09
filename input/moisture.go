package input

import (
	"github.com/warthog618/gpiod/device/rpi"
	"gitlab.com/r1chjames/go-aquarium-sensors/mqttBackend"
	"gitlab.com/r1chjames/go-aquarium-sensors/utils"
	"log"
	"time"
	"github.com/warthog618/gpiod"
)

var moistureStateTopic string

func InitMoisture(stateTopic string, readInterval time.Duration) {
	stateTopic = stateTopic
	utils.DoEvery(readInterval, processMoisture)
}

func processMoisture() {
	moistureValue := readRawMoisture()
	mqttBackend.Publish(moistureStateTopic, string(rune(moistureValue)))
}

func readRawMoisture() int {
	chip, err := gpiod.NewChip("gpiomem")
	defer chip.Close()
	if err != nil {
		log.Fatal("Unable to connect to GPIO")
	}

	line, err := chip.RequestLine(rpi.GPIO25, gpiod.AsInput)
	defer line.Close()
	if err != nil {
		log.Fatal("Unable to read moisture")
	}

	val, err := line.Value()
	if err != nil {
		log.Fatal("Unable to read moisture")
	}

	return val
}