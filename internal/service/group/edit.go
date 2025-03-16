package groupService

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"net/http"
)

func (s *GroupsService) Edit(ctx context.Context, group *models.Group) (*models.Group, error) {
	group, err := s.groupRepository.Edit(ctx, group)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return group, nil
}
