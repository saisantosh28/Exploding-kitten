package models

type Game struct {
	ID     string   `json:"id"`
	UserID string   `json:"userId"`
	Deck   []string `json:"deck"`
	Hand   []string `json:"hand"`
	Status string   `json:"status"` // "ongoing", "won", or "lost"
}
