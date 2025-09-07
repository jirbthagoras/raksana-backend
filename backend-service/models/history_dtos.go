package models

type ResponseHistory struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Category  string `json:"category"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"created_at"`
}
