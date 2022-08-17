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
	Timeout 	    time.Duration                            // timeout when waiting for a response from other party 
}


// CSMS acts as main handler for ChargePoints
type CSMS struct {
	ChargePoints sync.Map                               			// keeps track of all connected ChargePoints
	ActionHandlers map[string]func(*ChargePoint, Payload) Payload   // register implemented action handler functions
 	AfterHandlers map[string]func(*ChargePoint, Payload)            // register after-action habdler functions 
	Timeout time.Duration                                           // timeout when waiting for a response from other party
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
func NewCSMS(timeout int) *CSMS {
	csms = &CSMS{
		ChargePoints: sync.Map{},
		ActionHandlers: make(map[string]func(*ChargePoint, Payload) Payload),
		AfterHandlers: make(map[string]func(*ChargePoint, Payload)),
		Timeout: time.Duration(timeout) * time.Second,
	}
	return csms
}


// websocket reader to receive messages
func (cp *ChargePoint) reader() {
	defer func() {
		cp.Conn.Close()
	}()
	for {
		_, msg, err := cp.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WEBSOCKET | ERROR] %v", err)
				// delete charge point from ChargePoints map
				csms.ChargePoints.Delete(cp.Id)
			}
			break
		}
		call, callResult, callError, err := unpack(&msg)
		if err != nil {
			cp.Out <- call.createCallError(err)
			log.Printf("[ERROR | MSG] %v", err)
		}
		if call != nil {
			handler, ok := csms.ActionHandlers[call.Action]
			if ok {
				responsePayload := handler(cp, call.Payload)
				log.Printf("[WEBSOCKET | ERROR] %v", responsePayload)
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
			cp.Cr <- callResult
		}
		if callError != nil {
			cp.Ce <- callError
		}

	}
}

// websocket writer to send messages
func (cp *ChargePoint) writer() {
	for {
		message, ok := <-cp.Out
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
	deadline := time.Now().Add(cp.Timeout)
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
				Message: fmt.Sprintf("timeout of %s sec for response to Call with id: %s passed", cp.Timeout, uniqueId),
			}
		}
	}
}

// NewChargepoint creates a new ChargePoint
func NewChargePoint(conn *websocket.Conn, id string, proto string) *ChargePoint {
	cp := &ChargePoint{
		Proto:           proto,
		Conn:            conn,
		Id:              id,
		Out:             make(chan *[]byte),
		In:              make(chan *[]byte),
		Cr:              make(chan *CallResult),
		Ce:              make(chan *CallError),
		Extras: 		 make(map[string]interface{}),
		Timeout:         csms.Timeout,
	}
	go cp.reader()
	go cp.writer()
	csms.ChargePoints.Store(id, cp)
	return cp
}
