package repositoryUser

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository"

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
