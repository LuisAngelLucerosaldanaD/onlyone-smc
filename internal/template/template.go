package template

import (
	"bytes"
	"html/template"
	"onlyone_smc/internal/logger"
)

func GenerateTemplateMail(param map[string]string, path string) (string, error) {
	bf := &bytes.Buffer{}
	tpl := &template.Template{}

	tpl = template.Must(template.New("").ParseGlob("templates/*.gohtml"))
	err := tpl.ExecuteTemplate(bf, path, &param)
	if err != nil {
		logger.Error.Printf("couldn't generate template body email: %v", err)
		return "", err
	}
	return bf.String(), err
}
