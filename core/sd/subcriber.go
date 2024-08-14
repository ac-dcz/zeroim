package sd

import (
	"log"

	"github.com/zeromicro/go-zero/core/discov"
)

type Subcriber struct {
	*discov.Subscriber
}

func NewSubcriber(endpoints []string, key string, opts ...discov.SubOption) (*Subcriber, error) {
	sub, err := discov.NewSubscriber(endpoints, key, opts...)
	if err != nil {
		return nil, err
	}
	return &Subcriber{
		Subscriber: sub,
	}, nil
}

func (s *Subcriber) Values() []EdgeEndpoint {
	data := s.Subscriber.Values()
	var Edges []EdgeEndpoint
	for _, item := range data {
		e := &EdgeEndpoint{}
		if err := e.Decode([]byte(item)); err != nil {
			log.Printf("Endpoint Decode error: %v", err)
			return nil
		}
		Edges = append(Edges, *e)
	}
	return Edges
}
