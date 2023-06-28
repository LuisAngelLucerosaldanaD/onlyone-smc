package credentials

import (
	"github.com/google/uuid"
	"time"
)

type responseCreateCredential struct {
	Error bool   `json:"error"`
	Data  resTrx `json:"data"`
	Code  int    `json:"code"`
	Type  int    `json:"type"`
	Msg   string `json:"msg"`
}

type resTrx struct {
	Id        string  `json:"id,omitempty"`
	From      string  `json:"from,omitempty"`
	To        string  `json:"to,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
	TypeId    int32   `json:"type_id,omitempty"`
	Data      string  `json:"data,omitempty"`
	Block     int64   `json:"block,omitempty"`
	Files     string  `json:"files,omitempty"`
	CreatedAt string  `json:"created_at,omitempty"`
	UpdatedAt string  `json:"updated_at,omitempty"`
}

type requestCreateTransaction struct {
	From     string  `json:"from"`
	To       string  `json:"to"`
	Amount   float64 `json:"amount"`
	TypeId   int     `json:"type_id"`
	Data     string  `json:"data"`
	CipherId string  `json:"cipher_id"`
	Files    []*File `json:"files"`
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
	Files  []File  `json:"files"`
}

type Data struct {
	Category       string       `json:"category"`
	Type           int          `json:"type"`
	IdentityNumber string       `json:"identity_number"`
	Files          []*File      `json:"files"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	Identifiers    []Identifier `json:"identifiers"`
	ExpiresAt      *time.Time   `json:"expires_at"`
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

type DataTrx struct {
	Category    string       `json:"category"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Identifiers []Identifier `json:"identifiers"`
	Type        int32        `json:"type"`
	Id          string       `json:"id"`         // id de la credencial
	Status      string       `json:"status"`     // estado de la credencial
	CreatedAt   string       `json:"created_at"` // fecha de creación de la credencial
	ExpiresAt   *time.Time   `json:"expires_at"` // fecha de vencimiento
}

type reqSharedCredentials struct {
	Data             []AttributeShared `json:"data"`
	Password         string            `json:"password"`
	ExpiredAt        time.Time         `json:"expired_at"`
	MaxNumberQueries int               `json:"max_number_queries"`
}

type ResAnny struct {
	Error bool        `json:"error"`
	Data  interface{} `json:"data"`
	Code  int         `json:"code"`
	Type  int         `json:"type"`
	Msg   string      `json:"msg"`
}

type Credential struct {
	Attributes []AttributeShared `json:"attributes"`
	Entity     Entity            `json:"entity"`
}

type Entity struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AttributeShared struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

/*
TODO get fee miners, validators and nodes
}*/
