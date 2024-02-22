package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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
