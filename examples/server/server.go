package main

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	// _ "net/http/pprof"

	"github.com/aliml92/ocpp"
	v16 "github.com/aliml92/ocpp/v16"
)


var csms *ocpp.Server

// log replaces standard log
var log *zap.SugaredLogger



func main() {

	// go func() {
	// 	log.Debug(http.ListenAndServe(":6060", nil))
	// }()

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
	// csms.On("Authorize", AuthorizationHandler)
	csms.After("BootNotification", SendChangeConfigration)
	// go func() {
	// 	log.Debugln("sleeping 20s")
	// 	time.Sleep(20 * time.Second)
	// 	cp, ok := csms.Load("client00")
	// 	if !ok {
	// 		log.Debugln("charge point found")
	// 		return
	// 	}
	// 	var arr [6]ocpp.Payload
	// 	var req0 ocpp.Payload = v16.ChangeConfigurationReq{
	// 		Key: "WebSocketPingInterval",
	// 		Value: "30",
	// 	}
	// 	arr[0] = req0
	// 	var req1 ocpp.Payload = v16.ChangeConfigurationReq{
	// 		Key: "BlinkRepeat",
	// 		Value: "5",
	// 	}
	// 	arr[1] = req1
	// 	var req2 ocpp.Payload = v16.ChangeConfigurationReq{
	// 		Key: "ConnectionTimeOut",
	// 		Value: "300",
	// 	}
	// 	arr[2] = req2
	// 	arr[3] = req0
	// 	arr[4] = req1
	// 	arr[5] = req2
	// 	time.Sleep(2 * time.Second)
	// 	for i:=0; i < 5; i++ {
	// 		go func(idx int){
	// 			res, err := cp.Call("ChangeConfiguration", arr[idx], 10)
	// 			if err != nil {
	// 				log.Debug(err)
	// 			}
	// 			log.Debug(res)
	// 		}(i)
	// 	}
	// 	res, err := cp.Call("ChangeConfiguration", arr[5], 10)
	// 	if err != nil {
	// 		log.Debug(err)
	// 	}
	// 	// cp.ResetPingPong(30)
	// 	// cp.ResetPingPong(0)
	// 	// cp.EnableServerPing(10)
	// 	log.Debug(res)
	// }()
	csms.Start("0.0.0.0:8999", "/ws/", nil)
	

}

func SendChangeConfigration(cp *ocpp.ChargePoint, payload ocpp.Payload) {
	var arr [32]ocpp.Payload
	var req ocpp.Payload = v16.ChangeConfigurationReq{
		Key: "WebSocketPingInterval",
		Value: "30",
	}
	for i:=0; i < 30; i++ {
		arr[i] = req
		go func(idx int){
			res, err := cp.Call("ChangeConfiguration", arr[idx])
			if err != nil {
				log.Debug(err)
			}
			log.Debug(res)
		}(i)
	}
	arr[31] = req
	res, err := cp.Call("ChangeConfiguration", arr[31])
	if err != nil {
		log.Debug(err)
	}

	// cp.ResetPingPong(30)
	// cp.ResetPingPong(0)
	// cp.EnableServerPing(10)
	log.Debug(res)
}


// func SendGetLocalListVersion(cp *ocpp.ChargePoint, payload ocpp.Payload) {
// 	var req ocpp.Payload = v16.GetLocalListVersionReq{}
// 	res, err := cp.Call("GetLocalListVersion", req, 10)
// 	if err != nil {
// 		log.Debug(err)
// 	}
// 	log.Debug(res)
// }


// 
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

// func AuthorizationHandler(cp *ocpp.ChargePoint, p ocpp.Payload) ocpp.Payload {
// 	time.Sleep(time.Second * 2)
// 	req := p.(*v16.AuthorizeReq)
// 	log.Debugf("\nid: %s\nAuthorizeReq: %v", cp.Id, req)
// 	var res ocpp.Payload = &v16.AuthorizeConf{
// 		IdTagInfo: v16.IdTagInfo{
// 			Status: "Accepted",
// 		},
// 	}
// 	return res
// }