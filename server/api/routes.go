package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/users", CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/api/leaderboard", GetLeaderboard).Methods(http.MethodGet)
	router.HandleFunc("/api/score", UpdateScore).Methods(http.MethodPut)
	router.HandleFunc("/api/game/start", StartGame).Methods(http.MethodPost)
	router.HandleFunc("/api/game/draw", DrawCard).Methods(http.MethodPost)
	router.HandleFunc("/api/game/save", SaveGameState).Methods(http.MethodPost)
	router.HandleFunc("/api/game/restore", RestoreGameState).Methods(http.MethodPost)

	return router
}
