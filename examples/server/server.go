package main

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/aliml92/ocpp"
	ocpplog "github.com/aliml92/ocpp/logger"
	v16 "github.com/aliml92/ocpp/v16"
	"github.com/gorilla/websocket"
	_ "net/http/pprof"

)

// csms serves as a main hub to set default configurations,
// keep track of connected charge points, register handlers
// for charge point initiated actions
var csms *ocpp.Server


// cp is a ChangePoint which handles single websocket connection for
// for a connected charge point 
// var cp *ocpp.ChargePoint


// upgrader for websocket connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"ocpp1.6"},
}


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
	
	http.HandleFunc("/", wsHandler)
	
	// start csms server
	http.ListenAndServe("0.0.0.0:8999", nil)


}


func SendGetLocalListVersion(cp *ocpp.ChargePoint, payload ocpp.Payload) {
	var req ocpp.Payload = v16.GetLocalListVersionReq{}
	res, err := cp.Call("GetLocalListVersion", req)
	if err != nil {
		log.Debug(err)
	}
	log.Debug(res)
}


// websocket handler
// handles incoming websocket connections from charge points
func wsHandler(w http.ResponseWriter, r *http.Request) {
	
	// ocpp protocol supported by current csms server
	ssp := "ocpp1.6"

	// check if charge point supports current ocpp protocol
	subProtocol := r.Header.Get("Sec-WebSocket-Protocol")
	if subProtocol == "" {
		log.Debug("client hasn't requested any Subprotocol. Closing Connection")
		return
	}
	if !strings.Contains(subProtocol, ssp) {
		log.Debug("client has requested an unsupported Subprotocol. Closing Connection")
		return
	}


	// check if charge point is providing basic authentication
	u, p, ok := r.BasicAuth()
	if !ok {
		log.Debug("error parsing basic auth")
		w.WriteHeader(401)
		return
	}
	id := strings.Split(r.URL.Path, "/")[2]
	log.Debugf("%s is trying to connect with %s:%s", id, u, p)
	if u != id {
		log.Debug("username provided is correct: %s", u)
		w.WriteHeader(401)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("cannot upgrade to websocket")
		log.Error(err)
		return
	}


	// create charge point
	// with successfully connected websocket connection
	// unique id 
	// choice of ocpp protocol between charge point and csms server
	// boolean value meaning that this charge point represents client or server side charge point
	_ = ocpp.NewChargePoint(c, id, ssp, true)
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

