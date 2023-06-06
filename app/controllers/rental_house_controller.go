package controllers

import (
	"ekira-backend/app/errs"
	"ekira-backend/app/models"
	"ekira-backend/app/queries"
	"ekira-backend/pkg/utils"
	"ekira-backend/platform/database"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CreateRentalHouse method
// @Description Create a new rental house
// @Summary Create a new rental house
// @Tags Rental House
// @Accept json
// @Produce json
// @Param rentalHouse body controllers.CreateRentalHouse.Request true "Rental House"
// @Success 200 {object} models.ResponseOK{result=controllers.GetDetails.Result{creator=controllers.GetDetails.Creator}}
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /rental-house/create [post]
func CreateRentalHouse(c *fiber.Ctx) error {
	// New database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection))
	}

	// Get user from request context.
	user2 := c.Locals("user").(models.User)

	type Request struct {
		Title         string               `json:"title" example:"My rental house" required:"true"`
		Description   string               `json:"description" example:"3 rooms, 2 bathrooms, 1 kitchen, 1 living room" required:"true"`
		QuarterID     int                  `json:"quarter_id" example:"1" required:"true"`
		RentPeriod    int                  `json:"rent_period" example:"1" summary:"1 = daily, 2 = monthly, 3 = yearly" required:"true,min=1,max=3"`
		Price         float64              `json:"price" example:"100.00" required:"true"`
		MinDay        *int                 `json:"min_day" example:"1" required:"false"`
		CommisionType models.CommisionType `json:"commision_type" example:"0" summary:"0 = renter pays, 1 = owner pays" required:"true,min=0,max=1"`
		Lat           string               `json:"g_coordinate,omitempty"`
		ImageUUIDs    []string             `json:"imageUUIDs" swaggertype:"array,string" example:""`
	}

	// Get rental house from request.
	request := Request{}
	if err := c.BodyParser(&request); err != nil {
		fmt.Println("Create rental house body parser error:", err)
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("rentalHouse", err.Error()))
	}

	rentalHouse := models.RentalHouse{}
	rentalHouse.CreatorID = user2.ID
	rentalHouse.Creator = user2
	rentalHouse.Title = request.Title
	rentalHouse.Description = request.Description
	rentalHouse.QuarterID = request.QuarterID
	rentalHouse.RentPeriod = request.RentPeriod
	rentalHouse.Price = request.Price
	rentalHouse.CommisionType = request.CommisionType
	if rentalHouse.RentPeriod == models.RentPeriodDay && request.MinDay != nil {
		rentalHouse.MinDay = *request.MinDay
	} else {
		rentalHouse.MinDay = 1
	}
	//rentalHouse.GCoordinate = rentalHouseBody.GCoordinate
	rentalHouse.CreatedAt = time.Now()
	rentalHouse.UpdatedAt = time.Now()

	x, y := 0.0, 0.0
	if _, err := fmt.Sscanf(rentalHouse.GCoordinate, "(%f,%f)", &x, &y); err != nil {
		x, y = 0.0, 0.0
	}
	rentalHouse.GCoordinate = fmt.Sprintf("(%f,%f)", x, y)

	// Validate rental house.
	validate := utils.NewValidator()

	if err := validate.Struct(rentalHouse); err != nil {
		// Return, if some fields are not valid.
		fmt.Println("Create rental house validator error:", err)
		return c.Status(errs.ErrValidate.StatusCode).JSON(models.NewResponseError(errs.ErrValidate).SetHeader("validate", utils.ValidatorErrors(err)))
	}

	// Create new rental house.
	err = db.NewRentalHouse(&rentalHouse)
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	if rentalHouse.ID == 0 {
		return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseError(errs.ErrNotFound).SetHeader("db", "Rental house not found"))
	}

	// Create new rental house images.
	images, err := db.GetRentalHouseImagesByIDs(request.ImageUUIDs)
	if err == nil {
		for _, image := range images {
			image.RentalHouseID = &rentalHouse.ID
			err := db.UpdateRentalHouseImage(image)
			if err != nil {
				fmt.Println(err)
				continue
			}
			rentalHouse.Images = append(rentalHouse.Images, image)
		}
	}

	type Creator struct {
		ID           string  `json:"id"`
		FirstName    string  `json:"first_name"`
		LastName     string  `json:"last_name"`
		Email        string  `json:"email"`
		Phone        string  `json:"phone"`
		ProfileImage *string `json:"profile_image"`
	}

	type Result struct {
		ID            string                          `json:"id"`
		Title         string                          `json:"title"`
		Price         float64                         `json:"price"`
		MinDay        int                             `json:"min_day"`
		RentPeriod    int                             `json:"rent_period"`
		CommisionType string                          `json:"commision_type"`
		Address       models.Quarter                  `json:"address"`
		Images        [][]models.RentalHouseImageInfo `json:"images"`
		Description   string                          `json:"description"`
		Published     bool                            `json:"published"`
		Creator       interface{}                     `json:"creator"`
	}

	creator := Creator{
		ID:        rentalHouse.CreatorID.String(),
		FirstName: rentalHouse.Creator.FirstName,
		LastName:  rentalHouse.Creator.LastName,
		Email:     rentalHouse.Creator.Email,
		Phone:     rentalHouse.Creator.PhoneNumber,
	}

	result := Result{
		ID:            rentalHouse.UID.String(),
		Title:         rentalHouse.Title,
		Price:         rentalHouse.Price,
		MinDay:        rentalHouse.MinDay,
		RentPeriod:    rentalHouse.RentPeriod,
		CommisionType: rentalHouse.CommisionTypeInfo(),
		Address:       rentalHouse.Quarter,
		Images:        make([][]models.RentalHouseImageInfo, len(rentalHouse.Images)),
		Description:   rentalHouse.Description,
		Published:     rentalHouse.Published,
	}

	if rentalHouse.Creator.ProfileImage != nil {
		profileImage := rentalHouse.Creator.ProfileImage
		images := profileImage.Images
		if len(images) > 0 {
			lastItem := images[len(images)-1]
			creator.ProfileImage = &lastItem.URL
		}
	}
	result.Creator = creator

	for i, images := range rentalHouse.Images {
		result.Images[i] = make([]models.RentalHouseImageInfo, len(images.Images))
		for j, image := range images.Images {
			result.Images[i][j].URL = image.URL
			result.Images[i][j].Width = image.Width
			result.Images[i][j].Height = image.Height
		}
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&result))
}

// GetOwnedList method
// @Description Get user's owned rental houses
// @Summary Get user's owned rental houses
// @Tags Rental House
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit number" default(10)
// @Param sort query string false "Sort by" default(created_at:desc)
// @Param search query string false "Search by title" default()
// @Success 200 {object} models.ResponseOK{result=controllers.GetOwnedList.Response{results=[]controllers.GetOwnedList.Result}}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /rental-house/owned-list [get]
func GetOwnedList(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Query parameters.
	// @Param page query int false "Page number (default: 1)"
	// @Param limit query int false "Limit number of items per page"
	// @Param search query string false "Search string (title)"
	// @Param sort query string false "Sort string (field:asc/desc)"
	limit := 5
	page := 1
	//search := ""
	sort := "created_at:desc"

	if _limit, err := strconv.Atoi(c.Query("limit")); err == nil {
		limit = _limit
	}
	if _page, err := strconv.Atoi(c.Query("page")); err == nil {
		page = _page
	}
	//if _search := strings.TrimSpace(c.Query("search")); _search != "" {
	//	search = _search
	//}
	r := regexp.MustCompile(`^[a-z_A-Z0-9]:(asc|desc)$`)
	if _sort := strings.TrimSpace(c.Query("sort")); _sort != "" && r.MatchString(_sort) {
		sort = _sort
	}
	field := strings.Split(sort, ":")[0]
	order := strings.Split(sort, ":")[1]
	sort = fmt.Sprintf("%s %s", field, order)

	// Pagination.
	pagination := models.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	// Get all rental houses.
	rentalHouseList, err := db.GetRentalHouseOwnedList(user.ID, &pagination)
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}
	type Result struct {
		ID            string  `json:"id"`
		Title         string  `json:"title"`
		Price         float64 `json:"price"`
		MinDay        int     `json:"min_day"`
		RentPeriod    int     `json:"rent_period"`
		CommisionType string  `json:"commision_type"`
		Address       struct {
			QuarterID    int    `json:"quarter_id"`
			QuarterName  string `json:"quarter_name"`
			DistrictID   int    `json:"district_id"`
			DistrictName string `json:"district_name"`
			TownID       int    `json:"town_id"`
			TownName     string `json:"town_name"`
			CityID       int    `json:"city_id"`
			CityName     string `json:"city_name"`
			CountryID    int    `json:"country_id"`
			CountryName  string `json:"country_name"`
		} `json:"address"`
		Images    [][]models.RentalHouseImageInfo `json:"images"`
		Published bool                            `json:"published"`
	}
	type Response struct {
		Pagination struct {
			TotalCount int64 `json:"total_count"`
			FullCount  int64 `json:"full_count"`
			NextPage   bool  `json:"next_page"`
			PrevPage   bool  `json:"prev_page"`
		} `json:"pagination"`
		Results interface{} `json:"results"`
	}
	res := Response{}
	res.Pagination.FullCount = rentalHouseList.FullCount
	res.Pagination.TotalCount = rentalHouseList.TotalCount
	res.Pagination.NextPage = rentalHouseList.NextPage
	res.Pagination.PrevPage = rentalHouseList.PrevPage
	result := make([]Result, rentalHouseList.TotalCount)
	res.Results = &result
	if rentalHouseList.TotalCount == 0 {
		return c.JSON(models.NewResponseOK(&res))
	}
	i := 0
	for _, rentalHouse := range rentalHouseList.Houses {
		result[i].ID = rentalHouse.UID.String()
		result[i].Title = rentalHouse.Title
		result[i].Price = rentalHouse.Price
		result[i].MinDay = rentalHouse.MinDay
		result[i].RentPeriod = rentalHouse.RentPeriod
		result[i].Address.QuarterID = rentalHouse.Quarter.ID
		result[i].Address.QuarterName = rentalHouse.Quarter.Name
		result[i].Address.DistrictID = rentalHouse.Quarter.District.ID
		result[i].Address.DistrictName = rentalHouse.Quarter.District.Name
		result[i].Address.TownID = rentalHouse.Quarter.District.Town.ID
		result[i].Address.TownName = rentalHouse.Quarter.District.Town.Name
		result[i].Address.CityID = rentalHouse.Quarter.District.Town.City.ID
		result[i].Address.CityName = rentalHouse.Quarter.District.Town.City.Name
		result[i].Address.CountryID = rentalHouse.Quarter.District.Town.City.Country.ID
		result[i].Address.CountryName = rentalHouse.Quarter.District.Town.City.Country.Name
		result[i].Images = [][]models.RentalHouseImageInfo{}
		result[i].Published = rentalHouse.Published
		result[i].CommisionType = rentalHouse.CommisionTypeInfo()
		for _, image := range rentalHouse.Images {
			var dt []models.RentalHouseImageInfo
			for _, imageInfo := range image.Images {
				dt = append(dt, imageInfo)
			}
			result[i].Images = append(result[i].Images, dt)
		}
		i++
		if i == limit {
			break
		}
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&res))
}

// GetPublicList method
// @Description Get all public rental houses
// @Summary Get all public rental houses
// @Tags Rental House
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit number" default(10)
// @Param sort query string false "Sort by" default(created_at:desc)
// @Param search query string false "Search by title" default()
// @Param favorite query bool false "Only get favorite rental houses" default()
// @Success 200 {object} models.ResponseOK{result=controllers.GetPublicList.Response{results=[]controllers.GetPublicList.Result}}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /rental-house/list [get]
func GetPublicList(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	// Query parameters.
	// @Param page query int false "Page number (default: 1)"
	// @Param limit query int false "Limit number of items per page"
	// @Param search query string false "Search string (title)"
	// @Param sort query string false "Sort string (field:asc/desc)"
	// @Param favorite query bool false "Only get favorite rental houses"
	limit := 5
	page := 1
	favorite := false
	//search := ""
	sort := "created_at:desc"

	if c.Query("favorite") == "true" {
		favorite = true
	}

	if _limit, err := strconv.Atoi(c.Query("limit")); err == nil {
		limit = _limit
	}
	if _page, err := strconv.Atoi(c.Query("page")); err == nil {
		page = _page
	}
	//if _search := strings.TrimSpace(c.Query("search")); _search != "" {
	//	search = _search
	//}
	r := regexp.MustCompile(`^[a-z_A-Z0-9]:(asc|desc)$`)
	if _sort := strings.TrimSpace(c.Query("sort")); _sort != "" && r.MatchString(_sort) {
		sort = _sort
	}
	field := strings.Split(sort, ":")[0]
	order := strings.Split(sort, ":")[1]
	sort = fmt.Sprintf("%s %s", field, order)

	// Pagination.
	pagination := models.Pagination{
		Page:    page,
		Limit:   limit,
		Sort:    sort,
		Filters: models.Filter{},
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	var rentalHouseList queries.RentalHouseList

	if favorite {
		pagination.Filters["published"] = true
		// Get favorite rental houses.
		rentalHouseList, err = db.GetFavoriteRentalHouseList(user.ID, &pagination)
		if err != nil {
			return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
		}
	} else {
		// Get all rental houses.
		rentalHouseList, err = db.GetRentalHouseList(&pagination)
		if err != nil {
			return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
		}
	}

	type Result struct {
		ID            string  `json:"id"`
		Title         string  `json:"title"`
		Price         float64 `json:"price"`
		MinDay        int     `json:"min_day"`
		RentPeriod    int     `json:"rent_period"`
		CommisionType string  `json:"commision_type"`
		Address       struct {
			QuarterID    int    `json:"quarter_id"`
			QuarterName  string `json:"quarter_name"`
			DistrictID   int    `json:"district_id"`
			DistrictName string `json:"district_name"`
			TownID       int    `json:"town_id"`
			TownName     string `json:"town_name"`
			CityID       int    `json:"city_id"`
			CityName     string `json:"city_name"`
			CountryID    int    `json:"country_id"`
			CountryName  string `json:"country_name"`
		} `json:"address"`
		Images   [][]models.RentalHouseImageInfo `json:"images"`
		Favorite bool                            `json:"favorite"`
	}
	type Response struct {
		Pagination struct {
			TotalCount int64 `json:"total_count"`
			FullCount  int64 `json:"full_count"`
			NextPage   bool  `json:"next_page"`
			PrevPage   bool  `json:"prev_page"`
		} `json:"pagination"`
		Results interface{} `json:"results"`
	}
	res := Response{}
	res.Pagination.FullCount = rentalHouseList.FullCount
	res.Pagination.TotalCount = rentalHouseList.TotalCount
	res.Pagination.NextPage = rentalHouseList.NextPage
	res.Pagination.PrevPage = rentalHouseList.PrevPage
	result := make([]Result, rentalHouseList.TotalCount)
	res.Results = &result
	if rentalHouseList.TotalCount == 0 {
		return c.JSON(models.NewResponseOK(&res))
	}
	i := 0
	for _, rentalHouse := range rentalHouseList.Houses {
		result[i].ID = rentalHouse.UID.String()
		result[i].Title = rentalHouse.Title
		result[i].Price = rentalHouse.Price
		result[i].MinDay = rentalHouse.MinDay
		result[i].RentPeriod = rentalHouse.RentPeriod
		result[i].Address.QuarterID = rentalHouse.Quarter.ID
		result[i].Address.QuarterName = rentalHouse.Quarter.Name
		result[i].Address.DistrictID = rentalHouse.Quarter.District.ID
		result[i].Address.DistrictName = rentalHouse.Quarter.District.Name
		result[i].Address.TownID = rentalHouse.Quarter.District.Town.ID
		result[i].Address.TownName = rentalHouse.Quarter.District.Town.Name
		result[i].Address.CityID = rentalHouse.Quarter.District.Town.City.ID
		result[i].Address.CityName = rentalHouse.Quarter.District.Town.City.Name
		result[i].Address.CountryID = rentalHouse.Quarter.District.Town.City.Country.ID
		result[i].Address.CountryName = rentalHouse.Quarter.District.Town.City.Country.Name
		result[i].Images = [][]models.RentalHouseImageInfo{}
		result[i].CommisionType = rentalHouse.CommisionTypeInfo()
		fav, _ := db.IsUsersFavorite(user.ID, rentalHouse.ID)
		result[i].Favorite = fav
		for _, image := range rentalHouse.Images {
			var dt []models.RentalHouseImageInfo
			for _, imageInfo := range image.Images {
				dt = append(dt, imageInfo)
			}
			result[i].Images = append(result[i].Images, dt)
		}
		i++
		if i == limit {
			break
		}
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&res))
}

// UploadRentalHouseImage method
// @Description Upload rental house image
// @Summary Upload rental house image
// @Tags Rental House
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file"
// @Success 200 {object} models.ResponseOK{result=controllers.UploadRentalHouseImage.Result}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /rental-house/upload-image [post]
func UploadRentalHouseImage(c *fiber.Ctx) error {
	type Result struct {
		ID     uuid.UUID `json:"id"`
		Images []models.RentalHouseImageInfo
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	// Get user from request context.
	user2 := c.Locals("user").(models.User)

	if form, err := c.MultipartForm(); err == nil {
		// Get all files from form field "images".
		file := form.File["image"][0]

		// Is file image?
		if !strings.HasPrefix(file.Header["Content-Type"][0], "image/") {
			return c.Status(errs.ErrInvalidImage.StatusCode).JSON(models.NewResponseError(errs.ErrInvalidImage))
		}

		// Get file extension.
		f, ee := file.Open()
		if ee != nil {
			// File not found or cannot be opened.
			return c.Status(errs.ErrUploadImage.StatusCode).JSON(models.NewResponseError(errs.ErrUploadImage))
		}

		// Get file content type and check if it is a valid image.
		contentType, err := utils.GetFileContentType(f)
		if err != nil || contentType != file.Header.Get("Content-type") {
			return c.Status(errs.ErrInvalidImage.StatusCode).JSON(models.NewResponseError(errs.ErrInvalidImage))
		}
		if !utils.ValidImageContentType(file.Header.Get("Content-type")) {
			return c.Status(errs.ErrInvalidImage.StatusCode).JSON(models.NewResponseError(errs.ErrInvalidImage))
		}

		// Close file.
		f.Close()

		// Create files.
		images, err := utils.CreateImageFile(file, "public/photos/rental-house")
		if err != nil || len(images) == 0 {
			return c.Status(errs.ErrImageProc.StatusCode).JSON(models.NewResponseError(errs.ErrImageProc).SetHeader("image", err.Error()))
		}

		// Create result
		result := Result{
			ID:     images[0].ID,
			Images: make([]models.RentalHouseImageInfo, len(images)),
		}
		for i, image := range images {
			result.Images[i].URL = fmt.Sprintf("%s/photos/rental-house/%s", os.Getenv("API_URL"), image.Filename)
			result.Images[i].Width = image.Width
			result.Images[i].Height = image.Height
		}

		// Insert image to database.
		rentalHouseImage := models.RentalHouseImage{
			ID:        images[0].ID,
			CreatorID: user2.ID,
			Expire:    time.Now().AddDate(0, 0, 1),
			Images:    result.Images,
		}

		// Insert rental house image to database.
		err = db.CreateRentalHouseImage(&rentalHouseImage)
		if err != nil {
			// Return status 500 and error message.
			return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
		}

		// Return status 200 OK.
		return c.JSON(models.NewResponseOK(&result))
	}
	return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest))
}

// GetDetails method
// @Description Get rental house details
// @Summary Get rental house details
// @Tags Rental House
// @Accept json
// @Produce json
// @Param id path string true "Rental house ID"
// @Success 200 {object} models.ResponseOK{result=controllers.GetDetails.Result{creator=controllers.GetDetails.Creator}}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /rental-house/{id} [get]
func GetDetails(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	id := c.Params("id")
	validate := validator.New()
	err := validate.Var(id, "required,uuid4")
	if err != nil {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("id", err.Error()))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	rentalHouse, err := db.GetRentalHouseWithUid(uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseError(errs.ErrNotFound))
		}
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	if rentalHouse.Published == false && rentalHouse.CreatorID != user.ID {
		return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseError(errs.ErrNotFound))
	}

	type Creator struct {
		ID           string  `json:"id"`
		FirstName    string  `json:"first_name"`
		LastName     string  `json:"last_name"`
		Email        string  `json:"email"`
		Phone        string  `json:"phone"`
		ProfileImage *string `json:"profile_image"`
	}

	type Result struct {
		ID            string                          `json:"id"`
		Title         string                          `json:"title"`
		Price         float64                         `json:"price"`
		MinDay        int                             `json:"min_day"`
		RentPeriod    int                             `json:"rent_period"`
		CommisionType string                          `json:"commision_type"`
		Address       models.Quarter                  `json:"address"`
		Images        [][]models.RentalHouseImageInfo `json:"images"`
		Description   string                          `json:"description"`
		Published     bool                            `json:"published"`
		Creator       interface{}                     `json:"creator"`
	}

	creator := Creator{
		ID:        rentalHouse.CreatorID.String(),
		FirstName: rentalHouse.Creator.FirstName,
		LastName:  rentalHouse.Creator.LastName,
		Email:     rentalHouse.Creator.Email,
		Phone:     rentalHouse.Creator.PhoneNumber,
	}

	result := Result{
		ID:            rentalHouse.UID.String(),
		Title:         rentalHouse.Title,
		Price:         rentalHouse.Price,
		MinDay:        rentalHouse.MinDay,
		RentPeriod:    rentalHouse.RentPeriod,
		CommisionType: rentalHouse.CommisionTypeInfo(),
		Address:       rentalHouse.Quarter,
		Images:        make([][]models.RentalHouseImageInfo, len(rentalHouse.Images)),
		Description:   rentalHouse.Description,
		Published:     rentalHouse.Published,
	}

	if rentalHouse.Creator.ProfileImage != nil {
		profileImage := rentalHouse.Creator.ProfileImage
		images := profileImage.Images
		if len(images) > 0 {
			lastItem := images[len(images)-1]
			creator.ProfileImage = &lastItem.URL
		}
	}
	result.Creator = creator

	for i, images := range rentalHouse.Images {
		result.Images[i] = make([]models.RentalHouseImageInfo, len(images.Images))
		for j, image := range images.Images {
			result.Images[i][j].URL = image.URL
			result.Images[i][j].Width = image.Width
			result.Images[i][j].Height = image.Height
		}
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&result))
}

// EditRentalHouse method
// @Description Edit rental house
// @Summary Edit rental house
// @Tags Rental House
// @Accept json
// @Produce json
// @Param id path string true "Rental house ID"
// @Param rentalHouse body controllers.EditRentalHouse.Request true "Rental House"
// @Success 200 {object} models.ResponseOK{result=models.RentalHouse}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /rental-house/{id} [put]
func EditRentalHouse(c *fiber.Ctx) error {
	id := c.Params("id")
	validate := validator.New()
	err := validate.Var(id, "required,uuid4")
	if err != nil {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("id", err.Error()))
	}

	user := c.Locals("user").(models.User)

	type Request struct {
		Title       *string  `json:"title" example:"My rental house" required:"false"`
		Description *string  `json:"description" example:"3 rooms, 2 bathrooms, 1 kitchen, 1 living room" required:"false"`
		QuarterID   *int     `json:"quarter_id" example:"1" required:"false"`
		RentPeriod  *int     `json:"rent_period" example:"1" summary:"1 = daily, 2 = monthly, 3 = yearly" required:"false,min=1,max=3"`
		Price       *float64 `json:"price" example:"100.00" required:"true"`
		MinDay      *int     `json:"min_day" example:"1"`
		Lat         *string  `json:"g_coordinate,omitempty"`
		Published   *bool    `json:"published" example:"true" required:"false"`
	}

	// Parse request body.
	var body Request
	err = c.BodyParser(&body)
	if err != nil {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("body", err.Error()))
	}

	// Validate request body.
	err = validate.Struct(body)
	if err != nil {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("body", err.Error()))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	// Get all rental houses.
	rentalHouse, err := db.GetRentalHouseWithUid(uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseError(errs.ErrNotFound))
		}
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	if rentalHouse.CreatorID != user.ID {
		return c.Status(errs.ErrForbidden.StatusCode).JSON(models.NewResponseError(errs.ErrForbidden))
	}

	// Check not nil fields and non empty fields.
	if body.Title != nil && *body.Title != "" {
		rentalHouse.Title = *body.Title
	}
	if body.Published != nil {
		rentalHouse.Published = *body.Published
	}
	if body.Description != nil && *body.Description != "" {
		rentalHouse.Description = *body.Description
	}
	if body.QuarterID != nil && *body.QuarterID != 0 {
		rentalHouse.QuarterID = *body.QuarterID
	}
	if body.RentPeriod != nil && *body.RentPeriod != 0 {
		rentalHouse.RentPeriod = *body.RentPeriod
	}
	if body.Price != nil && *body.Price != 0 {
		rentalHouse.Price = *body.Price
	}
	if body.MinDay != nil && *body.MinDay > 0 {
		rentalHouse.MinDay = *body.MinDay
	}
	if body.Lat != nil && *body.Lat != "" {
		rentalHouse.GCoordinate = *body.Lat
	}

	// Update rental house.
	err = db.UpdateRentalHouse(&rentalHouse)
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	updatedRentalHouse, err := db.GetRentalHouseWithUid(rentalHouse.UID)
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK(&updatedRentalHouse))
}

// GetReservedDates method
// @Description Get rental house not available dates
// @Summary Get rental house not available dates
// @Tags Rental House
// @Accept json
// @Produce json
// @Param id path string true "Rental house ID"
// @Success 200 {object} models.ResponseOK{result=controllers.GetReservedDates.Response}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /rental-house/{id}/reserved-dates [get]
func GetReservedDates(c *fiber.Ctx) error {
	id := c.Params("id")
	validate := validator.New()
	err := validate.Var(id, "required,uuid4")
	if err != nil {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("id", err.Error()))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	// Get all rental houses.
	rentalHouse, err := db.GetRentalHouseWithUid(uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseError(errs.ErrNotFound))
		}
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Get reserved dates.
	reservedDates, err := db.GetReservedDatesByRentalHouseID(rentalHouse.ID)
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Return status 200 OK.
	type Response []string
	res := make(Response, len(reservedDates))

	for i := range reservedDates {
		res[i] = reservedDates[i].Format("2006-01-02")
	}

	return c.JSON(models.NewResponseOK(&res))
}

// FavoriteRentalHouse method
// @Description Favorite a rental house
// @Summary Favorite a rental house
// @Tags Rental House
// @Accept json
// @Produce json
// @Param id path string true "Rental house ID"
// @Success 200 {object} models.ResponseOK{result=string}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /rental-house/{id}/favorite [get]
func FavoriteRentalHouse(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	id := c.Params("id")
	validate := validator.New()
	err := validate.Var(id, "required,uuid4")
	if err != nil {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("id", err.Error()))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	// Get all rental houses.
	rentalHouse, err := db.GetRentalHouseWithUid(uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseError(errs.ErrNotFound))
		}
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	if rentalHouse.Published == false && rentalHouse.CreatorID != user.ID {
		return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseError(errs.ErrNotFound))
	}

	// Check if user already favorite this rental house.
	var rentalHouseFavorite models.RentalHouseFavorite
	err = db.Model(&models.RentalHouseFavorite{}).Where("creator = ? AND rental_house_id = ?", user.ID, rentalHouse.ID).First(&rentalHouseFavorite).Error
	if err == nil && rentalHouseFavorite.ID != 0 {
		return c.JSON(models.NewResponseOK("OK"))
	}

	// Create rental house favorite.
	rentalHouseFavorite = models.RentalHouseFavorite{
		CreatorID:     user.ID,
		RentalHouseID: rentalHouse.ID,
	}

	err = db.Create(&rentalHouseFavorite).Error
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK("OK"))
}

// UnfavoriteRentalHouse method
// @Description Unfavorite a rental house
// @Summary Unfavorite a rental house
// @Tags Rental House
// @Accept json
// @Produce json
// @Param id path string true "Rental house ID"
// @Success 200 {object} models.ResponseOK{result=string}
// @Failure 404 {object} models.ResponseErr
// @Failure 400 {object} models.ResponseErr
// @Failure 500 {object} models.ResponseErr
// @Security Authentication
// @Router /rental-house/{id}/unfavorite [get]
func UnfavoriteRentalHouse(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	id := c.Params("id")
	validate := validator.New()
	err := validate.Var(id, "required,uuid4")
	if err != nil {
		return c.Status(errs.ErrBadRequest.StatusCode).JSON(models.NewResponseError(errs.ErrBadRequest).SetHeader("id", err.Error()))
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(errs.ErrDatabaseConnection.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseConnection).SetHeader("db", err.Error()))
	}

	// Get all rental houses.
	rentalHouse, err := db.GetRentalHouseWithUid(uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseError(errs.ErrNotFound))
		}
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	if rentalHouse.Published == false && rentalHouse.CreatorID != user.ID {
		return c.Status(errs.ErrNotFound.StatusCode).JSON(models.NewResponseError(errs.ErrNotFound))
	}

	// Check if user already favorite this rental house.
	var rentalHouseFavorite models.RentalHouseFavorite
	err = db.Model(&models.RentalHouseFavorite{}).Where("creator = ? AND rental_house_id = ?", user.ID, rentalHouse.ID).First(&rentalHouseFavorite).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(models.NewResponseOK("OK"))
		}
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Delete rental house favorite.
	err = db.Delete(&rentalHouseFavorite).Error
	if err != nil {
		return c.Status(errs.ErrDatabaseQuery.StatusCode).JSON(models.NewResponseError(errs.ErrDatabaseQuery).SetHeader("db", err.Error()))
	}

	// Return status 200 OK.
	return c.JSON(models.NewResponseOK("OK"))
}
