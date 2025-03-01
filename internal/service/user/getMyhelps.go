package userService

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *UsersService) GetMyHelps(ctx context.Context, userID uuid.UUID) ([]*models.HelpRequest, error) {
	return s.usersRepository.GetMyRequests(ctx, userID)
}
