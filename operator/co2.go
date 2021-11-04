package operator

import (
	"strconv"
	"strings"

	sensor "github.com/abserari/jx-co2-101-sensor/device/co2_sensor"
	"github.com/abserari/jx-co2-101-sensor/util/log"
)

type CO2Operator struct {
	sensor *sensor.CO2Sensor
}

func NewOperator(sensor *sensor.CO2Sensor) *CO2Operator {
	err := sensor.SendActiveModeChange()
	if err != nil {
		log.Error(err)
	}

	return &CO2Operator{
		sensor,
	}
}

func (o *CO2Operator) QueryCO2() (int, error) {
	raw, _, err := o.sensor.ReadLine()
	if err != nil {
		log.Error(err)
	}

	// bytes data            4444       ppm
	//           space space 4444 space ppm
	strs := strings.Split(string(raw), " ")
	log.Debug("receive raw", raw, "and split to ", strs)

	// convert 4444 string to int
	return strconv.Atoi(strs[2])

}
