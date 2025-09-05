package models

type ResponseChallenge struct {
	ID          int    `json:"id,omitempty"`
	Day         int    `json:"day"`
	Difficulty  string `json:"difficulty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PointGain   int    `json:"point_gain,omitempty"`
}
