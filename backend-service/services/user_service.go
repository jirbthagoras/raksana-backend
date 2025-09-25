package services

import (
	"context"
	"errors"
	"fmt"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type UserService struct {
	Repository *repositories.Queries
	*StreakService
	*LeaderboardService
}

func NewUserService(
	r *repositories.Queries,
	ss *StreakService,
	ls *LeaderboardService,
) *UserService {
	return &UserService{
		Repository:         r,
		StreakService:      ss,
		LeaderboardService: ls,
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

	streak, err := s.StreakService.GetCurrentStreak(context.Background(), res.UserID)
	if err != nil {
		return profile, err
	}

	rank, err := s.LeaderboardService.GetUserRank(strconv.Itoa(int(res.UserID)))
	if err != nil {
		return profile, err
	}

	levelBefore := res.Level - 2
	neededExpBefore := helpers.CalculateExpNeeded(int(levelBefore))
	if res.Level == 1 {
		neededExpBefore = 0
	}

	profile = models.ResponseGetUserProfileStatistic{
		Id:                     int(res.UserID),
		Name:                   res.Name,
		Username:               res.Username,
		Email:                  res.Email,
		Rank:                   int(rank),
		CurrentExp:             res.CurrentExp,
		ExpNeeded:              res.ExpNeeded,
		Level:                  res.Level,
		Points:                 res.Points,
		ProfileUrl:             bucketUrl + res.ProfileKey,
		Challenges:             res.Challenges,
		Events:                 res.Events,
		Quests:                 res.Quests,
		Treasures:              res.Treasures,
		TaskCompletionRate:     stringCompletionRate,
		NeededExpPreviousLevel: neededExpBefore,
		CompletedTask:          int32(tasks.CompletedTask),
		AssignedTask:           int32(tasks.AssignedTask),
		LongestStreak:          res.LongestStreak,
		CurrentStreak:          streak,
		Badges:                 s.CheckBadges(res),
	}

	return profile, nil
}

func (s *UserService) CheckBadges(profile repositories.GetUserProfileStatisticRow) []models.Badge {

	var badges []models.Badge

	if profile.Challenges == 0 &&
		profile.Treasures == 0 &&
		profile.Events == 0 &&
		profile.Quests == 0 {
		return []models.Badge{
			{
				Category:  "nuetral",
				Name:      "Peasant",
				Frequency: 0,
			},
		}
	}

	GetBadge("challenge", int(profile.Challenges), "Challenger", &badges)
	GetBadge("quest", int(profile.Quests), "Adventurer", &badges)
	GetBadge("event", int(profile.Events), "Scholar", &badges)
	GetBadge("treasure", int(profile.Treasures), "Hunter", &badges)

	return badges
}

func GetBadge(category string, frequency int, role string, badges *[]models.Badge) {
	switch {
	case frequency == 0:
		return
	case frequency >= 1 && frequency <= 5:
		*badges = append(*badges, models.Badge{
			Category:  category,
			Name:      "Beginner " + role,
			Frequency: frequency,
		})
	case frequency >= 6 && frequency <= 15:
		*badges = append(*badges, models.Badge{
			Category:  category,
			Name:      "Novice " + role,
			Frequency: frequency,
		})
	case frequency >= 16:
		*badges = append(*badges, models.Badge{
			Category:  category,
			Name:      "Expert " + role,
			Frequency: frequency,
		})
	}

	return
}
