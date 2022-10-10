package main

import (
	"fmt"
	"net/http"

	"github.com/aliml92/ocpp"
	ocpplog "github.com/aliml92/ocpp/logger"
	"github.com/aliml92/ocpp/v16"
	"go.uber.org/zap"
	_ "net/http/pprof"
)

var confData map[string]string

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
	ocpplog.SetLogger(log)
	confData = make(map[string]string)
	// create client
	client = ocpp.NewClient()
	id := "client00"
	client.SetID(id)
	client.AddSubProtocol("ocpp1.6")
	client.AddBasicAuth(id, "dummypass")
	client.On("ChangeAvailability", ChangeAvailabilityHandler)
	client.On("GetLocalListVersion", GetLocalListVersionHandler)
	
	cp, err := client.Start("ws://localhost:8999", "/ws")
	if err != nil {
		fmt.Printf("error dialing: %v\n", err)
		return
	}
	sendBootNotification(cp)
	defer cp.Shutdown()
	select {}
}




func ChangeConfigurationHandler(cp *ocpp.ChargePoint, p ocpp.Payload) ocpp.Payload {
	req := p.(*v16.ChangeConfigurationReq)
	log.Debugf("ChangeConfigurationReq: %v\n", req)
	confData[req.Key] = req.Value
	var res ocpp.Payload = &v16.ChangeConfigurationConf{
		Status: "Accepted",
	}
	return res
}


// Later use
func ChangeAvailabilityHandler(cp *ocpp.ChargePoint, p ocpp.Payload) ocpp.Payload {
	req := p.(*v16.ChangeAvailabilityReq)
	log.Debugf("ChangeAvailability: %v\n", req)
	var res ocpp.Payload = &v16.ChangeAvailabilityConf{
		Status: "Accepted",
	}
	return res
}


func GetLocalListVersionHandler(cp *ocpp.ChargePoint, p ocpp.Payload) ocpp.Payload {
	req :=p.(*v16.GetLocalListVersionReq)
	log.Debugf("GetLocalListVersionReq: %v\n", req)
	var res ocpp.Payload = &v16.GetLocalListVersionConf{
		ListVersion: 1,		
	}
	return res
}











func sendBootNotification(c *ocpp.ChargePoint){
	req := &v16.BootNotificationReq{
		ChargePointModel: "client00",
		ChargePointVendor: "VendorX",
	}
	res, err := c.Call("BootNotification", req)
	if err != nil {
		fmt.Printf("error dialing: %v\n", err)
		return 
	}
	fmt.Printf("BootNotificationConf: %v\n", res)
} 

