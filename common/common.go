package common

import (
	"encoding/json"
	"fmt"
)

const (
	SERVERADDR = "192.168.1.116:20000"
	CLIENTADDR = "192.168.1.116:20000"
)

type HeartBeat struct {
	UserID uint32 `json:"userid"`
	Type   uint32 `json:"type"`
}

func SendHeart(gameType uint32) []byte {
	heart := new(HeartBeat)
	heart.UserID = 12665
	heart.Type = gameType
	data, err := json.Marshal(heart)
	if err != nil {
		fmt.Errorf("sendHeart json marhsal error %v", err)
		return []byte{}
	}
	return data
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
