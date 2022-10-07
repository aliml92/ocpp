package ocpp

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	log "github.com/aliml92/ocpp/logger"
	"github.com/aliml92/ocpp/v16"
	"github.com/aliml92/ocpp/v201"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)



const (

	ocppV16 = "ocpp1.6"
	ocppV201 = "ocpp2.0.1"

	// Time allowed to wait until corresponding ocpp call result received
	ocppWait = 30 * time.Second

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pigWait = 60 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

)


var validateV16 = v16.Validate
var validateV201 = v201.Validate
var server *Server
var client *Client

// add more

// Payload used as a container is for both Call and CallResult' Payload
type Payload interface{}

type TimeoutConfig struct {

	// ocpp response timeout in seconds
	ocppWait 	time.Duration
	
	// time allowed to write a message to the peer
	writeWait   time.Duration 

	// time allowed to read the next pong message from the peer
	pingWait    time.Duration 

	// pong wait in seconds
	pongWait    time.Duration 

	// ping period in seconds
	pingPeriod  time.Duration 
}

// ChargePoint Represents a connected ChargePoint (also known as a Charging Station)
type ChargePoint struct {
	// OCPP protocol version
	proto         string     
	
	// the websocket connection
	conn          *websocket.Conn 

	// chargePointId
	Id            string  
	
	// outgoing message channel
	out           chan *[]byte 
	
	// incoming message channel
	in            chan *[]byte 
	
	// mutex ensures that only one message is sent at a time
	mu            sync.Mutex      
	cr            chan *CallResult
	ce            chan *CallError
	Extras        map[string]interface{}
	tc 	  	  TimeoutConfig
	isServer      bool

	unmarshalRes func(a string, r *json.RawMessage) (Payload, error)
	
	// ping in channel           
	pingIn        chan []byte
	closeC        chan websocket.CloseError
	forceWClose   chan error
}

func (cp *ChargePoint) Shutdown() {
	cp.closeC <- websocket.CloseError{Code: websocket.CloseNormalClosure, Text: ""}
}

func (cp *ChargePoint) SetTimeoutConfig(config TimeoutConfig) {
	cp.tc = config
}





func (c *ChargePoint) getReadTimeout() time.Time {
	if c.tc.pingWait == 0 {
		return time.Time{}
	}
	return time.Now().Add(c.tc.pingWait)
}



// websocket reader to receive messages
func (cp *ChargePoint) clientReader() {
	cp.conn.SetPongHandler(func(appData string) error {
		log.L.Debug("pong <- ")
		return cp.conn.SetReadDeadline(cp.getReadTimeout())
	})
	for {
		messageType, msg, err := cp.conn.ReadMessage()
		log.L.Debugf("messageType: %d", messageType)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.L.Error(err)
			}
			cp.forceWClose <- err
			return
		}
		call, callResult, callError, err := unpack(&msg, cp.proto)
		if err != nil {
			cp.out <- call.createCallError(err)
			log.L.Error(err)
		}
		if call != nil {
			handler, ok := client.actionHandlers[call.Action]
			if ok {
				responsePayload := handler(cp, call.Payload)
				switch cp.proto {
				case ocppV16:
					err = validateV16.Struct(responsePayload)
				case ocppV201:
					err = validateV201.Struct(responsePayload)	
				}
				if err != nil {
					log.L.Error(err)
				} else {
					cp.out <- call.createCallResult(responsePayload)
					time.Sleep(time.Second)
					if afterHandler, ok := client.afterHandlers[call.Action]; ok {
						go afterHandler(cp, call.Payload)
					}
				}
			} else {
				var err error = &OCPPError{
					id:    call.UniqueId,
					code:  "NotSupported",
					cause: fmt.Sprintf("Action %s is not supported", call.Action),
				}
				cp.out <- call.createCallError(err)
				log.L.Errorf("No handler for action %s", call.Action)
			}
		}
		if callResult != nil {
			log.L.Debug(callResult)
			cp.cr <- callResult
		}
		if callError != nil {
			log.L.Debug(callError)
			cp.ce <- callError
		}

	}
}

// websocket writer to send messages
func (cp *ChargePoint) clientWriter() {
	defer func() {
		_ = cp.conn.Close()
	}()
	switch cp.tc.pingPeriod {
	case 0:
		for {
			select {
			case <-cp.forceWClose:
				return
			case closeErr := <-cp.closeC:
				b := websocket.FormatCloseMessage(closeErr.Code, closeErr.Text)
				err := cp.conn.WriteControl(websocket.CloseMessage, b, time.Now().Add(time.Second))
				if err != nil && err != websocket.ErrCloseSent {
					log.L.Error(err)
				}
				return
			case message, ok := <-cp.out:
				_ = cp.conn.SetWriteDeadline(time.Now().Add(cp.tc.writeWait))
				if !ok {
					cp.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				w, err := cp.conn.NextWriter(websocket.TextMessage)
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
		ticker := time.NewTicker(client.pingPeriod)
		defer ticker.Stop()
		for {
			select {
			case message, ok := <-cp.out:
				_ = cp.conn.SetWriteDeadline(time.Now().Add(cp.tc.writeWait))
				if !ok {
					cp.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				w, err := cp.conn.NextWriter(websocket.TextMessage)
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
				_ = cp.conn.SetWriteDeadline(time.Now().Add(cp.tc.writeWait))
				if err := cp.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					log.L.Error(err)
					return
				}
			case <-cp.forceWClose:
				return
			case closeErr := <-cp.closeC:
				b := websocket.FormatCloseMessage(closeErr.Code, closeErr.Text)
				err := cp.conn.WriteControl(websocket.CloseMessage, b, time.Now().Add(time.Second))
				if err != nil && err != websocket.ErrCloseSent {
					log.L.Error(err)
				}
				return
			}
		}
	}
}

func (cp *ChargePoint) serverReader() {
	// var count int
	cp.conn.SetPingHandler(func(appData string) error {
		cp.pingIn <- []byte(appData)
		log.L.Debug("ping <- ")
		i := cp.getReadTimeout()
		return cp.conn.SetReadDeadline(i)
	})
	defer func() {
		_ = cp.conn.Close()
	}()
	for {
		messageType, msg, err := cp.conn.ReadMessage()
		log.L.Debugf("messageType: %d ", messageType)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.L.Error(err)
			}
			cp.forceWClose <- err
			return
		}
		call, callResult, callError, err := unpack(&msg, cp.proto)
		if err != nil {
			cp.out <- call.createCallError(err)
			log.L.Error(err)
		}
		if call != nil {
			handler, ok := server.actionHandlers[call.Action]
			if ok {
				responsePayload := handler(cp, call.Payload)
				switch cp.proto {
				case ocppV16:
					err = validateV16.Struct(responsePayload)
				case ocppV201:
					err = validateV201.Struct(responsePayload)	
				}
				if err != nil {
					log.L.Error(err)
				} else {
					cp.out <- call.createCallResult(responsePayload)
					time.Sleep(time.Second)
					if afterHandler, ok := server.afterHandlers[call.Action]; ok {
						go afterHandler(cp, call.Payload)
					}
				}
			} else {
				var err error = &OCPPError{
					id:    call.UniqueId,
					code:  "NotSupported",
					cause: fmt.Sprintf("action %s is not supported", call.Action),
				}
				cp.out <- call.createCallError(err)
				log.L.Debugf("no handler for action %s", call.Action)
			}
		}
		if callResult != nil {
			log.L.Debugf("call result <-: %v", callResult)
			cp.cr <- callResult
		}
		if callError != nil {
			log.L.Debugf("call err <-: %v", callError)
			cp.ce <- callError
		}
	}
}

// websocket writer to send messages
func (cp *ChargePoint) serverWriter() {
	for {
		select {
		case message, ok := <-cp.out:
			err := cp.conn.SetWriteDeadline(time.Now().Add(cp.tc.writeWait))
			if err != nil {
				log.L.Error(err)
			}
			if !ok {
				err := cp.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.L.Error(err)
				}
				log.L.Debug("close msg ->")
				server.Chargepoints.Delete(cp.Id)
				return
			}
			w, err := cp.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.L.Error(err)
				server.Chargepoints.Delete(cp.Id)
				return
			}
			n, err := w.Write(*message)
			if err != nil {
				log.L.Error(err)
				server.Chargepoints.Delete(cp.Id)
				return
			}
			if err := w.Close(); err != nil {
				server.Chargepoints.Delete(cp.Id)
				log.L.Error(err)
				server.Chargepoints.Delete(cp.Id)
				return
			}
			log.L.Debugf("text msg -> %d", n)
		case <-cp.pingIn:
			err := cp.conn.SetWriteDeadline(time.Now().Add(cp.tc.writeWait))
			if err != nil {
				log.L.Error(err)
			}
			err = cp.conn.WriteMessage(websocket.PongMessage, []byte{})
			if err != nil {
				log.L.Error(err)
				server.Chargepoints.Delete(cp.Id)
				return
			}
			log.L.Debug("pong ->")
		case <-cp.forceWClose:
			server.Chargepoints.Delete(cp.Id)
			return
		case closeErr := <-cp.closeC:
			server.Chargepoints.Delete(cp.Id)
			b := websocket.FormatCloseMessage(closeErr.Code, closeErr.Text)
			err := cp.conn.WriteControl(websocket.CloseMessage, b, time.Now().Add(time.Second))
			if err != nil && err != websocket.ErrCloseSent {
				log.L.Error(err)
			}
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
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.out <- &raw
	callResult, callError, err := cp.waitForResponse(id)
	if callResult != nil {
		resPayload, err := cp.unmarshalRes(action, callResult.Payload)
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
	deadline := time.Now().Add(cp.tc.ocppWait)
	for {
		select {
		case r1 := <-cp.cr:
			if r1.UniqueId == uniqueId {
				return r1, nil, nil
			}
		case r2 := <-cp.ce:
			if r2.UniqueId == uniqueId {
				return nil, r2, nil
			}
		case <-time.After(time.Until(deadline)):
			return nil, nil, &TimeoutError{
				Message: fmt.Sprintf("timeout of %s sec for response to Call with id: %s passed", cp.tc.ocppWait, uniqueId),
			}
		}
	}
}

// NewChargepoint creates a new ChargePoint
func NewChargePoint(conn *websocket.Conn, id, proto string, isServer bool) *ChargePoint {
	cp := &ChargePoint{
		proto:       proto,
		conn:        conn,
		Id:          id,
		out:         make(chan *[]byte),
		in:          make(chan *[]byte),
		cr:          make(chan *CallResult),
		ce:          make(chan *CallError),
		Extras:      make(map[string]interface{}),
		closeC:      make(chan websocket.CloseError, 1),
		forceWClose: make(chan error, 1),
	}

	switch isServer {
	case true:
		cp.tc.ocppWait = server.ocppWait
		cp.tc.writeWait = server.writeWait
		cp.tc.pingWait = server.pingWait
		cp.pingIn = make(chan []byte)
		cp.isServer = true
		server.Chargepoints.Store(id, cp)
		go cp.serverReader()
		go cp.serverWriter()
	case false:
		cp.tc.ocppWait = client.ocppWait
		cp.tc.writeWait = client.writeWait
		cp.tc.pongWait = client.pongWait
		cp.tc.pingPeriod = client.pingPeriod
		go cp.clientReader()
		go cp.clientWriter()					
	}

	switch proto {
	case ocppV16:
		cp.unmarshalRes = unmarshalResV16
	case ocppV201:
		cp.unmarshalRes = unmarshalResV201	
	}
	return cp
}
