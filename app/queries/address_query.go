package queries

import (
	"ekira-backend/app/models"
	"gorm.io/gorm"
)

// AddressQueries struct
type AddressQueries struct {
	*gorm.DB
}

// GetCountries method for getting countries
func (q *AddressQueries) GetCountries(order string) (models.Countries, error) {
	countries := models.Countries{}

	q.Model(&models.Country{}).Order(order).Find(&countries)

	return countries, nil
}

// GetCities method for getting cities
func (q *AddressQueries) GetCities(countryId int, order string) (models.Cities, error) {
	cities := models.Cities{}

	err := q.Model(&models.City{}).Where("country_id = ?", countryId).Order(order).Preload("Country").Find(&cities).Error
	if err != nil {
		return cities, err
	}

	return cities, nil
}

// GetTowns method for getting towns
func (q *AddressQueries) GetTowns(cityId int, order string) (models.Towns, error) {
	towns := models.Towns{}

	q.Model(&models.Town{}).Where("city_id = ?", cityId).Order(order).Preload("City.Country").Find(&towns)

	return towns, nil
}

// GetDistricts method for getting districts
func (q *AddressQueries) GetDistricts(townId int, order string) (models.Districts, error) {
	districts := models.Districts{}

	q.Model(&models.District{}).Where("town_id = ?", townId).Order(order).Preload("Town.City.Country").Find(&districts)

	return districts, nil
}

// GetQuarters method for getting quarters
func (q *AddressQueries) GetQuarters(districtId int, order string) (models.Quarters, error) {
	quarters := models.Quarters{}

	q.Model(&models.Quarter{}).Where("district_id = ?", districtId).Order(order).Preload("District.Town.City.Country").Find(&quarters)

	return quarters, nil
}
