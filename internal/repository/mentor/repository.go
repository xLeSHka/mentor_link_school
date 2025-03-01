package repositoryMentor

import (
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository"

	"gorm.io/gorm"
)

type MentorRepository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) repository.MentorRepository {
	return &MentorRepository{DB: db}
}
