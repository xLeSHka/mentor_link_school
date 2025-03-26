package studentService

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (r *StudentService) GetRequest(ctx context.Context, studentID, mentorID, groupID uuid.UUID) (*models.HelpRequest, error) {
	req, err := r.studentRepository.GetRequest(ctx, studentID, mentorID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return req, err
}
