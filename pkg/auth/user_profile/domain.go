package user_profile

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de UserProfile
type UserProfile struct {
	ID        int64     `json:"id" db:"id" valid:"-"`
	UserId    string    `json:"user_id" db:"user_id" valid:"required"`
	Name      string    `json:"name" db:"name" valid:"required"`
	Path      string    `json:"path" db:"path" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewUserProfile(id int64, userId string, name string, path string) *UserProfile {
	return &UserProfile{
		ID:     id,
		UserId: userId,
		Name:   name,
		Path:   path,
	}
}

func NewCreateUserProfile(userId string, name string, path string) *UserProfile {
	return &UserProfile{
		UserId: userId,
		Name:   name,
		Path:   path,
	}
}

func (m *UserProfile) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
