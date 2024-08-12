package sd

import (
	"encoding/json"
	"fmt"
)

type Service struct {
	Name     string
	NetWotk  string
	Addr     string
	MetaData map[string]any
}

func (s *Service) Key() string {
	return fmt.Sprintf("%s/%s@%s", s.Name, s.NetWotk, s.Addr)
}

func (s *Service) Encode() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Service) Decode(data []byte) error {
	return json.Unmarshal(data, s)
}

type Register interface {
	Registry(*Service) error
	DisRegistry(*Service) error
}
