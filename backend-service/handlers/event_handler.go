package handlers

import (
	"context"
	"errors"
	"fmt"
	"jirbthagoras/raksana-backend/exceptions"
	"jirbthagoras/raksana-backend/helpers"
	"jirbthagoras/raksana-backend/models"
	"jirbthagoras/raksana-backend/repositories"
	"jirbthagoras/raksana-backend/services"
	"log/slog"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type EventHandler struct {
	Validator  *validator.Validate
	Repository *repositories.Queries
	*services.PointService
	*services.JournalService
	*services.StreakService
	Mu sync.Mutex
}

func NewEventHandler(
	v *validator.Validate,
	r *repositories.Queries,
	ps *services.PointService,
	js *services.JournalService,
	ss *services.StreakService,
) *EventHandler {
	return &EventHandler{
		Validator:      v,
		Repository:     r,
		PointService:   ps,
		JournalService: js,
		StreakService:  ss,
	}
}

func (h *EventHandler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/event")
	g.Use(helpers.TokenMiddleware)
	g.Post("/:id", h.handlerRegisterEvent)
	g.Get("/", h.handleGetEvents)
	g.Get("/pending", h.handleGetAllPendingAttendance)
	g.Get("/:id", h.handleGetAttendanceDetail)
}

func (h *EventHandler) handlerRegisterEvent(c *fiber.Ctx) error {
	req := &models.RequestRegisterAttendance{}

	err := c.BodyParser(req)
	if err != nil {
		slog.Error("Failed to parse payload", "err", err.Error())
		return err
	}

	err = h.Validator.Struct(req)
	if err != nil && errors.As(err, &validator.ValidationErrors{}) {
		return exceptions.NewFailedValidationError(*req, err.(validator.ValidationErrors))
	}

	eventId, err := c.ParamsInt("id")
	if err != nil {
		slog.Error("Failed to take event id", "err", err)
		return err
	}

	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	ctx := context.Background()

	event, err := h.Repository.GetEventById(ctx, int64(eventId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "Event does not exists")
		}
		slog.Error("Failed to get the event", "err", err)
		return err
	}

	pendingAttendance, err := h.Repository.GetUserAttendances(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get user pending attendance", "err", err)
		return err
	}

	for _, attendance := range pendingAttendance {
		if attendance.EventID == event.ID {
			return fiber.NewError(fiber.StatusBadRequest, "Kamu sudah berpartisipasi pada event ini!")
		}
	}

	_, err = h.Repository.CreateAttendance(ctx, repositories.CreateAttendanceParams{
		UserID:        int64(userId),
		EventID:       event.ID,
		ContactNumber: req.ContactNumber,
	})
	if err != nil {
		slog.Error("Failed to create attendance", "err", err)
		return err
	}

	logMsg := fmt.Sprintf("Baru saja mendaftar di event: %s! Tunggu kabar saya!", event.Name)
	err = h.JournalService.AppendLog(&models.PostLogAppend{
		Text:      logMsg,
		IsSystem:  true,
		IsPrivate: false,
	}, userId)
	if err != nil {
		return err
	}

	err = h.StreakService.UpdateStreak(ctx, int64(userId))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"message": "success",
		},
	})
}

func (h *EventHandler) handleGetEvents(c *fiber.Ctx) error {
	ctx := context.Background()
	res, err := h.Repository.GetAllEvents(ctx)
	if err != nil {
		slog.Error("Failed to get events", "err", err)
		return err
	}

	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	attendances, err := h.Repository.GetUserAttendances(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get all attendance", "err", err)
		return err
	}

	var events []models.ResponseEvent
	for _, event := range res {
		var participated bool = false
		if len(attendances) > 0 {
			for _, p := range attendances {
				if p.EventID == event.ID {
					participated = true
				}
			}
		}
		events = append(events, models.ResponseEvent{
			ID:           event.ID,
			Name:         event.DetailName,
			Description:  event.DetailDescription,
			Location:     event.Location,
			Longitude:    event.Longitude,
			Latitude:     event.Latitude,
			PointGain:    event.PointGain,
			CreatedAt:    event.DetailCreatedAt.Time.Format("2006-01-02 15:04"),
			StartsAt:     event.StartsAt.Time.Format("2006-01-02 15:04"),
			EndsAt:       event.EndsAt.Time.Format("2006-01-02 15:04"),
			Contact:      event.Contact,
			IsEnded:      event.Ended,
			Participated: participated,
			CoverUrl:     "https://raksana-admin.s3.ap-southeast-2.amazonaws.com/" + event.CoverKey.String,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"events": events,
		},
	})
}

func (h *EventHandler) handleGetAllPendingAttendance(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	res, err := h.Repository.GetUserPendingAttendances(context.Background(), int64(userId))
	if err != nil {
		slog.Error("Failed get attendance", "err", err)
		return err
	}

	var events []models.ResponseEvent
	for _, event := range res {
		events = append(events, models.ResponseEvent{
			ID:          event.AttendanceID,
			Name:        event.DetailName,
			Description: event.DetailDescription,
			Location:    event.Location,
			Longitude:   event.Longitude,
			Latitude:    event.Latitude,
			PointGain:   event.PointGain,
			CreatedAt:   event.RegisteredAt.Time.Format("2006-01-02 15:04"),
			StartsAt:    event.StartsAt.Time.Format("2006-01-02 15:04"),
			EndsAt:      event.EndsAt.Time.Format("2006-01-02 15:04"),
			Contact:     event.Contact,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"events": events,
		},
	})

}

func (h *EventHandler) handleAttend(c *fiber.Ctx) error {
	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}
	req := &models.ActivityRequest{}

	err = c.BodyParser(req)
	if err != nil {
		slog.Error("Failed to parse payload", "err", err)
		return err
	}

	err = h.Validator.Struct(req)
	if err != nil && errors.As(err, &validator.ValidationErrors{}) {
		return exceptions.NewFailedValidationError(*req, err.(validator.ValidationErrors))
	}

	_, payload, err := helpers.ValidateActivityToken(req.Token)
	if err != nil {
		slog.Error("Failed to validate token")
		return fiber.NewError(fiber.StatusBadRequest, "Token Invalid")
	}

	ctx := context.Background()

	event, err := h.Repository.GetEventByCodeId(ctx, payload.Subject)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "Event not found")
		}
		slog.Error("Failed to get event", "err", err)
		return err
	}

	attendance, err := h.Repository.GetUserAttendanceByUserId(ctx, repositories.GetUserAttendanceByUserIdParams{
		UserID: int64(userId),
		ID:     event.ID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "Sepertinya anda belum terdaftar dalam event")
		}
		slog.Error("Failed to get the attendance")
		return err
	}

	err = h.Repository.Attend(ctx, attendance.AttendanceID)
	if err != nil {
		slog.Error("Failed to finish the attend", "err", err)
		return err
	}

	logMsg := fmt.Sprintf("Baru saja menghadiri event: %s!", event.Name)
	err = h.JournalService.AppendLog(&models.PostLogAppend{
		Text:      logMsg,
		IsSystem:  true,
		IsPrivate: false,
	}, userId)
	if err != nil {
		return err
	}

	profile, err := h.Repository.GetUserProfile(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to get user profile", "err", err)
		return err
	}

	historyMsg := fmt.Sprintf("Mendapat poin event: %s", event.Name)
	_, err = h.PointService.UpdateUserPoint(int64(userId), event.PointGain, historyMsg, "event", int(profile.Level))
	if err != nil {
		return err
	}

	err = h.Repository.UpdaAttendedAt(ctx, attendance.AttendanceID)
	if err != nil {
		slog.Error("Failed to update attended_at", "err", err)
		return err
	}

	err = h.StreakService.UpdateStreak(ctx, int64(userId))
	if err != nil {
		return err
	}

	_, err = h.Repository.IncreaseEventsFieldByOne(ctx, int64(userId))
	if err != nil {
		slog.Error("Failed to increase", "er", err)
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"message": "success",
			"type":    "event",
			"event": fiber.Map{
				"name":        event.Name,
				"description": event.Description,
				"point_gain":  event.PointGain,
			},
		},
	})
}

func (h *EventHandler) handleGetAttendanceDetail(c *fiber.Ctx) error {
	attendanceId, err := c.ParamsInt("id")
	if err != nil {
		slog.Error("Failed to take attendance id", "err", err)
		return err
	}

	userId, err := helpers.GetSubjectFromToken(c)
	if err != nil {
		return err
	}

	attendance, err := h.Repository.GetUserAttendance(context.Background(), repositories.GetUserAttendanceParams{
		ID:     int64(attendanceId),
		UserID: int64(userId),
	})
	if err != nil {
		if errors.As(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusBadRequest, "attendance tidak ditemukan")
		}
		slog.Error("err", "failed to get attendance", err)
		return err
	}

	bucketUrl := "https://raksana-admin.s3.ap-southeast-2.amazonaws.com/"

	var resAttendance = models.ResponseAttendance{
		ID:                attendance.AttendanceID,
		RegisteredAt:      attendance.CreatedAt.Time.Format("2006-01-02 15:04"),
		AttendedAt:        attendance.AttendedAt.Time.Format("2006-01-02 15:04"),
		Location:          attendance.Location,
		Latitude:          attendance.Latitude,
		Longitude:         attendance.Longitude,
		StartsAt:          attendance.StartsAt.Time.Format("2006-01-02 15:04"),
		EndsAt:            attendance.EndsAt.Time.Format("2006-01-02 15:04"),
		CoverUrl:          bucketUrl + attendance.CoverKey.String,
		ContactNumber:     attendance.ContactNumber,
		DetailName:        attendance.DetailName,
		DetailDescription: attendance.DetailDescription,
		PointGain:         attendance.PointGain,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": resAttendance,
	})
}
