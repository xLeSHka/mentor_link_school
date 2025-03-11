package repositoryUser

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"

	"gorm.io/gorm"
)

type UsersRepository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) repository.UsersRepository {
	return &UsersRepository{
		DB: db,
	}
}
