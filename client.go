package ocpp

import "time"




type ClientTimeoutConfig struct {
	// ocpp response timeout in seconds
	OcppWait 	time.Duration 

	// time allowed to write a message to the peer
	WriteWait   time.Duration 

	// pong wait in seconds
	PongWait    time.Duration 

	// ping period in seconds
	PingPeriod  time.Duration 
}



type Client struct {
	// register implemented action handler functions
	actionHandlers map[string]func(*ChargePoint, Payload) Payload 
	// register after-action habdler functions
	afterHandlers  map[string]func(*ChargePoint, Payload)
	// timeout configuration
	ocppWait  time.Duration
	
	writeWait time.Duration

	pongWait     time.Duration

	pingPeriod  time.Duration
}


// create new Client instance
func NewClient() *Client {
	client = &Client{
		actionHandlers: make(map[string]func(*ChargePoint, Payload) Payload),
		afterHandlers:  make(map[string]func(*ChargePoint, Payload)),
		ocppWait: ocppWait,
		writeWait: writeWait,
		pongWait: pongWait,
		pingPeriod: pingPeriod,
	}
	return client
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