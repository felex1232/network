package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

func main() {

}

func StartWebSocketClient() {
	addr := "192.168.1.116:9998"
	conn, _, err := websocket.Dialer{}.Dial(addr, nil)
	if err != nil {
		fmt.Errorf("websocket dial add error %v", err)
		return
	}

}
