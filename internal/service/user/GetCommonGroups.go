package userService

import (
	"errors"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"gorm.io/gorm"
	"net/http"
)

func (s *UsersService) GetCommonGroups(userID, mentorID uuid.UUID) ([]uuid.UUID, error) {
	group, err := s.usersRepository.GetCommonGroups(userID, mentorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []uuid.UUID{}, nil
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return group, nil
}
