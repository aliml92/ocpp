package ocpp

import (
	"sync"
	"time"
)





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

}

// create new CSMS instance actin1`g as main handler for ChargePoints
func NewServer() *Server {
	server = &Server{
		chargepoints:   make(map[string]*ChargePoint),
		actionHandlers: make(map[string]func(*ChargePoint, Payload) Payload),
		afterHandlers:  make(map[string]func(*ChargePoint, Payload)),
		ocppWait: ocppWait,
		writeWait: writeWait,
		pingWait: pigWait,
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
	delete(s.chargepoints, id)
	s.mu.Unlock()
}


func (s *Server) AddConn(cp *ChargePoint) {
	s.mu.Lock()
	server.chargepoints[cp.Id] = cp
	s.mu.Unlock()
}