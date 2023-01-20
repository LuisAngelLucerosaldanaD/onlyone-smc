package shared_credential

import (
	"github.com/jmoiron/sqlx"

	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesSharedCredentialRepository interface {
	create(m *SharedCredential) error
	update(m *SharedCredential) error
	delete(id int64) error
	getByID(id int64) (*SharedCredential, error)
	getAll() ([]*SharedCredential, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesSharedCredentialRepository {
	var s ServicesSharedCredentialRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newSharedCredentialPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
