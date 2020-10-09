package main

import (
	"gitlab.com/r1chjames/go-aquarium-sensors/input"
	"gitlab.com/r1chjames/go-aquarium-sensors/mqttBackend"
	"gitlab.com/r1chjames/go-aquarium-sensors/output"
	"gitlab.com/r1chjames/go-aquarium-sensors/utils"
	"net/url"
	"time"
)

func main() {
	initMqtt()
	initInput()
	initOutput()
}

func initMqtt() {
	clientId := utils.GetEnv("CLIENT_ID", "aquarium_sensors")
	brokerUrl := utils.GetEnv("BROKER_URL", "broker")
	mqttBackend.Connect(clientId, &url.URL{Host: brokerUrl})
}

func initInput() {

	if utils.GetEnv("TEMP_SENSOR_ENABLED", "false") != "false"{
		tempStateTopic := utils.GetEnv("INPUT_TEMP_TOPIC", "temperature_topic")
		tempDuration, _ := time.ParseDuration(utils.GetEnv("INPUT_TEMP_DURATION", "2m"))
		input.InitTemperature(tempStateTopic, tempDuration)
	}

	if utils.GetEnv("MOISTURE_SENSOR_ENABLED", "false") != "false" {
		moistureStateTopic := utils.GetEnv("INPUT_MOISTURE_TOPIC", "moisture_topic")
		moistureDuration, _ := time.ParseDuration(utils.GetEnv("INPUT_MOISTURE_DURATION", "2m"))
		input.InitMoisture(moistureStateTopic, moistureDuration)
	}
}

func initOutput() {
	if utils.GetEnv("DOSING_PUMP_ENABLED", "false") != "false" {
		dosingCommandTopic := utils.GetEnv("OUTPUT_DOSING_COMMAND_TOPIC", "")
		dosingStateTopic := utils.GetEnv("OUTPUT_DOSING_STATE_TOPIC", "")
		output.InitDosing(dosingCommandTopic, dosingStateTopic)
	}
}
