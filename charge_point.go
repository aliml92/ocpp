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
	// crOrce carries CallResult or CallError      
	crOrce		  chan interface{}
	// Extras is for future use to carry data between different actions	
	Extras        map[string]interface{}

	// tc timeout config ensures a ChargePoint has its unique timeout configuration
	tc 	  	  	  TimeoutConfig

	// isServer defines if a ChargePoint at server or client side
	isServer      bool

	// TODO:
	unmarshalRes func(a string, r json.RawMessage) (Payload, error)
	
	// ping in channel           
	pingIn        chan []byte
	
	// closeC used to close the websocket connection by user
	closeC        chan websocket.CloseError
	// TODO
	forceWClose   chan error
	connected 	  bool
	// Followings used for sending ping messages
	ticker        *time.Ticker
	tickerC       <-chan time.Time 	
	
	// serverPing defines if ChargePoint is in server initiated ping mode
	serverPing	  bool

	stopC 		  chan struct{}
	dispatcherIn  chan *CallReq	 	
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



func (cp *ChargePoint) ResetPingPong(t int) (err error) {
	if t < 0 {
		err = errors.New("interval cannot be less than 0")
		return
	}
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if cp.isServer {
		log.Debug("ping/pong reconfigured")
		cp.tc.pingWait = time.Duration(t) * time.Second
		cp.conn.SetPingHandler(func(appData string) error {
			cp.pingIn <- []byte(appData)
			log.Debug("<- ping")
			return cp.conn.SetReadDeadline(cp.getReadTimeout())
		})
		return 
	}
	log.Debug("ping/pong reconfigured")
	cp.tc.pongWait = time.Duration(t) * time.Second
	cp.tc.pingPeriod = (cp.tc.pongWait * 9) / 10
	cp.conn.SetPongHandler(func(appData string) error {
		log.Debug("<- pong")
		return cp.conn.SetReadDeadline(cp.getReadTimeout())
	})
	if t == 0 {
		cp.ticker.Stop()
	} else {
		cp.ticker.Reset(cp.tc.pingPeriod)
	}
	return
}

func (cp *ChargePoint) EnableServerPing(t int) (err error){
	if t <= 0 {
		err = errors.New("interval must be greater than 0")
		return
	}
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.serverPing = true
	if cp.isServer {
		log.Debug("server ping enabled")
		cp.tc.pongWait = time.Duration(t) * time.Second
		cp.tc.pingPeriod = (cp.tc.pongWait * 9) / 10
		cp.conn.SetPingHandler(nil)
		cp.ticker = time.NewTicker(cp.tc.pingPeriod)
		cp.tickerC = cp.ticker.C
		cp.conn.SetPongHandler(func(appData string) error {
			log.Debug("<- pong")
			return cp.conn.SetReadDeadline(cp.getReadTimeout())
		})
		return 
	}
	log.Debug("server ping enabled")
	cp.ticker.Stop()
	cp.tickerC = nil
	cp.conn.SetPongHandler(nil)
	cp.tc.pingWait = time.Duration(t) * time.Second
	cp.pingIn = make(chan []byte)
	cp.conn.SetPingHandler(func(appData string) error {
		cp.pingIn <- []byte(appData)
		log.Debug("<- ping")
		return cp.conn.SetReadDeadline(cp.getReadTimeout())
	})
	return
}




func (cp *ChargePoint) getReadTimeout() time.Time {
	if cp.serverPing {
		if cp.isServer {
			if cp.tc.pongWait == 0 {
				return time.Time{}
			}
			return time.Now().Add(cp.tc.pongWait)
		}
		if cp.tc.pingWait == 0 {
			return time.Time{}
		}
		return time.Now().Add(cp.tc.pingWait)
	}
	if cp.isServer {
		if cp.tc.pingWait == 0 {
			return time.Time{}
		}
		return time.Now().Add(cp.tc.pingWait)
	}
	if cp.tc.pongWait == 0 {
		return time.Time{}
	}
	return time.Now().Add(cp.tc.pongWait)
	
}




func (cp *ChargePoint) processIncoming(peer Peer)  bool {
	messageType, msg, err := cp.conn.ReadMessage()
	log.Debugf("messageType: %d", messageType)
	if err != nil {
		log.Debug(err)
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
			log.Debug(err)
		}
		cp.forceWClose <- err
		cp.stopC <- struct{}{}
		return true
	}
	call, callResultOrcallError, err := unpack(msg, cp.proto)
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
				if afterHandler := peer.getAfterHandler(call.Action); afterHandler != nil {
					time.Sleep(time.Second)
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
	if callResultOrcallError != nil {
		select {
		case cp.crOrce <- callResultOrcallError:
		default:	
		}
	}
	return false
}




// websocket reader to receive messages
func (cp *ChargePoint) clientReader() {
	defer func ()  {
		cp.connected = false
	}()
	cp.conn.SetPongHandler(func(appData string) error {
		log.Debug("<- pong")
		return cp.conn.SetReadDeadline(cp.getReadTimeout())
	})
	for {
		if cp.processIncoming(client) { break }
	}
}

func (cp *ChargePoint) processOutgoing() (br bool) {
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
		return true
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
		return true	
	case <-cp.tickerC:
		_ = cp.conn.SetWriteDeadline(time.Now().Add(cp.tc.writeWait))
		if err := cp.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			log.Error(err)
			return
		}
		log.Debug("ping ->")
		return true
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

func (cp *ChargePoint) clientWriter() {
	defer func() {
		_ = cp.conn.Close()
	}()
	if cp.tc.pingPeriod != 0 {
		cp.ticker = time.NewTicker(cp.tc.pingPeriod)
		cp.tickerC = cp.ticker.C
		defer cp.ticker.Stop()
	}
	for {
		if !cp.processOutgoing() { break }
	}
}




func (cp *ChargePoint) serverReader() {
	log.Debug("pass once")
	cp.conn.SetPingHandler(func(appData string) error {
		cp.pingIn <- []byte(appData)
		log.Debug("<- ping")
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

func (cp *ChargePoint) serverWriter() {
	defer server.Delete(cp.Id)
	for {
		if !cp.processOutgoing() { break }
	}
}

type CallReq struct {
	id 		 	string
	data  	 	[]byte
	recvChan 	chan interface{}  
}


// Call sends a message to peer
func (cp *ChargePoint) Call(action string, p Payload) (Payload, error) {
	// check if cahrge point is connected
	if !cp.IsConnected(){
		return nil, errors.New("charge point not connected") 
	}
	// validate payload 
	var err error
	switch cp.proto{
	case ocppV16:
		err = validateV16.Struct(p)
	case ocppV201:
		err = validateV201.Struct(p)	
	}  
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()
	call := [4]interface{}{
		2,
		id,
		action,
		p,
	}
	raw, _ := json.Marshal(call)
	recvChan := make(chan interface{}, 1)
	cr := &CallReq{
		id: id, 
		data: raw,
		recvChan: recvChan,
	}
	select {
    case cp.dispatcherIn <- cr:
		log.Debug("call request added to dispatcher")
    default:
		return nil, errors.New("call queque is full")
    }
	r, ok := <-recvChan
	if !ok {
		return nil, errors.New("charge point disconnected")
	}
	log.Debugf("callback : %v", &r)
	if callResult, ok := r.(*CallResult); ok {
		resPayload, err := cp.unmarshalRes(action, callResult.Payload)
		if err != nil {
			return nil, err
		}
		return resPayload, nil
	}
	if callError, ok := r.(*CallError); ok {
		return nil, callError
	}
	return nil, r.(*TimeoutError)
}

func (cp *ChargePoint) callDispatcher(){
	cleanUp := make(chan struct{},1)
	for {
		select {
		case callReq := <-cp.dispatcherIn:
				log.Debug("dispatcher in")
				cp.out <- callReq.data
				deadline := time.Now().Add(cp.tc.ocppWait)
				in:
				for {
					select {
					case <- cleanUp:
						log.Debug("clean up")
						break in
					case <- cp.stopC:
						log.Debug("charge point is closed")
						for ch := range cp.dispatcherIn{
							log.Debug("cancel remaning call requests in queque")
							close(ch.recvChan)
						}
						break in		
					case r := <-cp.crOrce:
						log.Debug("call result in dispatcher")
						switch v := r.(type) {
						case *CallResult:
							if v.UniqueId == callReq.id {
								callReq.recvChan <- v
								break in
							}
						case *CallError:
							if v.UniqueId == callReq.id {
								callReq.recvChan <- v
								break in
							}	
						}
					case <-time.After(time.Until(deadline)):
						log.Debug("ocpp timeout occured")
						callReq.recvChan <-  &TimeoutError{
							Message: fmt.Sprintf("timeout of %s sec for response to Call with id: %s passed", cp.tc.ocppWait, callReq.id),
						}
						break in		
					} 					
				}
				log.Debug("broke from loop") 
		case <- cp.stopC:
			log.Debug("charge point is closed")
			select {
			case cleanUp <- struct{}{}:
			default:	
			}
			for ch := range cp.dispatcherIn{
				log.Debug("cancel remaning call requests in queque")
				close(ch.recvChan)
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
		crOrce:      make(chan interface{}),
		Extras:      make(map[string]interface{}),
		closeC:      make(chan websocket.CloseError, 1),
		forceWClose: make(chan error, 1),
		stopC: 		 make(chan struct{}),
		connected:   true,
	}
	go cp.callDispatcher()
	switch isServer {
	case true:
		cp.dispatcherIn = make(chan *CallReq, server.getCallQueueSize())
		cp.tc.ocppWait = server.ocppWait
		cp.tc.writeWait = server.writeWait
		cp.tc.pingWait = server.pingWait
		cp.pingIn = make(chan []byte)
		cp.isServer = true
		cp.tickerC = nil
		go cp.serverReader()
		go cp.serverWriter()
	case false:
		cp.dispatcherIn = make(chan *CallReq, client.callQuequeSize)
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
