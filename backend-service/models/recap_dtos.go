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

type AIResponsWeeklyRecap struct {
	GrowthRating string `json:"growth_rating"`
	Summary      string `json:"summary"`
	Tips         string `json:"tips"`
}
