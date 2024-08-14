package session

import (
	"fmt"
	"strconv"
	"strings"
)

const sessionFieldNums = 2

type SessionID struct {
	Uid    uint64 `zeroim:"uid"`
	EdgeId uint64 `zeroim:"edgeid"`
}

func (s SessionID) String() string {
	return fmt.Sprintf("%d-%d", s.Uid, s.EdgeId)
}

func NewSessionIDFromString(key string) (SessionID, error) {
	temp := strings.Split(key, "-")
	if len(temp) != sessionFieldNums {
		return SessionID{}, fmt.Errorf("sessionID fotmat error: %s", key)
	}

	u, err := strconv.ParseUint(temp[0], 10, 64)
	if err != nil {
		return SessionID{}, fmt.Errorf("parse Uid error: %v", err)
	}
	e, err := strconv.ParseUint(temp[1], 10, 64)
	if err != nil {
		return SessionID{}, fmt.Errorf("parse Edgeid error: %v", err)
	}
	return SessionID{Uid: u, EdgeId: e}, nil
}

func (s *SessionID) EtcdEdgeKey() string {
	return fmt.Sprintf("Edge#%d", s.EdgeId)
}
