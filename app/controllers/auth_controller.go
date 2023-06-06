package controllers

import (
	"ekira-backend/app/errs"
	"ekira-backend/app/models"
	"ekira-backend/pkg/utils"
	"ekira-backend/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

// Login method for user login.
// @Description Login the user account.
// @Summary login the user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body controllers.Login.Request true "login params"
// @Success 200 {object} models.ResponseOK{result=controllers.Login.Result{user=controllers.Login.ResultUser}}
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	// Create login body params struct
	type Request struct {
		Email    string `json:"email" example:"alperen@e-kira.tk" required:"true"`
		Password string `json:"password" example:"123" required:"true"`
	}
	params := &Request{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(params); err != nil {
		// Return status 400 and error message.
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection))
	}

	// Create a new validator for request params
	validate := utils.NewValidator()

	// Validate request params
	if err := validate.Struct(params); err != nil {
		// Return, if some fields are not valid.
		return c.Status(errs.ErrValidate.StatusCode).JSON(models.NewResponseError(errs.ErrValidate).SetHeader("validate", utils.ValidatorErrors(err)))
	}

	// Check user by email
	user, err := db.GetUserByEmail(params.Email)
	if err != nil {
		// Return status 500 and error message.
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}
	if user.Email == "" {
		// Return status 200 and error message.
		return c.Status(errs.ErrLogin.StatusCode).JSON(models.NewResponseError(errs.ErrLogin))
	}
	if user.Password != utils.HashPassword(user.ID, params.Password) {
		// Return status 200 and error message.
		return c.Status(errs.ErrLogin.StatusCode).JSON(models.NewResponseError(errs.ErrLogin))
	}

	// Create a session id
	sessionID := uuid.New()

	data := map[string]string{
		"Session": sessionID.String(),
		"UserId":  user.ID.String(),
		"Email":   user.Email,
	}
	token, claims, err := utils.GenerateNewAuthAccessToken(data)
	if err != nil {
		return c.Status(errs.ErrJwt.StatusCode).JSON(models.NewResponseError(errs.ErrJwt).SetHeader("jwt", err.Error()))
	}

	// Create a session
	session := models.Session{}
	session.ID = sessionID
	session.UserID = user.ID
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()
	session.ExpiresAt = claims["exp"].(int64)
	session.UserAgent = c.Get("User-Agent")
	session.IP = c.IP()

	// Save session to database
	if _, err := db.CreateSession(session); err != nil {
		// Return status 500 and error message.
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Result Structs
	type ResultUser struct {
		ID           uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
		Email        string    `json:"email" example:"john.doe@e-kira.tk"`
		FirstName    string    `json:"firstName" example:"John"`
		LastName     string    `json:"lastName" example:"Doe"`
		Validated    bool      `json:"validated" example:"false"`
		PhoneNumber  string    `json:"phoneNumber" example:"+905555555555"`
		ProfileImage *string   `json:"profileImage" example:"https://api.e-kira.tk/photos/profile/00000000-0000-0000-0000-000000000000.jpg"`
		RegisterDate string    `json:"registerDate" example:"2006-01-02T15:04:05Z07:00"`
	}

	type Result struct {
		User        interface{} `json:"user"`
		AccessToken string      `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"` // JWT access token
	}

	authUser := ResultUser{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Validated:    user.Validated,
		PhoneNumber:  user.PhoneNumber,
		RegisterDate: user.CreatedAt.Format(time.RFC3339),
	}
	if user.ProfileImage != nil {
		if len(user.ProfileImage.Images) > 0 {
			authUser.ProfileImage = &user.ProfileImage.Images[len(user.ProfileImage.Images)-1].URL
		}
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&Result{
		User:        authUser,
		AccessToken: token,
	}))
}

// Register method for user registration.
// @Description Create a user account.
// @Summary create a user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body controllers.Register.Request true "register params"
// @Success 200 {object} models.ResponseOK{result=controllers.Register.Result{user=controllers.Register.ResultUser}}
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Router /auth/register [post]
func Register(c *fiber.Ctx) error {
	// Create register body params struct
	type Request struct {
		Email    string `json:"email" example:"alperen@e-kira.tk" required:"true"`
		Password string `json:"password" example:"123" required:"true"`
	}
	params := &Request{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(params); err != nil {
		// Return status 400 and error message.
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection))
	}

	// Create a new validator for a User model.
	validate := utils.NewValidator()

	// Validate request params
	if err := validate.Struct(params); err != nil {
		// Return, if some fields are not valid.
		return c.Status(errs.ErrValidate.StatusCode).JSON(models.NewResponseError(errs.ErrValidate).SetHeader("validate", utils.ValidatorErrors(err)))
	}

	// Create new user
	user := models.User{}

	// Set initialized default data for user:
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.Validated = false
	user.Email = params.Email
	user.Password = utils.HashPassword(user.ID, params.Password)

	// Validate user fields.
	if err := validate.Struct(user); err != nil {
		// Return, if some fields are not valid.
		return c.Status(errs.ErrValidate.StatusCode).JSON(models.NewResponseError(errs.ErrValidate).SetHeader("validate", utils.ValidatorErrors(err)))
	}

	if user2, err := db.GetUserByEmail(user.Email); err == nil && user2.Email == user.Email {
		// Return status 500 and error message.
		return c.Status(errs.ErrRegister.StatusCode).JSON(models.NewResponseError(errs.ErrRegister))
	}

	// Create new user
	if err := db.NewUser(&user); err != nil {
		// Return status 500 and error message.
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Create a session id
	sessionID := uuid.New()

	data := map[string]string{
		"Session": sessionID.String(),
		"UserId":  user.ID.String(),
		"Email":   user.Email,
	}
	token, claims, err := utils.GenerateNewAuthAccessToken(data)
	if err != nil {
		return c.Status(errs.ErrJwt.StatusCode).JSON(models.NewResponseError(errs.ErrJwt).SetHeader("jwt", err.Error()))
	}

	// Create a session
	session := models.Session{}
	session.ID = sessionID
	session.UserID = user.ID
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()
	session.ExpiresAt = claims["exp"].(int64)
	session.UserAgent = c.Get("User-Agent")
	session.IP = c.IP()

	// Save session to database
	if _, err := db.CreateSession(session); err != nil {
		// Return status 500 and error message.
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Result Structs
	type ResultUser struct {
		ID           uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
		Email        string    `json:"email" example:"john.doe@e-kira.tk"`
		FirstName    string    `json:"firstName" example:"John"`
		LastName     string    `json:"lastName" example:"Doe"`
		Validated    bool      `json:"validated" example:"false"`
		PhoneNumber  string    `json:"phoneNumber" example:"+905555555555"`
		ProfileImage *string   `json:"profileImage" example:"https://api.e-kira.tk/photos/profile/00000000-0000-0000-0000-000000000000.jpg"`
		RegisterDate string    `json:"registerDate" example:"2006-01-02T15:04:05Z07:00"`
	}

	type Result struct {
		User        interface{} `json:"user"`
		AccessToken string      `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"` // JWT access token
	}

	authUser := ResultUser{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Validated:    user.Validated,
		PhoneNumber:  user.PhoneNumber,
		RegisterDate: user.CreatedAt.Format(time.RFC3339),
	}

	if user.ProfileImage != nil {
		if len(user.ProfileImage.Images) > 0 {
			authUser.ProfileImage = &user.ProfileImage.Images[len(user.ProfileImage.Images)-1].URL
		}
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&Result{
		User:        authUser,
		AccessToken: token,
	}))
}

// AuthCheck Authentication method for user.
// @Description Check logged-in
// @Summary Check user's logged-in
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseOK{result=controllers.AuthCheck.Result}
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /auth/check [get]
func AuthCheck(c *fiber.Ctx) error {
	// Get user from request context.
	user2 := c.Locals("user").(models.User)

	// Result Structs
	type Result struct {
		ID           uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
		Email        string    `json:"email" example:"john.doe@e-kira.tk"`
		FirstName    string    `json:"firstName" example:"John"`
		LastName     string    `json:"lastName" example:"Doe"`
		Validated    bool      `json:"validated" example:"false"`
		PhoneNumber  string    `json:"phoneNumber" example:"+905555555555"`
		ProfileImage *string   `json:"profileImage" example:"https://api.e-kira.tk/photos/profile/00000000-0000-0000-0000-000000000000.jpg"`
		RegisterDate string    `json:"registerDate" example:"2006-01-02T15:04:05Z07:00"`
	}
	userResult := Result{
		ID:           user2.ID,
		Email:        user2.Email,
		FirstName:    user2.FirstName,
		LastName:     user2.LastName,
		Validated:    user2.Validated,
		PhoneNumber:  user2.PhoneNumber,
		RegisterDate: user2.CreatedAt.Format(time.RFC3339),
	}

	if user2.ProfileImage != nil {
		profileImage := user2.ProfileImage
		images := profileImage.Images
		if len(images) > 0 {
			lastItem := images[len(images)-1]
			userResult.ProfileImage = &lastItem.URL
		}
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&userResult))
}

// Logout method for user.
// @Description Logout from a session, revoke authorization token.
// @Summary Logout from a session
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseOK{result=string}
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /auth/logout [get]
func Logout(c *fiber.Ctx) error {
	// Get user from request context.
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection))
	}
	sessionId := claims["Session"].(string)

	// Revoke session.
	if err := db.DeleteSessionBySessionID(sessionId); err != nil {
		// Return status 500 and error message.
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseString("OK"))
}
