package shared_credential

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// SharedCredential  Model struct SharedCredential
type SharedCredential struct {
	ID               int64     `json:"id" db:"id" valid:"-"`
	Data             string    `json:"data" db:"data" valid:"required"`
	UserId           string    `json:"user_id" db:"user_id" valid:"required"`
	Password         string    `json:"password" db:"password" valid:"required"`
	ExpiredAt        time.Time `json:"expired_at" db:"expired_at" valid:"required"`
	MaxNumberQueries int       `json:"max_number_queries" db:"max_number_queries" valid:"-"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

func NewSharedCredential(id int64, data string, userId string, password string, expiredAt time.Time, maxNumberQueries int) *SharedCredential {
	return &SharedCredential{
		ID:               id,
		Data:             data,
		UserId:           userId,
		Password:         password,
		ExpiredAt:        expiredAt,
		MaxNumberQueries: maxNumberQueries,
	}
}

func NewCreateSharedCredential(data string, userId string, password string, expiredAt time.Time, maxNumberQueries int) *SharedCredential {
	return &SharedCredential{
		Data:             data,
		UserId:           userId,
		Password:         password,
		ExpiredAt:        expiredAt,
		MaxNumberQueries: maxNumberQueries,
	}
}

func (m *SharedCredential) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
