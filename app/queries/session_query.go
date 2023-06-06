package queries

import (
	"ekira-backend/app/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionQueries struct {
	*gorm.DB
}

func (q *SessionQueries) UpdateSession(session *models.Session) error {
	return q.Save(session).Error
}

func (q *SessionQueries) GetSessionBySessionID(sessionID string) (models.Session, error) {
	// Define a new object.
	var session models.Session

	// Send query to database.
	err := q.Table("sessions").Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return session, nil
		}
		// Return empty object and error.
		return session, err
	}

	// Return object and nil error.
	return session, nil
}

func (q *SessionQueries) GetSessionsByUserID(userID uuid.UUID, page, limit, itemPerPage int) ([]models.Session, error) {
	// Define a new object.
	var sessions []models.Session

	offset := (page - 1) * itemPerPage

	// Send query to database.
	err := q.Unscoped().Table("sessions").Order("updated_at desc").Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&sessions).Error
	if err != nil {
		// Return empty object and error.
		return sessions, err
	}

	// Return object and nil error.
	return sessions, nil
}

func (q *SessionQueries) CreateSession(session models.Session) (models.Session, error) {
	// Send query to database.
	err := q.Table("sessions").Create(&session).Error
	if err != nil {
		// Return empty object and error.
		return session, err
	}

	// Return object and nil error.
	return session, nil
}

func (q *SessionQueries) DeleteSessionBySessionID(sessionID string) error {
	// Send query to database.
	err := q.Table("sessions").Where("session_id = ?", sessionID).Delete(&models.Session{}).Error
	if err != nil {
		// Return only error
		return err
	}
	return nil
}

func (q *SessionQueries) DeleteUserSessionBySessionID(userId uuid.UUID, sessionID string) (int64, error) {
	// Send query to database.
	tx := q.Table("sessions").Where("user_id = ? AND session_id = ?", userId.String(), sessionID).Delete(&models.Session{})
	if tx.Error != nil {
		// Return only error
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

func (q *SessionQueries) DeleteSessionByUserID(userId uuid.UUID) error {
	// Send query to database.
	err := q.Table("sessions").Where("user_id = ?", userId.String()).Delete(&models.Session{}).Error
	if err != nil {
		// Return only error
		return err
	}
	return nil
}

func (q *SessionQueries) DeleteSessionExByUserID(userId uuid.UUID, except []string) error {
	// Send query to database.
	err := q.Table("sessions").Where("user_id = ? AND session_id NOT IN ?", userId.String(), except).Delete(&models.Session{}).Error
	if err != nil {
		// Return only error
		return err
	}
	return nil
}
