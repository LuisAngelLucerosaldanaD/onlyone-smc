package auth

import (
	"github.com/jmoiron/sqlx"
	"onlyone_smc/internal/models"
	"onlyone_smc/pkg/auth/files"
	"onlyone_smc/pkg/auth/user_profile"
	"onlyone_smc/pkg/auth/users_credential"
)

type Server struct {
	SrvUserProfile     user_profile.PortsServerUserProfile
	SrvFiles           files.PortsServerFile
	SrvUsersCredential users_credential.PortsServerUsersCredential
}

func NewServerAuth(db *sqlx.DB, user *models.User, txID string) *Server {

	repoS3File := files.FactoryFileDocumentRepository(user, txID)
	srvFiles := files.NewFileService(repoS3File, user, txID)

	repoUserProfile := user_profile.FactoryStorage(db, user, txID)
	srvUserProfile := user_profile.NewUserProfileService(repoUserProfile, user, txID)

	repoUsersCredential := users_credential.FactoryStorage(db, user, txID)
	srvUsersCredential := users_credential.NewUsersCredentialService(repoUsersCredential, user, txID)

	return &Server{
		SrvFiles:           srvFiles,
		SrvUserProfile:     srvUserProfile,
		SrvUsersCredential: srvUsersCredential,
	}
}
