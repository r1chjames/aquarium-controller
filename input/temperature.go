package input

import (
	"fmt"
	"github.com/yryz/ds18b20"
	"gitlab.com/r1chjames/aquarium-controller/mqttBackend"
	"gitlab.com/r1chjames/aquarium-controller/utils"
	"log"
	"strconv"
	"time"
)

var temperatureStateTopic string

func InitTemperature(stateTopic string, readInterval time.Duration) {
	temperatureStateTopic = stateTopic
	utils.DoEvery(readInterval, processTemperature)
}

func processTemperature() {
	sensors, err := ds18b20.Sensors()
	if err != nil {
		log.Fatal("No temperature sensors found")
	}

	log.Print(fmt.Sprintf("Sensor IDs found: %v\n", sensors))

	for _, sensor := range sensors {
		value, err := ds18b20.Temperature(sensor)
		if err != nil {
			log.Fatal("Unable to read temperature from sensor")
		}
		log.Print(fmt.Sprintf("Sensor: %s temperature: %.2f°C\n", sensor, value))
		mqttBackend.Publish(temperatureStateTopic, strconv.FormatFloat(value, 'E', -1, 64))
	}
}
