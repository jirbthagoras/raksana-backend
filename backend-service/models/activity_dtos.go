package models

type ActivityRequest struct {
	Token string `json:"token" validate:"required"`
}
