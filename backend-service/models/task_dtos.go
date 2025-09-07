package models

type ResponseGetTask struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"created_at"`
	CompletedAt string `json:"completed_at,omitempty"`
}

type InputTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"created_at"`
	CompletedAt string `json:"completed_at"`
}
