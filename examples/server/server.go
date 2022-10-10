package main

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/aliml92/ocpp"
	ocpplog "github.com/aliml92/ocpp/logger"
	v16 "github.com/aliml92/ocpp/v16"
	_ "net/http/pprof"

)

// csms serves as a main hub to set default configurations,
// keep track of connected charge points, register handlers
// for charge point initiated actions
var csms *ocpp.Server





// log replaces standard log
var log *zap.SugaredLogger


// initialize zap logger
// for deveplopment only
func initLogger() {
	logger, _ := zap.NewDevelopment()
	log = logger.Sugar()
}



// main function to bootstrap csms server
func main() {

	go func() {
		log.Debug(http.ListenAndServe(":6060", nil))
	}()

	// initialize logger
	initLogger()
	defer log.Sync()

	// set ocpp library's logger to zap logger
	ocpplog.SetLogger(log)





	// start csms server with default configurations
	csms = ocpp.NewServer()

	// custom timeout configuration 
	// config := ocpp.ServerTimeoutConfig{
	// 	OcppWait: 30 * time.Second,
	// 	WriteWait: 10 * time.Second,
	// 	PingWait: 30 * time.Second,
	// }
	// // set timeout configuration
	// csms.SetTimeoutConfig(config)


	// register charge-point-initiated action handlers
	csms.On("BootNotification", BootNotificationHandler)
	csms.After("BootNotification", SendGetLocalListVersion)
	csms.AddSubProtocol("ocpp1.6")
	csms.SetCheckOriginHandler(func(r *http.Request) bool {
		return true
	})
	csms.SetPreUpgradeHandler(customPreUpgradeHandler)
	csms.Start("0.0.0.0:8999", "/ws/", nil)

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
	// check if charge point is providing basic authentication
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

