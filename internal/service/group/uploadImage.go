package groupService

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *GroupsService) UploadImage(ctx context.Context, file *models.File, groupID, personID uuid.UUID) (string, *httpError.HTTPError) {
	exists, err := s.groupRepository.CheckGroupExists(ctx, personID, groupID)
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exists {
		return "", httpError.New(http.StatusForbidden, "group does not exist")
	}

	url, err := s.minioRepository.UploadImage(file)
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	url = strings.Split(url, "?X-Amz-Algorithm=AWS4-HMAC-SHA256")[0] + "?X-Amz-Algorithm=AWS4-HMAC-SHA256"
	_, err = s.groupRepository.Edit(ctx, groupID, map[string]any{"avatar_url": file.Filename})
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	return url, nil
}
