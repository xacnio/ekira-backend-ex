package middleware

import (
	"ekira-backend/app/errs"
	"ekira-backend/app/models"
	"ekira-backend/platform/database"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	jwtWare "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

type Handler func(c *fiber.Ctx) error

// JWTProtected func for specify routes group with JWT authentication.
// See: https://github.com/gofiber/jwt
func JWTProtectedMaintenance(handle Handler) []func(c *fiber.Ctx) error {
	// Create config for JWT authentication middleware.
	jwtHandler := jwtWare.New(jwtWare.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET_KEY")),
		ContextKey: "jwt",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Next()
		},
	})

	return []func(c *fiber.Ctx) error{jwtHandler, func(c *fiber.Ctx) error {

		db, err := database.OpenDBConnection()
		if err != nil {
			// Return status 500 and database connection error.
			return c.Next()
		}

		// Get JWT claims from context.
		if _, ok := c.Locals("jwt").(*jwt.Token); !ok {
			// Return status 401 and failed authentication error.
			return c.Next()
		}
		user := c.Locals("jwt").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		sessionId := claims["Session"].(string)

		if len(sessionId) < 36 {
			// Return status 401 and failed authentication error.
			return c.Next()
		}

		session, err := db.GetSessionBySessionID(sessionId)
		if err != nil {
			// Return status 500 and database query error.
			return c.Next()
		}

		if session.ID == uuid.Nil {
			// Return status 401 and failed authentication error.
			return c.Next()
		}

		if session.ExpiresAt < time.Now().Unix() {
			// Return status 401 and failed authentication error.
			return c.Next()
		}

		session.UserAgent = c.Get("User-Agent")
		session.IP = c.IP()
		session.UpdatedAt = time.Now()
		err = db.UpdateSession(&session)
		if err != nil {
			fmt.Println(err)
		}

		c.Locals("user-id", session.UserID)
		c.Locals("session", session)

		return c.Next()
	}, handle}
}

// JWTProtected func for specify routes group with JWT authentication.
// See: https://github.com/gofiber/jwt
func JWTProtected(handlers ...fiber.Handler) []fiber.Handler {
	// Create config for JWT authentication middleware.
	jwtHandler := jwtWare.New(jwtWare.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
		ContextKey:   "jwt",
		ErrorHandler: jwtError,
	})

	sfMiddleware := func(c *fiber.Ctx) error {
		db, err := database.OpenDBConnection()
		if err != nil {
			// Return status 500 and database connection error.
			return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection))
		}

		// Get JWT claims from context.
		if _, ok := c.Locals("jwt").(*jwt.Token); !ok {
			// Return status 401 and failed authentication error.
			return c.Status(errs.ErrNoAuth.StatusCode).JSON(models.NewResponseError(errs.ErrNoAuth).SetHeader("jwt", "invalid ctx").SetHeader("token", c.Get("Authorization")))
		}
		user := c.Locals("jwt").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		sessionId := claims["Session"].(string)

		if len(sessionId) < 36 {
			// Return status 401 and failed authentication error.
			return c.Status(errs.ErrNoAuth.StatusCode).JSON(models.NewResponseError(errs.ErrNoAuth).SetHeader("jwt", "invalid session id").SetHeader("token", c.Get("Authorization")))
		}

		session, err := db.GetSessionBySessionID(sessionId)
		if err != nil {
			// Return status 500 and database query error.
			return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery))
		}

		if session.ID == uuid.Nil {
			// Return status 401 and failed authentication error.
			return c.Status(errs.ErrNoAuth.StatusCode).JSON(models.NewResponseError(errs.ErrNoAuth).SetHeader("jwt", "invalid session").SetHeader("token", c.Get("Authorization")))
		}

		if session.ExpiresAt < time.Now().Unix() {
			// Return status 401 and failed authentication error.
			return c.Status(errs.ErrNoAuth.StatusCode).JSON(models.NewResponseError(errs.ErrNoAuth).SetHeader("jwt", "expired").SetHeader("token", c.Get("Authorization")))
		}

		user1, err := db.GetUserById(session.UserID)
		if err != nil || user1.Email == "" {
			// Return status 500 and database query error.
			return c.Status(errs.ErrNoAuth.StatusCode).JSON(models.NewResponseError(errs.ErrNoAuth).SetHeader("jwt", "expired").SetHeader("token", c.Get("Authorization")))
		}

		db.Model(&models.User{}).Where("id = ?", user1.ID.String()).Update("last_access", time.Now())

		session.UserAgent = c.Get("User-Agent")
		session.IP = c.IP()
		session.UpdatedAt = time.Now()
		err = db.UpdateSession(&session)
		if err != nil {
			fmt.Println(err)
		}

		c.Locals("user", user1)
		c.Locals("user-id", session.UserID)
		c.Locals("session", session)

		return c.Next()
	}

	rHandlers := []fiber.Handler{jwtHandler, sfMiddleware}
	rHandlers = append(rHandlers, handlers...)

	return rHandlers
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(errs.ErrNoAuth.StatusCode).JSON(models.NewResponseError(errs.ErrNoAuth).SetHeader("jwt", err.Error()).SetHeader("token", c.Get("Authorization")))
	}

	// Return status 401 and failed authentication error.
	return c.Status(errs.ErrNoAuth.StatusCode).JSON(models.NewResponseError(errs.ErrNoAuth).SetHeader("jwt", err.Error()).SetHeader("token", c.Get("Authorization")))
}
