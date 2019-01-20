package main

import (
	//"bytes"
	//"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
	"network/common"
	"encoding/json"
	"sync/atomic"
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

var sConnetId uint32

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

func StartWebSocketSerer() {
	addr := "192.168.0.105:9998"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return

		}
		connid := atomic.AddUint32(&sConnetId, 1)
		webc := &WebSocketConn{
			Conn: c,
		}

		go handleWebsocketConnection(webc, connid)
	})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("websocket ListenAndServe: ", err)
	}
}

func handleWebsocketConnection(webc *WebSocketConn, connid uint32) {
	c := webc.Conn
	defer c.Close()

	var seq int

	for {
		c.SetReadDeadline(time.Now().Add(1 * time.Minute))
		mt, message, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("handleWebsocketConnection read message error %v", err)
			break
		}

		if mt != websocket.BinaryMessage {
			break
		}

		seq ++
		var heart common.HeartBeat
		//rd := bytes.NewReader(message[0:16])
		json.Unmarshal(message, &heart)
		//binary.Read(rd, binary.BigEndian, &heart)
		fmt.Printf("recieve seq %v client addr %v userid %v client heart beat\r\n", seq,c.RemoteAddr(),heart.UserID)
		if heart.UserID > 0 {
			time.Sleep(1 * time.Second)
			data := common.SendHeart(2)
			c.WriteMessage(websocket.BinaryMessage, data)
		}
	}
}
