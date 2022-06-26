package categories

import "onlyone_smc/pkg/cfg/categories"

type responseAllCategories struct {
	Error bool                     `json:"error"`
	Data  []*categories.Categories `json:"data"`
	Code  int                      `json:"code"`
	Type  int                      `json:"type"`
	Msg   string                   `json:"msg"`
}
