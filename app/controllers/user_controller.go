package controllers

import (
	"ekira-backend/app/errs"
	"ekira-backend/app/models"
	"ekira-backend/pkg/utils"
	"ekira-backend/platform/database"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"os"
	"regexp"
	"strings"
	"time"
)

// SetProfile set user's profile
// @Description Set user's profile
// @Summary Set user's profile
// @Tags User
// @Accept json
// @Produce json
// @Param payload body controllers.SetProfile.Request true "set profile params"
// @Success 200 {object} models.ResponseOK{result=controllers.SetProfile.Response}
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /user/set-profile [post]
func SetProfile(c *fiber.Ctx) error {
	// Get user from request context.
	user := c.Locals("user").(models.User)

	type Request struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
	}

	var params Request
	if err := c.BodyParser(&params); err != nil {
		// Return status 400 and error message.
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest))
	}

	// Create a new validator for request params
	validate := utils.NewValidator()

	// Validate request params
	if err := validate.Struct(params); err != nil {
		// Return, if some fields are not valid.
		return c.Status(errs.ErrValidate.StatusCode).JSON(models.NewResponseError(errs.ErrValidate).SetHeader("validate", utils.ValidatorErrors(err)))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	type Response struct {
		FirstName string `json:"first_name" validate:"required,lte=50"`
		LastName  string `json:"last_name" validate:"required,lte=50"`
	}

	nameRegex := regexp.MustCompile(`^[a-zA-ZÜĞİŞÇÖüğışöç]+$`)
	if !nameRegex.MatchString(params.FirstName) {
		return c.Status(errs.ErrValidate.StatusCode).JSON(models.NewResponseError(errs.ErrValidate).SetHeader("validate", "first_name is not valid"))
	}

	if !nameRegex.MatchString(params.LastName) {
		return c.Status(errs.ErrValidate.StatusCode).JSON(models.NewResponseError(errs.ErrValidate).SetHeader("validate", "last_name is not valid"))
	}

	validator2 := validator.New()
	if err := validator2.Struct(params); err != nil {
		return c.Status(errs.ErrValidate.StatusCode).JSON(models.NewResponseError(errs.ErrValidate).SetHeader("validate", utils.ValidatorErrors(err)))
	}

	// Update user
	updates := map[string]interface{}{
		"first_name": params.FirstName,
		"last_name":  params.LastName,
	}
	err = db.Model(&models.User{}).Where("id = ?", user.ID.String()).Updates(updates).Error
	if err != nil {
		// Return status 500 and error message.
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Return status 200 OK.
	res := Response{
		FirstName: params.FirstName,
		LastName:  params.LastName,
	}

	return c.JSON(models.NewResponseOK(&res))
}

// SetPhone set user's phone
// @Description Set user's phone
// @Summary Set user's phone
// @Tags User
// @Accept json
// @Produce json
// @Param payload body controllers.SetPhone.Request true "set phone params"
// @Success 200 {object} models.ResponseOK{result=controllers.SetPhone.Response}
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /user/set-phone [post]
func SetPhone(c *fiber.Ctx) error {
	// Get user from request context.
	user := c.Locals("user").(models.User)

	type Request struct {
		Phone string `json:"phone" validate:"required" example:"05555555555"`
	}

	var params Request
	if err := c.BodyParser(&params); err != nil {
		// Return status 400 and error message.
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest))
	}

	// regex, starts with 0 and 10 digits
	regex := regexp.MustCompile(`^0[0-9]{10}$`)
	if !regex.MatchString(params.Phone) || params.Phone == "05555555555" {
		return c.Status(errs.ErrValidate.StatusCode).JSON(models.NewResponseError(errs.ErrValidate).SetHeader("validate", "phone is not valid"))
	}

	type Response struct {
		Timeout int `json:"timeout"`
	}

	// Open redis connection.
	con := database.NewRConnectionDB(database.RedisDatabasePhoneVerification)
	defer con.RClose()

	vfRes, err := utils.StartVerifyKitOTP(c.IP(), "9"+params.Phone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewResponseErr(errors.New("verification sms not send")).SetHeader("verifykit", err.Error()))
	}

	// Set phone verification code to redis
	key := fmt.Sprintf("%s:%s", user.ID.String(), "9"+params.Phone)
	err = con.RSetTTL(key, vfRes.Reference, time.Duration(vfRes.Timeout)*time.Second)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewResponseErr(errors.New("verification sms not send")).SetHeader("redis", err.Error()))
	}

	res := Response{
		Timeout: vfRes.Timeout,
	}

	return c.JSON(models.NewResponseOK(&res))
}

// VerifyPhone verify user's phone
// @Description Verify user's phone
// @Summary Verify user's phone
// @Tags User
// @Accept json
// @Produce json
// @Param payload body controllers.VerifyPhone.Request true "verify phone params"
// @Success 200 {object} models.ResponseOK{result=controllers.VerifyPhone.Response}
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /user/verify-phone [post]
func VerifyPhone(c *fiber.Ctx) error {
	// Get user from request context.
	user := c.Locals("user").(models.User)

	type Request struct {
		PhoneNumber string `json:"phone" validate:"required" example:"05555555555"`
		Code        string `json:"code" validate:"required,len=6,numeric" example:"123456"`
	}

	var params Request
	if err := c.BodyParser(&params); err != nil {
		// Return status 400 and error message.
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest))
	}

	fmt.Println(params)

	// regex, starts with 0 and 10 digits
	regex := regexp.MustCompile(`^0[0-9]{10}$`)
	if !regex.MatchString(params.PhoneNumber) || params.PhoneNumber == "05555555555" {
		return c.Status(errs.ErrValidate.StatusCode).JSON(models.NewResponseError(errs.ErrValidate).SetHeader("validate", "phone is not valid"))
	}

	validator2 := validator.New()
	if err := validator2.Struct(params); err != nil {
		return c.Status(errs.ErrValidate.StatusCode).JSON(models.NewResponseError(errs.ErrValidate).SetHeader("validate", utils.ValidatorErrors(err)))
	}

	// Open redis connection.
	con := database.NewRConnectionDB(database.RedisDatabasePhoneVerification)
	defer con.RClose()

	// Open database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and error message.
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	// Get phone verification code from redis
	key := fmt.Sprintf("%s:%s", user.ID.String(), "9"+params.PhoneNumber)
	reference, err := con.RGet(key)
	if err != nil || reference == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewResponseErr(errors.New("verification code not found")))
	}

	// Verify phone number
	err = utils.CheckVerifyKitOTP(c.IP(), "9"+params.PhoneNumber, reference, params.Code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewResponseErr(errors.New("verification code not valid")).SetHeader("verifykit", err.Error()))
	}

	// Delete phone verification code from redis
	con.RDel(key)

	// Update user
	updates := map[string]interface{}{
		"phone_number": "+9" + params.PhoneNumber,
		"validated":    true,
	}

	err = db.Model(&models.User{}).Where("id = ?", user.ID.String()).Updates(updates).Error
	if err != nil {
		// Return status 500 and error message.
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	type Response struct {
		PhoneNumber string `json:"phone"`
	}
	return c.JSON(models.NewResponseOK(&Response{
		PhoneNumber: "+9" + params.PhoneNumber,
	}))
}

// SetProfileImage method
// @Description Upload profile image
// @Summary Upload profile image
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file"
// @Success 200 {object} models.ResponseOK{result=controllers.SetProfileImage.Result}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /user/set-profile-image [post]
func SetProfileImage(c *fiber.Ctx) error {
	type Result struct {
		ID     uuid.UUID `json:"id"`
		Images []models.UserProfileImageInfo
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	// Get user from request context.
	user := c.Locals("user").(models.User)

	form, err := c.MultipartForm()
	if err == nil {
		// Get all files from form field "images".
		file := form.File["image"][0]

		// Is file image?
		if !strings.HasPrefix(file.Header["Content-Type"][0], "image/") {
			fmt.Println("hata2, file is not image")
			return c.Status(errs.ErrInvalidImage.StatusCode).JSON(models.NewResponseError(errs.ErrInvalidImage).SetHeader("image", "file is not image"))
		}

		// Get file extension.
		f, ee := file.Open()
		if ee != nil {
			// File not found or cannot be opened.
			fmt.Println("hata2, file not found or cannot be opened")
			return c.Status(errs.ErrUploadImage.StatusCode).JSON(models.NewResponseError(errs.ErrUploadImage).SetHeader("image", ee.Error()))
		}

		// Get file content type and check if it is a valid image.
		contentType, err := utils.GetFileContentType(f)
		if err != nil || contentType != file.Header.Get("Content-type") {
			fmt.Println("hata2, content type is not valid")
			return c.Status(errs.ErrInvalidImage.StatusCode).JSON(models.NewResponseError(errs.ErrInvalidImage).SetHeader("image", "content type is not valid").SetHeader("content-type", contentType).SetHeader("req-content-type", file.Header.Get("Content-type")))
		}
		if !utils.ValidImageContentType(file.Header.Get("Content-type")) {
			fmt.Println("hata2, content type is not valid")
			return c.Status(errs.ErrInvalidImage.StatusCode).JSON(models.NewResponseError(errs.ErrInvalidImage).SetHeader("image", "content type is not valid").SetHeader("content-type", contentType))
		}

		// Close file.
		f.Close()

		// Create files.
		images, err := utils.CreateProfileImageFile(file, "public/photos/profile")
		if err != nil || len(images) == 0 {
			fmt.Println("hata2, file cannot be created")
			return c.Status(errs.ErrImageProc.StatusCode).JSON(models.NewResponseError(errs.ErrImageProc).SetHeader("image", err.Error()))
		}

		// Create result
		result := Result{
			ID:     images[0].ID,
			Images: make([]models.UserProfileImageInfo, len(images)),
		}
		for i, image := range images {
			result.Images[i].URL = fmt.Sprintf("%s/photos/profile/%s", os.Getenv("API_URL"), image.Filename)
			result.Images[i].Width = image.Width
			result.Images[i].Height = image.Height
		}

		// Insert image to database.
		userProfileImage := models.UserProfileImage{
			ID:     images[0].ID,
			UserID: user.ID,
			Images: result.Images,
		}

		// Insert rental house image to database.
		err = db.Model(&models.UserProfileImage{}).Create(&userProfileImage).Error
		if err != nil {
			// Return status 500 and error message.
			return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
		}

		// Update user profile image id.
		err = db.Model(&models.User{}).Where("id = ?", user.ID.String()).Update("profile_image_id", images[0].ID).Error
		if err != nil {
			// Return status 500 and error message.
			return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
		}

		// Return status 200 OK.
		return c.JSON(models.NewResponseOK(&result))
	}
	fmt.Println("hata1", form, err)
	return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("form", err.Error()))
}
