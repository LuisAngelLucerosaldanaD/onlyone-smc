package shared_credential

import (
	"fmt"
	"onlyone_smc/internal/pwd"
	"time"

	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

type PortsServerSharedCredential interface {
	CreateSharedCredential(data string, userId string, password string, expiredAt time.Time, maxNumberQueries int) (*SharedCredential, int, error)
	UpdateSharedCredential(id int64, data string, userId string, password string, expiredAt time.Time, maxNumberQueries int) (*SharedCredential, int, error)
	DeleteSharedCredential(id int64) (int, error)
	GetSharedCredentialByID(id int64) (*SharedCredential, int, error)
	GetAllSharedCredential() ([]*SharedCredential, error)
}

type service struct {
	repository ServicesSharedCredentialRepository
	user       *models.User
	txID       string
}

func NewSharedCredentialService(repository ServicesSharedCredentialRepository, user *models.User, TxID string) PortsServerSharedCredential {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateSharedCredential(data string, userId string, password string, expiredAt time.Time, maxNumberQueries int) (*SharedCredential, int, error) {
	m := NewCreateSharedCredential(data, userId, password, expiredAt, maxNumberQueries)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	m.Password = pwd.Encrypt(m.Password)

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create SharedCredential :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateSharedCredential(id int64, data string, userId string, password string, expiredAt time.Time, maxNumberQueries int) (*SharedCredential, int, error) {
	m := NewSharedCredential(id, data, userId, password, expiredAt, maxNumberQueries)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update SharedCredential :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteSharedCredential(id int64) (int, error) {
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return 15, fmt.Errorf("id is required")
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

func (s *service) GetSharedCredentialByID(id int64) (*SharedCredential, int, error) {
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return nil, 15, fmt.Errorf("id is required")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllSharedCredential() ([]*SharedCredential, error) {
	return s.repository.getAll()
}
