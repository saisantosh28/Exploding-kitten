package main

import (
	"log"
	"net/http"

	"github.com/saisantosh28/exploding-kitten/api"
	"github.com/saisantosh28/exploding-kitten/services"
)

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
