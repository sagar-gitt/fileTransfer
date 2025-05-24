package repository

import "fileTransfer/internal/models"

type UserDbRepo interface {
	FindOrCreateUser(user *models.GoogleUser) (*models.GoogleUser, error)
}
