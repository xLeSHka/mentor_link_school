package userService

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *UsersService) GetMentors(ctx context.Context, userID uuid.UUID) ([]*models.User, error) {
	return nil, nil
}
