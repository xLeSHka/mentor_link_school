package group

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
	"gorm.io/gorm"
)

type GroupRepository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) repository.GroupRepository {
	return &GroupRepository{
		DB: db,
	}
}
