package user_profile

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

func newUserProfilePsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *UserProfile) error {
	const psqlInsert = `INSERT INTO auth.user_profile (user_id, name, path) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.UserId,
		m.Name,
		m.Path,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *UserProfile) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE auth.user_profile SET user_id = :user_id, name = :name, path = :path, updated_at = :updated_at WHERE id = :id `
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
func (s *psql) delete(id int64) error {
	const psqlDelete = `DELETE FROM auth.user_profile WHERE id = :id `
	m := UserProfile{ID: id}
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
func (s *psql) getByID(id int64) (*UserProfile, error) {
	const psqlGetByID = `SELECT id , user_id, name, path, created_at, updated_at FROM auth.user_profile WHERE id = $1 `
	mdl := UserProfile{}
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
func (s *psql) getAll() ([]*UserProfile, error) {
	var ms []*UserProfile
	const psqlGetAll = ` SELECT id , user_id, name, path, created_at, updated_at FROM auth.user_profile `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getByUserID(userID string) (*UserProfile, error) {
	const psqlGetByUserID = `SELECT id , user_id, name, path, created_at, updated_at FROM auth.user_profile WHERE user_id = $1 `
	mdl := UserProfile{}
	err := s.DB.Get(&mdl, psqlGetByUserID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}
