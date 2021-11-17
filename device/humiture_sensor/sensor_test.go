package humiture_sensor

import (
	"fmt"
	"testing"
)

func TestExample(t *testing.T) {
	sensor, err := Connect(&Config{Bus: "1", Addr: 0x38})
	if err != nil {
		panic(err)
	}
	if err = sensor.InitRegister(); err != nil {
		panic(err)
	}
	if err = sensor.Start(); err != nil {
		panic(err)
	}
	data, err := sensor.Query()
	if err != nil {
		panic(err)
	}

	fmt.Printf("RH = %.2f, T = %.2f\n", data.RH, data.T)
}
