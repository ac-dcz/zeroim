package types

import "time"

type ConnectMsg struct {
	AccessToken string `json:"access_token"`
}

type ChatMsg struct {
	Id          int64     `json:"msg_id"`
	From        int64     `json:"from"` //UserId/RoomId
	To          int64     `json:"to"`
	ChatType    uint8     `json:"chat_type"` //私聊or群聊
	ContentType uint32    `json:"content_type"`
	Data        []byte    `json:"data"`
	TimeAt      time.Time `json:"time_at"`
	Status      uint8     `json:"status"`
}
