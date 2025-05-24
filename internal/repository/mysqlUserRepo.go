package repository

import (
	"database/sql"
	"errors"
	"fileTransfer/internal/models"
	"fmt"
	"github.com/google/uuid"
)

type MysqlUserRepo struct {
	db *sql.DB
}

func (m *MysqlUserRepo) FindOrCreateUser(user *models.GoogleUser) (*models.GoogleUser, error) {
	// 1. Check if user exists
	q1 := `SELECT Id, Email, Name, Avatar, IsEmailVerified  FROM user WHERE Email = ?`
	row := m.db.QueryRow(q1, user.Email)

	var existingUser models.GoogleUser
	err := row.Scan(&existingUser.ID, &existingUser.Email, &existingUser.Name, &existingUser.Avatar, &existingUser.IsEmailVerified)
	if err == nil {
		// User found
		return &existingUser, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		// Some DB error
		return nil, err
	}

	// 2. User not found, insert new user
	userInsert := `
		INSERT INTO user (Id, Email, Name, Avatar, IsEmailVerified)
		VALUES (?, ?, ?, ?, ?)
	`
	id := uuid.New().String()
	_, err = m.db.Exec(userInsert, id, user.Email, user.Name, user.Avatar, user.IsEmailVerified)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	user.ID = id
	return user, nil
}

func NewMysqlUserRepo(db *sql.DB) UserDbRepo {
	return &MysqlUserRepo{db: db}
}
