package userService

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
)

func (s *UsersService) GetMentors(ctx context.Context, userID uuid.UUID) ([]*models.Role, error) {
	mentors, err := s.usersRepository.GetMentors(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*models.Role{}, nil
		}
		return nil, err
	}
	return mentors, nil
}
