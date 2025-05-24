package repository

import (
	"fileTransfer/internal/models"
	"time"
)

type FileDbRepo interface {
	AddFile(file *models.File) error
	GetExpiredFiles(time time.Time) ([]models.File, error)
	DeleteFileByID(id string) error
	IncreaseDownloadCount(key string) error
}
