
# ocpp

[![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](https://github.com/tterb/atomic-design-ui/blob/master/LICENSEs)

Golang package implementing the JSON version of the Open Charge Point Protocol (OCPP). Currently OCPP 1.6 and 2.0.1 is supported.
The project is initially inspired by [mobility/ocpp](https://github.com/mobilityhouse/ocpp)

## Installation

Go version 1.18+ is required

```bash
  go get github.com/aliml92/ocpp
```

 ## Features

- [x] ocpp1.6 and ocpp2.0.1 support
- [x] logging
- [x] ping/pong customization on `WebSocketPingInterval`
- [x] server initiated ping activation 

## Roadmap

- [ ]   add unit/integration tests
- [x]   improve logging
- [ ]   add validation disabling feature
- [ ]   add better queque implementation
 

## Usage

### Cental System  (Server)
```go
package main

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/aliml92/ocpp"
	v16 "github.com/aliml92/ocpp/v16"
)


var csms *ocpp.Server

// log replaces standard log
var log *zap.SugaredLogger



func main() 
	logger, _ := zap.NewDevelopment()
	log = logger.Sugar()
	defer log.Sync()

	// set ocpp library's logger to zap logger
	ocpp.SetLogger(log)

	// start csms server with default configurations
	csms = ocpp.NewServer()

	csms.AddSubProtocol("ocpp1.6")
	csms.SetCheckOriginHandler(func(r *http.Request) bool { return true })
	csms.SetPreUpgradeHandler(customPreUpgradeHandler)
	csms.SetCallQueueSize(32)	

	// register charge-point-initiated action handlers
	csms.On("BootNotification", BootNotificationHandler)
	csms.After("BootNotification", SendChangeConfigration)
	csms.On("Authorize", AuthorizationHandler)
	csms.Start("0.0.0.0:8999", "/ws/", nil)
	
}

func SendChangeConfigration(cp *ocpp.ChargePoint, payload ocpp.Payload) {
	var req ocpp.Payload = v16.ChangeConfigurationReq{
		Key: "WebSocketPingInterval",
		Value: "30",
	}
	res, err := cp.Call("ChangeConfiguration", req)
	if err != nil {
		log.Debug(err)
	}
	log.Debug(res)
}


func customPreUpgradeHandler(w http.ResponseWriter, r *http.Request) bool {
	u, p, ok := r.BasicAuth()
	if !ok {
		log.Debug("error parsing basic auth")
		w.WriteHeader(401)
		return false
	}
	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-1]
	log.Debugf("%s is trying to connect with %s:%s", id, u, p)
	if u != id {
		log.Debug("username provided is correct: %s", u)
		w.WriteHeader(401)
		return false
	}
	return true
}



func BootNotificationHandler(cp *ocpp.ChargePoint, p ocpp.Payload) ocpp.Payload {
	req := p.(*v16.BootNotificationReq)
	log.Debugf("\nid: %s\nBootNotification: %v", cp.Id, req)
	
    var res ocpp.Payload = &v16.BootNotificationConf{
		CurrentTime: time.Now().Format("2006-01-02T15:04:05.000Z"),
		Interval:    60 ,
		Status:      "Accepted",
	}
	return res
}

func AuthorizationHandler(cp *ocpp.ChargePoint, p ocpp.Payload) ocpp.Payload {
	req := p.(*v16.AuthorizeReq)
	log.Debugf("\nid: %s\nAuthorizeReq: %v", cp.Id, req)
	
    var res ocpp.Payload = &v16.AuthorizeConf{
		IdTagInfo: v16.IdTagInfo{
			Status: "Accepted",
		},
	}
	return res
}
```
`ChargePoint` represents a single Charge Point (CP) connected to Central System
and after initializing `*ocpp.Server` , register CP initiated call handlers using `csms.On` method.
Making a Call can be done by excuting `cp.Call` method.



### Charge Point (Client)
```go
package main

import (
	"fmt"
	"time"
	"github.com/aliml92/ocpp"
	v16 "github.com/aliml92/ocpp/v16"
	"go.uber.org/zap"
)



var client *ocpp.Client

// log replaces standard log
var log *zap.SugaredLogger

// initialize zap logger
// for deveplopment only
func initLogger() {
	logger, _ := zap.NewDevelopment()
	log = logger.Sugar()
}

func main() {
	initLogger()
	defer log.Sync()

	// set ocpp library's logger to zap logger
	ocpp.SetLogger(log)

	// create client
	client = ocpp.NewClient()
	id := "client00"
	client.SetID(id)
	client.AddSubProtocol("ocpp1.6")
	client.SetBasicAuth(id, "dummypass")
	client.SetCallQueueSize(32)
	client.On("ChangeAvailability", ChangeAvailabilityHandler)
	client.On("GetLocalListVersion", GetLocalListVersionHandler)
	client.On("ChangeConfiguration", ChangeConfigurationHandler)
	
	cp, err := client.Start("ws://localhost:8999", "/ws")
	if err != nil {
		fmt.Printf("error dialing: %v\n", err)
		return
	}
	sendBootNotification(cp)
	defer cp.Shutdown()
	log.Debugf("charge point status %v", cp.IsConnected())
	select {}
}

func ChangeConfigurationHandler(cp *ocpp.ChargePoint, p ocpp.Payload) ocpp.Payload {
	req := p.(*v16.ChangeConfigurationReq)
	log.Debugf("ChangeConfigurationReq: %v\n", req)
	var res ocpp.Payload = &v16.ChangeConfigurationConf{
		Status: "Accepted",
	}
	return res
}

Later use
func ChangeAvailabilityHandler(cp *ocpp.ChargePoint, p ocpp.Payload) ocpp.Payload {
	req := p.(*v16.ChangeAvailabilityReq)
	log.Debugf("ChangeAvailability: %v\n", req)
	var res ocpp.Payload = &v16.ChangeAvailabilityConf{
		Status: "Accepted",
	}
	return res
}

func GetLocalListVersionHandler(cp *ocpp.ChargePoint, p ocpp.Payload) ocpp.Payload {
	req := p.(*v16.GetLocalListVersionReq)
	log.Debugf("GetLocalListVersionReq: %v\n", req)
	var res ocpp.Payload = &v16.GetLocalListVersionConf{
		ListVersion: 1,
	}
	return res
}

func sendBootNotification(c *ocpp.ChargePoint) {
	req := &v16.BootNotificationReq{
		ChargePointModel:  "client00",
		ChargePointVendor: "VendorX",
	}
	res, err := c.Call("BootNotification", req)
	if err != nil {
		fmt.Printf("error dialing: %v\n", err)
		return
	}
	fmt.Printf("BootNotificationConf: %v\n", res)
}


func sendAuthorize(c *ocpp.ChargePoint) {
	req := &v16.AuthorizeReq{
		IdTag: "safdasdfdsa",
	}
	res, err := c.Call("Authorize", req, 10)
	if err != nil {
		fmt.Printf("error dialing: %v\n", err)
		return
	}
	fmt.Printf("AuthorizeConf: %v\n", res)
}
```
After creating `*ocpp.Client` instance, register CS (Central System) initiated call handlers.
Making a call to CS is same as the above snippet where just call `cp.Call` method.
## Contributing

Contributions are always welcome!
Implementing higher versions of ocpp is highly appreciated!

See `CONTRIBUTING.md` for ways to get started.
