package models

import (
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/repositories"
	"time"
)

type ResponseMemory struct {
	MemoryID        int64     `json:"memory_id"`
	FileURL         string    `json:"file_url"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UserID          int64     `json:"user_id"`
	UserName        string    `json:"user_name"`
	IsParticipation bool      `json:"is_participation"`
	ChallengeID     *int64    `json:"challenge_id,omitempty"`
	Day             *int32    `json:"day,omitempty"`
	Difficulty      *string   `json:"difficulty,omitempty"`
	ChallengeName   *string   `json:"challenge_name,omitempty"`
	PointGain       *int64    `json:"point_gain,omitempty"`
}

type PostMemoryCreate struct {
	FileKey     string `json:"file_key" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func ToResponseMemory(row repositories.GetMemoryWithParticipationRow) ResponseMemory {

	cnf := helpers.NewConfig()
	bucketUrl := cnf.GetString("AWS_URL")

	resp := ResponseMemory{
		MemoryID:        row.MemoryID,
		FileURL:         bucketUrl + row.FileKey,
		Description:     row.MemoryDescription,
		CreatedAt:       row.MemoryCreatedAt.Time,
		UserID:          row.UserID,
		UserName:        row.UserName,
		IsParticipation: row.IsParticipation,
	}

	if row.ChallengeID.Valid {
		resp.ChallengeID = &row.ChallengeID.Int64
	}
	if row.Day.Valid {
		resp.Day = &row.Day.Int32
	}
	if row.Difficulty.Valid {
		resp.Difficulty = &row.Difficulty.String
	}
	if row.ChallengeName.Valid {
		resp.ChallengeName = &row.ChallengeName.String
	}
	if row.PointGain.Valid {
		resp.PointGain = &row.PointGain.Int64
	}

	return resp
}
