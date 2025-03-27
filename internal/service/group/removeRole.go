package groupService

import (
	"context"
	"errors"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
)

func (s *GroupsService) RemoveRole(ctx context.Context, role *models.Role) error {
	_, err := s.userRepository.GetByID(ctx, role.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpError.New(http.StatusNotFound, "user not found")
		}
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	err = s.groupRepository.RemoveRole(ctx, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpError.New(http.StatusNotFound, "User Not Found")
		}
		if errors.Is(err, gorm.ErrInvalidTransaction) {
			return httpError.New(http.StatusUnprocessableEntity, "Can not delete last users role")
		}
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	err = s.cacheRepository.RemoveRoles(ctx, []*models.Role{role})
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}
