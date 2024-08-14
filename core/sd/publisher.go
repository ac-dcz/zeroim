package sd

import "github.com/zeromicro/go-zero/core/discov"

type Publisher struct {
	*discov.Publisher
}

func NewPublisher(endpoints []string, Edge EdgeEndpoint, opts ...discov.PubOption) (*Publisher, error) {
	key := Edge.Key()
	value, err := Edge.Encode()
	if err != nil {
		return nil, err
	}
	pub := &Publisher{}
	pub.Publisher = discov.NewPublisher(endpoints, key, string(value), opts...)
	return pub, nil
}
