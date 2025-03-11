package userService

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"net/http"
)

func (s *UsersService) Edit(ctx context.Context, personID uuid.UUID, updates map[string]any) error {
	exist, err := s.usersRepository.CheckExists(ctx, personID)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !exist {
		return httpError.New(http.StatusNotFound, "User Not Found")
	}

	_, err = s.usersRepository.EditUser(ctx, personID, updates)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}
