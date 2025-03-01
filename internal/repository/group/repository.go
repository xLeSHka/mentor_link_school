package group

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository"
	"gorm.io/gorm"
)

type GroupRepository struct {
	DB *gorm.DB
}

func NewUsersRepository(db *gorm.DB) repository.GroupRepository {
	return &GroupRepository{
		DB: db,
	}
}
