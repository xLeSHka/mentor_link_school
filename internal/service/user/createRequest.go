package userService

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *UsersService) CreateRequest(ctx context.Context, request *models.HelpRequest) error {
	err := s.usersRepository.CreateRequest(ctx, request)
	if err != nil {
		return err
	}
	return nil
}
