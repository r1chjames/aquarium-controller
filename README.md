# Aquarium Controller written in GO
---

Application that reads sensor data from connected sensors. Included is a dosing module that allows attached dosing pumps to be actuated. All communication is via MQTT.

Inputs:
- Temperature
- Moisture

Outputs:
- Dosing pump

## Installation
For use on a Raspberry Pi, additional modules need to be loaded:

```
$ sudo vi /etc/modules
```
Add the following lines to the end of the file
```
w1-gpio
w1-therm
```

