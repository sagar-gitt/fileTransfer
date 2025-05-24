package utils

import (
	"bytes"
	"html/template"
)

type EmailData struct {
	DownloadLink string
	LinkValidity string
}

func RenderEmailHTML(downloadLink, linkValidity string) (string, error) {
	tmpl, err := template.ParseFiles("../internal/utils/templates/email_template.html")
	if err != nil {
		return "", err
	}

	data := EmailData{
		DownloadLink: downloadLink,
		LinkValidity: linkValidity,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
