package userService

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"net/http"
)

func (s *UsersService) UploadImage(ctx context.Context, file *models.File, personID uuid.UUID) (string, *httpError.HTTPError) {
	exist, err := s.usersRepository.CheckExists(ctx, personID)
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exist {
		return "", httpError.New(http.StatusNotFound, "User Not Found")
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
