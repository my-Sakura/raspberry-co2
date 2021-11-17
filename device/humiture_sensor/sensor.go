package humiture_sensor

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

type HumitureSensor struct {
	Data *Data

	i2cDev *i2c.Dev
}

type Data struct {
	// temperature
	T float64
	// relative humidity
	RH float64
}

type Config struct {
	Bus  string
	Addr uint16
}

func Connect(c *Config) (*HumitureSensor, error) {
	_, err := host.Init()
	if err != nil {
		return nil, err
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open(c.Bus)
	if err != nil {
		return nil, err
	}

	return &HumitureSensor{
		i2cDev: &i2c.Dev{Addr: c.Addr, Bus: b},
	}, nil
}

func (h *HumitureSensor) InitRegister() error {
	time.Sleep(time.Millisecond * 100)
	result := make([]byte, 1)
	if err := h.writeAndRead([]byte{0x71}, result); err != nil {
		return err
	}

	if result[0]&0x18 != 0x18 {
		// TODO: init 0x1B、0x1C、0x1E register
		return errors.New("init 0x1B、0x1C、0x1E register failed")
	}

	return nil
}

func (h *HumitureSensor) Start() error {
	time.Sleep(time.Millisecond * 10)

	if err := h.write([]byte{0xAC}); err != nil {
		return err
	}
	if err := h.write([]byte{0x33}); err != nil {
		return err
	}
	if err := h.write([]byte{0x00}); err != nil {
		return err
	}

	return nil
}

func (h *HumitureSensor) Query() (*Data, error) {
	if b, ok := h.i2cDev.Bus.(i2c.BusCloser); ok {
		defer b.Close()
	}

	time.Sleep(time.Millisecond * 80)
	r := make([]byte, 1)
	if err := h.read(r); err != nil {
		return nil, err
	}

	retry := 0
	for r[0]&0x80 == 0x80 {
		retry++
		time.Sleep(time.Millisecond * 80)
		if err := h.read(r); err != nil {
			return nil, err
		}

		fmt.Printf("read state retry: %d\n", retry)
		if retry >= 20 {
			fmt.Println("retry many times")
			os.Exit(1)
		}
	}

	result := make([]byte, 6)
	if err := h.writeAndRead([]byte{0x71}, result); result != nil {
		return nil, err
	}

	var humidity, t uint32
	humidity = (humidity | uint32(result[1])) << 8
	humidity = (humidity | uint32(result[2])) << 8
	humidity = humidity | uint32(result[3])
	humidity >>= 4
	t = (t | uint32(result[3])) << 8
	t = (t | uint32(result[4])) << 8
	t = (t | uint32(result[5]))
	t = t & 0xfffff
	RH, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float32(humidity)*100/1024/1024), 32)
	T, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float32(t)*200/1024/1024-50), 32)
	data := &Data{
		T:  T,
		RH: RH,
	}
	h.Data = data

	return data, nil
}

func (h *HumitureSensor) write(w []byte) error {
	return h.i2cDev.Tx(w, nil)
}

func (h *HumitureSensor) read(r []byte) error {
	return h.i2cDev.Tx(nil, r)
}

func (h *HumitureSensor) writeAndRead(w, r []byte) error {
	return h.i2cDev.Tx(w, r)
}
