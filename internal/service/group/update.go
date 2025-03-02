package groupService

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gorm.io/gorm"
)

func (r *GroupsService) UpdateInviteCode(ctx context.Context, groupID uuid.UUID, ownerID uuid.UUID) (string, error) {
	_, err := r.groupRepository.GetGroup(ctx, ownerID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", httpError.New(http.StatusBadRequest, "group not found")
		}
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}

	inviteCode, _ := generateInviteCode(10)

	err = r.groupRepository.UpdateInviteCode(ctx, groupID, inviteCode)
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
	return inviteCode[:length], nil
}
