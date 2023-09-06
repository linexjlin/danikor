package danikor

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type DanikorTCPConnection struct {
	address         string
	conn            net.Conn
	receiveCallBack func(AnsData)
}

func NewDanikorTCPConnection(addr string, receiveCallBack func(AnsData)) *DanikorTCPConnection {
	dc := &DanikorTCPConnection{
		address:         addr,
		receiveCallBack: receiveCallBack,
	}
	return dc
}

func (dc *DanikorTCPConnection) Dial() {
	for {
		conn, err := net.Dial("tcp", dc.address)
		if err != nil {
			fmt.Printf("Failed to dial: %v\n", err)
			time.Sleep(time.Second)
		} else {
			dc.conn = conn
			break
		}
	}
}

func showData(data []byte) {
	// Unmarshal the binary data into the AnsData struct
	var ansData AnsData
	if err := ansData.UnmarshalBinary(data); err != nil {
		panic(err)
	}

	// Marshal the AnsData struct to JSON
	jsonData, err := json.Marshal(ansData)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonData))
}

func parseData(data []byte) AnsData {
	// Unmarshal the binary data into the AnsData struct
	var ansData AnsData
	if err := ansData.UnmarshalBinary(data); err != nil {
		panic(err)
	}
	return ansData
}

func (dc *DanikorTCPConnection) Establish() {
	data := []byte{0x02, 0x00, 0x00, 0x00, 0x05, 0x52, 0x30, 0x30, 0x30, 0x31, 0x03} //mid 001 建立通信数据包

	_, err := dc.conn.Write(data)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	// Receive and print the response
	response := make([]byte, 1024)
	n, err := dc.conn.Read(response)
	if err != nil {
		fmt.Println("Error receiving response:", err)
		return
	}
	showData(response[:n])
}

func (dc *DanikorTCPConnection) SubscribeRealTimeData() {
	data := []byte{0x02, 0x00, 0x00, 0x00, 0x05, 0x52, 0x30, 0x32, 0x30, 0x33, 0x03} //mid 0203 实时曲线数据

	_, err := dc.conn.Write(data)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	// Receive and print the response
	response := make([]byte, 1024)
	n, err := dc.conn.Read(response)
	if err != nil {
		fmt.Println("Error receiving response:", err)
		return
	}
	showData(response[:n])
}

func (dc *DanikorTCPConnection) ForwardTurn() {
	data := []byte{0x02, 0x00, 0x00, 0x00, 0x0A, 0x57, 0x30, 0x33, 0x30, 0x31, 0x30, 0x31, 0x3D, 0x31, 0x3B, 0x03} //mid 正转
	_, err := dc.conn.Write(data)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	// Receive and print the response
	response := make([]byte, 1024)
	n, err := dc.conn.Read(response)
	if err != nil {
		fmt.Println("Error receiving response:", err)
		return
	}
	showData(response[:n])
}

func (dc *DanikorTCPConnection) StartReceiveData() {
	response := make([]byte, 1024)
	// Continuously receive and print data
	for {
		n, err := dc.conn.Read(response)
		if err != nil {
			fmt.Println("Error receiving response:", err)
			return
		}
		//showData(response[:n])
		ansData := parseData(response[:n])
		dc.receiveCallBack(ansData)
		//fmt.Println("xxxx", hex.EncodeToString(response[:n]))
	}
}
