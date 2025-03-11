package userService

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
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
	// if alredy exists
	_, err = s.usersRepository.GetRequest(ctx, request.UserID, request.MentorID, request.GroupID)
	if err == nil {
		return httpError.New(http.StatusBadRequest, "request already exists")
	}
	err = s.usersRepository.CreateRequest(ctx, request)
	if err != nil {
		return err
	}
	return nil
}
