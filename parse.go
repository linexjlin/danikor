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

/*
0=无定意；
1=OK，拧紧合格；
2=NG，拧紧不合格；
*/
func (dr *DanitorTorqueResult) ShowFinalStatus() string {
	switch dr.FinalStatus {
	case "0":
		return "无定意"
	case "1":
		return "OK，拧紧合格"
	case "2":
		return "NG，拧紧不合格"
	default:
		return "未知状态"
	}
}

/*
0x00=无定意；
0x01=最终扭矩过大；
0x02=最终扭矩过大；
0x03=最终角度过大；
0x04=最终角度过小；
0xn1=第 n 步扭矩过大；
0xn2=第 n 步扭矩过大；
1<n<5
0x90=总时间超限；
*/
func (dr *DanitorTorqueResult) ShowNgCode() string {
	r := ""
	switch dr.NgCode {
	case "00":
		r = "无定意"
	case "01":
		r = "最终扭矩过大"
	case "02":
		r = "最终扭矩过大"
	case "03":
		r = "最终角度过大"
	case "04":
		r = "最终角度过小"
	case "90":
		r = "总时间超限"
	default:
		if strings.HasSuffix(dr.NgCode, "1") {
			r = fmt.Sprintf("第 %s 步扭矩过大", string(dr.NgCode[0]))
		}
		if strings.HasSuffix(dr.NgCode, "2") {
			r = fmt.Sprintf("第 %s 步扭矩过大", string(dr.NgCode[0]))
		}
	}
	return r
}

/*
0=无定意；
1=OK
2=扭矩过大；
3=扭矩过小；
4=角度过大；
5=角度过小；
6=时间过长；
7=时间过短；
*/
func (dr *DanitorTorqueResult) ShowStageStatus(code string) string {
	r := ""
	switch code {
	case "0":
		r = "无定意"
	case "1":
		r = "OK"
	case "2":
		r = "扭矩过大"
	case "3":
		r = "扭矩过小"
	case "4":
		r = "角度过大"
	case "5":
		r = "角度过小"
	case "6":
		r = "时间过长"
	case "7":
		r = "时间过短"
	}
	return r
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
