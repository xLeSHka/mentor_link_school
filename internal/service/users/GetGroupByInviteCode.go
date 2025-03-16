package usersService

import (
	"context"
	"errors"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (s *UserService) GetGroupByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error) {
	group, err := s.usersRepository.GetGroupByInviteCode(ctx, inviteCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpError.New(http.StatusNotFound, err.Error())
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return group, nil
}
