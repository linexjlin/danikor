package danikor

import (
	"encoding/binary"
)

// AnsData represents the structure of the AnsData packet
type AnsData struct {
	Header  byte
	DataLen uint32
	AnsMode byte
	MID     string
	Data    []byte
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
	return nil
}
