package repositoryUser

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (r *UsersRepository) CreateRequest(ctx context.Context, request *models.HelpRequest) error {
	return r.DB.Create(request).Error
}
