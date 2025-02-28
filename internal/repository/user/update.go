package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"prodapp/internal/models"

	"gorm.io/gorm/clause"
)

func (r *UsersRepository) EditUser(ctx context.Context, userID uuid.UUID, updates map[string]any) (*models.User, error) {
	user := models.User{}
	err := r.DB.Model(&user).WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", userID).Updates(updates).Error
	return &user, err
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
