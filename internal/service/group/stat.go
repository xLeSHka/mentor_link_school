package groupService

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
)

func (s *GroupsService) GetStat(ctx context.Context, groupID uuid.UUID) (*models.GroupStat, error) {

	stat, err := s.groupRepository.GetStat(ctx, groupID)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return stat, nil
}
