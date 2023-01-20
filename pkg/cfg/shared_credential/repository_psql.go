package shared_credential

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

func newSharedCredentialPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *SharedCredential) error {
	const psqlInsert = `INSERT INTO cfg.shared_credential (data, user_id, password, expired_at, max_number_queries) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.Data,
		m.UserId,
		m.Password,
		m.ExpiredAt,
		m.MaxNumberQueries,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *SharedCredential) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE cfg.shared_credential SET data = :data, user_id = :user_id, password = :password, expired_at = :expired_at, max_number_queries = :max_number_queries, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM cfg.shared_credential WHERE id = :id `
	m := SharedCredential{ID: id}
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
func (s *psql) getByID(id int64) (*SharedCredential, error) {
	const psqlGetByID = `SELECT id , data, user_id, password, expired_at, max_number_queries, created_at, updated_at FROM cfg.shared_credential WHERE id = $1 `
	mdl := SharedCredential{}
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
func (s *psql) getAll() ([]*SharedCredential, error) {
	var ms []*SharedCredential
	const psqlGetAll = ` SELECT id , data, user_id, password, expired_at, max_number_queries, created_at, updated_at FROM cfg.shared_credential `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}
