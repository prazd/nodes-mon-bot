package state

import (
	"sync"
)

type NodesState struct {
	sync.Mutex
	Result map[string]bool
}

func New() *NodesState {
	return &NodesState{
		Result: make(map[string]bool),
	}
}

func (ds *NodesState) set(key string, value bool) {
	ds.Result[key] = value
}

func (ds *NodesState) Set(key string, value bool) {
	ds.Lock()
	defer ds.Unlock()
	ds.set(key, value)
}
