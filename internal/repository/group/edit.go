package group

import (
	"context"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"

	"gorm.io/gorm/clause"
)

func (r *GroupRepository) Edit(ctx context.Context, group *models.Group) (*models.Group, error) {
	g := models.Group{}
	err := r.DB.Model(&g).WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", group.ID).Updates(group).Error
	return &g, err
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
