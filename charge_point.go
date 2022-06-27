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


// Payload used as a container is for both Call and CallResult' Payload
type Payload interface{}

// ChargePoint Represents a connected ChargePoint (also known as a Charging Station)
type ChargePoint struct {
	Proto           string
	Conn            *websocket.Conn                          // the websocket connection
	Id              string                                   // chargePointId
	Out             chan *[]byte                             // channel to send messages to the ChargePoint
	In              chan *[]byte                             // channel to receive messages from the ChargePoint
	Mu              sync.Mutex                               // mutex ensuring that only one message is sent at a time
	Cr              chan *CallResult
	Ce              chan *CallError
	Extras          map[string]interface{}
	Timeout 	    time.Duration
}


// DRAFTING CSMS
type CSMS struct {
	ChargePoints sync.Map
	ActionHandlers map[string]func(*ChargePoint, Payload) Payload
	AfterHandlers map[string]func(*ChargePoint, Payload)
	Timeout time.Duration
}

func (csms *CSMS) On(action string, f func(*ChargePoint, Payload) Payload) *CSMS {
	csms.ActionHandlers[action] = f
	return csms
}

func (csms *CSMS) After(action string, f func(*ChargePoint,  Payload)) *CSMS {
	csms.AfterHandlers[action] = f
	return csms
}

func NewCSMS(timeout int) *CSMS {
	csms = &CSMS{
		ChargePoints: sync.Map{},
		ActionHandlers: make(map[string]func(*ChargePoint, Payload) Payload),
		AfterHandlers: make(map[string]func(*ChargePoint, Payload)),
		Timeout: time.Duration(timeout) * time.Second,
	}
	return csms
}


// Websocket reader to read messages from ChargePoint
func (cp *ChargePoint) reader() {
	defer func() {
		cp.Conn.Close()
	}()
	for {
		_, msg, err := cp.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WEBSOCKET][ERROR][READER] %v", err)
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
			// TODO: check if this is causing a deadlock
			handler, ok := csms.ActionHandlers[call.Action]
			if ok {
				responsePayload := handler(cp, call.Payload)
				// log response
				log.Printf("[WEBSOCKET][RESPONSE] %v", responsePayload)
				// TODO check if validation works as expected / CP <-
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
						afterHandler(cp, call.Payload)
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

// Websocket writer to send messages to the ChargePoint
func (cp *ChargePoint) writer() {
	for {
		message, ok := <-cp.Out
		if !ok {
			log.Printf("[WEBSOCKET][ERROR][WRITER1] Channel closed")
			cp.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		w, err := cp.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Printf("[WEBSOCKET][ERROR][WRITER2] %v", err)
			return
		}
		i, err := w.Write(*message)
		if err != nil {
			log.Printf("[WEBSOCKET][ERROR][WRITER3] %v", err)
			return
		}
		log.Printf("[WEBSOCKET][SENT] %v", i)
		if err := w.Close(); err != nil {
			csms.ChargePoints.Delete(cp.Id)
			log.Printf("[WEBSOCKET][ERROR][WRITER4] %v", err)
			return
		}
	}
}

// // On method used by the implementers to register action handlers
// func (cp *ChargePoint) On(action string, f func(string, Payload) Payload) *ChargePoint {
// 	cp.MessageHandlers[action] = f
// 	return cp
// }

// // After method used by the implementers to register functions to be called after a CP initiated action
// func (cp *ChargePoint) After(action string, f func(string, Payload)) *ChargePoint {
// 	cp.AfterHandlers[action] = f
// 	return cp
// }

// Call method   used by the implementers to execute a CSMS initiated action
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

// Waits for a response to a certain Call
func (cp *ChargePoint) waitForResponse(uniqueId string) (*CallResult, *CallError, error) {
	deadline := time.Now().Add(cp.Timeout)
	for {
		select {
		case r1 := <-cp.Cr:
			log.Printf("[WEBSOCKET][RECEIVED] %v", r1.UniqueId)
			if r1.UniqueId == uniqueId {
				return r1, nil, nil
			}
		case r2 := <-cp.Ce:
			log.Printf("[WEBSOCKET][RECEIVED] %v", r2.UniqueId)
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

// NewChargePoint creates a new ChargePoint
func NewChargePoint(conn *websocket.Conn, id string, proto string) *ChargePoint {
	cp := &ChargePoint{
		Proto:           proto,
		Conn:            conn,
		Id:              id,
		Out:             make(chan *[]byte),
		In:              make(chan *[]byte),
		Cr:              make(chan *CallResult,2),
		Ce:              make(chan *CallError,2),
		Extras: 		 make(map[string]interface{}),
		Timeout:         time.Second * 10,
	}
	go cp.reader()
	go cp.writer()
	csms.ChargePoints.Store(id, cp)
	return cp
}
