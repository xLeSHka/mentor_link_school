package groupService

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *GroupsService) Create(ctx context.Context, group *models.Group, userID uuid.UUID) error {

	err := s.groupRepository.Create(ctx, group, userID)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}
