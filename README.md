Aquarium Controller written in GO

---

Inputs:
- Temperature
- Moisture

Outputs:
- Dosing pump

# Installation
For use on a Raspberry Pi, additional modules need to be loaded:

```
$ sudo vi /etc/modules
```
Add the following lines to the end of the file
```
w1-gpio
w1-therm
```

