package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/saisantosh28/Exploding-kitten/server/models"
)

var redisClient *redis.Client

// InitRedis initializes the Redis connection.
func InitRedis() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis-11499.c264.ap-south-1-1.ec2.cloud.redislabs.com:11499",
		Password: "NZRLvX2RMLB0BEbqpVdPiZMvhDGGNhs1",
		Username: "default",
		DB:       0,
	})

	// Ping the Redis server to check the connection
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("error connecting to Redis: %v", err)
	}

	fmt.Println("Connected to Redis")
	return nil
}

// CreateUser creates a new user in the Redis database.
func CreateUser(user *models.User) error {
	if redisClient == nil {
		return errors.New("Redis client is not initialized")
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = redisClient.Set(redisClient.Context(), fmt.Sprintf("user:%s", user.ID), userJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// StartGame starts a new game session for a user.
func StartGame(game *models.Game) error {
	if redisClient == nil {
		return errors.New("Redis client is not initialized")
	}

	// Create a deck of cards
	deck := []string{"Cat", "Cat", "Cat", "Cat", "Exploding Kitten", "Defuse", "Shuffle"}

	// Shuffle the deck
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })

	// Set the deck and status in the game
	game.Deck = deck
	game.Status = "ongoing"

	// Save the game state
	err := SaveGame(game)
	if err != nil {
		return err
	}

	return nil
}

// DrawCard draws a card from the game deck for the user.
func DrawCard(game *models.Game) (string, error) {
	if redisClient == nil {
		return "", errors.New("Redis client is not initialized")
	}

	// Check if the game is ongoing
	if game.Status != "ongoing" {
		return "", errors.New("game is not ongoing")
	}

	// Check if the deck is empty
	if len(game.Deck) == 0 {
		return "", errors.New("deck is empty")
	}

	// Draw a card from the deck
	card := game.Deck[0]

	// Remove the drawn card from the deck
	game.Deck = game.Deck[1:]

	// Check the drawn card
	switch card {
	case "Cat":
		// Cat card: Do nothing, just continue the game
	case "Exploding Kitten":
		// Exploding Kitten: Player loses the game
		game.Status = "lost"
	case "Defuse":
		// Defuse card: Do nothing, can be used to defuse Exploding Kitten
	case "Shuffle":
		// Shuffle card: Shuffle the deck again
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(game.Deck), func(i, j int) { game.Deck[i], game.Deck[j] = game.Deck[j], game.Deck[i] })
	}

	// Save the updated game state
	err := SaveGame(game)
	if err != nil {
		return "", err
	}

	return card, nil
}

// SaveGame saves the game state to the Redis database.
func SaveGame(game *models.Game) error {
	if redisClient == nil {
		return errors.New("Redis client is not initialized")
	}

	gameJSON, err := json.Marshal(game)
	if err != nil {
		return err
	}

	err = redisClient.Set(redisClient.Context(), fmt.Sprintf("game:%s", game.ID), gameJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetGame retrieves a game state from the Redis database by ID.
func GetGame(id string) (*models.Game, error) {
	if redisClient == nil {
		return nil, errors.New("Redis client is not initialized")
	}

	gameJSON, err := redisClient.Get(redisClient.Context(), fmt.Sprintf("game:%s", id)).Bytes()
	if err != nil {
		return nil, err
	}

	var game models.Game
	err = json.Unmarshal(gameJSON, &game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

// GetLeaderboard retrieves the current leaderboard from the Redis database.
func GetLeaderboard() ([]models.User, error) {
	if redisClient == nil {
		return nil, errors.New("Redis client is not initialized")
	}

	var users []models.User
	keys, err := redisClient.Keys(redisClient.Context(), "user:*").Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		userJSON, err := redisClient.Get(redisClient.Context(), key).Bytes()
		if err != nil {
			return nil, err
		}

		var user models.User
		err = json.Unmarshal(userJSON, &user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// UpdateUserScore updates the score of a user in the Redis database.
func UpdateUserScore(id string, score int) error {
	if redisClient == nil {
		return errors.New("Redis client is not initialized")
	}

	user, err := GetUser(id)
	if err != nil {
		return err
	}

	user.Score = score

	// After updating the user, save it back to Redis
	return SetUser(user)
}

// SetUser sets a user in the Redis database.
func SetUser(user *models.User) error {
	if redisClient == nil {
		return errors.New("Redis client is not initialized")
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = redisClient.Set(redisClient.Context(), fmt.Sprintf("user:%s", user.ID), userJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUser gets a user from the Redis database by ID.
func GetUser(id string) (*models.User, error) {
	if redisClient == nil {
		return nil, errors.New("Redis client is not initialized")
	}

	userJSON, err := redisClient.Get(redisClient.Context(), fmt.Sprintf("user:%s", id)).Bytes()
	if err != nil {
		return nil, err
	}

	var user models.User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
