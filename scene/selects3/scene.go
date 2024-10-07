package selects3

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"

	"github.com/coder/websocket"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
)

// embed all template files
//
//go:embed *.html
var Templates embed.FS

// TODO: Currently, Scene does not handle any dbs or opts.
type Scene struct {
	res  *scene.Resources
	opts *scene.Options
	dbs  *scene.Databases
	ws   *websocket.Conn // WebSocket connection

	msgChan  chan scene.PlayArgs // Channel to send parsed PlayArgs to the Update method
	quitChan chan struct{}       // Channel to signal the goroutine to stop
}

// Suppose the web socket has already been established.
// The scene will listen to the web socket.
// When the web socket receives a message, the scene will parse the json
// and return the PlayArgs to the game.
// NewScene initializes a new scene with the provided resources, options, handlers, and databases.
func NewScene(res *scene.Resources, opts *scene.Options, dbs *scene.Databases, ws *websocket.Conn) (*Scene, error) {
	s := &Scene{
		res:  res,
		opts: opts,
		dbs:  dbs,
		ws:   ws,

		msgChan:  make(chan scene.PlayArgs), // Create the message channel
		quitChan: make(chan struct{}),       // Create the quit channel
	}

	// Start the goroutine to listen for WebSocket messages
	go s.listen()

	return s, nil
}

// listen listens to the WebSocket and sends parsed PlayArgs to the message channel.
func (s *Scene) listen() {
	for {
		ctx := context.Background()

		// Read messages from the WebSocket using ReadMessage method
		_, msg, err := s.ws.Read(ctx)
		if err != nil {
			fmt.Println("Error reading message:", err)
			return // Exit the loop if an error occurs
		}

		// Parse the JSON message
		var playArgs scene.PlayArgs
		err = json.Unmarshal(msg, &playArgs)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			continue // Skip to the next message if parsing fails
		}

		// Send parsed PlayArgs to the channel
		select {
		case s.msgChan <- playArgs: // Send to the channel
		case <-s.quitChan: // Exit if quit signal received
			return
		}
	}
}

// Update processes the latest PlayArgs received from the WebSocket.
func (s *Scene) Update() any {
	select {
	case playArgs := <-s.msgChan: // Non-blocking receive of the latest PlayArgs
		// Process playArgs as needed
		// For example, update game state or transition scenes
		return playArgs // Return the processed PlayArgs to the game
	default:
		// No new messages received; proceed with the normal update
		return nil
	}
}

func (s Scene) Draw(dst draws.Image) {
	return
}

func (s Scene) DebugString() string {
	return "selects3"
}
