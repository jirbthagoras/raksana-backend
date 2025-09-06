package models

import (
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/repositories"
	"time"
)

type ResponseMemory struct {
	MemoryID        int64   `json:"memory_id,omitempty"`
	FileURL         string  `json:"file_url"`
	Description     string  `json:"description"`
	CreatedAt       string  `json:"created_at"`
	UserID          int64   `json:"user_id"`
	UserName        string  `json:"user_name"`
	IsParticipation bool    `json:"is_participation,omitempty"`
	ChallengeID     *int64  `json:"challenge_id,omitempty"`
	Day             *int32  `json:"day,omitempty"`
	Difficulty      *string `json:"difficulty,omitempty"`
	ChallengeName   *string `json:"challenge_name,omitempty"`
	PointGain       *int64  `json:"point_gain,omitempty"`
}

type PostMemoryCreate struct {
	ContentType string `json:"content_type" validate:"required"`
	FileName    string `json:"filename" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func ToResponseMemory(row repositories.GetMemoryWithParticipationRow) ResponseMemory {

	cnf := helpers.NewConfig()
	bucketUrl := cnf.GetString("AWS_URL")
	loc, _ := time.LoadLocation("Asia/Jakarta")
	resp := ResponseMemory{
		MemoryID:        row.MemoryID,
		FileURL:         bucketUrl + row.FileKey,
		Description:     row.MemoryDescription,
		CreatedAt:       row.MemoryCreatedAt.Time.In(loc).Format("2006-01-02 15:04"),
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

type PostCreateParticipation struct {
	Description string `json:"description" validate:"required"`
	FileName    string `json:"filename" validate:"required"`
	ContentType string `json:"content_type" validate:"required"`
}
