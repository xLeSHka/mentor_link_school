package group

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"
	"gorm.io/gorm"
	"net/http"

	"github.com/xLeSHka/mentorLinkSchool/internal/app/httpError"
)

func (r *GroupRepository) RemoveRole(ctx context.Context, role *models.Role) error {
	tx := r.DB.Begin()
	var c int64
	err := tx.Model(&models.Role{}).Where("user_id = ? AND group_id = ?", role.UserID, role.GroupID).Count(&c).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if c == 1 {
		var role models.Role
		err := tx.Model(&models.Role{}).Where("user_id = ? AND group_id = ? AND role = ?", role.UserID, role.GroupID, role.Role).First(&role).Error
		if err == nil {
			return gorm.ErrInvalidTransaction
		}
	}
	if role.Role == "mentor" {
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
	res := tx.WithContext(ctx).Delete(role, &models.Role{UserID: role.UserID, GroupID: role.GroupID, Role: role.Role})
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
