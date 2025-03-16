package repositoryStudent

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *StudentRepository) CreateRequest(ctx context.Context, request *models.HelpRequest) error {
	return r.DB.Create(request).Error
}
