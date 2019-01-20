package main

import (
	"encoding/json"
	"fmt"
	"net"
	"network/common"
	"sync"
	"time"
)

var ConnLocker sync.Mutex

func main() {
	StartClient()
}

func StartClient() {
	addr := "192.168.0.105:9999"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("dial server error %v", err)
		return
	}
	//defer conn.Close()

	// send heart
	time.Sleep(1 * time.Second)
	data := common.SendHeart(1)
	conn.Write(data)

	var connid int64
	handleClient(conn, connid)

	select {}
}

func handleClient(conn net.Conn, connid int64) {
	for {
		buffer := make([]byte, 1024)
		connid++

		n, err := conn.Read(buffer)
		if err != nil || n == 0 {
			fmt.Printf("handleClient error %v n %v", err, n)
			return
		}

		// todo 设置包头 包体 包长度（基础 所以不写 只传数据）
		heart := new(common.HeartBeat)
		if err := json.Unmarshal(buffer[:n], &heart); err != nil {
			fmt.Printf("handleBuffer json unmarshal error %v", err)
			return
		}

		fmt.Printf(" *** handleBuffer reciev from server connid %v heart %+v\r\n", connid, heart)
		if heart.UserID > 0 {
			time.Sleep(2 * time.Second)
			data := common.SendHeart(1)
			conn.Write(data)
		}
	}
}
