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

func (r *StudentService) GetRequestByID(ctx context.Context, reqID, groupID uuid.UUID) (*models.HelpRequest, error) {
	req, err := r.studentRepository.GetRequestByID(ctx, reqID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.HelpRequest{}, httpError.New(http.StatusNotFound, "That request does not exists in this group")
		}
		return &models.HelpRequest{}, httpError.New(http.StatusInternalServerError, err.Error())
	}
	return req, err
}
