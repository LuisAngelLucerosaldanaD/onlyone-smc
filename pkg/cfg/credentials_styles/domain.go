package credentials_styles

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de CredentialStyles
type CredentialStyles struct {
	ID           string    `json:"id" db:"id" valid:"required,uuid"`
	Type         int       `json:"type" db:"type" valid:"required"`
	Background   string    `json:"background" db:"background" valid:"required"`
	Logo         string    `json:"logo" db:"logo" valid:"required"`
	Front        string    `json:"front" db:"front" valid:"required"`
	Back         string    `json:"back" db:"back" valid:"required"`
	CredentialId string    `json:"credential_id" db:"credential_id" valid:"required"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func NewCredentialStyles(id string, typeStyle int, background string, logo string, front string, back string, credentialId string) *CredentialStyles {
	return &CredentialStyles{
		ID:           id,
		Type:         typeStyle,
		Background:   background,
		Logo:         logo,
		Front:        front,
		Back:         back,
		CredentialId: credentialId,
	}
}

func (m *CredentialStyles) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
