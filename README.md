
# ocpp

[![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](https://github.com/tterb/atomic-design-ui/blob/master/LICENSEs)

Golang package implementing the JSON version of the Open Charge Point Protocol (OCPP). Currently OCPP 1.6 is supported


## Installation

Go version 1.18+ is required

```bash
  go get github.com/aliml92/ocpp
```
    
## Usage

### Cental System  
```go
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

	
	// create a ChargePoint
	cp = ocpp.NewChargePoint(c, chargePointId)
	

	// register handlers for CP initiated calls
	cp.On("BootNotification", BootNotificationHandler)

	// make CS initiated calls
	var req ocpp.Payload = &v16.ChangeAvailabilityReq{
		ConnectorId: 1,
		Type: "Operative",
	}
	res, err := cp.Call("ChangeAvailability", req)
	if err != nil {
		fmt.Printf("error calling: %v", err)
		return
	}
	fmt.Printf("ChangeAvailabilityRes: %v\n", res)

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
```
`ChargePoint` represents a single Charge Point (CP) connected to Central System
and after it is created, register CP initiated call handlers using `cp.On` method
Making Central System initiated call can be created using `cp.Call` method.
To make a Call to multiple charge points concurrently refer to `examples/` folder.

### Charge Point
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/aliml92/ocpp"
	"github.com/aliml92/ocpp/v16"
	"github.com/gorilla/websocket"
)


var cp *ocpp.ChargePoint

func main() {
	
	chargePointId := "client_01"
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
	cp = ocpp.NewChargePoint(c, chargePointId)


	// register handlers for CS initiated calls
	cp.On("ChangeAvailability", ChangeAvailabilityHandler)


	// make a BootNotification Call to Central System
	req := &v16.BootNotificationReq{
		ChargePointModel: "ModelX",
		ChargePointVendor: "VendorX",
	} 
	res, err := cp.Call("BootNotification", req)
	if err != nil {
		fmt.Printf("error calling: %v", err)
		return
	}
	fmt.Printf("BootNotificationRes: %v\n", res)

}


func ChangeAvailabilityHandler(p ocpp.Payload) ocpp.Payload {
	req := p.(*v16.ChangeAvailabilityReq)
	fmt.Printf("ChangeAvailabilityReq: %v\n", req)
	
	var res ocpp.Payload = &v16.ChangeAvailabilityConf{
		Status: "Accepted",
	}
	
	return res
}
```
After creating `ChargePoint` register CS (Central System) initiated call handlers.
Making a call to CS is same as the above snippet where just call `cp.Call` method.
## Contributing

Contributions are always welcome!
Implementing higher versions of ocpp is highly appreciated

See `CONTRIBUTING.md` for ways to get started.
