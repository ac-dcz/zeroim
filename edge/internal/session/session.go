package session

import (
	"sync"
	"zeroim/common/libnet"
)

type Session struct {
	EdgeID     string
	KafkaID    string
	UserID     int64
	codec      libnet.Codec
	sendChan   chan *libnet.Message
	closeFlag  int32
	closeChan  chan int
	closeMutex sync.Mutex
}
