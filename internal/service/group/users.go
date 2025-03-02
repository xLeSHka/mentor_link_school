package groupService

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
)

func (r *GroupsService) UpdateRole(ctx context.Context, ownerID, groupID, userID uuid.UUID, role string) error {
	exists, err := r.groupRepository.CheckGroupExists(ctx, ownerID, groupID)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exists {
		httpError.New(http.StatusNotFound, "group does not exist")
	}
	err = r.groupRepository.UpdateRole(ctx, groupID, userID, role)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}
