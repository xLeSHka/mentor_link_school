package groupService

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"net/http"
)

func (s *GroupsService) UploadImage(ctx context.Context, file *models.File, groupdID, personID uuid.UUID) (string, *httpError.HTTPError) {
	exist, err := s.userRepository.CheckExists(ctx, personID)
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exist {
		return "", httpError.New(http.StatusForbidden, "User Not Found")
	}

	url, err := s.minioRepository.UploadImage(file)
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	_, err = s.groupRepository.Edit(ctx, groupdID, map[string]any{"avatar_url": file.Filename})
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	return url, nil
}
