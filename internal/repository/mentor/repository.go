package repositoryMentor

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"

	"gorm.io/gorm"
)

type MentorRepository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) repository.MentorRepository {
	return &MentorRepository{DB: db}
}
