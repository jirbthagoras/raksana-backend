package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type PostUserRegister struct {
	Name                 string `json:"name" validate:"required"`
	Username             string `json:"username" validate:"required,min=3,max=20"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type PostUserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type PostLogAppend struct {
	Text      string `json:"text" validate:"required"`
	IsSystem  bool   `json:"is_system" validate:"boolean"`
	IsPrivate bool   `json:"is_private" validate:"boolean"`
}

type PostPacketCreate struct {
	Target      string `json:"target" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type ResponseGetLogs struct {
	Text      string           `json:"text"`
	IsMarked  bool             `json:"is_marked"`
	IsSystem  bool             `json:"is_system"`
	IsPrivate bool             `json:"is_private"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type EcoachHabitResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
}

type EcoachCreatePacketResponse struct {
	Name         string `json:"name"`
	ExpectedTask int    `json:"expected_task"`
	TaskPerDay   int    `json:"task_per_day"`
	Habits       []EcoachHabitResponse
}

type ResponseGetTask struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Difficulty  string    `json:"difficulty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}
