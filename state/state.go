package state

import (
	"sync"
)

type SingleState struct {
	sync.Mutex
	Result map[string]bool
}

func NewSingleState() *SingleState {
	return &SingleState{
		Result: make(map[string]bool),
	}
}

func (ds *SingleState) Set(key string, value bool) {
	ds.Lock()
	defer ds.Unlock()
	ds.Result[key] = value
}
