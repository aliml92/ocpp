package ocpp

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/aliml92/ocpp/logger"
	"github.com/aliml92/ocpp/v16"
	"github.com/aliml92/ocpp/v201"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func init(){
	log = &logger.EmptyLogger{}
}

func SetLogger(logger logger.Logger) {
	if logger == nil {
		panic("logger cannot be nil")
	}
	log = logger
}


type Peer interface {
	getHandler(string) func(*ChargePoint, Payload) Payload
	getAfterHandler(string) func(*ChargePoint, Payload) 
}

const (

	ocppV16 = "ocpp1.6"
	ocppV201 = "ocpp2.0.1"

	// Time allowed to wait until corresponding ocpp call result received
	ocppWait = 20 * time.Second

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pingWait = 30 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

)


var validateV16 = v16.Validate
var validateV201 = v201.Validate
var log logger.Logger




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
	out           chan []byte 
	
	// incoming message channel
	in            chan []byte 
	
	// mutex ensures that only one message is sent at a time
	mu            sync.Mutex      
	cr            chan *CallResult
	ce            chan *CallError
	Extras        map[string]interface{}
	tc 	  	  	  TimeoutConfig
	isServer      bool

	unmarshalRes func(a string, r json.RawMessage) (Payload, error)
	
	// ping in channel           
	pingIn        chan []byte
	closeC        chan websocket.CloseError
	forceWClose   chan error
	connected 	  bool
	ticker        *time.Ticker	
}

func (cp *ChargePoint) IsConnected() bool {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	return cp.connected
}


func (cp *ChargePoint) Shutdown() {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.closeC <- websocket.CloseError{Code: websocket.CloseNormalClosure, Text: ""}
}

func (cp *ChargePoint) SetTimeoutConfig(config TimeoutConfig) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.tc = config
}


// func (cp *ChargePoint) DisablePingPong() {
// 	cp.mu.Lock()
// 	defer cp.mu.Unlock()
// 	if cp.isServer {
// 		log.Debug("server ping handler is disabled")
// 		cp.tc.pingWait = 0
// 		cp.conn.SetPingHandler(func(appData string) error {
// 			cp.pingIn <- []byte(appData)
// 			log.Debug("ping <- ")
// 			return cp.conn.SetReadDeadline(cp.getReadTimeout())
// 		})
// 		return 
// 	}
// 	log.Debug("client pong handler is disabled")
// 	cp.ticker.Stop()
// 	cp.tc.pongWait = 0
// 	cp.conn.SetPongHandler(func(appData string) error {
// 		log.Debug("pong <- ")
// 		return cp.conn.SetReadDeadline(cp.getReadTimeout())
// 	})
// }


func (cp *ChargePoint) ResetPingPong(t int) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if cp.isServer {
		log.Debug("server ping handler is reconfigured")
		cp.tc.pingWait = time.Duration(t) * time.Second
		cp.conn.SetPingHandler(func(appData string) error {
			cp.pingIn <- []byte(appData)
			log.Debug("ping <- ")
			return cp.conn.SetReadDeadline(cp.getReadTimeout())
		})
		return 
	}
	log.Debug("client pong handler is reconfigured")
	cp.tc.pongWait = time.Duration(t) * time.Second
	cp.tc.pingPeriod = (cp.tc.pongWait * 9) / 10
	cp.conn.SetPongHandler(func(appData string) error {
		log.Debug("pong <- ")
		return cp.conn.SetReadDeadline(cp.getReadTimeout())
	})
	if t == 0 {
		cp.ticker.Stop()
	} else {
		cp.ticker.Reset(cp.tc.pingPeriod)
	}
}

func (c *ChargePoint) getReadTimeout() time.Time {
	if c.isServer {
		if c.tc.pingWait == 0 {
			return time.Time{}
		}
		return time.Now().Add(c.tc.pingWait)
	}
	if c.tc.pongWait == 0 {
		return time.Time{}
	}
	return time.Now().Add(c.tc.pongWait)
	
}




func (cp *ChargePoint) processIncoming(peer Peer) bool {
	messageType, msg, err := cp.conn.ReadMessage()
	log.Debugf("messageType: %d", messageType)
	if err != nil {
		log.Debug(err)
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
			log.Debug(err)
		}
		cp.forceWClose <- err
		return true
	}
	call, callResult, callError, err := unpack(msg, cp.proto)
	if err != nil {
		cp.out <- call.createCallError(err)
		log.Error(err)
	}
	if call != nil {
		handler := peer.getHandler(call.Action)
		if handler != nil {
			responsePayload := handler(cp, call.Payload)
			switch cp.proto {
			case ocppV16:
				err = validateV16.Struct(responsePayload)
			case ocppV201:
				err = validateV201.Struct(responsePayload)	
			}
			if err != nil {
				log.Error(err)
			} else {
				cp.out <- call.createCallResult(responsePayload)
				time.Sleep(time.Second)
				if afterHandler := peer.getAfterHandler(call.Action); afterHandler != nil {
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
			log.Errorf("No handler for action %s", call.Action)
		}
	}
	if callResult != nil {
		log.Debug(callResult)
		cp.cr <- callResult
	}
	if callError != nil {
		log.Debug(callError)
		cp.ce <- callError
	}
	return false
}




// websocket reader to receive messages
func (cp *ChargePoint) clientReader() {
	log.Debug("pass once")
	defer func ()  {
		cp.connected = false
	}()
	cp.conn.SetPongHandler(func(appData string) error {
		log.Debug("pong <- ")
		return cp.conn.SetReadDeadline(cp.getReadTimeout())
	})
	for {
		if cp.processIncoming(client) { break }
	}
}

// websocket writer to send messages
func (cp *ChargePoint) clientWriter() {
	defer func() {
		_ = cp.conn.Close()
	}()
	if cp.tc.pingPeriod != 0 {
		cp.ticker = time.NewTicker(cp.tc.pingPeriod)
		defer cp.ticker.Stop()
	}
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
				log.Error(err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Error(err)
				return
			}
			if err := w.Close(); err != nil {
				log.Error(err)
				return
			}
		case <-cp.ticker.C:
			_ = cp.conn.SetWriteDeadline(time.Now().Add(cp.tc.writeWait))
			if err := cp.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Error(err)
				return
			}
			log.Debug("ping ->")
		case <-cp.forceWClose:
			return
		case closeErr := <-cp.closeC:
			b := websocket.FormatCloseMessage(closeErr.Code, closeErr.Text)
			err := cp.conn.WriteControl(websocket.CloseMessage, b, time.Now().Add(time.Second))
			if err != nil && err != websocket.ErrCloseSent {
				log.Error(err)
			}
			return
		}
	}	
}


func (cp *ChargePoint) serverReader() {
	log.Debug("pass once")
	cp.conn.SetPingHandler(func(appData string) error {
		cp.pingIn <- []byte(appData)
		log.Debug("ping <- ")
		i := cp.getReadTimeout()
		return cp.conn.SetReadDeadline(i)
	})
	defer func() {
		_ = cp.conn.Close()
		server.Delete(cp.Id)
	}()
	for {
		if cp.processIncoming(server) { break }
	}
}



// websocket writer to send messages
func (cp *ChargePoint) serverWriter() {
	defer server.Delete(cp.Id)
	for {
		select {
		case message, ok := <-cp.out:
			err := cp.conn.SetWriteDeadline(time.Now().Add(cp.tc.writeWait))
			if err != nil {
				log.Error(err)
			}
			if !ok {
				err := cp.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Error(err)
				}
				log.Debug("close msg ->")
				return
			}
			w, err := cp.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Debug(err)
				return
			}
			n, err := w.Write(message)
			if err != nil {
				log.Error(err)
				return
			}
			if err := w.Close(); err != nil {
				log.Error(err)
				return
			}
			log.Debugf("text msg -> %d", n)
		case <-cp.pingIn:
			err := cp.conn.SetWriteDeadline(time.Now().Add(cp.tc.writeWait))
			if err != nil {
				log.Error(err)
			}
			err = cp.conn.WriteMessage(websocket.PongMessage, []byte{})
			if err != nil {
				log.Error(err)
				return
			}
			log.Debug("pong ->")
		case <-cp.forceWClose:
			return
		case closeErr := <-cp.closeC:
			b := websocket.FormatCloseMessage(closeErr.Code, closeErr.Text)
			err := cp.conn.WriteControl(websocket.CloseMessage, b, time.Now().Add(time.Second))
			if err != nil && err != websocket.ErrCloseSent {
				log.Error(err)
			}
			return
		}
	}
}



// Call sends a message to peer
func (cp *ChargePoint) Call(action string, p Payload) (Payload, error) {
	if !cp.connected {
		return nil, errors.New("charge point not connected") 
	} 
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
	cp.out <- raw
	callResult, callError, err := cp.waitForResponse(id)
	if callResult != nil {
		resPayload, err := cp.unmarshalRes(action, callResult.Payload)
		if err != nil {
			log.Error(err)
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
		out:         make(chan []byte),
		in:          make(chan []byte),
		cr:          make(chan *CallResult),
		ce:          make(chan *CallError),
		Extras:      make(map[string]interface{}),
		closeC:      make(chan websocket.CloseError, 1),
		forceWClose: make(chan error, 1),
		connected:   true ,
	}

	switch isServer {
	case true:
		cp.tc.ocppWait = server.ocppWait
		cp.tc.writeWait = server.writeWait
		cp.tc.pingWait = server.pingWait
		cp.pingIn = make(chan []byte)
		cp.isServer = true
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
