package main

import (
	"fmt"

	. "github.com/linexjlin/danikor"
)

func main() {
	dc := NewDanikorTCPConnection("192.168.2.5:5000", func(ansData AnsData) {
		fmt.Println("ansMid:", string(ansData.MID))
		switch ansData.MID {
		case "0203":
			fmt.Println(ansData.Torque.Pset, ansData.Torque.IsCurveStart, ansData.Torque.IsCurveEnd)
		case "0202":
			fmt.Println("torque result:", ansData.TorqueResult.FinalAngleFinal)
		}
	})

	dc.Dial()
	dc.Establish()
	dc.SubscribeRealTimeData()
	dc.SubscribeResultData()
	dc.ForwardTurn()
	dc.StartReceiveData()
}
