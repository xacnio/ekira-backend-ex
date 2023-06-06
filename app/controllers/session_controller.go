package controllers

import (
	"ekira-backend/app/errs"
	"ekira-backend/app/models"
	"ekira-backend/platform/database"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mileusna/useragent"
	"strconv"
)

// GetUserSessions method for user
// @Description user get sessions.
// @Summary user get sessions.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param page query int false "page number"
// @Param limit query int false "limit number"
// @Param itemPerPage query int false "item per page"
// @Security Authentication
// @Success 200 {object} models.ResponseOK{result=[]controllers.GetUserSessions.ResultSession}
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Router /auth/sessions [get]
func GetUserSessions(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	itemPerPage, _ := strconv.Atoi(c.Query("itemPerPage", "10"))

	// Database connection
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(errs.ErrDatabaseConnection)
	}

	// Get user id from context.
	userId := c.Locals("user-id").(uuid.UUID)

	// Get sessions from database.
	sessions, err := db.GetSessionsByUserID(userId, page, limit, itemPerPage)
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(errs.ErrDatabaseQuery)
	}

	type ResultSession struct {
		ID        uuid.UUID `json:"id"`
		IpAddress string    `json:"ipAddress"`
		Device    string    `json:"device"`
		Os        string    `json:"os"`
		Browser   string    `json:"browser"`
		ExpiresAt int64     `json:"expiresAt"`
		LastLogin int64     `json:"lastLogin"`
		DeletedAt int64     `json:"deletedAt"`
		Valid     bool      `json:"valid"`
		Current   bool      `json:"current"`
	}
	resultSessions := make([]ResultSession, 0)
	for session := range sessions {
		current := false
		ua := useragent.Parse(sessions[session].UserAgent)
		if sessions[session].ID == c.Locals("session").(models.Session).ID {
			current = true
		}

		resultSessions = append(resultSessions, ResultSession{
			ID:        sessions[session].ID,
			IpAddress: sessions[session].IP,
			Device:    ua.Device,
			Os:        ua.OS + " " + ua.OSVersion,
			Browser:   ua.Name,
			ExpiresAt: sessions[session].ExpiresAt,
			LastLogin: sessions[session].UpdatedAt.Unix(),
			DeletedAt: sessions[session].DeletedAt.Time.Unix(),
			Valid:     !sessions[session].DeletedAt.Valid,
			Current:   current,
		})
	}

	// Return status 200 and sessions.
	return c.Status(fiber.StatusOK).JSON(models.NewResponseOK(resultSessions))
}

// EndSession method for user
// @Description user end session.
// @Summary user end session.
// @Tags Authentication
// @Accept json
// @Produce json
// @Security Authentication
// @Param payload body controllers.EndSession.Request true "end session params"
// @Success 200 {string} string "OK"
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Router /auth/end-session [post]
func EndSession(c *fiber.Ctx) error {
	type Request struct {
		SessionID uuid.UUID `json:"sessionId" required:"true,uuid"`
	}

	req := Request{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest))
	}

	validator := validator.New()

	if err := validator.Struct(req); err != nil {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseErr(errors.New("invalid session id")).SetHeader("validator", err))
	}

	// Database connection
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(errs.ErrDatabaseConnection)
	}

	// Get user id from context.
	userId := c.Locals("user-id").(uuid.UUID)

	// Get sessions from database.
	affected, err := db.DeleteUserSessionBySessionID(userId, req.SessionID.String())
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(errs.ErrDatabaseQuery)
	}

	if affected == 0 {
		return c.Status(errs.ErrNotFound.StatusCode).JSON(errs.ErrNotFound)
	}

	// Return status 200 and sessions.
	return c.Status(fiber.StatusOK).JSON(models.NewResponseOK("OK"))
}
