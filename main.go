package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// Define the WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
	openAIapiKey := r.URL.Query().Get("api_key")
	if openAIapiKey == "" {
		http.Error(w, "API key is required", http.StatusUnauthorized)
		return
	}

	err := validateAPIKey(openAIapiKey)
	if err != nil {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	// Upgrade HTTP request to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	fmt.Printf("A client has connected, starting game: %s\n", conn.RemoteAddr().String())

	defer conn.Close()

	// Start a game
	game := NewGame(conn, openAIapiKey)

	// Ping the client periodically to see if the connection is still alive
	go func() {
		for {
			fmt.Printf("Pinging client\n")
			err := conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				fmt.Println("Failed to ping client:", err)
				game.End()
				return
			}
			<-time.After(1 * time.Second)
		}
	}()

	// If the client disconnects, close the connection
	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("Client disconnected: %d %s\n", code, text)
		game.End()
		return nil
	})

	RunGame(game)
}

func main() {
	// WebSocket endpoint
	http.HandleFunc("/ws", wsHandler)

	port := os.Getenv("WEBSOCKET_PORT")

	fmt.Printf("Server running on port %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
