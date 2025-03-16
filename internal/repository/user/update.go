package repositoryUser

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/mentorLinkSchool/internal/models"

	"gorm.io/gorm/clause"
)

func (r *UsersRepository) EditUser(ctx context.Context, userID uuid.UUID, user *models.User) (*models.User, error) {
	usr := models.User{}
	err := r.DB.Model(&usr).WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", userID).Updates(&user).Error
	return &usr, err
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
