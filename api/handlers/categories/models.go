package categories

import "onlyone_smc/pkg/cfg/categories"

type responseAllCategories struct {
	Error bool             `json:"error"`
	Data  []*resCategories `json:"data"`
	Code  int              `json:"code"`
	Type  int              `json:"type"`
	Msg   string           `json:"msg"`
}

type requestCreateStyle struct {
	Type       int    `json:"type"`
	Background string `json:"background"`
	Logo       string `json:"logo"`
	Front      string `json:"front"`
	Back       string `json:"back"`
	CategoryID string `json:"category_id"`
}

type resCategories struct {
	Category *categories.Categories `json:"category"`
	Styles   []Styles               `json:"styles"`
}

type Styles struct {
	Type        int           `json:"type"`
	Background  string        `json:"background"`
	Logo        string        `json:"logo"`
	Identifiers []identifiers `json:"identifiers"`
}

type identifiers struct {
	Type       string `json:"type"`
	Attributes string `json:"attributes"`
}

type resAny struct {
	Error bool        `json:"error"`
	Data  interface{} `json:"data"`
	Code  int         `json:"code"`
	Type  int         `json:"type"`
	Msg   string      `json:"msg"`
}
