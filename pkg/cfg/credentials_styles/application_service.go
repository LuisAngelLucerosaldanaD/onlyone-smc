package credentials_styles

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

type PortsServerCredentialStyles interface {
	CreateCredentialStyles(id string, typeStyle int, background string, logo string, front string, back string, credentialId string) (*CredentialStyles, int, error)
	UpdateCredentialStyles(id string, typeStyle int, background string, logo string, front string, back string, credentialId string) (*CredentialStyles, int, error)
	DeleteCredentialStyles(id string) (int, error)
	GetCredentialStylesByID(id string) (*CredentialStyles, int, error)
	GetAllCredentialStyles() ([]*CredentialStyles, error)
	GetCredentialStylesByCredentialID(credentialId string) ([]*CredentialStyles, int, error)
}

type service struct {
	repository ServicesCredentialStylesRepository
	user       *models.User
	txID       string
}

func NewCredentialStylesService(repository ServicesCredentialStylesRepository, user *models.User, TxID string) PortsServerCredentialStyles {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateCredentialStyles(id string, typeStyle int, background string, logo string, front string, back string, credentialId string) (*CredentialStyles, int, error) {
	m := NewCredentialStyles(id, typeStyle, background, logo, front, back, credentialId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create CredentialStyles :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateCredentialStyles(id string, typeStyle int, background string, logo string, front string, back string, credentialId string) (*CredentialStyles, int, error) {
	m := NewCredentialStyles(id, typeStyle, background, logo, front, back, credentialId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update CredentialStyles :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteCredentialStyles(id string) (int, error) {
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

func (s *service) GetCredentialStylesByID(id string) (*CredentialStyles, int, error) {
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

func (s *service) GetCredentialStylesByCredentialID(credentialId string) ([]*CredentialStyles, int, error) {
	if !govalidator.IsUUID(credentialId) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("credential id isn't uuid"))
		return nil, 15, fmt.Errorf("credential id isn't uuid")
	}
	m, err := s.repository.getByCredentialID(credentialId)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByCredentialID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllCredentialStyles() ([]*CredentialStyles, error) {
	return s.repository.getAll()
}
