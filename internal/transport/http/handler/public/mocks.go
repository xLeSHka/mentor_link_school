package publicRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *Route) mocks(c *gin.Context) {
	BIO := "*Информация о студенте*"
	swaggerStudent := &models.User{
		ID:       uuid.New(),
		Name:     "Свагер студент",
		BIO:      &BIO,
		Telegram: "t.me/t_prodano",
	}
	BIO2 := "*Информация о менторе*"
	swaggerMentor := &models.User{
		ID:       uuid.New(),
		Name:     "Свагер ментор",
		BIO:      &BIO2,
		Telegram: "t.me/t_prodano",
	}
	BIO3 := "*Информация о владельце*"
	swaggerOwner := &models.User{
		ID:       uuid.New(),
		Name:     "Свагер владелец",
		BIO:      &BIO3,
		Telegram: "t.me/t_prodano",
	}
	inviteCode := "gnjkf"
	swaggerGroup := &models.Group{
		ID:         uuid.New(),
		Name:       "Свагер огранизация",
		InviteCode: &inviteCode,
	}
	swaggerStudentRole := &models.Role{
		UserID:  swaggerStudent.ID,
		GroupID: swaggerGroup.ID,
		Role:    "student",
	}
	swaggerMentorRole := &models.Role{
		UserID:  swaggerMentor.ID,
		GroupID: swaggerGroup.ID,
		Role:    "mentor",
	}
	swaggerOwnerRole := &models.Role{
		UserID:  swaggerOwner.ID,
		GroupID: swaggerGroup.ID,
		Role:    "owner",
	}
	r.DB.Create(&swaggerStudent)
	r.DB.Create(&swaggerMentor)
	r.DB.Create(&swaggerOwner)
	r.DB.Create(&swaggerGroup)
	r.DB.Create(&swaggerStudentRole)
	r.DB.Create(&swaggerMentorRole)
	r.DB.Create(&swaggerOwnerRole)

	owBio := "Информация о владельце 1"
	owner := &models.User{
		ID:        uuid.New(),
		Name:      "Владелец 1",
		AvatarURL: nil,
		BIO:       &owBio,
		Telegram:  "t.me/t_prodano",
	}
	ownerGroup := &models.Group{
		ID:        uuid.New(),
		AvatarURL: nil,
		Name:      "Группа 1",
	}
	ow2Bio := "Информация о владельце 2"
	owner2 := &models.User{
		ID:        uuid.New(),
		Name:      "Владелец 2",
		AvatarURL: nil,
		BIO:       &ow2Bio,
		Telegram:  "t.me/t_prodano",
	}
	ownerGroup2 := &models.Group{
		ID:        uuid.New(),
		AvatarURL: nil,
		Name:      "Группа 2",
	}
	sBio := "Информация о студенте"
	student1 := &models.User{
		ID:        uuid.New(),
		Name:      "Студент 1",
		AvatarURL: nil,
		BIO:       &sBio,
		Telegram:  "t.me/t_prodano",
	}
	s2BIO := "Информауия о студенте 2"
	student2 := &models.User{
		ID:        uuid.New(),
		Name:      "Студент 2",
		AvatarURL: nil,
		BIO:       &s2BIO,
		Telegram:  "t.me/t_prodano",
	}

	bio := "new bio"
	mentor1 := &models.User{
		ID:        uuid.New(),
		Name:      "Ментор 1",
		AvatarURL: nil,
		BIO:       &bio,
		Telegram:  "t.me/t_prodano",
	}
	bio2 := "new bio"
	mentor2 := &models.User{
		ID:        uuid.New(),
		Name:      "Ментор 2",
		AvatarURL: nil,
		BIO:       &bio2,
		Telegram:  "t.me/t_prodano",
	}
	r.DB.Create(&owner)
	r.DB.Create(&owner2)
	r.DB.Create(&ownerGroup)
	r.DB.Create(&ownerGroup2)
	r.DB.Create(&student1)
	r.DB.Create(&student2)
	r.DB.Create(&mentor1)
	r.DB.Create(&mentor2)
	r.DB.Create(&models.Role{
		UserID:  owner.ID,
		GroupID: ownerGroup.ID,
		Role:    "owner",
	})
	r.DB.Create(&models.Role{
		UserID:  owner2.ID,
		GroupID: ownerGroup2.ID,
		Role:    "owner",
	})
	r.DB.Create(&models.Role{
		UserID:  mentor1.ID,
		GroupID: ownerGroup2.ID,
		Role:    "mentor",
	})
	r.DB.Create(&models.Role{
		UserID:  mentor2.ID,
		GroupID: ownerGroup2.ID,
		Role:    "mentor",
	})
	r.DB.Create(&models.Role{
		UserID:  mentor1.ID,
		GroupID: ownerGroup.ID,
		Role:    "mentor",
	})
	r.DB.Create(&models.Role{
		UserID:  student2.ID,
		GroupID: ownerGroup2.ID,
		Role:    "student",
	})
	r.DB.Create(&models.Role{
		UserID:  student1.ID,
		GroupID: ownerGroup2.ID,
		Role:    "student",
	})
	r.DB.Create(&models.Role{
		UserID:  student1.ID,
		GroupID: ownerGroup.ID,
		Role:    "student",
	})
	c.JSON(200, gin.H{
		"message": "PROOOOOOOOOOOOOOOOOD",
	})
}
