package auth

import (
	"github.com/jmoiron/sqlx"
	"onlyone_smc/internal/models"
	"onlyone_smc/pkg/auth/files"
	"onlyone_smc/pkg/auth/persons"
	"onlyone_smc/pkg/auth/user_profile"
)

type Server struct {
	SrvUserProfile user_profile.PortsServerUserProfile
	SrvFiles       files.PortsServerFile
	SrvPersons     persons.PortsServerPersons
}

func NewServerAuth(db *sqlx.DB, user *models.User, txID string) *Server {

	repoS3File := files.FactoryFileDocumentRepository(user, txID)
	srvFiles := files.NewFileService(repoS3File, user, txID)

	repoUserProfile := user_profile.FactoryStorage(db, user, txID)
	srvUserProfile := user_profile.NewUserProfileService(repoUserProfile, user, txID)

	repoPersons := persons.FactoryStorage(db, user, txID)
	srvPersons := persons.NewPersonsService(repoPersons, user, txID)

	return &Server{
		SrvFiles:       srvFiles,
		SrvUserProfile: srvUserProfile,
		SrvPersons:     srvPersons,
	}
}
