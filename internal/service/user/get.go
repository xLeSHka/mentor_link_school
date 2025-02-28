package userService

import (
	"context"
	"prodapp/internal/models"

	"github.com/google/uuid"
)

func (s *UsersService) GetByID(ctx context.Context, id uuid.UUID) (person *models.User, err error) {
	return s.usersRepository.GetByID(ctx, id)
}
