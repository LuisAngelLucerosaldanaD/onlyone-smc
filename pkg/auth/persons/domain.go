package persons

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de Persons
type Persons struct {
	ID           int64     `json:"id" db:"id" valid:"-"`
	NumeroCedula string    `json:"numero_cedula" db:"numero_cedula" valid:"required"`
	NombreUno    string    `json:"nombre1" db:"nombre1" valid:"required"`
	NombreDos    string    `json:"nombre2" db:"nombre2" valid:"-"`
	ApellidoUno  string    `json:"apellido1" db:"apellido1" valid:"required"`
	ApellidoDos  string    `json:"apellido2" db:"apellido2" valid:"-"`
	Particula    string    `json:"particula" db:"particula" valid:"required"`
	Vigencia     string    `json:"vigencia" db:"vigencia" valid:"required"`
	FechaNac     string    `json:"fechaNac" db:"fechaNac" valid:"required"`
	FechaExp     string    `json:"fechaExp" db:"fechaExp" valid:"required"`
	Genero       string    `json:"genero" db:"genero" valid:"required"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func NewPerson(id int64, numeroCedula string, primerNombre string, segundoNombre string, primerApellido string, segundoApellido string, particula string, vigencia string, fechaNacimiento string, fechaExpedicion string, genero string) *Persons {
	return &Persons{
		ID:           id,
		NumeroCedula: numeroCedula,
		NombreUno:    primerNombre,
		NombreDos:    segundoNombre,
		ApellidoUno:  primerApellido,
		ApellidoDos:  segundoApellido,
		Particula:    particula,
		Vigencia:     vigencia,
		FechaNac:     fechaNacimiento,
		FechaExp:     fechaExpedicion,
		Genero:       genero,
	}
}

func NewCreatePerson(numeroCedula string, primerNombre string, segundoNombre string, primerApellido string, segundoApellido string, particula string, vigencia string, fechaNacimiento string, fechaExpedicion string, genero string) *Persons {
	return &Persons{
		NumeroCedula: numeroCedula,
		NombreUno:    primerNombre,
		NombreDos:    segundoNombre,
		ApellidoUno:  primerApellido,
		ApellidoDos:  segundoApellido,
		Particula:    particula,
		Vigencia:     vigencia,
		FechaNac:     fechaNacimiento,
		FechaExp:     fechaExpedicion,
		Genero:       genero,
	}
}

func (m *Persons) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
