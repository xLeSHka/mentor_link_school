package repositoryUser

import (
	"prodapp/internal/repository"

	"gorm.io/gorm"
)

type UsersRepository struct {
	DB *gorm.DB
}

func NewUsersRepository(db *gorm.DB) repository.UsersRepository {
	return &UsersRepository{
		DB: db,
	}
}
