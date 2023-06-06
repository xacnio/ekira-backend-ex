package queries

import (
	"ekira-backend/app/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// RentalHouseQueries struct
type RentalHouseQueries struct {
	*gorm.DB
}

// NewRentalHouse method for create new rental house.
func (q *RentalHouseQueries) NewRentalHouse(rh *models.RentalHouse) error {
	// Send query to database.
	err := q.Model(models.RentalHouse{}).Create(rh).Preload("Images").Preload("Quarter.District.Town.City.Country").First(rh).Error
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateRentalHouse method for create new rental house image.
func (q *RentalHouseQueries) UpdateRentalHouse(rh *models.RentalHouse) error {
	// Insert query to database.
	err := q.Save(rh).Error
	if err != nil {
		// Return only error.
		return err
	}
	return nil
}

// GetRentalHouseWithUid method for get rental house with uid.
func (q *RentalHouseQueries) GetRentalHouseWithUid(uid uuid.UUID) (models.RentalHouse, error) {
	// Define user variable.
	house := models.RentalHouse{}

	// Send query to database.
	err := q.Debug().Model(models.RentalHouse{}).Where("uid = ?", uid.String()).Preload("Images").Preload("Quarter.District.Town.City.Country").Preload("Creator." + clause.Associations).First(&house).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return house, nil
		}
		// Return empty object and error.
		return house, err
	}

	// Return query result.
	return house, nil
}

type RentalHouseList struct {
	TotalCount int64
	FullCount  int64
	NextPage   bool
	PrevPage   bool
	Houses     []models.RentalHouse
}

// IsUsersFavorite method the rental house is favorite of user.
func (q *RentalHouseQueries) IsUsersFavorite(userId uuid.UUID, id int) (bool, error) {
	// Send query to database.
	var count int64
	err := q.Model(models.RentalHouseFavorite{}).Where("creator = ? AND rental_house_id = ?", userId.String(), id).Count(&count).Error
	if err != nil {
		// Return empty object and error.
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	// Return query result.
	return false, nil
}

// GetFavoriteRentalHouseList method for get rental house list. (results, prev, next, error)
func (q *RentalHouseQueries) GetFavoriteRentalHouseList(userId uuid.UUID, pagination *models.Pagination) (RentalHouseList, error) {
	// Define variables.
	data := RentalHouseList{}
	offset := (pagination.Page - 1) * pagination.Limit
	if pagination.Page > 1 {
		data.PrevPage = true
	}

	filterArgs := pagination.Filters.ToPreloadQuery()
	var favoriteHouses []models.RentalHouseFavorite

	// Get full count over()
	q.Model(models.RentalHouseFavorite{}).Where("creator = ?", userId.String()).Order(pagination.Sort).Preload("RentalHouse", filterArgs...).Preload(clause.Associations).Select("count(*) over() as full_count").Count(&data.FullCount)
	if data.FullCount > 0 {
		if int(data.FullCount) > offset+pagination.Limit {
			data.NextPage = true
		}
	}

	// Send query to database.
	err := q.Model(models.RentalHouseFavorite{}).Where("creator = ?", userId.String()).Limit(pagination.Limit).Offset(offset).Order(pagination.Sort).Preload("RentalHouse", filterArgs...).Preload(clause.Associations).Find(&favoriteHouses).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, nil
		}
		// Return empty object and error.
		return data, err
	}
	data.TotalCount = int64(len(favoriteHouses))

	for _, favoriteHouse := range favoriteHouses {
		data.Houses = append(data.Houses, favoriteHouse.RentalHouse)
	}

	// Return query result.
	return data, nil
}

// GetRentalHouseList method for get rental house list. (results, prev, next, error)
func (q *RentalHouseQueries) GetRentalHouseList(pagination *models.Pagination) (RentalHouseList, error) {
	// Define variables.
	data := RentalHouseList{}
	offset := (pagination.Page - 1) * pagination.Limit
	if pagination.Page > 1 {
		data.PrevPage = true
	}

	filterQuery, filterArgs := pagination.Filters.ToWhereQuery()

	// Get full count over()
	q.Model(models.RentalHouse{}).Order(pagination.Sort).Where("published IS TRUE").Where(filterQuery, filterArgs...).Select("count(*) OVER()").Scan(&data.FullCount)
	if data.FullCount > 0 {
		if int(data.FullCount) > offset+pagination.Limit {
			data.NextPage = true
		}
	}

	// Send query to database.
	err := q.Model(models.RentalHouse{}).Limit(pagination.Limit).Offset(offset).Order(pagination.Sort).Where("published IS TRUE").Preload("Images").Preload("Quarter.District.Town.City.Country").Find(&data.Houses).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, nil
		}
		// Return empty object and error.
		return data, err
	}
	data.TotalCount = int64(len(data.Houses))

	// Return query result.
	return data, nil
}

// GetRentalHouseOwnedList method for get rental house, only for owner. (results, prev, next, error)
func (q *RentalHouseQueries) GetRentalHouseOwnedList(creatorId uuid.UUID, pagination *models.Pagination) (RentalHouseList, error) {
	// Define variables.
	data := RentalHouseList{}
	offset := (pagination.Page - 1) * pagination.Limit
	if pagination.Page > 1 {
		data.PrevPage = true
	}

	// Get full count over()
	q.Model(models.RentalHouse{}).Order(pagination.Sort).Where("creator = ?", creatorId).Select("count(*) OVER()").Scan(&data.FullCount)
	if data.FullCount > 0 {
		if int(data.FullCount) > offset+pagination.Limit {
			data.NextPage = true
		}
	}

	// Send query to database.
	err := q.Model(models.RentalHouse{}).Limit(pagination.Limit).Offset(offset).Order(pagination.Sort).Where("creator = ?", creatorId).Preload("Images").Preload("Quarter.District.Town.City.Country").Find(&data.Houses).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, nil
		}
		// Return empty object and error.
		return data, err
	}
	data.TotalCount = int64(len(data.Houses))

	// Return query result.
	return data, nil
}

// CreateRentalHouseImage method for create new rental house image.
func (q *RentalHouseQueries) CreateRentalHouseImage(rhi *models.RentalHouseImage) error {
	// Insert query to database.
	err := q.Model(models.RentalHouseImage{}).Create(rhi).Error
	if err != nil {
		// Return only error.
		return err
	}
	return nil
}

// GetRentalHouseImagesByRentalHouseID method for create new rental house image.
func (q *RentalHouseQueries) GetRentalHouseImagesByRentalHouseID(id int) ([]models.RentalHouseImage, error) {
	// Init data
	var data []models.RentalHouseImage

	// Insert query to database.
	err := q.Model(models.RentalHouseImage{}).Where("rental_house_id = ?", id).Preload("RentalHouse").Preload("Creator").Find(&data).Error
	if err != nil {
		// Return only error.
		return data, err
	}
	return data, nil
}

// GetRentalHouseImageByID method for create new rental house image.
func (q *RentalHouseQueries) GetRentalHouseImageByID(id string) (models.RentalHouseImage, error) {
	// Init data
	var data models.RentalHouseImage

	// Insert query to database.
	err := q.Model(models.RentalHouseImage{}).Where("id = ?", id).Preload("RentalHouse").Preload("Creator").First(&data).Error
	if err != nil {
		// Return only error.
		return data, err
	}
	return data, nil
}

// GetRentalHouseImagesByIDs method for create new rental house image.
func (q *RentalHouseQueries) GetRentalHouseImagesByIDs(id []string) ([]models.RentalHouseImage, error) {
	fmt.Println(id)
	// Init data
	var data []models.RentalHouseImage

	// Insert query to database.
	err := q.Model(models.RentalHouseImage{}).Where("id IN ?", id).Where("rental_house_id IS NULL").Preload("RentalHouse").Preload("Creator").Find(&data).Error
	if err != nil {
		// Return only error.
		return []models.RentalHouseImage{}, err
	}
	return data, nil
}

// UpdateRentalHouseImage method for create new rental house image.
func (q *RentalHouseQueries) UpdateRentalHouseImage(rh models.RentalHouseImage) error {
	// Insert query to database.
	err := q.Model(models.RentalHouseImage{}).Where("id = ?", rh.ID).Updates(rh).Error
	if err != nil {
		// Return only error.
		return err
	}
	return nil
}

// GetReservedDatesByRentalHouseID method for create new rental house image.
func (q *RentalHouseQueries) GetReservedDatesByRentalHouseID(id int) ([]time.Time, error) {
	// Define user variable.
	reservations := []models.Reservation{}

	// Send query to database.
	err := q.Debug().Model(models.Reservation{}).Where("rental_house_id = ? AND status NOT IN (4,5) AND (status != 1 OR (expire > ? AND status = 1))", id, time.Now()).Find(&reservations).Error
	if err != nil {
		// Return empty object and error.
		return []time.Time{}, err
	}

	// Define dates variable.
	var dates []time.Time

	// Loop through reservations.
	for _, reservation := range reservations {
		startDate := reservation.StartDate
		endDate := reservation.EndDate
		// Loop through every day between start and end date.
		for startDate.Before(endDate) {
			// Append date to dates.
			dates = append(dates, startDate)
			// Add one day to start date.
			startDate = startDate.AddDate(0, 0, 1)
		}
	}

	// Return query result.
	return dates, nil
}
