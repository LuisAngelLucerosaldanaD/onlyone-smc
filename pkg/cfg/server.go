package cfg

import (
	"onlyone_smc/internal/models"
	"onlyone_smc/pkg/cfg/categories"
	"onlyone_smc/pkg/cfg/credentials_styles"
	"onlyone_smc/pkg/cfg/dictionaries"
	"onlyone_smc/pkg/cfg/messages"
	"onlyone_smc/pkg/cfg/shared_credential"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvDictionaries     dictionaries.PortsServerDictionaries
	SrvMessage          messages.PortsServerMessages
	SrvCategories       categories.PortsServerCategories
	SrvStyles           credentials_styles.PortsServerCredentialStyles
	SrvSharedCredential shared_credential.PortsServerSharedCredential
}

func NewServerCfg(db *sqlx.DB, user *models.User, txID string) *Server {

	repoDictionaries := dictionaries.FactoryStorage(db, user, txID)
	srvDictionaries := dictionaries.NewDictionariesService(repoDictionaries, user, txID)

	repoCategories := categories.FactoryStorage(db, user, txID)
	srvCategories := categories.NewCategoriesService(repoCategories, user, txID)

	repoMessage := messages.FactoryStorage(db, user, txID)
	srvMessage := messages.NewMessagesService(repoMessage, user, txID)

	repoStyles := credentials_styles.FactoryStorage(db, user, txID)
	srvStyles := credentials_styles.NewCredentialStylesService(repoStyles, user, txID)

	repoSharedCredential := shared_credential.FactoryStorage(db, user, txID)
	srvSharedCredential := shared_credential.NewSharedCredentialService(repoSharedCredential, user, txID)

	return &Server{
		SrvDictionaries:     srvDictionaries,
		SrvCategories:       srvCategories,
		SrvMessage:          srvMessage,
		SrvStyles:           srvStyles,
		SrvSharedCredential: srvSharedCredential,
	}
}
