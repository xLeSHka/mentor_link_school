package groupService

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
)

func (r *GroupsService) GetGroup(ctx context.Context, userID, groupID uuid.UUID) (*models.Group, error) {
	gr, err := r.groupRepository.GetGroup(ctx, userID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpError.New(http.StatusNotFound, "group not found")
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return gr, nil
}
