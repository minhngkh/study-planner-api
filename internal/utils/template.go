package utils

import (
	"bytes"
	"html/template"
)

func CreateHtml(template *template.Template, data any) (string, error) {
	var content bytes.Buffer
	err := template.Execute(&content, data)
	if err != nil {
		return "", err
	}

	return content.String(), nil
}
