package models

import (
	"time"
)

type File struct {
	ID             string    `json:"id"`
	S3Key          string    `json:"s3_key"`
	Name           string    `json:"name"`
	Size           int64     `json:"size"`
	ExpirationDate time.Time `json:"expiration_date"`
	UserId         string    `json:"user_id"`
	DownloadLink   string    `json:"download_link"`
	UploadedAt     time.Time `json:"uploaded_at"`
	DownloadCount  int       `json:"download_count"`
}

func NewFile(id string, s3Key string, name string, size int64, expiry time.Time, userId string, downloadLink string, uploadTime time.Time, downloadCount int) *File {
	return &File{
		ID:             id,
		S3Key:          s3Key,
		Name:           name,
		Size:           size,
		ExpirationDate: expiry,
		UserId:         userId,
		DownloadLink:   downloadLink,
		UploadedAt:     uploadTime,
		DownloadCount:  downloadCount,
	}
}
