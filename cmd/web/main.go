package main

import (
  "flag"
  "log/slog"
  "net/http"
  "os"
  "github.com/gorilla/websocket"
  "github.com/atmask/mooz/internal/models"
)

type application struct {
  logger *slog.Logger
  rooms *models.RoomMap
  broadcastChannel chan models.BroadcastMsg
  wsUpgrader *websocket.Upgrader
}


func main() {
  addr := flag.String("addr", ":8080", "http service address")
  flag.Parse()

  logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
  }))
  
  // Create RoomMap and init
  rooms := &models.RoomMap{}
  rooms.Init()
  
  // Creat channel for broadcasting messages to all clients
  var broadcast = make(chan models.BroadcastMsg)

  var upgrader = &websocket.Upgrader{
  	CheckOrigin: func(r *http.Request) bool {
  		return true
  	},
  }

  app := &application{
    logger: logger,
    rooms: rooms,
    broadcastChannel: broadcast,
    wsUpgrader: upgrader,
  }

  logger.Info("Starting server on", "addr", *addr)
  err := http.ListenAndServe(*addr, app.routes())
  logger.Error(err.Error())
  os.Exit(1)
}
