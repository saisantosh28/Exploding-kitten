package services

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/saisantosh28/Exploding-kitten/server/models"
)

var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Update with your Redis connection details
		Password: "",               // Set password if required
		DB:       0,                // Use the desired Redis database
	})
}

func SetUser(user *models.User) error {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = redisClient.Set(fmt.Sprintf("user:%s", user.ID), userJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetUser(id string) (*models.User, error) {
	userJSON, err := redisClient.Get(fmt.Sprintf("user:%s", id)).Bytes()
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

func UpdateUserScore(id string, score int) error {
	user, err := GetUser(id)
	if err != nil {
		return err
	}

	user.Score = score
	return SetUser(user)
}

func GetLeaderboard() ([]models.User, error) {
	var users []models.User
	keys, err := redisClient.Keys("user:*").Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		userJSON, err := redisClient.Get(key).Bytes()
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

func SaveGame(game *models.Game) error {
	gameJSON, err := json.Marshal(game)
	if err != nil {
		return err
	}

	err = redisClient.Set(fmt.Sprintf("game:%s", game.ID), gameJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetGame(id string) (*models.Game, error) {
	gameJSON, err := redisClient.Get(fmt.Sprintf("game:%s", id)).Bytes()
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
