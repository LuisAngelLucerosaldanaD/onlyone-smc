package credential_page

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// CredentialPage  Model struct CredentialPage
type CredentialPage struct {
	ID        string    `json:"id" db:"id" valid:"required,uuid"`
	Url       string    `json:"url" db:"url" valid:"required"`
	Ttl       int       `json:"ttl" db:"ttl" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewCredentialPage(id string, url string, ttl int) *CredentialPage {
	return &CredentialPage{
		ID:  id,
		Url: url,
		Ttl: ttl,
	}
}

func (m *CredentialPage) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
