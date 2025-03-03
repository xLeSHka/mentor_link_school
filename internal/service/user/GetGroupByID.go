package userService

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *UsersService) GetGroupByID(ctx context.Context, ID uuid.UUID) (*models.Group, error) {
	return s.usersRepository.GetGroupByID(ctx, ID)
}
