package groupService

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (s *GroupsService) GetStat(ctx context.Context, ownerID, groupID uuid.UUID) (*models.GroupStat, error) {
	exist, err := s.groupRepository.CheckGroupExists(ctx, ownerID, groupID)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exist {
		return nil, httpError.New(http.StatusForbidden, "Specified Group not found")
	}
	stat, err := s.groupRepository.GetStat(ctx, groupID)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return stat, nil
}
