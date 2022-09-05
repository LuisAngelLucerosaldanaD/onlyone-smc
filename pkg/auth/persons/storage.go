package persons

import (
	"github.com/jmoiron/sqlx"

	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesPersonsRepository interface {
	getByIdentityNumber(identityNumber string) (*Persons, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesPersonsRepository {
	var s ServicesPersonsRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newPersonsPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
