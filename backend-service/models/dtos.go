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
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Difficulty  string    `json:"difficulty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ResponseGetPacket struct {
	Id             int32            `json:"id"`
	Name           string           `json:"name"`
	Target         string           `json:"target"`
	Description    string           `json:"description"`
	CompletedTask  int32            `json:"completed_task"`
	ExpectedTask   int32            `json:"expected_task"`
	CompletionRate string           `json:"completion_rate"`
	TaskPerDay     int32            `json:"task_per_day"`
	Completed      bool             `json:"completed"`
	CreatedAt      pgtype.Timestamp `json:"created_at"`
}

type ResponseGetUserProfileStatistic struct {
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	CurrentExp         int64  `json:"current_exp"`
	ExpNeeded          int64  `json:"exp_needed"`
	Level              int32  `json:"level"`
	Points             int64  `json:"points"`
	ProfileUrl         string `json:"profile_url"`
	Challenges         int32  `json:"challenges"`
	Events             int32  `json:"events"`
	Quests             int32  `json:"quests"`
	Treasures          int32  `json:"treasures"`
	LongestStreak      int32  `json:"longest_streak"`
	TreeGrown          int32  `json:"tree_grown"`
	CompletedTask      int32  `json:"completed_task"`
	AssignedTask       int32  `json:"assigend_task"`
	TaskCompletionRate string `json:"task_completion_rate"`
}

type ResponsePacketDetail struct {
	PacketID           int64                       `json:"packet_id"`
	Username           string                      `json:"username"`
	PacketName         string                      `json:"packet_name"`
	Target             string                      `json:"target"`
	Description        string                      `json:"description"`
	CompletedTask      int32                       `json:"completed_task"`
	ExpectedTask       int32                       `json:"expected_task"`
	TaskCompletionRate string                      `json:"task_completion_rate"`
	TaskPerDay         int32                       `json:"task_per_day"`
	Completed          bool                        `json:"completed"`
	CreatedAt          time.Time                   `json:"created_at"`
	Habits             []ResponsePacketDetailHabit `json:"habits"`
}

type ResponsePacketDetailHabit struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	Locked      bool   `json:"locked"`
	ExpGain     int32  `json:"point_gain"`
}

type PostFilePresigned struct {
	Filename    string `json:"filename" validate:"required"`
	ContentType string `json:"content_type" validate:"required"`
}
