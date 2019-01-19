package main

import (
	"encoding/json"
	"fmt"
	"net"
	"network/common"
	"time"
)

func main() {
	startServer()
}

func startServer() {
	addr := "192.168.1.116:9999"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Errorf("listen addr %v error %v\r\n", addr, err)
		return
	}
	defer listener.Close()
	fmt.Printf("listen addr %v ok !!!\r\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Errorf("listenet accept data error %v\r\n", err)
			return
		}

		conn.(*net.TCPConn).SetNoDelay(true)
		conn.(*net.TCPConn).SetKeepAlive(true)
		conn.(*net.TCPConn).SetLinger(8)

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	fmt.Printf(" *** handleConn addr %v\r\n", conn.RemoteAddr())

	// 设置1分钟超时
	if err := conn.SetReadDeadline(time.Now().Add(1 * time.Minute)); err != nil {

	}
	defer conn.Close()

	// 另一种方式
	//buffReader := bufio.NewReader(conn)
	var connid int64
	for {
		buffer := make([]byte, 1024)
		var bufferLen = 0

		n, err := conn.Read(buffer)

		//n, err := buffReader.Read(buffer)
		if n == 0 || err != nil {
			fmt.Errorf("conn read data len is 0 or read data error %v", err)
			break
		}
		bufferLen += n
		connid++
		go handleBuffer(buffer, bufferLen, conn, connid)
	}
}

func handleBuffer(buffer []byte, bufferLEen int, conn net.Conn, connid int64) {
	heart := new(common.HeartBeat)
	if err := json.Unmarshal(buffer[:bufferLEen], &heart); err != nil {
		fmt.Errorf("handleBuffer json unmarshal error %v", err)
		return
	}

	fmt.Printf("*** handleBuffer recievr from client connid %v heart %+v\r\n", connid, heart)
	if heart.UserID > 0 {
		time.Sleep(4 * time.Second)
		data := common.SendHeart(2)
		conn.Write(data)
	}
}
