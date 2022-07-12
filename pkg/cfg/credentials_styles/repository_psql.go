package credentials_styles

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"onlyone_smc/internal/models"
	"time"
)

// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newCredentialStylesPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *CredentialStyles) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	const psqlInsert = `INSERT INTO cfg.credentials_styles (id ,type, background, logo, front, back, credential_id, created_at, updated_at) VALUES (:id ,:type, :background, :logo, :front, :back, :credential_id,:created_at, :updated_at) `
	rs, err := s.DB.NamedExec(psqlInsert, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *CredentialStyles) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE cfg.credentials_styles SET type = :type, background = :background, logo = :logo, front = :front, back = :back, credential_id = :credential_id, updated_at = :updated_at WHERE id = :id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *psql) delete(id string) error {
	const psqlDelete = `DELETE FROM cfg.credentials_styles WHERE id = :id `
	m := CredentialStyles{ID: id}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *psql) getByID(id string) (*CredentialStyles, error) {
	const psqlGetByID = `SELECT id , type, background, logo, front, back, credential_id, created_at, updated_at FROM cfg.credentials_styles WHERE id = $1 `
	mdl := CredentialStyles{}
	err := s.DB.Get(&mdl, psqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) getAll() ([]*CredentialStyles, error) {
	var ms []*CredentialStyles
	const psqlGetAll = ` SELECT id , type, background, logo, front, back, credential_id, created_at, updated_at FROM cfg.credentials_styles `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// getByCredentialID consulta un registro por el id de la credencial
func (s *psql) getByCredentialID(credentialId string) ([]*CredentialStyles, error) {
	var ms []*CredentialStyles
	const psqlGetAllByCategoryID = ` SELECT id , type, background, logo, front, back, credential_id, created_at, updated_at FROM cfg.credentials_styles where credential_id = '%s' `

	err := s.DB.Select(&ms, fmt.Sprintf(psqlGetAllByCategoryID, credentialId))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}
