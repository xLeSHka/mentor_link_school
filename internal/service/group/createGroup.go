package groupService

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (s *GroupsService) Create(ctx context.Context, group *models.Group, userID uuid.UUID) (string, error) {
	inviteCode, _ := generateInviteCode(5)
	group.InviteCode = &inviteCode
	err := s.groupRepository.Create(ctx, group, userID)
	if err != nil {
		return "", httpError.New(http.StatusInternalServerError, err.Error())
	}
	return inviteCode, nil
}
