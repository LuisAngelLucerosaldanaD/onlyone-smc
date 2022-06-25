package cfg

import (
	"onlyone_smc/internal/models"
	"onlyone_smc/pkg/cfg/categories"
	"onlyone_smc/pkg/cfg/dictionaries"
	"onlyone_smc/pkg/cfg/messages"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvDictionaries dictionaries.PortsServerDictionaries
	SrvMessage      messages.PortsServerMessages
	SrvCategories   categories.PortsServerCategories
}

func NewServerCfg(db *sqlx.DB, user *models.User, txID string) *Server {

	repoDictionaries := dictionaries.FactoryStorage(db, user, txID)
	srvDictionaries := dictionaries.NewDictionariesService(repoDictionaries, user, txID)

	repoCategories := categories.FactoryStorage(db, user, txID)
	srvCategories := categories.NewCategoriesService(repoCategories, user, txID)

	repoMessage := messages.FactoryStorage(db, user, txID)
	srvMessage := messages.NewMessagesService(repoMessage, user, txID)

	return &Server{
		SrvDictionaries: srvDictionaries,
		SrvCategories:   srvCategories,
		SrvMessage:      srvMessage,
	}
}
