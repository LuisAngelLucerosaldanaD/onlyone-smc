package persons

import (
	"fmt"

	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

type PortsServerPersons interface {
	GetPersonByIdentityNumber(identityNumber string) (*Persons, int, error)
}

type service struct {
	repository ServicesPersonsRepository
	user       *models.User
	txID       string
}

func NewPersonsService(repository ServicesPersonsRepository, user *models.User, TxID string) PortsServerPersons {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) GetPersonByIdentityNumber(identityNumber string) (*Persons, int, error) {
	if identityNumber == "" {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("identityNumber is required"))
		return nil, 15, fmt.Errorf("identityNumber is required")
	}
	m, err := s.repository.getByIdentityNumber(identityNumber)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByIdentityNumber row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
