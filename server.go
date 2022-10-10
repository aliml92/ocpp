package ocpp

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)



var server *Server

type ServerTimeoutConfig struct {
	// ocpp response timeout in seconds
	OcppWait 	time.Duration 

	// time allowed to write a message to the peer
	WriteWait   time.Duration 

	// time allowed to read the next pong message from the peer
	PingWait    time.Duration 
}



// Server type representes csms server
type Server struct {
	// keeps track of all connected ChargePoints
	chargepoints   map[string]*ChargePoint
	
	// register implemented action handler functions
	actionHandlers map[string]func(*ChargePoint, Payload) Payload
	
	// register after-action habdler functions
	afterHandlers  map[string]func(*ChargePoint, Payload)
	
	// timeout configuration
	ocppWait 	time.Duration
	
	
	writeWait   time.Duration


	pingWait    time.Duration

	mu 			sync.Mutex
	
	upgrader    websocket.Upgrader

	preUpgradeHandler func(w http.ResponseWriter, r *http.Request) bool

	returnError func(err error)
}

// create new CSMS instance acting as main handler for ChargePoints
func NewServer() *Server {
	server = &Server{
		chargepoints:   make(map[string]*ChargePoint),
		actionHandlers: make(map[string]func(*ChargePoint, Payload) Payload),
		afterHandlers:  make(map[string]func(*ChargePoint, Payload)),
		ocppWait: ocppWait,
		writeWait: writeWait,
		pingWait: pigWait,
		upgrader: websocket.Upgrader{
			Subprotocols: []string{},
		},
	}
	return server
}


func (s *Server) SetTimeoutConfig(config ServerTimeoutConfig) {
	s.ocppWait = config.OcppWait
	s.writeWait = config.WriteWait
	s.pingWait = config.PingWait
}


// register action handler function
func (s *Server) On(action string, f func(*ChargePoint, Payload) Payload) *Server {
	s.actionHandlers[action] = f
	return s
}

// register after-action handler function
func (s *Server) After(action string, f func(*ChargePoint, Payload)) *Server {
	s.afterHandlers[action] = f
	return s
}

func (s *Server) getHandler(action string) func(*ChargePoint, Payload) Payload {
	return s.actionHandlers[action]
}

func (s *Server) getAfterHandler(action string) func(*ChargePoint, Payload) {
	return s.afterHandlers[action]
}


func (s *Server) DeleteConn(id string) {
	s.mu.Lock()
	// check if DeleteConn deletes 
	if _, ok := s.chargepoints[id]; ok {
		fmt.Printf("ChargePoint with id: %s exist\n", id)
	}
	delete(s.chargepoints, id)
	if _, ok := s.chargepoints[id]; !ok {
		fmt.Printf("ChargePoint with id: %s deleted\n", id)
	}
	s.mu.Unlock()
}


func (s *Server) AddConn(cp *ChargePoint) {
	s.mu.Lock()
	if _, ok := s.chargepoints[cp.Id]; !ok {
		fmt.Printf("ChargePoint with id: %s does not exist\n", cp.Id)
	}
	server.chargepoints[cp.Id] = cp
	if _, ok := s.chargepoints[cp.Id]; ok {
		fmt.Printf("ChargePoint with id: %s added\n", cp.Id)
	}
	s.mu.Unlock()
}


func (s *Server) AddSubProtocol(protocol string) {
	for _, p := range server.upgrader.Subprotocols {
		if p == protocol {
			return
		}
	}
	s.upgrader.Subprotocols = append(s.upgrader.Subprotocols, protocol)
}


func (s *Server) SetCheckOriginHandler(f func(r *http.Request) bool) {
	s.upgrader.CheckOrigin = f
}


func (s *Server) SetPreUpgradeHandler(f func(w http.ResponseWriter, r *http.Request) bool) {
	s.preUpgradeHandler = f
}


// TODO: add more functionality
func (s *Server) Start(addr string, path string, handler func(http.ResponseWriter, *http.Request)) {
	if handler != nil {
		http.HandleFunc(path, handler)
	} else {
		http.HandleFunc(path, defaultWebsocketHandler)
	}
	http.ListenAndServe(addr, nil)
}

func defaultWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	preCheck := server.preUpgradeHandler
	if preCheck != nil {
		if preCheck(w, r) {
			upgrade(w, r)
		} else {
			server.returnError(errors.New("cannot start server"))
		}
	} else {
		upgrade(w, r)
	}
} 


func upgrade(w http.ResponseWriter, r *http.Request) {
	c, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		server.returnError(err)
		return
	}
	p := strings.Split(r.URL.Path, "/")
	id := p[len(p)-1]
	cp := NewChargePoint(c, id, c.Subprotocol(), true)
	server.AddConn(cp)
}

