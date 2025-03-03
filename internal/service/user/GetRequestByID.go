package userService

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *UsersService) GetRequestByID(ctx context.Context, reqID uuid.UUID) (models.HelpRequest, error) {
	return r.usersRepository.GetRequestByID(ctx, reqID)
}
