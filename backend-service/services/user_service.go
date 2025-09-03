package services

import (
	"context"
	"errors"
	"fmt"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type UserService struct {
	Repository *repositories.Queries
}

func NewUserService(
	r *repositories.Queries,
) *UserService {
	return &UserService{
		Repository: r,
	}
}

func (s *UserService) GetUserDetail(id int) (models.ResponseGetUserProfileStatistic, error) {
	var profile models.ResponseGetUserProfileStatistic

	res, err := s.Repository.GetUserProfileStatistic(context.Background(), int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return profile, fiber.NewError(fiber.StatusBadRequest, "User not found")
		}
		slog.Error("Failed to get user profile and statistics")
		return profile, err
	}

	tasks, err := s.Repository.CountUserTask(context.Background(), int64(id))
	if err != nil {
		slog.Error("Failed to count", "err", err)
		return profile, err
	}
	var completionRate float64 = 0.0

	if tasks.AssignedTask != 0 || tasks.CompletedTask != 0 {
		completionRate = float64(tasks.CompletedTask) * 100.0 / float64(tasks.AssignedTask)
	}

	stringCompletionRate := fmt.Sprintf("%v", completionRate) + "%"

	cnf := helpers.NewConfig()
	bucketUrl := cnf.GetString("AWS_URL")

	profile = models.ResponseGetUserProfileStatistic{
		Id:                 int(res.UserID),
		Name:               res.Name,
		Username:           res.Username,
		Email:              res.Email,
		CurrentExp:         res.CurrentExp,
		ExpNeeded:          res.ExpNeeded,
		Level:              res.Level,
		Points:             res.Points,
		ProfileUrl:         bucketUrl + res.ProfileKey,
		Challenges:         res.Challenges,
		Events:             res.Events,
		Quests:             res.Quests,
		Treasures:          res.Treasures,
		TaskCompletionRate: stringCompletionRate,
		CompletedTask:      int32(tasks.CompletedTask),
		AssignedTask:       int32(tasks.AssignedTask),
		LongestStreak:      res.LongestStreak,
	}

	return profile, nil
}
