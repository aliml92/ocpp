package ocpp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	v16 "github.com/aliml92/ocpp/v16"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	Subprotocols: []string{"ocpp1.6"},
}

var validate = v16.Validate
var ChargePoints = make(map[string]*ChargePoint)

type ReqPayload interface{}

type ResPayload interface{}



type ChargePoint struct {
	Conn 			*websocket.Conn
	Id 				string
	Out 			chan []byte
	In 				chan []byte
	MessageHandlers map[string]func(ReqPayload) ResPayload
	AfterHandlers   map[string]func(ReqPayload) 	
	Mu 				sync.Mutex
	Cr              chan *CallResult
	Ce              chan *CallError
	Timeout 	    time.Duration 
}



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
		p := &msg
		call, callResult, callError, err := unpack(p)
		if err != nil {
			cp.SendCallError(err)
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
					cp.Out <- *call.Marshal(&call.UniqueId, &responsePayload)
					if afterHandler, ok := cp.AfterHandlers[call.Action]; ok {
						afterHandler(call.Payload)
					}
				}
			} else {
				// TODO: send CallError with NotSupported error
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



// 6. ChargePoint.Writer
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
		w.Write(message)
		w.Close()
	}
}


func (cp *ChargePoint) SendCallError(err error) {
	var id string
	var code string
	var cause string
	var ocppErr OCPPError
	if errors.Is(err, &ocppErr) {
		id = ocppErr.id
		code = ocppErr.code
		cause = ocppErr.cause
	}
	if id == "" {
		id = uuid.New().String()
	}
	callError := &CallError{
			UniqueId: id,
			ErrorCode: code,
			ErrorDescription: "",
			ErrorDetails: cause,
	}
	switch code {
	case "ProtocolError":
		callError.ErrorDescription = "Payload for Action is incomplete"
	case "PropertyConstraintViolation":
		callError.ErrorDescription = "Payload is syntactically correct but at least one field contains an invalid value"
	case "NotImplemented":
		callError.ErrorDescription = "Requested Action is not known by receiver"
	case "TypeConstraintViolationError":
		callError.ErrorDescription = "Payload for Action is syntactically correct but at least one of the fields violates data type constraints (e.g. “somestring”: 12)"		
	case "PropertyConstraintViolationError":
		callError.ErrorDescription = "Payload is syntactically correct but at least one field contains an invalid value"	
	default:
		callError.ErrorDescription = "Unknown error"	
	}
	cp.Out <- *callError.Marshal()
}

func (cp *ChargePoint) On(action string, f func(ReqPayload) ResPayload) *ChargePoint {
	cp.MessageHandlers[action] = f
	return cp
}


func (cp *ChargePoint) After(action string, f func(ReqPayload)) *ChargePoint {
	cp.AfterHandlers[action] = f
	return cp
}


func (cp *ChargePoint) Call(action string, p ReqPayload) (ResPayload, error) {
	id := uuid.New().String()
	// TODO: check if validation works as expected / CS -> 
	call := [4]interface{}{
		2,
		id,
		action,
		p,
	}
	raw, _ := json.Marshal(call)
	// use a lock to make sure we don't send two messages at the same time
	cp.Mu.Lock()
	defer cp.Mu.Unlock()
	cp.Out <- raw
	callResult, callError, err := cp.WaitForResponse(&id)
	if callResult != nil {
		resPayload, err := unmarshall_call_result_payload_from_cp(&action, callResult.Payload)
		// TODO just return the error
		if err != nil {
			log.Printf("[ERROR | MSG] %v", err)
			return nil, err
		}
		return resPayload, nil
	}
	if callError != nil {
		return nil, errors.New("CallError")
	}
	// TODO return timeout error
	return nil, err
}

func (cp *ChargePoint) WaitForResponse(uniqueId *string) (*CallResult, *CallError, error) {
	wait_until := time.Now().Add(cp.Timeout)
	for {
		select {
		case r1 := <-cp.Cr:
			fmt.Println("Received CallResult: ", r1)
			if r1.UniqueId == *uniqueId {
				fmt.Println("CallResult matches UniqueId")
				return r1, nil, nil
			}
		case r2 := <-cp.Ce:
			fmt.Println("Received CallError: ", r2)
			if r2.UniqueId == *uniqueId {
				fmt.Println("CallError matches UniqueId")
				return nil, r2, nil
			}	
		case <-time.After(time.Until(wait_until)):
			fmt.Println("Timed out")
			return nil,nil, fmt.Errorf("timed out waiting for response")
		}
	}
}




func NewChargePoint(w http.ResponseWriter, r *http.Request) (*ChargePoint, error) {
	id := strings.Split(r.URL.Path, "/")[2]
	log.Printf("[INFO] New connection from %s", id)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR | SOCKET CONNECT] %v", err)
		return nil, err
	}
	cp := &ChargePoint{
		Conn:   			conn,
		Id:     			id,
		Out:    			make(chan []byte),
		In:     			make(chan []byte),
		MessageHandlers: 	make(map[string]func(ReqPayload) ResPayload),
		AfterHandlers: 		make(map[string]func(ReqPayload)),
		Cr: 				make(chan *CallResult, 1),
		Ce: 				make(chan *CallError, 1),
		Timeout: 			time.Second * 10,
	}
	go cp.Reader()
	go cp.Writer()
	// add chargepoint to list of chargepoints
	cp.Mu.Lock()
	defer cp.Mu.Unlock()
	ChargePoints[cp.Id] = cp
	return cp, nil
}

