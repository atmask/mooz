package main

import (
	"net/http"
	"runtime/debug"

  "github.com/gorilla/websocket"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri = r.URL.RequestURI()
		trace = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)

    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
    http.Error(w, http.StatusText(status), status)
}

// Send a message to the user on the connection. If there's an error, log it and close the connection.
func (app *application) sendWsMessage(conn *websocket.Conn, message interface{}) {
  err := conn.WriteJSON(message)
  if err != nil {
    app.logger.Error("Error sending message to client. Closing their connection", "error", err.Error())
    conn.Close()
  }
}

func (app *application) WsError(conn *websocket.Conn, err error) {
  app.logger.Error("Sending ws error to user", "error", err.Error())
  app.sendWsMessage(conn, map[string]interface{}{ "error": err.Error() })
}


func (app *application) startNewBroadcaster() {
	for {
    // Listen on the channel for any messages and send to all clients
		msg := <- app.broadcastChannel
		for _, client := range app.rooms.Map[msg.RoomID] {
			if client.Conn != msg.Client {
        app.sendWsMessage(client.Conn, msg.Message)
			}
		}
	}
}


