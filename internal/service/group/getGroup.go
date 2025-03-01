package mentorService

import (
	"context"
	"errors"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (r *GroupsService) GetGroup(ctx context.Context, group *models.Group) (*models.Group, error) {
	gr, err := r.groupRepository.GetGroup(ctx, group)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, httpError.New(http.StatusNotFound, "group not found")
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return gr, nil
}
