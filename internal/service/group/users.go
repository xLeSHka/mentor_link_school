package groupService

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"net/http"

	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
)

func (s *GroupsService) UpdateRole(ctx context.Context, ownerID, groupID, userID uuid.UUID, role string) error {
	exists, err := s.groupRepository.CheckGroupExists(ctx, ownerID, groupID)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exists {
		return httpError.New(http.StatusForbidden, "specified group does not exist")
	}
	exist, err := s.userRepository.CheckExists(ctx, userID)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exist {
		return httpError.New(http.StatusNotFound, "User Not Found")
	}
	err = s.groupRepository.UpdateRole(ctx, groupID, userID, role)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpError.New(http.StatusNotFound, "User Not Found")
		}
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}
