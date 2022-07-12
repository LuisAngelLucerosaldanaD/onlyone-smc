package credentials_styles

import (
	"github.com/jmoiron/sqlx"

	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesCredentialStylesRepository interface {
	create(m *CredentialStyles) error
	update(m *CredentialStyles) error
	delete(id string) error
	getByID(id string) (*CredentialStyles, error)
	getAll() ([]*CredentialStyles, error)
	getByCredentialID(credentialId string) ([]*CredentialStyles, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesCredentialStylesRepository {
	var s ServicesCredentialStylesRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newCredentialStylesPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
