package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
}
