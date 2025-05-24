package handlers

import (
	"fileTransfer/internal/repository"
	"fileTransfer/internal/utils"
)

type Handlers struct {
	UserDbRepo repository.UserDbRepo
	FileDbRepo repository.FileDbRepo
	JWT        *utils.JWTService
	AwsS3      *utils.AwsS3
}

func NewHandlers(mysqlUserRepo repository.UserDbRepo, FileDbRepo repository.FileDbRepo, jwt *utils.JWTService, awsS3 *utils.AwsS3) *Handlers {
	return &Handlers{
		UserDbRepo: mysqlUserRepo,
		FileDbRepo: FileDbRepo,
		JWT:        jwt,
		AwsS3:      awsS3,
	}
}
