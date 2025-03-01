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

func (r *GroupsService) GetGroups(ctx context.Context, userID uuid.UUID) ([]*models.Group, error) {
	gr, err := r.groupRepository.GetGroups(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, httpError.New(http.StatusNotFound, "groups not found")
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return gr, nil
}
