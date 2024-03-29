package main

import (
	"fmt"


	"net/http"
	_ "net/http/pprof"

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
	go func() {
		log.Debugln(http.ListenAndServe("localhost:5050", nil))
	}()
	// initialize logger
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


