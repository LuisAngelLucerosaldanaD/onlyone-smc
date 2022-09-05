package persons

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"onlyone_smc/internal/models"
)

// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newPersonsPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

func (s *psql) getByIdentityNumber(identityNumber string) (*Persons, error) {
	const psqlGetByUserID = `select id, numero_cedula, nombre1, nombre2, apellido1, apellido2, particula, vigencia, fechaNac, fechaExp, genero from dbo.persons WHERE numero_cedula = $1`
	mdl := Persons{}
	err := s.DB.Get(&mdl, psqlGetByUserID, identityNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}
