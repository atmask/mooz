package models 

import (
  "github.com/gorilla/websocket"
  "sync"
)


type Participant struct {
  Host bool
  Mutex sync.RWMutex // Prevent concurrent writes to the same connection
  Conn *websocket.Conn
}

func (p *Participant) Close() {
  p.Mutex.Lock()
  defer p.Mutex.Unlock()
  p.Conn.Close()
}

func (p *Participant) SendJSON(v interface{}) error {
  p.Mutex.Lock()
  defer p.Mutex.Unlock()
  return p.Conn.WriteJSON(v)
}
