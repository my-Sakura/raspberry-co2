package sensor

import (
	"testing"
	"time"

	"github.com/tarm/serial"
)

// pi3 should open uart and communicate with device: /dev/ttyAMA0 | /dev/serial0
func TestExample(t *testing.T) {
	sensor, err := Connect(&serial.Config{Name: "COM1", Baud: 9600, ReadTimeout: time.Second * 5})
	if err != nil {
		t.Fatal(err)
	}
	sensor.SendActiveModeChange()

	go func() {
		for {
			sensor.SendQuery()
			time.Sleep(3 * time.Second)
		}
	}()

	for {
		data, err := sensor.ReadLine()
		if err != nil {
			t.Error(err)
		}

		t.Log(data)
	}

}

func TestCRC(t *testing.T) {
	CRC(ActiveModeChange, 0xfc)
}
