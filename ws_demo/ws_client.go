package ws_demo

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"network/common"
	"time"
)

func main() {
	StartWebSocketClient()
}

var (
	dialer websocket.Dialer
)

func NewWebDialer() {
	dialer = websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}
}

func init() {
	NewWebDialer()
}

func StartWebSocketClient() {
	addr := "ws://192.168.0.105:9998"
	conn, _, err := dialer.Dial(addr, nil)
	if err != nil {
		fmt.Printf("websocket dial add error %v", err)
		return
	}
	defer conn.Close()
	data := common.SendHeart(1)
	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		fmt.Printf("conn writeMessage error %v", err)
	}

	var seq int

	for {
		conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("handleWebsocketConnection read message error %v", err)
			break
		}

		if mt != websocket.BinaryMessage {
			break
		}

		seq++
		var heart common.HeartBeat
		//rd := bytes.NewReader(message[0:16])
		json.Unmarshal(message, &heart)
		//binary.Read(rd, binary.BigEndian, &heart)
		fmt.Printf("recieve server seq %v addr %v userid %v client heart beat\r\n", seq, conn.RemoteAddr(), heart.UserID)
		if heart.UserID > 0 {
			time.Sleep(1 * time.Second)
			data := common.SendHeart(1)
			conn.WriteMessage(websocket.BinaryMessage, data)
		}
	}

}
