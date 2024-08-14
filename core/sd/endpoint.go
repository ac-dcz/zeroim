package sd

import (
	"encoding/json"
	"fmt"
)

type EdgeEndpoint struct {
	EdgeId  int64     `json:"edge_id"`
	Address string    `json:"address"`
	KqInfo  KafkaInfo `json:"kq_info"`
}

type KafkaInfo struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}

func (e *EdgeEndpoint) Key() string {
	return fmt.Sprintf("Edge#%d", e.EdgeId)
}

func (e *EdgeEndpoint) Encode() ([]byte, error) {
	return json.Marshal(e)
}

func (e *EdgeEndpoint) Decode(data []byte) error {
	return json.Unmarshal(data, e)
}
