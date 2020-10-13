package main

import (
	"gitlab.com/r1chjames/aquarium-controller/input"
	"gitlab.com/r1chjames/aquarium-controller/mqttBackend"
	"gitlab.com/r1chjames/aquarium-controller/output"
	"gitlab.com/r1chjames/aquarium-controller/utils"
	"net/url"
	"time"
)

func main() {
	done := make(chan bool)
	initMqtt()
	go initTemperatureSensorModule()
	go initMoistureSensorModule()
	go initDosingPumpModule()
	<-done
}

func initMqtt() {
	clientId := utils.GetEnv("MQTT_CLIENT_ID", "aquarium_sensors")
	brokerUrl := utils.GetEnv("MQTT_BROKER_URL", "broker")
	user := utils.GetEnv("MQTT_USER", "user")
	password := utils.GetEnv("MQTT_PASSWORD", "password")
	mqttBackend.Connect(clientId, &url.URL{Host: brokerUrl, User: url.UserPassword(user, password)})
}

func initTemperatureSensorModule() {
	if utils.GetEnv("TEMP_SENSOR_ENABLED", "false") != "false"{
		tempStateTopic := utils.GetEnv("INPUT_TEMP_TOPIC", "temperature_topic")
		tempDuration, _ := time.ParseDuration(utils.GetEnv("INPUT_TEMP_DURATION", "2m"))
		input.InitTemperature(tempStateTopic, tempDuration)
	}
}

func initMoistureSensorModule() {
	if utils.GetEnv("MOISTURE_SENSOR_ENABLED", "false") != "false" {
		moistureStateTopic := utils.GetEnv("INPUT_MOISTURE_TOPIC", "moisture_topic")
		moistureDuration, _ := time.ParseDuration(utils.GetEnv("INPUT_MOISTURE_DURATION", "2m"))
		gpioChip := utils.GetEnv("GPIO_CHIP", "gpiochip0")
		gpioPin := utils.GetEnv("INPUT_MOISTURE_PIN", "25")
		input.InitMoisture(moistureStateTopic, moistureDuration, gpioChip, gpioPin)
	}
}

func initDosingPumpModule() {
	if utils.GetEnv("DOSING_PUMP_ENABLED", "false") != "false" {
		dosingCommandTopic := utils.GetEnv("OUTPUT_DOSING_COMMAND_TOPIC", "")
		dosingStateTopic := utils.GetEnv("OUTPUT_DOSING_STATE_TOPIC", "")
		gpioChip := utils.GetEnv("GPIO_CHIP", "gpiochip0")
		timezone := utils.GetEnv("TZ", "Europe/London")
		output.InitDosing(dosingCommandTopic, dosingStateTopic, gpioChip, timezone)
	}
}
