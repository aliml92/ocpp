package ocpp

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/aliml92/ocpp/log"
	"github.com/aliml92/ocpp/v16"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)


var validate = v16.Validate
var server *Server
var client *Client

// add more



// Payload used as a container is for both Call and CallResult' Payload
type Payload interface{}


type TimeoutConfig struct {
	OcppTimeout			time.Duration	// ocpp response timeout in seconds
	WriteWait			time.Duration	// time allowed to write a message to the peer
	PingWait			time.Duration   // time allowed to read the next pong message from the peer
	PongWait       		time.Duration   // pong wait in seconds
	PingPeriod      	time.Duration   // ping period in seconds
}




// ChargePoint Represents a connected ChargePoint (also known as a Charging Station)
type ChargePoint struct {
	Proto           string                                   // OCPP protocol version
	Conn            *websocket.Conn                          // the websocket connection
	Id              string                                   // chargePointId
	Out             chan *[]byte                             // outgoing message channel
	In              chan *[]byte                             // incoming message channel
	Mu              sync.Mutex                               // mutex ensures that only one message is sent at a time
	Cr              chan *CallResult
	Ce              chan *CallError
	Extras          map[string]interface{}
	TimeoutConfig   TimeoutConfig
	isServer        bool 
	PingIn          chan []byte                              // ping in channel              #server specific
	srvRClose		chan struct{}							 // channel to close connection
	srvWClose   	chan struct{}
	cliRClose    	chan struct{}
	cliWClose     	chan struct{}   
	closeAck     	chan struct{}
	forceRClose		chan struct{}
	forceWClose		chan struct{}	
}


func (cp *ChargePoint) Shutdown(timeout time.Duration, force bool) error {
	b := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	err := cp.Conn.WriteControl(websocket.CloseMessage, b, time.Now().Add(time.Second))
	if err != nil && err != websocket.ErrCloseSent {
		log.L.Error(err)
		return cp.Conn.Close()
	}
	if force {
		cp.forceRClose <- struct{}{}
		cp.forceWClose <-struct{}{}
		return cp.Conn.Close()
	}

	if cp.isServer {
		cp.srvWClose <- struct{}{}
		select {
		case <- cp.closeAck:
			log.L.Debug("Close Singal Received")
			cp.srvRClose <- struct{}{}	
		case <- time.After(timeout):
			log.L.Debug("Timeout Happened")
			cp.srvRClose <- struct{}{}
		}
	} else {
		cp.cliWClose <- struct{}{}
		select {
		case <- cp.closeAck:
			log.L.Debug("Close Singal Received")
			cp.cliRClose <- struct{}{}	
		case <- time.After(timeout):
			log.L.Debug("Timeout Happened")
			cp.cliRClose <- struct{}{}
		}
	} 
	return cp.Conn.Close()
	
}





func (cp *ChargePoint) SetTimeoutConfig(config TimeoutConfig) {
	cp.TimeoutConfig = config
}


// CSMS acts as main handler for ChargePoints
type Server struct {
	Chargepoints sync.Map                               			// keeps track of all connected ChargePoints
	ActionHandlers map[string]func(*ChargePoint, Payload) Payload   // register implemented action handler functions
 	AfterHandlers map[string]func(*ChargePoint, Payload)            // register after-action habdler functions 
	TimeoutConfig TimeoutConfig            							// timeout configuration                            
}


// create new CSMS instance acting as main handler for ChargePoints
func NewServer() *Server {
	server = &Server{
		Chargepoints: 	sync.Map{},
		ActionHandlers: make(map[string]func(*ChargePoint, Payload) Payload),
		AfterHandlers: 	make(map[string]func(*ChargePoint, Payload)),
		TimeoutConfig: TimeoutConfig{
			OcppTimeout: 30 * time.Second,
			WriteWait: 10 * time.Second,
			PingWait: 30 * time.Second,
		},
	}
	return server
}

func (c *Server) SetTimeoutConfig(config TimeoutConfig) {
	c.TimeoutConfig = config
}

func (c *ChargePoint) getReadTimeout() time.Time {
	if c.TimeoutConfig.PingWait == 0 {
		return time.Time{}
	}
	return time.Now().Add(c.TimeoutConfig.PingWait)
}

type Client struct {
	ActionHandlers map[string]func(*ChargePoint, Payload) Payload   // register implemented action handler functions
 	AfterHandlers map[string]func(*ChargePoint, Payload)            // register after-action habdler functions 
	TimeoutConfig 	TimeoutConfig							        // timeout configuration
}


// create new Client instance 
func NewClient() *Client {
	client = &Client{
		ActionHandlers: make(map[string]func(*ChargePoint, Payload) Payload),
		AfterHandlers: 	make(map[string]func(*ChargePoint, Payload)),
		TimeoutConfig: TimeoutConfig{
			OcppTimeout: 	30 * time.Second,  
			WriteWait: 		10 * time.Second,
			PongWait: 		30 * time.Second,
			PingPeriod:   	25 * time.Second,
		},
	}
	return client
}


func (c *Client) SetTimeoutConfig(config TimeoutConfig) {
	c.TimeoutConfig = config
}

func (c *Client) getReadTimeout() time.Time {
	if c.TimeoutConfig.PongWait == 0 {
		return time.Time{}
	}
	return time.Now().Add(c.TimeoutConfig.PongWait)
}




// register action handler function
func (csms *Server) On(action string, f func(*ChargePoint, Payload) Payload) *Server {
	csms.ActionHandlers[action] = f
	return csms
}


// register after-action handler function
func (csms *Server) After(action string, f func(*ChargePoint,  Payload)) *Server {
	csms.AfterHandlers[action] = f
	return csms
}




// register action handler function
func (c *Client) On(action string, f func(*ChargePoint, Payload) Payload) *Client {
	c.ActionHandlers[action] = f
	return c
}


// register after-action handler function
func (c *Client) After(action string, f func(*ChargePoint,  Payload)) *Client {
	c.AfterHandlers[action] = f
	return c
}




// websocket reader to receive messages
func (cp *ChargePoint) cliReader() {
	cp.Conn.SetPongHandler(func(appData string) error {
		log.L.Debugf("Pong received: %v", appData)
		return cp.Conn.SetReadDeadline(client.getReadTimeout())
	})
	// cp.Conn.SetCloseHandler(func(code int, text string) error {
	// 	log.L.Debugf("code received: %v", code)
	// 	log.L.Debugf("text received: %v", text)
	// 	return cp.Conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.), time.Now().Add(time.Second))
	// })
	defer func() {
		cp.Conn.Close()
	}()
	for {
		select {
		case <- cp.forceRClose:
			return
		case <- cp.cliRClose:
			// TODO
			return
		default:
			messageType, msg, err := cp.Conn.ReadMessage()
			log.L.Debugf("messageType: %d , err: %v", messageType, err)
			if err != nil {
				var closeErr *websocket.CloseError
				if errors.As(err, &closeErr){
					log.L.Error(closeErr)
					if closeErr.Code == 1000 {
						cp.closeAck <- struct{}{}
					}
				}
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
					log.L.Error(err)
				}				
				break
			}
			call, callResult, callError, err := unpack(&msg)
			if err != nil {
				cp.Out <- call.createCallError(err)
				log.L.Error(err)
			}
			if call != nil {
				handler, ok := client.ActionHandlers[call.Action]
				if ok {
					responsePayload := handler(cp, call.Payload)
					log.L.Debugf("[ocpp] %v", responsePayload)
					err = validate.Struct(responsePayload)
					if err != nil {
						log.L.Errorf("[ocpp] %v", err)
					} else {
						cp.Out <- call.createCallResult(responsePayload)
						time.Sleep(time.Second)
						if afterHandler, ok := client.AfterHandlers[call.Action]; ok {
							go afterHandler(cp, call.Payload)
						}
					}
				} else {
					var err error = &OCPPError{
						id:    call.UniqueId,
						code:  "NotSupported",
						cause: fmt.Sprintf("Action %s is not supported", call.Action),
					}
					cp.Out <- call.createCallError(err)
					log.L.Errorf("No handler for action %s", call.Action)
				}
			}
			if callResult != nil {
				log.L.Debug(callResult)
				cp.Cr <- callResult
			}
			if callError != nil {
				log.L.Debug(callError)
				cp.Ce <- callError
			}				  
		}

	}
}



// websocket writer to send messages
func (cp *ChargePoint) cliWriter() {
	switch client.TimeoutConfig.PingPeriod {
	case 0:
		for {
			select {
			case <- cp.forceWClose:
				return
			case <- cp.cliWClose:
				return	
			case message, ok := <-cp.Out:
				_ = cp.Conn.SetWriteDeadline(time.Now().Add(client.TimeoutConfig.WriteWait))			
				if !ok {
					cp.Conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				w, err := cp.Conn.NextWriter(websocket.TextMessage)
				if err != nil {
					log.L.Error(err)
					return
				}
				_, err = w.Write(*message)
				if err != nil {
					log.L.Error(err)
					return
				}
				if err := w.Close(); err != nil {
					log.L.Error(err)
					return
				}					
			}
		}
	default:
		defer cp.Conn.Close()
		ticker := time.NewTicker(client.TimeoutConfig.PingPeriod)
		defer ticker.Stop()
		for {
			select {
			case message, ok := <-cp.Out:
				_ = cp.Conn.SetWriteDeadline(time.Now().Add(client.TimeoutConfig.WriteWait))			
				if !ok {
					cp.Conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				w, err := cp.Conn.NextWriter(websocket.TextMessage)
				if err != nil {
					log.L.Error(err)
					return
				}
				_, err = w.Write(*message)
				if err != nil {
					log.L.Error(err)
					return
				}
				if err := w.Close(); err != nil {
					log.L.Error(err)
					return
				}
			case <-ticker.C:
				_ = cp.Conn.SetWriteDeadline(time.Now().Add(client.TimeoutConfig.WriteWait))
				if err := cp.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					log.L.Error(err)
					return
				}
			case <- cp.forceWClose:
				return
			case <- cp.cliWClose:
				return			
			}
		}					
	} 
}




// 
func (cp *ChargePoint) srvReader() {
	cp.Conn.SetPingHandler(func(appData string) error {
		cp.PingIn <- []byte(appData)
		log.L.Debugf("ping received: %v", appData)
		i := cp.getReadTimeout()
		log.L.Debugf("second read deadline: %v", i)
		return cp.Conn.SetReadDeadline(i)
	})
	// cp.Conn.SetCloseHandler(func(code int, text string) error {
	// 	log.L.Debugf("code received: %v", code)
	// 	log.L.Debugf("text received: %v", text)
	// 	return nil
	// })	
	defer func() {
		cp.Conn.Close()
	}()
	for {
		select {
		case <- cp.forceRClose:
			return	
		case <-cp.srvRClose:
			log.L.Debug("stop singal received")
			server.Chargepoints.Delete(cp.Id)
			return
		default:
			messageType, msg, err := cp.Conn.ReadMessage()
			log.L.Debugf("messageType: %d ", messageType)
			if err != nil {
				var closeErr *websocket.CloseError
				if errors.As(err, &closeErr) {
					log.L.Debugf("close error occured: code: %v, text: %v", closeErr.Code, closeErr.Text)
					if closeErr.Code == 1000 {
						log.L.Debug("closeErr received")
						cp.closeAck <- struct{}{}
					}
				}
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
					log.L.Error(err)
					server.Chargepoints.Delete(cp.Id)
					log.L.Debugf("charge point with id %s deleted", cp.Id)
				}
				break
			}
			call, callResult, callError, err := unpack(&msg)
			if err != nil {
				cp.Out <- call.createCallError(err)
				log.L.Error(err)
			}
			if call != nil {
				handler, ok := server.ActionHandlers[call.Action]
				if ok {
					responsePayload := handler(cp, call.Payload)
					err = validate.Struct(responsePayload)
					if err != nil {
						log.L.Error(err)
					} else {
						cp.Out <- call.createCallResult(responsePayload)
						time.Sleep(time.Second)
						if afterHandler, ok := server.AfterHandlers[call.Action]; ok {
							go afterHandler(cp, call.Payload)
						}
					}
				} else {
					var err error = &OCPPError{
						id:    call.UniqueId,
						code:  "NotSupported",
						cause: fmt.Sprintf("action %s is not supported", call.Action),
					}
					cp.Out <- call.createCallError(err)
					log.L.Debugf("no handler for action %s", call.Action)
				}
			}
			if callResult != nil {
				log.L.Debugf("call result received: %v", callResult)
				cp.Cr <- callResult
			}
			if callError != nil {
				log.L.Debugf("call error received: %v", callError)
				cp.Ce <- callError
			}		
		}	
	}
}



// websocket writer to send messages
func (cp *ChargePoint) srvWriter() {
	for {
		select {
		case message, ok := <-cp.Out:
			err := cp.Conn.SetWriteDeadline(time.Now().Add(cp.TimeoutConfig.WriteWait))
			if err != nil {
				log.L.Error(err)
			}
			if !ok {
				err := cp.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.L.Error(err)
				}
				log.L.Debug("close message sent")
				return
			}
			w, err := cp.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.L.Error(err)
				return
			}
			n, err := w.Write(*message)
			if err != nil {
				log.L.Error(err)
				return
			}
			if err := w.Close(); err != nil {
				server.Chargepoints.Delete(cp.Id)
				log.L.Error(err)
				return
			}
			log.L.Debugf("text message sent %d", n)
		case <-cp.PingIn:
			err := cp.Conn.SetWriteDeadline(time.Now().Add(cp.TimeoutConfig.WriteWait))
			if err != nil {
				log.L.Error(err)
			}
			err = cp.Conn.WriteMessage(websocket.PongMessage, []byte{})
			if err != nil {
				log.L.Error(err)
				return
			}
			log.L.Debug("Pong sent")
		case <- cp.forceWClose:
			return
		case <- cp.srvWClose:
			return			
		}
	}
}


// Call sends a message to peer
func (cp *ChargePoint) Call(action string, p Payload) (Payload, error) {
	id := uuid.New().String()
	call := [4]interface{}{
		2,
		id,
		action,
		p,
	}
	raw, _ := json.Marshal(call)
	cp.Mu.Lock()
	defer cp.Mu.Unlock()
	cp.Out <- &raw
	callResult, callError, err := cp.waitForResponse(id)
	if callResult != nil {
		resPayload, err := unmarshalConf(action, callResult.Payload)
		if err != nil {
			log.L.Error(err)
			return nil, err
		}
		return resPayload, nil
	}
	if callError != nil {
		return nil, callError
	}
	return nil, err
}


// waitForResponse waits for a response to a Call with id
func (cp *ChargePoint) waitForResponse(uniqueId string) (*CallResult, *CallError, error) {
	deadline := time.Now().Add(cp.TimeoutConfig.OcppTimeout)
	for {
		select {
		case r1 := <-cp.Cr:
			if r1.UniqueId == uniqueId {
				return r1, nil, nil
			}
		case r2 := <-cp.Ce:
			if r2.UniqueId == uniqueId {
				return nil, r2, nil
			}
		case <-time.After(time.Until(deadline)):
			return nil, nil, &TimeoutError{
				Message: fmt.Sprintf("timeout of %s sec for response to Call with id: %s passed", cp.TimeoutConfig.OcppTimeout, uniqueId),
			}
		}
	}
}


// NewChargepoint creates a new ChargePoint
func NewChargePoint(conn *websocket.Conn, id, proto string, isServer bool) *ChargePoint {
	cp := &ChargePoint{
		Proto:           proto,
		Conn:            conn,
		Id:              id,
		Out:             make(chan *[]byte),
		In:              make(chan *[]byte),
		Cr:              make(chan *CallResult),
		Ce:              make(chan *CallError),
		Extras: 		 make(map[string]interface{}),
		closeAck: 	 	 make(chan struct{}, 1),
		forceRClose:     make(chan struct{}, 1),
		forceWClose:     make(chan struct{}, 1),
	}

	if isServer {
		cp.TimeoutConfig = server.TimeoutConfig
		cp.PingIn = make(chan []byte)
		cp.isServer = true
		cp.srvRClose = 	make(chan struct{}, 1)
		cp.srvWClose = make(chan struct{}, 1)
		server.Chargepoints.Store(id, cp)
		go cp.srvReader()
		go cp.srvWriter()
	} else {
		cp.TimeoutConfig = client.TimeoutConfig
		cp.cliRClose = make(chan struct{},1)
		cp.cliWClose = make(chan struct{},1)
		go cp.cliReader()
		go cp.cliWriter()
	} 
	return cp
}
