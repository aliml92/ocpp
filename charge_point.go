package ocpp

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aliml92/ocpp/v16"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)


var validate = v16.Validate
var csms *CSMS
var client *Client
// add more


// Payload used as a container is for both Call and CallResult' Payload
type Payload interface{}



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
	ResTimeout 	    time.Duration                            // ocpp response timeout in seconds
	WriteWait       time.Duration                            // write wait in seconds  
	PingWait        time.Duration                            // ping wait in seconds         #server specific 
	PongWait        time.Duration                            // pong wait in seconds         #client specific
	PingPeriod      time.Duration                            // ping period in seconds       #client specific
	PingIn          chan []byte                              // ping in channel              #server specific
}




// CSMS acts as main handler for ChargePoints
type CSMS struct {
	ChargePoints sync.Map                               			// keeps track of all connected ChargePoints
	ActionHandlers map[string]func(*ChargePoint, Payload) Payload   // register implemented action handler functions
 	AfterHandlers map[string]func(*ChargePoint, Payload)            // register after-action habdler functions 
	ResTimeout time.Duration                                        // ocpp response timeout in seconds
	WriteWait time.Duration                                         // time allowed to write a message to the peer
	PingWait time.Duration                                          // time allowed to read the next pong message from the peer
}


func (c *CSMS) SetTimeoutConfig(resTimeout, writeWait, pingWait time.Duration) {
	c.ResTimeout = resTimeout
	c.WriteWait = writeWait
	c.PingWait = pingWait
}

func (c *CSMS) getReadTimeout() time.Time {
	if c.PingWait == 0 {
		log.Println("PingWait is 0")
		return time.Time{}
	}
	return time.Now().Add(c.PingWait)
}

type Client struct {
	ActionHandlers map[string]func(*ChargePoint, Payload) Payload   // register implemented action handler functions
 	AfterHandlers map[string]func(*ChargePoint, Payload)            // register after-action habdler functions 
	ResTimeout time.Duration                                        // ocpp response timeout in seconds
	WriteWait       time.Duration                            		// write wait in seconds
	PongWait        time.Duration                            		// pong wait in seconds
	PingPeriod      time.Duration                            		// ping period in seconds
}


func (c *Client) SetTimeoutConfig(resTimeout, writeWait, pongWait, pingPeriod time.Duration) {
	c.ResTimeout = resTimeout
	c.WriteWait = writeWait
	c.PongWait = pongWait
	c.PingPeriod = pingPeriod
}


func (c *Client) getReadTimeout() time.Time {
	if c.PongWait == 0 {
		log.Println("PongWait is 0")
		return time.Time{}
	}
	return time.Now().Add(c.PongWait)
}

// register action handler function
func (csms *CSMS) On(action string, f func(*ChargePoint, Payload) Payload) *CSMS {
	csms.ActionHandlers[action] = f
	return csms
}


// register after-action handler function
func (csms *CSMS) After(action string, f func(*ChargePoint,  Payload)) *CSMS {
	csms.AfterHandlers[action] = f
	return csms
}



// create new CSMS instance acting as main handler for ChargePoints
func NewCSMS() *CSMS {
	csms = &CSMS{
		ChargePoints: 	sync.Map{},
		ActionHandlers: make(map[string]func(*ChargePoint, Payload) Payload),
		AfterHandlers: 	make(map[string]func(*ChargePoint, Payload)),
		ResTimeout: 	30 * time.Second ,
		WriteWait:    	10 * time.Second ,
		PingWait:     	60 * time.Second ,
	}
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


// create new Client instance 
func NewClient() *Client {
	client = &Client{
		ActionHandlers: make(map[string]func(*ChargePoint, Payload) Payload),
		AfterHandlers: 	make(map[string]func(*ChargePoint, Payload)),
		ResTimeout: 	30 * time.Second ,
		WriteWait:    	10 * time.Second ,
		PongWait:     	60 * time.Second ,
		PingPeriod:   	54 * time.Second ,
	}
	return client
}


// websocket reader to receive messages
func (cp *ChargePoint) reader() {
	_ = cp.Conn.SetReadDeadline(client.getReadTimeout())
	defer func() {
		cp.Conn.Close()
	}()
	cp.Conn.SetPongHandler(func(appData string) error {
		// print appData
		log.Printf("Pong received: %v", appData)
		return cp.Conn.SetReadDeadline(client.getReadTimeout())
	})
	for {
		log.Printf("Waiting for message[client side]")
		_, msg, err := cp.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Printf("[WEBSOCKET | ERROR] %v", err)
			}
			break
		}
		log.Printf("[WEBSOCKET | RECEIVED | CLIENT ] %v", msg)
		call, callResult, callError, err := unpack(&msg)
		if err != nil {
			cp.Out <- call.createCallError(err)
			log.Printf("[ERROR | MSG] %v", err)
		}
		if call != nil {
			handler, ok := client.ActionHandlers[call.Action]
			if ok {
				responsePayload := handler(cp, call.Payload)
				log.Printf("[WEBSOCKET | PAYLOAD | CLIENT ] %v", responsePayload)
				err = validate.Struct(responsePayload)
				if err != nil {
					// TODO simply log the error
					log.Printf("[ERROR | VALIDATION] %v", err)
				} else {
					// lock to ensure that only one message is sent at a time
					cp.Out <- call.createCallResult(responsePayload)
					// sleep for a bit to make sure the message is sent
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
				log.Printf("[ERROR | MSG] No handler for action %s", call.Action)
			}
		}
		if callResult != nil {
			log.Printf("[WEBSOCKET | CALL_RESULT | CLIENT] %v", callResult)
			cp.Cr <- callResult
		}
		if callError != nil {
			log.Printf("[WEBSOCKET | CALL_ERROR | CLIENT ] %v", callError)
			cp.Ce <- callError
		}

	}
}



// websocket writer to send messages
func (cp *ChargePoint) writer() {
	if client.PingPeriod == 0 {
		for {
			message, ok := <-cp.Out
			_ = cp.Conn.SetWriteDeadline(time.Now().Add(client.WriteWait))			
			if !ok {
				cp.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := cp.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("[WEBSOCKET | ERROR] %v", err)
				return
			}
			i, err := w.Write(*message)
			if err != nil {
				return
			}
			log.Printf("[WEBSOCKET | SENT | CLIENT ] %v", i)
			if err := w.Close(); err != nil {
				return
			}	
		}
	} else {
		defer cp.Conn.Close()
		ticker := time.NewTicker(client.PingPeriod)
		defer ticker.Stop()
		for {
			select {
			case message, ok := <-cp.Out:
				_ = cp.Conn.SetWriteDeadline(time.Now().Add(client.WriteWait))			
				if !ok {
					cp.Conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				w, err := cp.Conn.NextWriter(websocket.TextMessage)
				if err != nil {
					log.Printf("[WEBSOCKET | ERROR] %v", err)
					return
				}
				i, err := w.Write(*message)
				if err != nil {
					return
				}
				log.Printf("[WEBSOCKET | SENT | CLIENT ] %v", i)
				if err := w.Close(); err != nil {
					return
				}
			case <-ticker.C:
				log.Printf("[WEBSOCKET | TICK | CLIENT ]")
				_ = cp.Conn.SetWriteDeadline(time.Now().Add(client.WriteWait))
				if err := cp.Conn.WriteMessage(websocket.PingMessage, []byte("i")); err != nil {
					log.Printf("[WEBSOCKET | PING | ERROR ] %v", err)
					return
				}	
			}
		}	
	}
}



func (cp *ChargePoint) readerCsms() {
	_ = cp.Conn.SetReadDeadline(csms.getReadTimeout())
	cp.Conn.SetPingHandler(func(appData string) error {
		log.Printf("Ping received:  %v", appData)
		cp.PingIn <- []byte(appData)
		return cp.Conn.SetReadDeadline(csms.getReadTimeout())
	})	
	defer func() {
		cp.Conn.Close()
	}()
	for {
		_ = cp.Conn.SetReadDeadline(csms.getReadTimeout())
		log.Printf("Waiting for message[server side]")
		_, msg, err := cp.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WEBSOCKET | ERROR] %v", err)
				// delete charge point from ChargePoints map
				csms.ChargePoints.Delete(cp.Id)
			}
			break
		}
		log.Printf("[WEBSOCKET | RECEIVED | SERVER ] %v", msg)
		call, callResult, callError, err := unpack(&msg)
		if err != nil {
			cp.Out <- call.createCallError(err)
			log.Printf("[ERROR | MSG] %v", err)
		}
		if call != nil {
			handler, ok := csms.ActionHandlers[call.Action]
			if ok {
				responsePayload := handler(cp, call.Payload)
				log.Printf("[WEBSOCKET | PAYLOAD | SERVER ] %v", responsePayload)
				err = validate.Struct(responsePayload)
				if err != nil {
					// TODO simply log the error
					log.Printf("[ERROR | VALIDATION] %v", err)
				} else {
					// lock to ensure that only one message is sent at a time
					cp.Out <- call.createCallResult(responsePayload)
					// sleep for a bit to make sure the message is sent
					time.Sleep(time.Second)
					if afterHandler, ok := csms.AfterHandlers[call.Action]; ok {
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
				log.Printf("[ERROR | MSG] No handler for action %s", call.Action)
			}
		}
		if callResult != nil {
			log.Printf("[WEBSOCKET | CALL_RESULT | SERVER] %v", callResult)
			cp.Cr <- callResult
		}
		if callError != nil {
			log.Printf("[WEBSOCKET | CALL_ERROR | SERVER ] %v", callError)
			cp.Ce <- callError
		}
	}
}



// websocket writer to send messages
func (cp *ChargePoint) writerCsms() {
	for {
		select {
		case message, ok := <-cp.Out:
			_ = cp.Conn.SetWriteDeadline(time.Now().Add(csms.WriteWait))
			if !ok {
				cp.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := cp.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("[WEBSOCKET | ERROR] %v", err)
				return
			}
			i, err := w.Write(*message)
			if err != nil {
				return
			}
			log.Printf("[WEBSOCKET | SENT] %v", i)
			if err := w.Close(); err != nil {
				csms.ChargePoints.Delete(cp.Id)
				return
			}
		case ping := <-cp.PingIn:
			fmt.Println(ping)
			_ = cp.Conn.SetWriteDeadline(time.Now().Add(csms.WriteWait))
			err := cp.Conn.WriteMessage(websocket.PongMessage, []byte("o"))
			if err != nil {
				log.Printf("[WEBSOCKET | PING | ERROR ] %v", err)
				return
			}
			log.Printf("[WEBSOCKET | PONG_SEND | SERVER ]")		
		}

	}
}


// Call sends a message to other party (eg., ChargePoint, Central System)
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
		log.Println("Action name:", action)
		resPayload, err := unmarshalConf(action, callResult.Payload)
		if err != nil {
			log.Printf("[ERROR | MSG] %v", err)
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
	deadline := time.Now().Add(cp.ResTimeout)
	for {
		select {
		case r1 := <-cp.Cr:
			log.Printf("[WEBSOCKET | RECEIVED] %v", r1.UniqueId)
			if r1.UniqueId == uniqueId {
				return r1, nil, nil
			}
		case r2 := <-cp.Ce:
			log.Printf("[WEBSOCKET | RECEIVED] %v", r2.UniqueId)
			if r2.UniqueId == uniqueId {
				return nil, r2, nil
			}
		case <-time.After(time.Until(deadline)):
			return nil, nil, &TimeoutError{
				Message: fmt.Sprintf("timeout of %s sec for response to Call with id: %s passed", cp.ResTimeout, uniqueId),
			}
		}
	}
}

// NewChargepoint creates a new ChargePoint
func NewChargePoint(conn *websocket.Conn, id, proto string, isClient bool) *ChargePoint {
	cp := &ChargePoint{
		Proto:           proto,
		Conn:            conn,
		Id:              id,
		Out:             make(chan *[]byte),
		In:              make(chan *[]byte),
		Cr:              make(chan *CallResult),
		Ce:              make(chan *CallError),
		Extras: 		 make(map[string]interface{}),
	}

	if isClient {
		cp.ResTimeout = client.ResTimeout
		cp.WriteWait = client.WriteWait
		cp.PongWait = client.PongWait
		cp.PingPeriod = client.PingPeriod
		go cp.reader()
		go cp.writer()
	} else {
		cp.ResTimeout = csms.ResTimeout
		cp.WriteWait = csms.WriteWait
		cp.PingWait = csms.PingWait
		cp.PingIn   = make(chan []byte)
		go cp.readerCsms()
		go cp.writerCsms()
		csms.ChargePoints.Store(id, cp)
	} 
	return cp
}
