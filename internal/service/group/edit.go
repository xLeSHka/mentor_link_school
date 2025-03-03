package groupService

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"net/http"
)

func (s *GroupsService) Edit(ctx context.Context, userID, groupID uuid.UUID, updates map[string]any) error {
	exists, err := s.groupRepository.CheckGroupExists(ctx, userID, groupID)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exists {
		return httpError.New(http.StatusForbidden, "group does not exist")
	}
	return nil
}
