package danikor

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

// AnsData represents the structure of the AnsData packet
type AnsData struct {
	Header       byte
	DataLen      uint32
	AnsMode      byte
	MID          string
	Data         []byte
	Torque       DanitorTorque
	TorqueResult *DanitorTorqueResult
	Tailer       byte
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

	if a.MID == "0202" {
		a.TorqueResult = parseTorqueResult(string(a.Data))
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

type DanitorTorqueResult struct {
	FinalTorqueValue  string
	FinalAngleMonitor string
	FinalTime         string
	FinalAngleFinal   string
	FinalStatus       string
	NgCode            string
	StageResults      map[string]StageResult
	Status            map[string]string
}

type StageResult struct {
	Torque float64
	Angle  float64
	Time   float64
	Status string
}

func parseTorqueResult(str string) *DanitorTorqueResult {
	result := &DanitorTorqueResult{
		StageResults: make(map[string]StageResult),
		Status:       make(map[string]string),
	}

	pairs := strings.Split(str, ";")
	for _, pair := range pairs {
		values := strings.Split(pair, "=")
		if len(values) == 2 {
			key := values[0]
			value := values[1]
			switch key {
			case "00010":
				fields := strings.Split(value, ",")
				if len(fields) >= 4 {
					result.FinalTorqueValue = fields[0]
					result.FinalAngleMonitor = fields[1]
					result.FinalTime = fields[2]
					result.FinalAngleFinal = fields[3]
				}
			case "00011":
				result.FinalStatus = value
			case "00012":
				result.NgCode = value
			default:
				fmt.Println("key:", key, "value:", value)
				//key: 01030 3 is stageKey value: 0.013,1257.069,3.000(Torque,Angle,Time),
				if len(key) >= 5 && key[:3] == "010" && strings.HasSuffix(key, "0") { //value
					stageKey := key[3:4]
					fmt.Println("stage:", stageKey)
					//split value to get Torque,Angle,Time
					values := strings.Split(value[5:], ",")
					if len(values) == 3 {
						torque, err := strconv.ParseFloat(values[0], 64)
						if err != nil {
							// handle error
						}
						angle, err := strconv.ParseFloat(values[1], 64)
						if err != nil {
							// handle error
						}
						time, err := strconv.ParseFloat(values[2], 64)
						if err != nil {
							// handle error
						}
						stageResult := StageResult{
							Torque: torque,
							Angle:  angle,
							Time:   time,
						}
						result.StageResults[stageKey] = stageResult
					}
				}

				if len(key) >= 5 && key[:3] == "010" && strings.HasSuffix(key, "1") {
					stageKey := key[3:4]
					fmt.Println("status stage:", stageKey)
					result.Status[stageKey] = value

				}

			}
		}
	}

	return result
}
