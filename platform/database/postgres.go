package database

import (
	"ekira-backend/app/models"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"strconv"
	"time"
)

// PostgreSQLConnection func for connection to PostgreSQL database.
func PostgreSQLConnection() (*gorm.DB, error) {
	// Define database connection settings.
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	// Define database connection for PostgreSQL.
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_SERVER_URL")),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error, not connected to database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error, not connected to database")
	}

	sqlDB.SetMaxOpenConns(maxConn)                                         // the default is 0 (unlimited)
	sqlDB.SetMaxIdleConns(maxIdleConn)                                     // defaultMaxIdleConns = 2
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetimeConn) * time.Second) // 0, connections are reused forever

	if err := sqlDB.Ping(); err != nil {
		defer sqlDB.Close() // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}

// Migrate func for migrate database.
func Migrate(db *gorm.DB) error {
	// Schemas
	var models1 = []interface{}{
		&models.User{},
		&models.UserProfileImage{},
		&models.RentalHouse{},
		&models.RentalHouseImage{},
		&models.RentalHouseFavorite{},
		&models.Session{},
		&models.Reservation{},
		&models.Payment{},
	}
	// Migrate the schema
	if err := db.Debug().AutoMigrate(models1...); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Country{}); err == nil && db.Migrator().HasTable(&models.Country{}) {
		if err := db.First(&models.Country{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			db.Create(&models.Country{ID: 1, Name: "Türkiye", Abbreviation: "TR", Language: "tr", DisplayOrder: 1, SortOrder: 1, PhoneCode: "+90", Alpha2Code: "TR", Alpha3Code: "TUR"})
		}
	}
	if err := db.AutoMigrate(&models.City{}); err == nil && db.Migrator().HasTable(&models.City{}) {
		if err := db.First(&models.City{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			for _, city := range CitiesData {
				db.Create(&models.City{ID: city.ID, CountryID: city.CountryID, Name: city.Name, Tag: city.Tag, DisplayOrder: city.DisplayOrder, SortOrder: city.SortOrder})
			}
		}
	}
	if err := db.AutoMigrate(&models.Town{}); err == nil && db.Migrator().HasTable(&models.Town{}) {
		if err := db.First(&models.Town{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			query, err := GetTownsSQLFile()
			if err != nil {
				fmt.Println(err)
			} else {
				db.Exec(query)
			}
		}
	}
	if err := db.AutoMigrate(&models.District{}); err == nil && db.Migrator().HasTable(&models.District{}) {
		if err := db.First(&models.District{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			query, err := GetDistrictsSQLFile()
			if err != nil {
				fmt.Println(err)
			} else {
				db.Exec(query)
			}
		}
	}
	if err := db.AutoMigrate(&models.Quarter{}); err == nil && db.Migrator().HasTable(&models.Quarter{}) {
		if err := db.First(&models.Quarter{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			query, err := GetQuartersSQLFile()
			if err != nil {
				fmt.Println(err)
			} else {
				db.Exec(query)
			}
		}
	}
	return nil
}
