package repository

import (
	"database/sql"
	"fileTransfer/internal/models"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type MysqlFileRepo struct {
	db *sql.DB
}

func (m *MysqlFileRepo) DeleteFileByID(id string) error {
	_, err := m.db.Exec("DELETE FROM file WHERE Id = ?", id)
	return err
}

func (m *MysqlFileRepo) IncreaseDownloadCount(key string) error {
	query := `UPDATE file SET DownloadCount = DownloadCount + 1 WHERE S3Key = ?`

	_, err := m.db.Exec(query, key)
	if err != nil {
		return fmt.Errorf("failed to increase download count: %w", err)
	}

	return nil
}

func (m *MysqlFileRepo) AddFile(file *models.File) error {
	id := uuid.New()
	q := `
		INSERT INTO file (Id, S3Key, Name, Size, ExpirationDate, UserId, DownloadLink, UploadedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := m.db.Exec(q, id, file.S3Key, file.Name, file.Size, file.ExpirationDate,
		file.UserId, file.DownloadLink, file.UploadedAt)
	if err != nil {
		return fmt.Errorf("failed to insert file info: %w", err)
	}

	return nil
}

func (m *MysqlFileRepo) GetExpiredFiles(time time.Time) ([]models.File, error) {
	rows, err := m.db.Query("SELECT Id, S3Key FROM file WHERE ExpirationDate IS NOT NULL AND ExpirationDate <= ?", time)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var f models.File
		rows.Scan(&f.ID, &f.S3Key)
		files = append(files, f)
	}
	return files, nil
}

func NewMysqlFileRepo(db *sql.DB) FileDbRepo {
	return &MysqlFileRepo{db: db}
}
