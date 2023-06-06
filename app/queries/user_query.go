package queries

import (
	"ekira-backend/app/models"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserQueries struct for queries from User model.
type UserQueries struct {
	*gorm.DB
}

// GetUserById method for getting one user by given ID.
func (q *UserQueries) GetUserById(id uuid.UUID) (models.User, error) {
	// Define user variable.
	user := models.User{}

	// Send query to database.
	err := q.Table("users").Preload("ProfileImage").Where("id = ?", id).First(&user).Error
	if err != nil {
		// Return empty object and error.
		return user, err
	}

	// Return query result.
	return user, nil
}

// GetUserByEmail method for getting one user by given Email
func (q *UserQueries) GetUserByEmail(email string) (models.User, error) {
	// Define user variable.
	user := models.User{}

	// Send query to database.
	err := q.Table("users").Preload("ProfileImage").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, nil
		}
		// Return empty object and error.
		return user, err
	}

	// Return query result.
	return user, nil
}

// GetUserByEmailAndPassword method for getting one user by given Email and Password (Hash).
func (q *UserQueries) GetUserByEmailAndPassword(email, password string) (models.User, error) {
	// Define user variable.
	user := models.User{}

	// Send query to database.
	err := q.Table("users").Preload("ProfileImage").Where("email = ? AND password = ?", email, password).First(&user).Error
	if err != nil {
		// Return empty object and error.
		return user, err
	}

	// Return query result.
	return user, nil
}

// NewUser method for creating user by given User object.
func (q *UserQueries) NewUser(u *models.User) error {
	// Send query to database.
	err := q.Table("users").Create(u).Error
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateUser method for updating user by given User object.
func (q *UserQueries) UpdateUser(id uuid.UUID, u *models.User) error {
	// Define query string.
	query := `UPDATE users SET first_name = $2, last_name = $3, email = $4, password = $5, validated = $6 WHERE id = $1`

	// Send query to database.
	err := q.Table("users").Exec(query, id, u.FirstName, u.LastName, u.Email, u.Password, u.Validated).Error
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// DeleteUser method for delete user by given ID.
func (q *UserQueries) DeleteUser(id uuid.UUID) error {
	// Define query string.
	query := `DELETE FROM users WHERE id = $1`

	// Send query to database.
	err := q.Table("users").Exec(query, id).Error
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
