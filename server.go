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

	callQuequeSize int
}

// create new CSMS instance acting as main handler for ChargePoints
func NewServer() *Server {
	server = &Server{
		chargepoints:   make(map[string]*ChargePoint),
		actionHandlers: make(map[string]func(*ChargePoint, Payload) Payload),
		afterHandlers:  make(map[string]func(*ChargePoint, Payload)),
		ocppWait: ocppWait,
		writeWait: writeWait,
		pingWait: pingWait,
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


func (s *Server) IsConnected(id string) bool {
	if cp, ok := s.chargepoints[id]; ok {
		return cp.connected
	}
	return false 
}  

func (s *Server) getHandler(action string) func(*ChargePoint, Payload) Payload {
	return s.actionHandlers[action]
}

func (s *Server) getAfterHandler(action string) func(*ChargePoint, Payload) {
	return s.afterHandlers[action]
}


func (s *Server) Delete(id string) {
	s.mu.Lock() 
	if cp, ok := s.chargepoints[id]; ok {
		cp.connected = false
	}
	delete(s.chargepoints, id)
	s.mu.Unlock()
}



func (s *Server) Store(cp *ChargePoint) {
	s.mu.Lock()
	server.chargepoints[cp.Id] = cp
	s.mu.Unlock()
}


func (s *Server) Load(id string) (*ChargePoint, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cp, ok := s.chargepoints[id]; ok {
		fmt.Printf("ChargePoint with id: %s exist\n", cp.Id)
		return cp, true
	}
	return nil, false
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
	server.Store(cp)
}

func (s *Server) SetCallQueueSize(size int) {
	s.callQuequeSize = size
}  

func (s *Server) getCallQueueSize() int {
	s.mu.Lock()
	size := s.callQuequeSize
	s.mu.Unlock()
	return size
}