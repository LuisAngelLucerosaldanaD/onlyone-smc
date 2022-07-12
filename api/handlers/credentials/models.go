package credentials

import "github.com/google/uuid"

type responseCreateCredential struct {
	Error bool   `json:"error"`
	Data  string `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type requestCreateTransaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
	TypeId int     `json:"type_id"`
	Data   Data    `json:"data"`
}

type responseAllCredentials struct {
	Error bool          `json:"error"`
	Data  []*credential `json:"data"`
	Code  int           `json:"code"`
	Type  int           `json:"type"`
	Msg   string        `json:"msg"`
}

type JwtTransactionResponse struct {
	Error bool   `json:"error"`
	Data  string `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type ResGetFiles struct {
	Error bool    `json:"error"`
	Data  []*File `json:"data"`
	Code  int     `json:"code"`
	Type  int     `json:"type"`
	Msg   string  `json:"msg"`
}

type JwtTransactionRequest struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Block         int       `json:"block"`
	From          uuid.UUID `json:"from"`
	To            uuid.UUID `json:"to"`
	Verify        string    `json:"verify"`
	Ttl           int       `json:"ttl"`
	AttributesID  []int     `json:"attributes_id"`
	FilesID       []int     `json:"files_id"`
}

type credential struct {
	Id     string  `json:"id"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
	TypeId int     `json:"type_id"`
	Data   string  `json:"data"`
}

type Data struct {
	Category       string       `json:"category"`
	IdentityNumber string       `json:"identity_number"`
	Files          []*File      `json:"files"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	Identifiers    []Identifier `json:"identifiers"`
}

type File struct {
	FileID     int    `json:"id_file"`
	Name       string `json:"name"`
	FileEncode string `json:"file_encode"`
}

type Identifier struct {
	Name       string      `json:"name"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
