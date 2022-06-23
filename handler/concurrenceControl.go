package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
)

type QueueMutex struct {
	mutexes sync.Map
}

func (m *QueueMutex) Lock(req *http.Request, bodyRequest []byte, query string) func() {
	header, _ := json.Marshal(req.Header)

	pathID := req.Method + "-" + req.URL.Path + query + "-" + string(header) + "-" + bytes.NewBuffer(bodyRequest).String()

	value, _ := m.mutexes.LoadOrStore(pathID, &sync.Mutex{})
	mtx := value.(*sync.Mutex)
	mtx.Lock()

	return func() { mtx.Unlock() }
}
