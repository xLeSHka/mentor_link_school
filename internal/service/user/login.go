package userService

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *UsersService) Login(ctx context.Context, name string) (*models.User, string, error) {
	if user, err := s.usersRepository.GetByName(ctx, name); err == nil {
		token, err := s.GenerateAccessToken(user.ID)
		if err != nil {
			return nil, "", httpError.New(http.StatusInternalServerError, err.Error())
		}
		return user, token, nil
	}
	bio := "i want sleep"
	user, err := s.usersRepository.Login(ctx, &models.User{
		Name:     name,
		ID:       uuid.New(),
		Telegram: "t_prodano",
		BIO:      &bio,
	})

	token, err := s.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, "", httpError.New(http.StatusInternalServerError, err.Error())
	}

	return user, token, nil
}
