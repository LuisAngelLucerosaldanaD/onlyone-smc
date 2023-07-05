package credential_page

import (
	"github.com/jmoiron/sqlx"

	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesCredentialPageRepository interface {
	create(m *CredentialPage) error
	update(m *CredentialPage) error
	delete(id string) error
	getByID(id string) (*CredentialPage, error)
	getAll() ([]*CredentialPage, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesCredentialPageRepository {
	var s ServicesCredentialPageRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newCredentialPagePsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
