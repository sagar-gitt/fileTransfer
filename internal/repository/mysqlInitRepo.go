package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type MySQLInitRepo struct {
	db *sql.DB
}

func (m *MySQLInitRepo) CreateUserTableIfNotExist() error {
	query := `CREATE TABLE IF NOT EXISTS user (
    	Id VARCHAR(255) PRIMARY KEY,
		Name VARCHAR(255),
		Email VARCHAR(255) UNIQUE NOT NULL,
		Avatar VARCHAR(255),
		IsEmailVerified BOOLEAN DEFAULT FALSE,
    	AuthProvider VARCHAR(255) DEFAULT 'google'
	)`

	_, err := m.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (m *MySQLInitRepo) CreateFileTableIfNotExist() error {
	query := `CREATE TABLE IF NOT EXISTS file (
    	Id VARCHAR(255) PRIMARY KEY,
    	S3Key VARCHAR(512) UNIQUE NOT NULL,
		Name VARCHAR(255) NOT NULL,
		Size INT NOT NULL,
		ExpirationDate DATETIME,
		UserId VARCHAR(255),
    	DownloadLink VARCHAR(255) UNIQUE NOT NULL,
    	UploadedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    	DownloadCount INT DEFAULT 0
	)` //--FOREIGN KEY (UserId) REFERENCES user(Id) ON DELETE SET NULL can also use CASCADE or RESTRICT

	_, err := m.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (m *MySQLInitRepo) InsertSampleData() error {
	// Sample insert into user
	userInsert := `
		INSERT IGNORE INTO user (Id, Name, Email, Avatar)
		VALUES (?, ?, ?, ?)
	`
	id := uuid.New().String()
	_, err := m.db.Exec(userInsert, id, "Test User", "testuser@example.com", "http://testImage")
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// Sample insert into files
	fileInsert := `
		INSERT IGNORE INTO file (Id, S3Key, Name, Size, UserId, DownloadLink)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	id2 := uuid.New().String()
	_, err = m.db.Exec(fileInsert, id2, "SomeKey", "sample_file.txt", 1048576, id, "https://SampleDownloadlink.com")
	if err != nil {
		return fmt.Errorf("failed to insert file: %w", err)
	}

	return nil
}

func (m *MySQLInitRepo) TruncateAllTables() error {
	// Disable foreign key checks temporarily as even after deleting rows of files table, it still holds relation with user table
	//_, err := m.db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	//if err != nil {
	//	return fmt.Errorf("failed to disable FK checks: %w", err)
	//}

	// Truncate the tables - order here matters - Always truncate child tables before parent tables in the foreign key hierarchy.
	tables := []string{"file", "user"}
	for _, table := range tables {
		if _, err := m.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)); err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	// Re-enable foreign key checks
	//_, err = m.db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	//if err != nil {
	//	return fmt.Errorf("failed to re-enable FK checks: %w", err)
	//}

	return nil
}

func NewMySQLInitRepo(db *sql.DB) InitDbRepo {
	return &MySQLInitRepo{
		db: db,
	}
}
