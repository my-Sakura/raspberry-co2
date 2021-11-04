package sensor

import (
	"bufio"

	"github.com/tarm/serial"

	"github.com/abserari/jx-co2-101-sensor/util/log"
)

var ActiveModeChange = []byte{0xff, 0x05, 0x03, 0x01, 0x00, 0x00, 0x00, 0x00, 0xf8}
var QueryModeChange = []byte{0xff, 0x05, 0x03, 0x02, 0x00, 0x00, 0x00, 0x00, 0xf7}
var QueryPPM = []byte{0xff, 0x05, 0x03, 0x03, 0x01, 0x00, 0x00, 0x00, 0xf5}
var Correct = []byte{0xff, 0x05, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf8}
var MODBUS_RTU = []byte{0x05, 0x03, 0x00, 0x05, 0x00, 0x01, 0x94, 0x07}

type CO2Sensor struct {
	*bufio.Reader
	serial *serial.Port
}

func Connect(config *serial.Config) (*CO2Sensor, error) {
	s, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	return &CO2Sensor{
		bufio.NewReader(s),
		s,
	}, nil
}

func (s *CO2Sensor) SendMODBUS_RTU() error {
	n, err := s.serial.Write(MODBUS_RTU)
	if err != nil {
		log.Error(err)
		return err
	}

	if n == len(MODBUS_RTU) {
		log.Info("send MODBUS_RTU successful")
	}

	data, _, _ := s.ReadLine()
	log.Info(string(data))

	return nil
}

func (s *CO2Sensor) SendCorrect() error {
	n, err := s.serial.Write(Correct)
	if err != nil {
		log.Error(err)
		return err
	}

	if n == len(Correct) {
		log.Info("send Correct successful")
	}

	data, _, _ := s.ReadLine()
	log.Info(string(data))

	return nil
}

func (s *CO2Sensor) SendActiveModeChange() error {
	n, err := s.serial.Write(ActiveModeChange)
	if err != nil {
		log.Error(err)
		return err
	}

	if n == len(ActiveModeChange) {
		log.Debug("😀 send ActiveModeChange successful")
	}

	// read response
	var b = make([]byte, 8)
	n, err = s.serial.Read(b)
	if err != nil || n != 8 {
		log.Error("😥 set ActiveModeChange failed", err, n, b)
		return err
	}

	var check = false
	// read check bit and crc8 check
	// check if crc[0] -1 != b[0] - b[1]- b[2]- b[3] - b[4] - b[5] - b[6] - b[7]
	var crc = make([]byte, 1)
	n, err = s.serial.Read(crc)
	if err != nil || n != 1 {
		log.Error("😥 set ActiveModeChange failed", err, n, crc)
		return err
	}

	// need to check
	check = true
	if !check {
		log.Error("😥 set ActiveModeChange failed", err, n, crc)
		return err
	}

	log.Info("😀 set QueryModeChange successful")

	return nil
}

func (s *CO2Sensor) SendQueryModeChange() error {
	n, err := s.serial.Write(QueryModeChange)
	if err != nil {
		log.Error(err)
		return err
	}
	if n == len(QueryModeChange) {
		log.Debug("😀 send QueryModeChange successful")
	}

	// read response
	var b = make([]byte, 8)
	n, err = s.serial.Read(b)
	if err != nil || n != 8 {
		log.Error("😥 set QueryModeChange failed", err, n, b)
		return err
	}

	var check bool = false
	// read check bit and crc8 check
	// check if crc[0] -1 != b[0] - b[1]- b[2]- b[3] - b[4] - b[5] - b[6] - b[7]
	var crc = make([]byte, 1)
	n, err = s.serial.Read(crc)
	if err != nil || n != 1 {
		log.Error("😥 set QueryModeChange failed", err, n, crc)
		return err
	}

	// need to check
	check = true
	if !check {
		log.Error("😥 set QueryModeChange failed", err, n, crc)
		return err
	}

	log.Info("😀 set QueryModeChange successful")

	return nil
}

func (s *CO2Sensor) SendQuery() error {
	n, err := s.serial.Write(QueryPPM)
	if err != nil {
		log.Error(err)
		return err
	}

	if n == len(QueryPPM) {
		log.Info("send QueryPPM successful")
	}

	return nil
}
