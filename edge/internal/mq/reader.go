package mq

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/threading"
)

type Reader struct {
	reader  *kafka.Reader
	msgChan chan *kafka.Message
}

func NewReader(brokers []string, topic, groupID string) *Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		ErrorLogger: log.New(os.Stderr, "[Kafka] ", log.LstdFlags),
	})
	r := &Reader{
		reader:  reader,
		msgChan: make(chan *kafka.Message, 1000),
	}

	threading.GoSafe(r.run)

	return r
}

func (r *Reader) run() {
	defer r.Close()
	for {
		msg, err := r.reader.ReadMessage(context.Background())
		if err == context.DeadlineExceeded || err == kafka.LeaderNotAvailable {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		if err != nil {
			return
		}
		r.msgChan <- &msg
	}
}

func (r *Reader) MessageChannel() <-chan *kafka.Message {
	return r.msgChan
}

func (r *Reader) Close() error {
	close(r.msgChan)
	return r.reader.Close()
}
