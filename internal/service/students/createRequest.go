package studentService

import (
	"context"
	"errors"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (s *StudentService) CreateRequest(ctx context.Context, request *models.HelpRequest) error {
	isMentor, err := s.mentorRepository.CheckIsMentor(ctx, request.MentorID, request.GroupID)
	if err != nil {
		return err
	}
	if !isMentor {
		return httpError.New(http.StatusBadRequest, "user is not mentor")
	}

	err = s.studentRepository.CreateRequest(ctx, request)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return httpError.New(http.StatusConflict, "You already send req to this mentor with this goal")
		}
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}
