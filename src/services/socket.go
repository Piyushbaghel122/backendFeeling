package services

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

// InitSocket initializes your Socket.io server and returns it so you can mount it in main.go
func InitSocket() *socketio.Server {
	server := socketio.NewServer(nil)

	// Handle standard connection event
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("Socket connected:", s.ID())
		return nil
	})

	// Handle a custom event (e.g., "chat message")
	server.OnEvent("/", "chat_message", func(s socketio.Conn, msg string) {
		log.Println("Received chat message:", msg)
		// Broadcast the message back to all connected clients
		server.BroadcastToRoom("/", "chat_room", "chat_message", msg)
	})

	// Handle standard error event
	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("Socket.io error:", e)
	})

	// Handle standard disconnect event
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("Socket.io closed:", reason)
	})

	go server.Serve()
	
	// IMPORTANT: You need to defer server.Close() in your main.go 
	// and mount it like this: http.Handle("/socket.io/", server)

	return server
}
