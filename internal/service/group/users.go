package groupService

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
)

func (r *GroupsService) UpdateToMentor(ctx context.Context, ownerID, groupID, userID uuid.UUID) error {
	exists, err := r.groupRepository.CheckGroupExists(ctx, ownerID, groupID)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exists {
		httpError.New(http.StatusNotFound, "group does not exist")
	}
	err = r.groupRepository.UpdateToMentor(ctx, groupID, userID)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}
