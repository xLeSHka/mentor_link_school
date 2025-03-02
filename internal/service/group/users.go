package groupService

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gorm.io/gorm"
)

func (r *GroupsService) UpdateToMentor(ctx context.Context, owserID, groupID, userID uuid.UUID) error {
	_, err := r.groupRepository.GetGroup(ctx, owserID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpError.New(http.StatusBadRequest, "group not found")
		}
		return httpError.New(http.StatusInternalServerError, err.Error())
	}

	err = r.groupRepository.UpdateToMentor(ctx, groupID, userID)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}
