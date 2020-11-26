// gpio.go
package gpio

import (
	"io/ioutil"
	"log"
	"os"
)

const (
	gpioBasePath     = "/sys/class/gpio"
	gpioExportPath   = "/sys/class/gpio/export"
	gpioUnexportPath = "/sys/class/gpio/unexport"
)

type GPIO struct {
	pin string
}

func NewGPIO() GPIO {
	return GPIO{}
}

func (g GPIO) Pin(pin string) GPIO {
	g.pin = pin
	if _, err := os.Stat(gpioBasePath + "/gpio" + g.pin); os.IsNotExist(err) {
		err := ioutil.WriteFile(gpioExportPath, []byte(g.pin), 0666)
		if err != nil {
			log.Println(err)
		}
	}
	return g
}

func (g GPIO) Out() GPIO {
	err := ioutil.WriteFile(gpioBasePath+"/gpio"+g.pin+"/direction", []byte("out"), 0666)
	if err != nil {
		log.Println(err)
	}
	return g
}

func (g GPIO) In() GPIO {
	err := ioutil.WriteFile(gpioBasePath+"/gpio"+g.pin+"/direction", []byte("in"), 0666)
	if err != nil {
		log.Println(err)
	}
	return g
}

func (g GPIO) High() bool {
	err := ioutil.WriteFile(gpioBasePath+"/gpio"+g.pin+"/value", []byte("1"), 0666)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (g GPIO) Low() bool {
	err := ioutil.WriteFile(gpioBasePath+"/gpio"+g.pin+"/value", []byte("0"), 0666)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (g GPIO) PinRead(pin string) byte {
	value, err := ioutil.ReadFile(gpioBasePath + "/gpio" + pin + "/value")
	if err != nil {
		log.Println(err)
	}

	return value[0] - 48
}

func (g GPIO) PinUnexport(pin string) bool {
	err := ioutil.WriteFile(gpioUnexportPath, []byte(pin), 0666)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
