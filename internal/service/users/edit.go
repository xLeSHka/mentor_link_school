package usersService

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"net/http"
)

func (s *UserService) Edit(ctx context.Context, userID uuid.UUID, user *models.User) (*models.User, error) {

	updated, err := s.usersRepository.EditUser(ctx, userID, user)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return updated, nil
}
