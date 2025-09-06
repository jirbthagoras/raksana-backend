package models

type ResponseTreasure struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	PointGain int    `json:"point_gain"`
	ClaimedAt string `json:"claimed_at"`
}
