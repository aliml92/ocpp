package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/aliml92/ocpp"
	ocpplog "github.com/aliml92/ocpp/logger"
	"github.com/aliml92/ocpp/v16"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	_ "net/http/pprof"
)

var confData map[string]string

var client ocpp.Client

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
	client = *ocpp.NewClient()
	client.On("ChangeAvailability", ChangeAvailabilityHandler)
	client.On("GetLocalListVersion", GetLocalListVersionHandler)
	
	cp, err := start()
	if err != nil {
		fmt.Printf("error dialing: %v\n", err)
	}
	defer cp.Shutdown()

	time.Sleep(time.Second * 3)
	sendBootNotification(cp)
	time.Sleep(3600 * time.Second)
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









// utility functions
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func start() (*ocpp.ChargePoint, error) {
	chargePointId := "client00" // also username for basic auth
	password, ok := confData["AuthorizationKey"]
	if !ok {
		password = "dummypass"
	}
	fmt.Printf("AuthorizationKey is set to %s\n", password)
	url := fmt.Sprintf("ws://localhost:8999/ws/%s", chargePointId)
	header := http.Header{
		"Sec-WebSocket-Protocol": []string{"ocpp1.6"},
		"Authorization":          []string{"Basic " + basicAuth(chargePointId, password)},
	}
	c, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		return nil, err
	}
	cp := ocpp.NewChargePoint(c, chargePointId, "ocpp1.6", false)
	return cp, nil

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

