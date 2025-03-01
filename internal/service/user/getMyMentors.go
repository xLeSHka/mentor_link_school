package userService

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *UsersService) GetMyMentors(ctx context.Context, userID uuid.UUID) ([]*models.Pair, error) {
	return s.usersRepository.GetMyMentors(ctx, userID)
}
