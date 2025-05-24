package dto

type EmailRequestBody struct {
	To           string `json:"to"`
	Bcc          string `json:"bcc"`
	Cc           string `json:"cc"`
	DownloadLink string `json:"link"`
	LinkValidity string `json:"linkValidity"`
}
