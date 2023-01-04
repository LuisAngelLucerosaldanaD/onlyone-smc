package users_credential

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"onlyone_smc/internal/models"
)

// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newUsersCredentialPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *UsersCredential) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	const psqlInsert = `INSERT INTO auth.users_credential (id ,private_key, identity_number, mnemonic, created_at, updated_at) VALUES (:id ,:private_key, :identity_number, :mnemonic,:created_at, :updated_at) `
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
func (s *psql) update(m *UsersCredential) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE auth.users_credential SET private_key = :private_key, identity_number = :identity_number, mnemonic = :mnemonic, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM auth.users_credential WHERE id = :id `
	m := UsersCredential{ID: id}
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
func (s *psql) getByID(id string) (*UsersCredential, error) {
	const psqlGetByID = `SELECT id , private_key, identity_number, mnemonic, created_at, updated_at FROM auth.users_credential WHERE id = $1 `
	mdl := UsersCredential{}
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
func (s *psql) getAll() ([]*UsersCredential, error) {
	var ms []*UsersCredential
	const psqlGetAll = ` SELECT id , private_key, identity_number, mnemonic, created_at, updated_at FROM auth.users_credential `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// GetByID consulta un registro por su ID
func (s *psql) getByIdentityNumber(identityNumber string) (*UsersCredential, error) {
	const psqlGetByIdentityNumber = `SELECT id , private_key, identity_number, mnemonic, created_at, updated_at FROM auth.users_credential WHERE identity_number = $1 `
	mdl := UsersCredential{}
	err := s.DB.Get(&mdl, psqlGetByIdentityNumber, identityNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// Delete elimina un registro de la BD
func (s *psql) deleteByIdentityNumber(identityNumber string) error {
	const psqlDelete = `DELETE FROM auth.users_credential WHERE identity_number = :id `
	m := UsersCredential{IdentityNumber: identityNumber}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}
