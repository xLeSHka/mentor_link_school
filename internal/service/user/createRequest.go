package userService

import (
	"context"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"net/http"
)

func (s *UsersService) CreateRequest(ctx context.Context, request *models.HelpRequest) error {
	isStudent, err := s.usersRepository.CheckIsStudent(ctx, request.UserID, request.GroupID)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	if !isStudent {
		return httpError.New(http.StatusBadRequest, "user is not student")
	}
	isMentor, err := s.mentorRepository.CheckIsMentor(ctx, request.MentorID, request.GroupID)
	if err != nil {
		return err
	}
	if !isMentor {
		return httpError.New(http.StatusBadRequest, "user is not mentor")
	}
	err = s.usersRepository.CreateRequest(ctx, request)
	if err != nil {
		return err
	}
	return nil
}
