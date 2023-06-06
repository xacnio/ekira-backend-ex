package database

import (
	"ekira-backend/app/queries"
	"gorm.io/gorm"
)

// Queries struct for collect all app queries.
type Queries struct {
	*gorm.DB
	*queries.UserQueries        // load queries from User model
	*queries.AddressQueries     // load queries from Address model
	*queries.RentalHouseQueries // load queries from RentalHouse model
	*queries.SessionQueries     // load queries from Session model
	*queries.ReservationQueries // load queries from Reservation model
	*queries.PaymentQueries     // load queries from Payment model
}

// OpenDBConnection func for opening database connection.
func OpenDBConnection() (*Queries, error) {
	// Define a new PostgreSQL connection.
	db, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return &Queries{
		DB: db,
		// Set queries from models:
		UserQueries:        &queries.UserQueries{DB: db},        // from User model
		AddressQueries:     &queries.AddressQueries{DB: db},     // from Address model
		RentalHouseQueries: &queries.RentalHouseQueries{DB: db}, // from RentalHouse model
		SessionQueries:     &queries.SessionQueries{DB: db},     // from Session model
		ReservationQueries: &queries.ReservationQueries{DB: db}, // from Reservation model
		PaymentQueries:     &queries.PaymentQueries{DB: db},     // from Payment model
	}, nil
}
