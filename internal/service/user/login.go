package userService

import (
	"context"
	"errors"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/helpers"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func (s *UsersService) Login(ctx context.Context, email string, password string) (*models.User, string, error) {
	user, err := s.usersRepository.GetByEMail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", httpError.New(http.StatusUnauthorized, "Неверный email или пароль.")
		}
		return nil, "", httpError.New(http.StatusInternalServerError, err.Error())
	}

	if helpers.Compare([]byte(password), user.Password, s.cryptoKey) != nil {
		return nil, "", httpError.New(http.StatusUnauthorized, "Неверный email или пароль.")
	}

	err = s.rdb.Set(context.Background(), "jwt:"+user.ID.String(), time.Now().UnixMicro(), 6*time.Hour+time.Minute).Err()
	if err != nil {
		return nil, "", httpError.New(http.StatusInternalServerError, err.Error())
	}

	token, err := s.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, "", httpError.New(http.StatusInternalServerError, err.Error())
	}

	return user, token, nil
}
