package models

import "time"

type ResponseGetTask struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Difficulty  string    `json:"difficulty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type InputTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"created_at"`
	CompletedAt string `json:"completed_at"`
}
