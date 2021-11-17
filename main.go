package main

import (
	"time"

	"github.com/dovics/raspberry-co2/exporter"
	"github.com/dovics/raspberry-co2/operator"
	"github.com/dovics/raspberry-co2/util/log"

	"github.com/dovics/raspberry-co2/device/co2_sensor"
	"github.com/dovics/raspberry-co2/device/humiture_sensor"
)

func main() {
	co2Sensor, err := co2_sensor.Connect(&co2_sensor.Config{Name: "/dev/serial0", Baud: 9600, ReadTimeout: time.Second * 5})
	if err != nil {
		log.Fatal(err)
	}
	humitureSensor, err := humiture_sensor.Connect(&humiture_sensor.Config{Bus: "1", Addr: 0x38})
	if err != nil {
		log.Fatal(err)
	}

	co2Operator := operator.NewCO2Operator(co2Sensor)
	humitureOperator := operator.NewHumitureOperator(humitureSensor)

	e := exporter.NewExporter()
	e.Register("co2", func() (interface{}, error) {
		return co2Operator.QueryCO2()
	})
	e.Register("humiture", func() (interface{}, error) {
		return humitureOperator.Query()
	})

	if err := e.Run(); err != nil {
		log.Fatal(err)
	}
}
