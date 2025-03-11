package group

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"net/http"

	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
)

func (r *GroupRepository) UpdateRole(ctx context.Context, groupID, userID uuid.UUID, role string) error {
	tx := r.DB.Begin()
	if role == "student" {
		err := tx.Where("mentor_id = ? AND group_id = ?", userID, groupID).Delete(&models.HelpRequest{}).Error

		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Where("mentor_id = ? AND group_id = ? ", userID, groupID).Delete(&models.Pair{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	res := tx.WithContext(ctx).Table("roles").Where("user_id = ? AND group_id = ?", userID, groupID).Update("role", role)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	if res.RowsAffected == 0 {
		tx.Rollback()
		return httpError.New(http.StatusBadRequest, "user not found")
	}
	tx.Commit()
	return nil
}
