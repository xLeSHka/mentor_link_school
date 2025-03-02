package publicRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
)

func (r *Route) mocks(c *gin.Context) {
	owner := &models.User{
		ID:        uuid.New(),
		Name:      "owner 1",
		AvatarURL: nil,
		BIO:       nil,
		Telegram:  "@owner",
	}
	ownerGroup := &models.Group{
		ID:        uuid.New(),
		AvatarURL: nil,
		Name:      "owner group 1",
	}
	student1 := &models.User{
		ID:        uuid.New(),
		Name:      "student 1",
		AvatarURL: nil,
		BIO:       nil,
		Telegram:  "@student",
	}
	student2 := &models.User{
		ID:        uuid.New(),
		Name:      "student 1",
		AvatarURL: nil,
		BIO:       nil,
		Telegram:  "@student",
	}
	group1 := &models.Group{
		ID:        uuid.New(),
		AvatarURL: nil,
		Name:      "group 1",
	}
	group2 := &models.Group{
		ID:        uuid.New(),
		AvatarURL: nil,
		Name:      "group 2",
	}
	bio := "new bio"
	mentor1 := &models.User{
		ID:        uuid.New(),
		Name:      "mentor 1",
		AvatarURL: nil,
		BIO:       &bio,
		Telegram:  "@mentor",
	}
	helpReq1 := &models.HelpRequest{
		ID:       uuid.New(),
		UserID:   student2.ID,
		MentorID: mentor1.ID,
		GroupID:  group2.ID,
		Goal:     "PROD Project",
		BIO:      student2.AvatarURL,
		Status:   "pending",
	}
	helpReq2 := &models.HelpRequest{
		ID:       uuid.New(),
		UserID:   student1.ID,
		MentorID: mentor1.ID,
		GroupID:  group1.ID,
		Goal:     "PRODANO Project",
		BIO:      student1.AvatarURL,
		Status:   "accepted",
	}
	r.DB.Create(group1)
	r.DB.Create(student1)
	r.DB.Create(student2)
	r.DB.Create(&models.Role{
		UserID:  student1.ID,
		GroupID: group1.ID,
		Role:    "student",
	})
	r.DB.Create(&models.Role{
		UserID:  student2.ID,
		GroupID: group1.ID,
		Role:    "student",
	})
	r.DB.Create(mentor1)
	r.DB.Create(&models.Role{
		UserID:  mentor1.ID,
		GroupID: group1.ID,
		Role:    "mentor",
	})
	r.DB.Create(&models.Role{
		UserID:  student2.ID,
		GroupID: group2.ID,
		Role:    "student",
	})
	r.DB.Create(&models.Role{
		UserID:  student1.ID,
		GroupID: group2.ID,
		Role:    "student",
	})
	r.DB.Create(&models.Role{
		UserID:  mentor1.ID,
		GroupID: group2.ID,
		Role:    "mentor",
	})
	r.DB.Create(helpReq1)
	r.DB.Create(helpReq2)
	r.DB.Create(&models.Pair{
		UserID:   student1.ID,
		MentorID: mentor1.ID,
		GroupID:  group1.ID,
		Goal:     "PRODANO Project",
	})

	r.DB.Create(owner)
	r.DB.Create(ownerGroup)
	r.DB.Create(&models.Role{
		UserID:  owner.ID,
		GroupID: ownerGroup.ID,
		Role:    "owner",
	})
	c.JSON(200, gin.H{
		"message": "PROOOOOOOOOOOOOOOOOD",
	})
}
