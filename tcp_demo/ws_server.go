package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	StartWebSocketSerer()
}

type WebSocketConn struct {
	Conn   *websocket.Conn
	locker sync.Mutex
}

type MessageHeader struct {
	PackageLen uint16
	MsgID      uint16
	Seq        uint32
	PlayerID   uint32
	RoomID     uint32
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

func StartWebSocketSerer() {
	addr := "192.168.1.116:9998"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return

		}
		//connid := atomic.AddUint32(&sConnetId, 1)
		webc := &WebSocketConn{
			Conn: c,
		}

		go handleWebsocketConnection(webc)
	})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("websocket ListenAndServe: ", err)
	}
}

func handleWebsocketConnection(webc *WebSocketConn) {
	c := webc.Conn
	defer c.Close()

	for {
		c.SetReadDeadline(time.Now().Add(1 * time.Minute))
		mt, message, err := c.ReadMessage()
		if err != nil {
			fmt.Errorf("handleWebsocketConnection read message error %v", err)
			break
		}

		var header MessageHeader
		rd := bytes.NewReader(message[0:16])
		binary.Read(rd, binary.BigEndian, &header)
	}
}
