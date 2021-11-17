package operator

import (
	sensor "github.com/dovics/raspberry-co2/device/humiture_sensor"
)

type HumitureOperator struct {
	sensor *sensor.HumitureSensor
}

func NewHumitureOperator(sensor *sensor.HumitureSensor) *HumitureOperator {
	return &HumitureOperator{sensor: sensor}
}

func (h *HumitureOperator) Query() (data *sensor.Data, err error) {
	if err = h.sensor.InitRegister(); err != nil {
		return
	}
	if err = h.sensor.Start(); err != nil {
		return
	}

	return h.sensor.Query()
}
