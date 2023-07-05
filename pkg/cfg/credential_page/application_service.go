package credential_page

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

type PortsServerCredentialPage interface {
	CreateCredentialPage(id string, url string, ttl int) (*CredentialPage, int, error)
	UpdateCredentialPage(id string, url string, ttl int) (*CredentialPage, int, error)
	DeleteCredentialPage(id string) (int, error)
	GetCredentialPageByID(id string) (*CredentialPage, int, error)
	GetAllCredentialPage() ([]*CredentialPage, error)
}

type service struct {
	repository ServicesCredentialPageRepository
	user       *models.User
	txID       string
}

func NewCredentialPageService(repository ServicesCredentialPageRepository, user *models.User, TxID string) PortsServerCredentialPage {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateCredentialPage(id string, url string, ttl int) (*CredentialPage, int, error) {
	m := NewCredentialPage(id, url, ttl)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create CredentialPage :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateCredentialPage(id string, url string, ttl int) (*CredentialPage, int, error) {
	m := NewCredentialPage(id, url, ttl)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update CredentialPage :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteCredentialPage(id string) (int, error) {
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

func (s *service) GetCredentialPageByID(id string) (*CredentialPage, int, error) {
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

func (s *service) GetAllCredentialPage() ([]*CredentialPage, error) {
	return s.repository.getAll()
}
