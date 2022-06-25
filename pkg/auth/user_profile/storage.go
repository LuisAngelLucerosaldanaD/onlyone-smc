package user_profile

import (
	"github.com/jmoiron/sqlx"

	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesUserProfileRepository interface {
	create(m *UserProfile) error
	update(m *UserProfile) error
	delete(id int64) error
	getByID(id int64) (*UserProfile, error)
	getAll() ([]*UserProfile, error)
	getByUserID(userID string) (*UserProfile, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesUserProfileRepository {
	var s ServicesUserProfileRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newUserProfilePsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
