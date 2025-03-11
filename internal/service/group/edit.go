package groupService

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"net/http"
)

func (s *GroupsService) Edit(ctx context.Context, userID, groupID uuid.UUID, updates map[string]any) (*models.Group, error) {
	exists, err := s.groupRepository.CheckGroupExists(ctx, userID, groupID)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exists {
		return nil, httpError.New(http.StatusForbidden, "group does not exist")
	}
	group, err := s.groupRepository.Edit(ctx, groupID, updates)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return group, nil
}
