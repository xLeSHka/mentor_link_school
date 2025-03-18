package groupService

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (s *GroupsService) UploadImage(ctx context.Context, file *models.File, groupID uuid.UUID) (string, *httpError.HTTPError) {

	url, err := s.minioRepository.UploadImage(file)
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	url = strings.Replace(url, "http://minio:9000", "https://localhost", 1)
	url = strings.Split(url, "?X-Amz-Algorithm=AWS4-HMAC-SHA256")[0] + "?X-Amz-Algorithm=AWS4-HMAC-SHA256"
	_, err = s.groupRepository.Edit(ctx, &models.Group{ID: groupID, AvatarURL: &file.Filename})
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	return url, nil
}
