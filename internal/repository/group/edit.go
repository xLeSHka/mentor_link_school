package group

import (
	"context"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"

	"gorm.io/gorm/clause"
)

func (r *GroupRepository) EditUser(ctx context.Context, groupID uuid.UUID, updates map[string]any) (*models.Group, error) {
	group := models.Group{}
	err := r.DB.Model(&group).WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", groupID).Updates(updates).Error
	return &group, err
}

//
//func (r *UsersRepository) EditUser(ctx context.Context, user *models.User) (*models.User, error) {
//	editedUser := models.User{}
//	err := r.DB.Model(&editedUser).Clauses(clause.Returning{}).Where("id = ?", user.ID).Updates(user).Error
//	if err != nil {
//		log.Println("error editing user", err.Error())
//		return nil, err
//	}
//	return &editedUser, nil
//}
