package userService

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"net/http"
)

func (s *UsersService) Login(ctx context.Context, name string) (*models.User, string, error) {
	bio := "i want sleep"
	user, err := s.usersRepository.Login(ctx, &models.User{
		Name:     name,
		ID:       uuid.New(),
		Telegram: "t.me/meow",
		BIO:      &bio,
	})

	token, err := s.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, "", httpError.New(http.StatusInternalServerError, err.Error())
	}

	return user, token, nil
}
