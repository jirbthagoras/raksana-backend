package models

type ResponseQuest struct {
	ID              int            `json:"id"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	PointGain       int            `json:"point_gain"`
	ContributedAt   string         `json:"contributed_at"`
	Location        string         `json:"location"`
	Latitude        float64        `json:"latitude"`
	Longitude       float64        `json:"Longitude"`
	MaxContributors int            `json:"max_contributors"`
	CreatedAt       string         `json:"created_at"`
	Contributors    []Contributors `json:"contributors"`
}

type Contributors struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
}
