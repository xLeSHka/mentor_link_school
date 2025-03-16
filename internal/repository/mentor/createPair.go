package repositoryMentor

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *MentorRepository) CreatePair(ctx context.Context, pair *models.Pair) error {
	err := r.DB.Create(pair).WithContext(ctx).Error
	return err
}
