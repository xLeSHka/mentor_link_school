package groupService

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
)

func (s *GroupsService) UpdateInviteCode(ctx context.Context, groupID uuid.UUID, ownerID uuid.UUID) (string, error) {
	exists, err := s.groupRepository.CheckGroupExists(ctx, ownerID, groupID)
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exists {
		return "", httpError.New(http.StatusForbidden, "group does not exist")
	}
	inviteCode, _ := generateInviteCode(5)

	err = s.groupRepository.UpdateInviteCode(ctx, groupID, inviteCode)
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	return inviteCode, nil
}

func generateInviteCode(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	encoded := base64.URLEncoding.EncodeToString(bytes)
	inviteCode := strings.TrimRight(encoded, "=")
	return strings.ToLower(inviteCode[:length]), nil
}
