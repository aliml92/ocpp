package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/aliml92/ocpp"
	"github.com/aliml92/ocpp/v16"
	"github.com/gorilla/websocket"
)


var cp *ocpp.ChargePoint

func main() {
	
	chargePointId := "client_02"
	url := fmt.Sprintf("ws://localhost:8080/ocpp/v16/%s", chargePointId)
	header := http.Header{
		"Sec-WebSocket-Protocol": []string{"ocpp1.6"},
	}

	c, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		fmt.Printf("error dialing: %v", err)
		return
	}
	defer c.Close()


	// create a ChargePoint
	cp = ocpp.NewChargePoint(c, chargePointId, "ocpp1.6")


	// register handlers for CS initiated calls
	cp.On("ChangeAvailability", ChangeAvailabilityHandler)


	// make a BootNotification Call to Central System
	req := &v16.BootNotificationReq{
		ChargePointModel: "ModelY",
		ChargePointVendor: "VendorY",
	} 
	res, err := cp.Call("BootNotification", req)
	if err != nil {
		fmt.Printf("error calling: %v", err)
		return
	}
	fmt.Printf("BootNotificationRes: %v\n", res)
	
	// prevent main() from exiting
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()

}


func ChangeAvailabilityHandler(id string, p ocpp.Payload) ocpp.Payload {
	req := p.(*v16.ChangeAvailabilityReq)
	fmt.Printf("ChangeAvailabilityReq: %v\n", req)
	
	var res ocpp.Payload = &v16.ChangeAvailabilityConf{
		Status: "Accepted",
	}
	
	return res
}