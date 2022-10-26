package main

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/aliml92/ocpp"
	v16 "github.com/aliml92/ocpp/v16"
	_ "net/http/pprof"

)


var csms *ocpp.Server

// log replaces standard log
var log *zap.SugaredLogger



func main() {

	go func() {
		log.Debug(http.ListenAndServe(":6060", nil))
	}()

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



	// register charge-point-initiated action handlers
	csms.On("BootNotification", BootNotificationHandler)
	csms.After("BootNotification", SendChangeConfigration)
	csms.Start("0.0.0.0:8999", "/ws/", nil)
	

}

func SendChangeConfigration(cp *ocpp.ChargePoint, payload ocpp.Payload) {
	var req ocpp.Payload = v16.ChangeConfigurationReq{
		Key: "WebSocketPingInterval",
		Value: "0",
	}
	time.Sleep(25 * time.Second)
	res, err := cp.Call("ChangeConfiguration", req)
	if err != nil {
		log.Debug(err)
	}
	cp.DisablePingPong()
	log.Debug(res)
}


func SendGetLocalListVersion(cp *ocpp.ChargePoint, payload ocpp.Payload) {
	var req ocpp.Payload = v16.GetLocalListVersionReq{}
	res, err := cp.Call("GetLocalListVersion", req)
	if err != nil {
		log.Debug(err)
	}
	log.Debug(res)
}


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
	time.Sleep(time.Second * 2)
	req := p.(*v16.BootNotificationReq)
	log.Debugf("\nid: %s\nBootNotification: %v", cp.Id, req)
	var res ocpp.Payload = &v16.BootNotificationConf{
		CurrentTime: time.Now().Format("2006-01-02T15:04:05.000Z"),
		Interval:    60 ,
		Status:      "Accepted",
	}
	return res
}

