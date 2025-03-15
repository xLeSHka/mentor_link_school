package avatar

import (
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"github.com/xLeSHka/mentorLinkSchool/internal/repository"
	"net/http"
	"strings"
)

func GetUserAvatar(user *models.User, minioRepository repository.MinioRepository) error {
	if user.AvatarURL != nil {
		avatarURL, err := minioRepository.GetImage(*user.AvatarURL)
		if err != nil {
			return httpError.New(http.StatusInternalServerError, err.Error())
		}
		avatarURL = strings.Split(avatarURL, "?X-Amz-Algorithm=AWS4-HMAC-SHA256")[0] + "?X-Amz-Algorithm=AWS4-HMAC-SHA256"
		user.AvatarURL = &avatarURL
	}
	return nil
}
func GetGroupAvatar(group *models.Group, minioRepository repository.MinioRepository) error {
	if group.AvatarURL != nil {
		avatarURL, err := minioRepository.GetImage(*group.AvatarURL)
		if err != nil {
			return httpError.New(http.StatusInternalServerError, err.Error())
		}
		avatarURL = strings.Split(avatarURL, "?X-Amz-Algorithm=AWS4-HMAC-SHA256")[0] + "?X-Amz-Algorithm=AWS4-HMAC-SHA256"
		group.AvatarURL = &avatarURL
	}
	return nil
}
