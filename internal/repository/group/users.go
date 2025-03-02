package group

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
)

func (r *GroupRepository) UpdateToMentor(ctx context.Context, groupID, userID uuid.UUID) error {
	res := r.DB.WithContext(ctx).Table("roles").Where("user_id = ? AND group_id = ?", userID, groupID).Update("role", "mentor")
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return httpError.New(http.StatusBadRequest, "user not found")
	}
	return nil
}
