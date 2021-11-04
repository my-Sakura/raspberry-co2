package main

import (
	"time"

	"github.com/abserari/jx-co2-101-sensor/exporter"
	"github.com/abserari/jx-co2-101-sensor/operator"
	"github.com/abserari/jx-co2-101-sensor/util/log"

	sensor "github.com/abserari/jx-co2-101-sensor/device/co2_sensor"
	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: "/dev/ttyAMA0", Baud: 2400, ReadTimeout: time.Second * 5}
	co2Sensor, err := sensor.Connect(c)
	if err != nil {
		log.Fatal(err)
	}

	co2Operator := operator.NewOperator(co2Sensor)

	e := exporter.NewExporter()
	e.Register("co2", func() (interface{}, error) {
		return co2Operator.QueryCO2()
	})

	if err := e.Run(); err != nil {
		log.Fatal(err)
	}
}
