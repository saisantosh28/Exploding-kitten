package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/saisantosh28/Exploding-kitten/server/models"
	"github.com/saisantosh28/Exploding-kitten/server/services"
)

// CreateUser creates a new user.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Creating user:", user.Username)

	err = services.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("User created successfully:", user.Username)

	w.WriteHeader(http.StatusCreated)
}

// GetLeaderboard retrieves the current leaderboard.
func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching leaderboard...")

	users, err := services.GetLeaderboard()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Leaderboard retrieved successfully.")

	json.NewEncoder(w).Encode(users)
}

// UpdateScore updates the score of a user.
func UpdateScore(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Updating score for user:", user.Username)

	err = services.UpdateUserScore(user.ID, user.Score)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Score updated successfully for user:", user.Username)

	w.WriteHeader(http.StatusOK)
}

// StartGame starts a new game session for a user.
func StartGame(w http.ResponseWriter, r *http.Request) {
	var game models.Game
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Starting game for user:", game.UserID)

	err = services.StartGame(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Game started successfully for user:", game.UserID)

	w.WriteHeader(http.StatusCreated)
}

// DrawCard draws a card from the game deck for the user.
func DrawCard(w http.ResponseWriter, r *http.Request) {
	var game models.Game
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Drawing card for user:", game.UserID)

	card, err := services.DrawCard(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Card drawn successfully for user:", game.UserID)

	json.NewEncoder(w).Encode(card)
}

// SaveGameState saves the current state of the game.
func SaveGameState(w http.ResponseWriter, r *http.Request) {
	var game models.Game
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Saving game state for user:", game.UserID)

	err = services.SaveGame(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Game state saved successfully for user:", game.UserID)

	w.WriteHeader(http.StatusOK)
}

// RestoreGameState restores the game state for a user.
func RestoreGameState(w http.ResponseWriter, r *http.Request) {
	var game models.Game
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Restoring game state for user:", game.UserID)

	savedGame, err := services.GetGame(game.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Game state restored successfully for user:", game.UserID)

	json.NewEncoder(w).Encode(savedGame)
}
