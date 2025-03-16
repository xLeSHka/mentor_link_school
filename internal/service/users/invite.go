package usersService

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
)

func (s *UserService) Invite(ctx context.Context, inviteCode string, userID uuid.UUID) (bool, error) {
	group, err := s.usersRepository.GetGroupByInviteCode(ctx, inviteCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, httpError.New(http.StatusNotFound, "Group not found")
		}
		return false, err
	}
	err = s.groupRepository.AddRole(ctx, &models.Role{
		UserID:  userID,
		GroupID: group.ID,
		Role:    "student",
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return false, httpError.New(http.StatusConflict, "User already invited")
		}
		return false, err
	}
	return true, nil
}
