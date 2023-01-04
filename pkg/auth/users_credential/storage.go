package users_credential

import (
	"github.com/jmoiron/sqlx"

	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesUsersCredentialRepository interface {
	create(m *UsersCredential) error
	update(m *UsersCredential) error
	delete(id string) error
	getByID(id string) (*UsersCredential, error)
	getAll() ([]*UsersCredential, error)
	getByIdentityNumber(identityNumber string) (*UsersCredential, error)
	deleteByIdentityNumber(identityNumber string) error
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesUsersCredentialRepository {
	var s ServicesUsersCredentialRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newUsersCredentialPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
