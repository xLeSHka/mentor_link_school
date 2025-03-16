package repositoryStudent

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"

	"gorm.io/gorm"
)

type StudentRepository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) repository.StudentRepository {
	return &StudentRepository{
		DB: db,
	}
}
