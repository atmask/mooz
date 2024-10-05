package main

import (
  "net/http"
  "encoding/json"
  "github.com/atmask/mooz/internal/models"
)


func (app *application) CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	roomID := app.rooms.CreateRoom()

	type resp struct {
		RoomID string `json:"room_id"`
	}

  app.logger.Info("Creating a new room with RoomID: ", "roomID", roomID)
  err := json.NewEncoder(w).Encode(resp{RoomID: roomID})
  if err != nil {
    app.logger.Error("Error encoding response", "error", err.Error())
    app.serverError(w, r, err)
    return
  }
}


func (app *application) JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomID"]
	if !ok {
		app.logger.Info("RoomID missing in URL Parameters")
		app.clientError(w, http.StatusBadRequest)
    return
	}

	ws, err := app.wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		app.logger.Error("Web Socket Upgrade Error", "error", err.Error())
    app.serverError(w, r, err)
    return
	}

  // Create a new participant and add them to the room
  participant := &models.Participant{Conn: ws, Host: false}

	app.rooms.InsertIntoRoom(roomID[0], participant)

	go app.startNewBroadcaster()

	for {
		var msg models.BroadcastMsg

		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			app.logger.Error("Read Error:", "error", err.Error())
      // app.WsError(ws, err)
      return
		}

		msg.Client = ws
		msg.RoomID = roomID[0]

		app.logger.Info("Received message", "message", msg.Message)

		app.broadcastChannel <- msg

	}

}
