package repositoryUser

//import (
//	"context"
//	"errors"
//	"github.com/google/uuid"
//	"github.com/xLeSHka/mentorLinkSchool/internal/models"
//	"gorm.io/gorm"
//)
//
//func (r *UsersRepository) CheckExists(ctx context.Context, id uuid.UUID) (bool, error) {
//	var user models.User
//	err := r.DB.Model(&models.User{}).Where("id = ?", id).First(&user).Error
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return false, nil
//		}
//		return false, err
//	}
//	return true, nil
//}
