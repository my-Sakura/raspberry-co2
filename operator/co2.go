package operator

import (
	"errors"
	"io"
	"strconv"
	"strings"

	sensor "github.com/dovics/raspberry-co2/device/co2_sensor"
	"github.com/dovics/raspberry-co2/util/log"
)

type CO2Operator struct {
	sensor *sensor.CO2Sensor
}

func NewCO2Operator(sensor *sensor.CO2Sensor) *CO2Operator {
	err := sensor.SendActiveModeChange()
	if err != nil {
		log.Info(err)
	}

	return &CO2Operator{
		sensor,
	}
}

func (o *CO2Operator) QueryCO2() (int, error) {
	retry := 20
	for retry > 0 {
		raw, err := o.sensor.ReadLine()
		if err != nil {
			if err == io.EOF && retry < 15 {
				if retry <= 10 {
					err := o.sensor.Reconnect()
					if err != nil {
						log.Info(err)
					}
				}
				err := o.sensor.SendActiveModeChange()
				if err != nil {
					log.Info(err)
				}
			}
			retry--
			log.Info(err)
			continue
		}

		// bytes data            4444       ppm
		//           space space 4444 space ppm
		strs := strings.Split(string(raw), " ")
		log.Debug("receive raw", raw, "and split to ", strs)
		if len(strs) < 3 {
			retry--
			log.Info("wrong data format")
			continue
		}

		// convert 4444 string to int
		return strconv.Atoi(strs[2])
	}

	return 0, errors.New("over max retry times")
}
