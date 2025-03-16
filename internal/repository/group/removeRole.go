package group

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
)

func (r *GroupRepository) RemoveRole(ctx context.Context, role *models.Role) error {
	tx := r.DB.Begin()
	if role.Role == "student" {
		err := tx.Where("mentor_id = ? AND group_id = ?", role.UserID, role.GroupID).Delete(&models.HelpRequest{}).Error

		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Where("mentor_id = ? AND group_id = ? ", role.UserID, role.GroupID).Delete(&models.Pair{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	res := tx.WithContext(ctx).Table("roles").Where("user_id = ? AND group_id = ?", role.UserID, role.GroupID).Delete("role", role.Role)
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
