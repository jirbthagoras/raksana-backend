package models

type ResponseEvent struct {
	ID           int64   `json:"id,omitempty"`
	Location     string  `json:"location,omitempty"`
	Latitude     float64 `json:"latitude,omitempty"`
	Longitude    float64 `json:"longitude,omitempty"`
	Contact      string  `json:"contact,omitempty"`
	StartsAt     string  `json:"starts_at,omitempty"`
	EndsAt       string  `json:"ends_at,omitempty"`
	CoverUrl     string  `json:"cover_url,omitempty"`
	Name         string  `json:"name,omitempty"`
	Description  string  `json:"description,omitempty"`
	PointGain    int64   `json:"point_gain,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
	IsEnded      bool    `json:"is_ended,omitempty"`
	Participated bool    `json:"participated"`
}

type RequestRegisterAttendance struct {
	ContactNumber string `json:"contact_number" validate:"required"`
}

type ResponseAttendance struct {
	ID                int64   `json:"id"`
	RegisteredAt      string  `json:"registered_at"`
	ContactNumber     string  `json:"contact_number"`
	Location          string  `json:"location"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Contact           string  `json:"contact"`
	StartsAt          string  `json:"starts_at"`
	EndsAt            string  `json:"ends_at"`
	CoverUrl          string  `json:"cover_url"`
	DetailName        string  `json:"name"`
	DetailDescription string  `json:"description"`
	PointGain         int64   `json:"point_gain"`
	AttendedAt        string  `json:"attended_at,omitempty"`
}
