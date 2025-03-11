package repositoryMentor

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *MentorRepository) UpdateRequest(ctx context.Context, request *models.HelpRequest) error {
	err := r.DB.Model(&models.HelpRequest{}).Where("id = ?", request.ID).Update("status", request.Status).Error
	return err
}
func (r *MentorRepository) CreatePair(ctx context.Context, pair *models.Pair) error {
	err := r.DB.Create(pair).Error
	return err
}
