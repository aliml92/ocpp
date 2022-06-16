package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)


func main() {
	

	chargePointId := "12345"
	url := fmt.Sprintf("ws://localhost:8080/ocpp/v16/%s", chargePointId)
	
	fmt.Println(url)
	
	header := http.Header{}
	header.Add("Sec-WebSocket-Protocol", "ocpp1.5")
	
	fmt.Printf("connecting to %s", url)

	c, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		fmt.Printf("error dialing: %v", err)
		return
	}
	fmt.Printf("connected to %s", url)
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		fmt.Printf("recv: %s", message)
	}

}