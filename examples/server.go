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
		fmt.Println("Client hasn't requested any Subprotocol. Closing Connection")
		c.Close()
	}
	if subProtocol != "ocpp1.6" {
		fmt.Println("Client has requested an unsupported Subprotocol. Closing Connection")
		c.Close()
	}
	
	chargePointId := strings.Split(r.URL.Path, "/")[3]
	log.Printf("chargePointId: %s", chargePointId)
	
	
	cp = ocpp.NewChargePoint(c, chargePointId)
	
	cp.On("BootNotification", BootNotificationHandler)

}


func BootNotificationHandler(p ocpp.Payload) ocpp.Payload {
	req := p.(*v16.BootNotificationReq)
	fmt.Printf("BootNotificationReq: %+v\n", req)

	var res ocpp.Payload = &v16.BootNotificationConf{
		CurrentTime: time.Now().Format(time.RFC3339),
		Interval:    60,
		Status:      "Accepted",
	}
	return res
}