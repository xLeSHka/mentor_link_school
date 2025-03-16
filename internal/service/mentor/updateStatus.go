package mentorService

import (
	"context"
	"errors"
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
)

func (s *MentorService) UpdateRequest(ctx context.Context, request *models.HelpRequest) error {
	err := s.mentorRepository.UpdateRequest(ctx, request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpError.New(http.StatusNotFound, err.Error())
		}
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	if request.Status == "accepted" {
		req, err := s.studentRepository.GetRequestByID(ctx, request.ID, request.GroupID)
		if err != nil {
			return httpError.New(http.StatusInternalServerError, err.Error())
		}
		err = s.mentorRepository.CreatePair(ctx, &models.Pair{
			UserID:   req.UserID,
			MentorID: req.MentorID,
			GroupID:  req.GroupID,
			Goal:     req.Goal,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return httpError.New(http.StatusConflict, "You already have this student")
			}
			return httpError.New(http.StatusInternalServerError, err.Error())
		}
	}
	return nil
}
