package models

type ActivityRequest struct {
	Token string `json:"token" validate:"required"`
}

type ResponseContributions struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	PointGain   float64 `json:"point_gain"`
}

type ResponseAttendances struct {
	Id          int64   `json:"id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	PointGain   int64   `json:"point_gain"`
}
