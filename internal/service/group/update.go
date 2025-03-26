package groupService

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
)

func (s *GroupsService) UpdateInviteCode(ctx context.Context, groupID uuid.UUID) (string, error) {
	inviteCode, _ := GenerateInviteCode(5)

	err := s.groupRepository.UpdateInviteCode(ctx, groupID, inviteCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", httpError.New(http.StatusNotFound, "Организация не найдена!")
		}
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	return inviteCode, nil
}

func GenerateInviteCode(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	encoded := base64.URLEncoding.EncodeToString(bytes)
	inviteCode := strings.TrimRight(encoded, "=")
	return strings.ToLower(inviteCode[:length]), nil
}
