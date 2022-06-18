package main

import (
	"fmt"
	"net/http"

	"github.com/aliml92/ocpp"
	"github.com/aliml92/ocpp/v16"
	"github.com/gorilla/websocket"
)


var chargePoint *ocpp.ChargePoint

func main() {
	
	chargePointId := "12345"
	url := fmt.Sprintf("ws://localhost:8080/ocpp/v16/%s", chargePointId)
	

	header := http.Header{
		"Sec-WebSocket-Protocol": []string{"ocpp1.6"},
	}

	fmt.Printf("connecting to %s", url)

	c, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		fmt.Printf("error dialing: %v", err)
		return
	}
	fmt.Printf("connected to %s", url)
	
	defer c.Close()
	chargePoint = ocpp.NewChargePoint(c, chargePointId)


	// make a BootNotification Call to Central System
	req := &v16.BootNotificationReq{
		ChargePointModel: "ModelX",
		ChargePointVendor: "VendorX",
	} 
	res, err := chargePoint.Call("BootNotification", req)
	if err != nil {
		fmt.Printf("error calling: %v", err)
		return
	}
	fmt.Printf("BootNotificationRes: %v\n", res)


	// register handlers for CSMS initiated calls
	chargePoint.On("ChangeAvailability", ChangeAvailabilityHandler)
}


func ChangeAvailabilityHandler(p ocpp.Payload) ocpp.Payload {
	req := p.(*v16.ChangeAvailabilityReq)
	fmt.Printf("ChangeAvailabilityReq: %v\n", req)
	
	var res ocpp.Payload = &v16.ChangeAvailabilityConf{
		Status: "Accepted",
	}
	
	return res
}