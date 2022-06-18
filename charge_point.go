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



// map to store websocket connections,
// keys are the chargepoint ids and values are the chargepoint structs which
// contain the websocket connection   
var ChargePoints = make(map[string]*ChargePoint)


/* 
Used as a container is for both Call and CallResult' Payload
*/
type Payload interface{}


 
// Represents a connected ChargePoint (also known as a Charging Station)
type ChargePoint struct {
	Conn 			*websocket.Conn   // the websocket connection
	Id 				string            // chargePointId 
	Out 			chan *[]byte      // channel to send messages to the ChargePoint 
	In 				chan *[]byte      // channel to receive messages from the ChargePoint 
	MessageHandlers map[string]func(Payload) Payload // map to store CP initiated actions
	AfterHandlers   map[string]func(Payload) 	// map to store functions to be called after a CP initiated action
	Mu 				sync.Mutex        // mutex ensuring that only one message is sent at a time
	Cr              chan *CallResult  
	Ce              chan *CallError
	Timeout 	    time.Duration     // timeout for waiting for a response
}



// Websocket reader to read messages from ChargePoint
func (cp *ChargePoint) Reader() {
	defer func() {
		cp.Conn.Close()
	}()
	for {
		_, msg, err := cp.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WEBSOCKET][ERROR] %v", err)
				delete(ChargePoints, cp.Id)
			}
			break
		}
		call, callResult, callError, err := unpack(&msg)
		if err != nil {
			cp.Out <- call.CreateCallError(err)
			log.Printf("[ERROR | MSG] %v", err)
		}
		if call != nil {
			if handler, ok := cp.MessageHandlers[call.Action]; ok {
				responsePayload := handler(call.Payload)
				// TODO check if validation works as expected / CP <-
				err = validate.Struct(responsePayload)
				if err != nil {
					// TODO simply log the error
					log.Printf("[ERROR | VALIDATION] %v", err)
				} else {
					cp.Out <- call.CreateCallResult(responsePayload)
					if afterHandler, ok := cp.AfterHandlers[call.Action]; ok {
						afterHandler(call.Payload)
					}
				}
			} else {
				var err error = &OCPPError{
					id:    call.UniqueId,
					code: "NotSupported",
					cause: fmt.Sprintf("Action %s is not supported", call.Action),
				}
				cp.Out <- call.CreateCallError(err)
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
func (cp *ChargePoint) Writer() {
	for {
		message, ok := <-cp.Out
		if !ok {
			cp.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		w, err := cp.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(*message)
		w.Close()
	}
}




// The function to be used by the implementers to register CP initiated actions
func (cp *ChargePoint) On(action string, f func(Payload) Payload) *ChargePoint {
	cp.MessageHandlers[action] = f
	return cp
}



// The function to be used by the implementers to register functions to be called after a CP initiated action
func (cp *ChargePoint) After(action string, f func(Payload)) *ChargePoint {
	cp.AfterHandlers[action] = f
	return cp
}



// function to be used by the implementers to execute a CSMS initiated action
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
	callResult, callError, err := cp.WaitForResponse(id)
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



// Waites for a response to a certain Call 
func (cp *ChargePoint) WaitForResponse(uniqueId string) (*CallResult, *CallError, error) {
	wait_until := time.Now().Add(cp.Timeout)
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
		case <-time.After(time.Until(wait_until)):
			return nil,nil, &TimeoutError{
				Message: fmt.Sprintf("timeout of %s sec for response to Call with id: %s passed", cp.Timeout, uniqueId),
			}
		}
	}
}




// Creates a new ChargePoint 
func NewChargePoint(conn *websocket.Conn, id string) *ChargePoint {
	cp := &ChargePoint{
		Conn:   			conn,
		Id:     			id,
		Out:    			make(chan *[]byte),
		In:     			make(chan *[]byte),
		MessageHandlers: 	make(map[string]func(Payload) Payload),
		AfterHandlers: 		make(map[string]func(Payload)),
		Cr: 				make(chan *CallResult, 1),
		Ce: 				make(chan *CallError, 1),
		Timeout: 			time.Second * 10,
	}
	go cp.Reader()
	go cp.Writer()
	// add the ChargePoint to the list of ChargePoints
	cp.Mu.Lock()
	defer cp.Mu.Unlock()
	// add if not already there
	if _, ok := ChargePoints[cp.Id]; !ok {
		ChargePoints[cp.Id] = cp
	}
	return cp
}

