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
			fmt.Printf("torque  %+v\n", ansData)
			fmt.Println("Pset:", ansData.Torque.Pset, ansData.Torque.IsCurveStart, ansData.Torque.IsCurveEnd)
		case "0202":
			fmt.Printf("torque result %+v\n", ansData.TorqueResult)
		}
	})

	dc.Dial()
	dc.Establish()
	dc.ChosePset(2)
	dc.SubscribeRealTimeData()
	dc.SubscribeResultData()
	dc.ForwardTurn()
	dc.StartReceiveData()
}
