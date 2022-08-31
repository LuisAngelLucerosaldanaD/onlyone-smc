package users

import (
	"onlyone_smc/internal/models"
	"time"
)

type requestCreateUser struct {
	Id              string    `json:"id"`
	Nickname        string    `json:"nickname"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirm_password"`
	Name            string    `json:"name"`
	Lastname        string    `json:"lastname"`
	IdType          int       `json:"id_type"`
	IdNumber        string    `json:"id_number"`
	Cellphone       string    `json:"cellphone"`
	BirthDate       time.Time `json:"birth_date"`
}

type responseCreateUser struct {
	Error bool   `json:"error"`
	Data  string `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type responseUser struct {
	Error bool         `json:"error"`
	Data  *models.User `json:"data"`
	Code  int          `json:"code"`
	Type  int          `json:"type"`
	Msg   string       `json:"msg"`
}

type requestActivateUser struct {
	Code string `json:"code"`
}

type responseActivateUser struct {
	Error bool   `json:"error"`
	Data  string `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type responseValidateUser struct {
	Error bool           `json:"error"`
	Data  *models.Wallet `json:"data"`
	Code  int            `json:"code"`
	Type  int            `json:"type"`
	Msg   string         `json:"msg"`
}

type responseUserValid struct {
	Error bool   `json:"error"`
	Data  string `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type responseValidateIdentity struct {
	Error bool   `json:"error"`
	Data  string `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type requestValidateIdentity struct {
	IdentityNumber string `json:"identity_number"`
	IdentityEncode string `json:"identity_encode"`
	ConfirmEncode  string `json:"confirm_encode"`
	Country        string `json:"country"`
}

type responseUpdateUser struct {
	Error bool   `json:"error"`
	Data  string `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type requestUserPhoto struct {
	FileEncode string `json:"file_encode"`
	FileName   string `json:"file_name"`
}

type requestUpdatePassword struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type responseGetWallets struct {
	Error bool             `json:"error"`
	Data  []*models.Wallet `json:"data"`
	Code  int              `json:"code"`
	Type  int              `json:"type"`
	Msg   string           `json:"msg"`
}
