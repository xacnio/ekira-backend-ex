package controllers

import (
	"ekira-backend/app/errs"
	"ekira-backend/app/models"
	"ekira-backend/platform/database"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// GetCountries method
// @Description Get all countries
// @Summary Get all countries
// @Tags Address
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseOK{result=models.Countries}
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Router /get-countries [get]
func GetCountries(c *fiber.Ctx) error {
	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection))
	}

	// Get Countries
	countries, err := db.GetCountries("display_order ASC")
	if err != nil {
		// Return status 500 and error message.
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}
	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&countries))
}

// GetCitiesByCountry method
// @Description Get cities by country id
// @Summary Get cities by country id
// @Tags Address
// @Accept json
// @Produce json
// @Param countryId path int true "Country ID"
// @Success 200 {object} models.ResponseOK{result=models.Cities}
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Router /get-cities/{countryId} [get]
func GetCitiesByCountry(c *fiber.Ctx) error {
	// Get country id from request.
	countryId, err := strconv.Atoi(c.Params("countryId"))
	if err != nil {
		// Return status 400 and error message.
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("countryId", err.Error()))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection))
	}

	// Get Cities
	cities, err := db.GetCities(countryId, "display_order ASC")
	if err != nil {
		// Return status 500 and error message.
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}
	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&cities))
}

// GetTownsByCity method
// @Description Get towns by city id
// @Summary Get towns by city id
// @Tags Address
// @Accept json
// @Produce json
// @Param cityId path int true "City ID"
// @Success 200 {object} models.ResponseOK{result=models.Towns}
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Router /get-towns/{cityId} [get]
func GetTownsByCity(c *fiber.Ctx) error {
	// Get city id from request.
	cityId, err := strconv.Atoi(c.Params("cityId"))
	if err != nil {
		// Return status 400 and error message.
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("cityId", err.Error()))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection))
	}

	// Get Towns
	towns, err := db.GetTowns(cityId, "display_order ASC")
	if err != nil {
		// Return status 500 and error message.
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&towns))
}

// GetDistrictsByTown method
// @Description Get districts by town id
// @Summary Get districts by town id
// @Tags Address
// @Accept json
// @Produce json
// @Param townId path int true "Town ID"
// @Success 200 {object} models.ResponseOK{result=models.Districts}
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Router /get-districts/{townId} [get]
func GetDistrictsByTown(c *fiber.Ctx) error {
	// Get town id from request.
	townId, err := strconv.Atoi(c.Params("townId"))
	if err != nil {
		// Return status 400 and error message.
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("townId", err.Error()))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection))
	}

	// Get Districts
	districts, err := db.GetDistricts(townId, "display_order ASC")
	if err != nil {
		// Return status 500 and error message.
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&districts))
}

// GetQuartersByDistrict method
// @Description Get quarters by district id
// @Summary Get quarters by district id
// @Tags Address
// @Accept json
// @Produce json
// @Param districtId path int true "District ID"
// @Success 200 {object} models.ResponseOK{result=models.Quarters}
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Router /get-quarters/{districtId} [get]
func GetQuartersByDistrict(c *fiber.Ctx) error {
	// Get district id from request.
	districtId, err := strconv.Atoi(c.Params("districtId"))
	if err != nil {
		// Return status 400 and error message.
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("districtId", err.Error()))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection))
	}

	// Get Quarters
	quarters, err := db.GetQuarters(districtId, "display_order ASC")
	if err != nil {
		// Return status 500 and error message.
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&quarters))
}
