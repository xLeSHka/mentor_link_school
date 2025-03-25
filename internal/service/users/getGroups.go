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

func (r *UserService) GetGroups(ctx context.Context, userID uuid.UUID, page, size int) ([]*models.GroupWithRoles, int64, error) {
	gr, total, err := r.usersRepository.GetGroups(ctx, userID, page, size)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return []*models.GroupWithRoles{}, 0, nil
		}
		return nil, 0, httpError.New(http.StatusInternalServerError, err.Error())
	}
	log.Println(gr)
	return gr, total, nil
}
