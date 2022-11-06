package ocpp

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var client *Client

type ClientTimeoutConfig struct {
	// ocpp response timeout in seconds
	OcppWait time.Duration

	// time allowed to write a message to the peer
	WriteWait time.Duration

	// pong wait in seconds
	PongWait time.Duration

	// ping period in seconds
	PingPeriod time.Duration
}

type Client struct {
	Id string
	// register implemented action handler functions
	actionHandlers map[string]func(*ChargePoint, Payload) Payload
	// register after-action habdler functions
	afterHandlers map[string]func(*ChargePoint, Payload)
	// timeout configuration
	ocppWait time.Duration

	writeWait time.Duration

	pongWait time.Duration

	pingPeriod time.Duration

	header http.Header

	returnError func(error)

	callQuequeSize int
}

// create new Client instance
func NewClient() *Client {
	client = &Client{
		actionHandlers: make(map[string]func(*ChargePoint, Payload) Payload),
		afterHandlers:  make(map[string]func(*ChargePoint, Payload)),
		ocppWait:       ocppWait,
		writeWait:      writeWait,
		pongWait:       pongWait,
		pingPeriod:     pingPeriod,
		header:         http.Header{},
	}
	return client
}

func (c *Client) SetCallQueueSize(size int) {
	c.callQuequeSize = size
}

func (c *Client) SetTimeoutConfig(config ClientTimeoutConfig) {
	c.ocppWait = config.OcppWait
	c.writeWait = config.WriteWait
	c.pongWait = config.PongWait
	c.pingPeriod = config.PingPeriod
}

// register action handler function
func (c *Client) On(action string, f func(*ChargePoint, Payload) Payload) *Client {
	c.actionHandlers[action] = f
	return c
}

// register after-action handler function
func (c *Client) After(action string, f func(*ChargePoint, Payload)) *Client {
	c.afterHandlers[action] = f
	return c
}

func (c *Client) getHandler(action string) func(*ChargePoint, Payload) Payload {
	return c.actionHandlers[action]
}

func (c *Client) getAfterHandler(action string) func(*ChargePoint, Payload) {
	return c.afterHandlers[action]
}

func (c *Client) AddSubProtocol(protocol string) {
	c.header.Add("Sec-WebSocket-Protocol", protocol)
}

func (c *Client) SetBasicAuth(username string, password string) {
	auth := username + ":" + password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	c.header.Set("Authorization", "Basic "+enc)
}

func (c *Client) Start(addr string, path string) (cp *ChargePoint, err error) {
	urlStr, err := url.JoinPath(addr, path, c.Id)
	if err != nil {
		c.returnError(err)
		return
	}
	conn, _, err := websocket.DefaultDialer.Dial(urlStr, c.header)
	if err != nil {
		return
	}
	cp = NewChargePoint(conn, c.Id, conn.Subprotocol(), false)
	return
}

func (c *Client) SetID(id string) {
	c.Id = id
}
