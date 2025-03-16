package usersService

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
)

func (r *UserService) GetGroups(ctx context.Context, userID uuid.UUID) ([]*models.GroupWithRoles, error) {
	gr, err := r.usersRepository.GetGroups(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return []*models.GroupWithRoles{}, nil
		}
		return nil, httpError.New(http.StatusInternalServerError, err.Error())
	}
	log.Println(gr)
	return gr, nil
}
