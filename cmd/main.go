package main

import (
	"fmt"

	. "github.com/linexjlin/danikor"
)

func main() {
	dc := NewDanikorTCPConnection("192.168.2.5:5000", func(ansData AnsData) {
		fmt.Println("ansData:", string(ansData.MID))
		switch ansData.MID {
		case "0203":
			fmt.Println(ansData.Torque.Pset, ansData.Torque.IsCurveStart, ansData.Torque.IsCurveEnd)
		}
	})

	dc.Dial()
	dc.Establish()
	dc.SubscribeRealTimeData()
	dc.ForwardTurn()
	dc.StartReceiveData()
}
