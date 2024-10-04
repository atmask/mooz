package main

import (
	"net/http"
	"runtime/debug"
  "github.com/pion/webrtc/v3"
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


// Init peer connections
func (app *application) createPeerConnection() (*webrtc.PeerConnection, error) {
  // define ICS servers for NAT traversal
  iceServers := []webrtc.ICEServer{
    {
      URLs: []string{"stun:stun.l.google.com:19302"},
    },
  }

  // Create an RTCPeerConnection
  config := webrtc.Configuration{
    ICEServers: iceServers,
  }
  peerConnection, err := webrtc.NewPeerConnection(config)
  if err != nil {
    return nil, err
  }

  // Handle ICE connection state changes
  peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
    app.logger.Info("ICE connection state has changed", "state", connectionState.String())
  })

  return peerConnection, nil

}
