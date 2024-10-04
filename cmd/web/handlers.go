package main

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/websocket"
)

// func to return hello to the request
func (app *application) hello(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello, World!"))


}

func (app *application) handleWebSocket(w http.ResponseWriter, r *http.Request) {
  upgrader := websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
  }
  
  // Upgrade the HTTP connection to a WebSocket connection
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil { 
    app.serverError(w, r, err)
    return
  }
  defer conn.Close()


  for {
    _, message, err := conn.ReadMessage()
    if err != nil {
      app.serverError(w, r, err)
      break
    }

    var msg map[string]interface{}
    if err := json.Unmarshal(message, &msg); err != nil {
      app.serverError(w, r, err)
      continue
    }
  
    switch msg["type"]{
    case "join":
      app.logger.Info("A user joined the session", "name", msg["name"])
    case "offer":
      //handle offer
      app.logger.Info("A user sent an offer", "name", msg["name"])
    case "answer":
      //handle answer
      app.logger.Info("A user sent an answer", "name", msg["name"])
    case "candidate":
      //handle ICE candidate
      app.logger.Info("A user sent an ICE candidate", "name", msg["name"])

    }

  }
}

