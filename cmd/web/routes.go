package main

import (
  "net/http"

  "github.com/justinas/alice"
)



func (app *application) routes() http.Handler {
  mux := http.NewServeMux()
  
  // define the application routes
  fileserver := http.FileServer(http.Dir("./ui/"))
  mux.Handle("/", http.StripPrefix("", fileserver))

  mux.HandleFunc("/ws", app.handleWebSocket)
  
  std := alice.New(app.recoverPanic, app.logRequest)
  return std.Then(mux)
}
