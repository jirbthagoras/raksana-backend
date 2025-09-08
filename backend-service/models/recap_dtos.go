package models

import (
	"jirbthagoras/raksana-backend/repositories"
)

type InputRecap struct {
	Date               string      `json:"date"`
	Tasks              []InputTask `json:"tasks"`
	AssignedTask       int         `json:"assigned_task"`
	CompletedTask      int         `json:"completed_task"`
	TaskCompletionRate string      `json:"task_completion_rate"`
}

type RequestGetRecap struct {
	PreviousRecap repositories.Recap `json:"previous_recap"`
	InputRecap    InputRecap         `json:"current_recap"`
}

type AIResponseRecap struct {
	GrowthRating string `json:"growth_rating"`
	Summary      string `json:"summary"`
	Tips         string `json:"tips"`
}

type ResponseRecap struct {
	Summary            string `json:"summary"`
	Tips               string `json:"tips"`
	AssignedTask       int32  `json:"assigned_task"`
	CompletedTask      int32  `json:"completed_task"`
	TaskCompletionRate string `json:"completion_rate"`
	GrowthRating       string `json:"growth_rating"`
	CreatedAt          string `json:"created_at"`
}

type RequestGetMonthlyRecap struct {
	Statistics RequestStatistics
	Logs       []ResponseGetLogs
	Histories  []ResponseHistory
}

type ResponseMonthlyRecap struct {
	RecapID        int64  `json:"recap_id,omitempty"`
	Summary        string `json:"summary,omitempty"`
	Tips           string `json:"tips,omitempty"`
	AssignedTask   int32  `json:"assigned_task,omitempty"`
	CompletedTask  int32  `json:"completed_task,omitempty"`
	CompletionRate string `json:"completion_rate,omitempty"`
	GrowthRating   string `json:"growth_rating,omitempty"`
	Type           string `json:"type,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	Challenges     int    `json:"challenges"`
	Events         int    `json:"events"`
	Quests         int    `json:"quests"`
	Treasures      int    `json:"treasures"`
	LongestStreak  int    `json:"longest_streak"`
}
