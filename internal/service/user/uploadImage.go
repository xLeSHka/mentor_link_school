package userService

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (s *UsersService) UploadImage(ctx context.Context, file *models.File, personID uuid.UUID) (string, *httpError.HTTPError) {
	_, err := s.usersRepository.GetByID(ctx, personID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", httpError.New(http.StatusUnauthorized, "пользователь с таким id не найден")
		}
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	url, err := s.minioRepository.UploadImage(file)
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	_, err = s.usersRepository.EditUser(ctx, personID, map[string]any{"avatar_url": file.Filename})
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	return url, nil
}
