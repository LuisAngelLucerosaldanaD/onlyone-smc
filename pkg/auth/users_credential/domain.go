package users_credential

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// UsersCredential  Model struct UsersCredential
type UsersCredential struct {
	ID             string    `json:"id" db:"id" valid:"required,uuid"`
	PrivateKey     string    `json:"private_key" db:"private_key" valid:"required"`
	IdentityNumber string    `json:"identity_number" db:"identity_number" valid:"required"`
	Mnemonic       string    `json:"mnemonic" db:"mnemonic" valid:"required"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

func NewUsersCredential(id string, privateKey string, identityNumber string, mnemonic string) *UsersCredential {
	return &UsersCredential{
		ID:             id,
		PrivateKey:     privateKey,
		IdentityNumber: identityNumber,
		Mnemonic:       mnemonic,
	}
}

func (m *UsersCredential) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
