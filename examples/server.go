package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aliml92/ocpp"
	"github.com/aliml92/ocpp/v16"

	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
	Subprotocols: []string{"ocpp1.6"},
}


var cp *ocpp.ChargePoint


func main(){
	go callExample()
	http.HandleFunc("/", wsHandler)
	http.ListenAndServe("localhost:8080", nil)
}



func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	subProtocol := r.Header.Get("Sec-WebSocket-Protocol")	
	if subProtocol== "" {
		fmt.Println("client hasn't requested any Subprotocol. Closing Connection")
		c.Close()
	}
	if subProtocol != "ocpp1.6" {
		fmt.Println("client has requested an unsupported Subprotocol. Closing Connection")
		c.Close()
	}
	chargePointId := strings.Split(r.URL.Path, "/")[3]

	// create a ChargePoint
	cp = ocpp.NewChargePoint(c, chargePointId)
	registerHandlers()
	
}




// register handlers for CP initiated calls
func registerHandlers(){
	cp.On("BootNotification", BootNotificationHandler)
}


func BootNotificationHandler(p ocpp.Payload) ocpp.Payload {
	req := p.(*v16.BootNotificationReq)
	fmt.Printf("BootNotificationReq: %v\n", req)

	var res ocpp.Payload = &v16.BootNotificationConf{
		CurrentTime: time.Now().Format("2006-01-02T15:04:05.000Z"),
		Interval:    60,
		Status:      "Accepted",
	}
	return res
}


// make a ChangeAvailability Call to all connected Charge Points
func callExample(){
	// sleep for a while to allow both charge points to be connected
	time.Sleep(time.Second * 15)

	var req ocpp.Payload = &v16.ChangeAvailabilityReq{
		ConnectorId: 1,
		Type: "Operative",
	}
	cps := ocpp.ChargePoints
	fmt.Printf("number of connected charge points: %v\n", len(cps))

	ch := make(chan ocpp.Payload, len(cps))
	for _, cp := range cps {
		cp := cp
		go func() {
			res, err := cp.Call("ChangeAvailability", req)
			if err != nil {
				fmt.Printf("error calling: %v", err)
				return
			}
			ch <- res
		}()
	}
	// get all results from channel
	for c := range ch {
		fmt.Printf("ChangeAvailabilityRes: %v\n", c)
	}
}

