package users_credential

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

type PortsServerUsersCredential interface {
	CreateUsersCredential(id string, privateKey string, identityNumber string, mnemonic string) (*UsersCredential, int, error)
	UpdateUsersCredential(id string, privateKey string, identityNumber string, mnemonic string) (*UsersCredential, int, error)
	DeleteUsersCredential(id string) (int, error)
	GetUsersCredentialByID(id string) (*UsersCredential, int, error)
	GetAllUsersCredential() ([]*UsersCredential, error)
	GetUsersCredentialByIdentityNumber(identityNumber string) (*UsersCredential, int, error)
	DeleteUsersCredentialByIdentityNumber(identityNumber string) (int, error)
}

type service struct {
	repository ServicesUsersCredentialRepository
	user       *models.User
	txID       string
}

func NewUsersCredentialService(repository ServicesUsersCredentialRepository, user *models.User, TxID string) PortsServerUsersCredential {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateUsersCredential(id string, privateKey string, identityNumber string, mnemonic string) (*UsersCredential, int, error) {
	m := NewUsersCredential(id, privateKey, identityNumber, mnemonic)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create UsersCredential :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateUsersCredential(id string, privateKey string, identityNumber string, mnemonic string) (*UsersCredential, int, error) {
	m := NewUsersCredential(id, privateKey, identityNumber, mnemonic)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update UsersCredential :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteUsersCredential(id string) (int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return 15, fmt.Errorf("id isn't uuid")
	}

	if err := s.repository.delete(id); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetUsersCredentialByID(id string) (*UsersCredential, int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllUsersCredential() ([]*UsersCredential, error) {
	return s.repository.getAll()
}

func (s *service) GetUsersCredentialByIdentityNumber(identityNumber string) (*UsersCredential, int, error) {
	m, err := s.repository.getByIdentityNumber(identityNumber)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByIdentityNumber row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) DeleteUsersCredentialByIdentityNumber(identityNumber string) (int, error) {
	if err := s.repository.deleteByIdentityNumber(identityNumber); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't deleteByIdentityNumber row:", err)
		return 20, err
	}
	return 28, nil
}
