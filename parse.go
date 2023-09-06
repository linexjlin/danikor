package danikor

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

// AnsData represents the structure of the AnsData packet
type AnsData struct {
	Header  byte
	DataLen uint32
	AnsMode byte
	MID     string
	Data    []byte
	Torque  DanitorTorque
	Tailer  byte
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (a *AnsData) UnmarshalBinary(data []byte) error {
	a.Header = data[0]
	a.DataLen = binary.BigEndian.Uint32(data[1:5])
	a.AnsMode = data[5]
	a.MID = string(data[6:10])
	a.Data = data[10 : len(data)-1]
	a.Tailer = data[len(data)-1]
	if a.MID == "0203" {
		a.Torque = parseTorqueData(string(a.Data))
	}
	return nil
}

type DanitorTorque struct {
	SampleFrequency string    // 0101
	Pset            string    // 0102
	IsCurveEnd      bool      // 0201
	IsCurveStart    bool      // 0202
	Torque          []float64 // 0301
	Angle           []float64 // 0302
	CurrentPset     []int     // 0401
}

func parseTorqueData(str string) DanitorTorque {
	fmt.Println("coming str:", str)

	data := DanitorTorque{}
	parts := strings.Split(str, ";")
	for _, part := range parts {
		keyValue := strings.Split(part, "=")
		if len(keyValue) != 2 {
			continue
		}
		key := keyValue[0]
		value := keyValue[1]

		switch key {
		case "0101":
			data.SampleFrequency = value
		case "0102":
			data.Pset = value
		case "0201":
			if value == "1" {
				data.IsCurveEnd = true
			} else {
				data.IsCurveEnd = false
			}
		case "0202":
			if value == "1" {
				data.IsCurveStart = true
			} else {
				data.IsCurveStart = false
			}
		case "0301":
			torqueValues := strings.Split(value, ",")
			data.Torque = make([]float64, len(torqueValues))
			for i, v := range torqueValues {
				torque, _ := strconv.ParseFloat(v, 64)
				data.Torque[i] = torque
			}
		case "0302":
			angleValues := strings.Split(value, ",")
			data.Angle = make([]float64, len(angleValues))
			for i, v := range angleValues {
				angle, _ := strconv.ParseFloat(v, 64)
				data.Angle[i] = angle
			}
		case "0401":
			psetValues := strings.Split(value, ",")
			data.CurrentPset = make([]int, len(psetValues))
			for i, v := range psetValues {
				pset, _ := strconv.Atoi(v)
				data.CurrentPset[i] = pset
			}
		}
	}

	return data
}
