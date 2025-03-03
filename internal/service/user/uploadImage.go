package userService

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
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
	url = strings.Split(url, "?X-Amz-Algorithm=AWS4-HMAC-SHA256")[0] + "?X-Amz-Algorithm=AWS4-HMAC-SHA256"
	_, err = s.usersRepository.EditUser(ctx, personID, map[string]any{"avatar_url": file.Filename})
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	return url, nil
}
