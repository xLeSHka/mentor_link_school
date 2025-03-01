package mentorService

import (
	"context"
	"errors"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (s *MentorService) GetMentorRequest(ctx context.Context, req *models.HelpRequest) error {
	user, err := s.usersRepository.GetByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpError.New(http.StatusNotFound, "Такой пользователь не найден")
		}
		return httpError.New(http.StatusInternalServerError, err.Error())

	}
	req.BIO = user.BIO
	err = s.mentorRepository.GetMentorRequest(ctx, req)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}
