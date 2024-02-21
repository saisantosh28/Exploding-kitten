package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/saisantosh28/Exploding-kitten/server/api"
	"github.com/saisantosh28/Exploding-kitten/server/services"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Your WebSocket logic goes here
}

func main() {
	// Initialize Redis
	services.InitRedis()

	// Set up API routes
	router := api.SetupRoutes()

	// Set up WebSocket server (for Bonus Feature)
	http.HandleFunc("/ws", handleWebSocket)

	// Start the server
	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
