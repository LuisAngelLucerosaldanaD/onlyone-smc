package categories

import (
	"github.com/jmoiron/sqlx"
	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesCategoriesRepository interface {
	create(c *Categories) error
	update(c *Categories) error
	delete(id string) error
	getByID(id string) (*Categories, error)
	getAll() ([]*Categories, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesCategoriesRepository {
	var s ServicesCategoriesRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newCategoriesPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
