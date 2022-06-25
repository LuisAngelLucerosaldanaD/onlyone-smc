package user_profile

import (
	"fmt"

	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
)

type PortsServerUserProfile interface {
	CreateUserProfile(userId string, name string, path string) (*UserProfile, int, error)
	UpdateUserProfile(id int64, userId string, name string, path string) (*UserProfile, int, error)
	DeleteUserProfile(id int64) (int, error)
	GetUserProfileByID(id int64) (*UserProfile, int, error)
	GetAllUserProfile() ([]*UserProfile, error)
	GetUserProfileByUserID(userId string) (*UserProfile, int, error)
}

type service struct {
	repository ServicesUserProfileRepository
	user       *models.User
	txID       string
}

func NewUserProfileService(repository ServicesUserProfileRepository, user *models.User, TxID string) PortsServerUserProfile {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateUserProfile(userId string, name string, path string) (*UserProfile, int, error) {
	m := NewCreateUserProfile(userId, name, path)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create UserProfile :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateUserProfile(id int64, userId string, name string, path string) (*UserProfile, int, error) {
	m := NewUserProfile(id, userId, name, path)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update UserProfile :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteUserProfile(id int64) (int, error) {
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

func (s *service) GetUserProfileByID(id int64) (*UserProfile, int, error) {
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

func (s *service) GetAllUserProfile() ([]*UserProfile, error) {
	return s.repository.getAll()
}

func (s *service) GetUserProfileByUserID(userId string) (*UserProfile, int, error) {
	if userId == "" {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("user id is required"))
		return nil, 15, fmt.Errorf("user id is required")
	}
	m, err := s.repository.getByUserID(userId)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getBy User ID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
