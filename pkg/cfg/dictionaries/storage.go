package dictionaries

import (
	"github.com/jmoiron/sqlx"

	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

const (
	Postgresql = "postgres"
	MongoDB    = "mongodb"
)

type ServicesDictionariesRepository interface {
	create(m *Dictionary) error
	update(m *Dictionary) error
	delete(id int) error
	getByID(id int) (*Dictionary, error)
	getAll() ([]*Dictionary, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesDictionariesRepository {
	var s ServicesDictionariesRepository
	var engine string
	if db != nil {
		engine = db.DriverName()
	} else {
		engine = MongoDB
	}

	switch engine {
	case Postgresql:
		return newDictionaryPsqlRepository(db, user, txID)

	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
