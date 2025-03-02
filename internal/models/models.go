package models

import (
	"io"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" `
	Name      string    `gorm:"unique;not null"`
	AvatarURL *string
	BIO       *string
	Telegram  string `gorm:"not null"`
	Role      *Role  `gorm:"foreignKey:user_id"`
}

func (_ *User) TableName() string {
	return "users"
}

type Group struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	AvatarURL  *string
	Name       string `gorm:"not null"`
	InviteCode *string
}

func (_ *Group) TableName() string {
	return "groups"
}

type HelpRequest struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID   uuid.UUID `gorm:"type:uuid;not null"`
	MentorID uuid.UUID `gorm:"type:uuid;not null"`
	GroupID  uuid.UUID `gorm:"type:uuid;not null"`
	Goal     string    `gorm:"not null"`
	BIO      *string
	Status   string `gorm:"not null"`
	Mentor   *User  `gorm:"foreignKey:mentor_id"`
	Student  *User  `gorm:"foreignKey:user_id"`
}

func (_ *HelpRequest) TableName() string {
	return "help_requests"
}
func (_ *HelpRequest) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{DoNothing: true})
	return nil
}

type Role struct {
	UserID  uuid.UUID `gorm:"type:uuid;not null"`
	GroupID uuid.UUID `gorm:"type:uuid;not null"`
	Role    string    `gorm:"not null"`
	Mentor  *User     `gorm:"foreignKey:user_id"`
}

func (_ *Role) TableName() string {
	return "roles"
}
func (_ *Role) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{DoNothing: true})
	return nil
}

type Pair struct {
	UserID   uuid.UUID `gorm:"type:uuid;not null"`
	MentorID uuid.UUID `gorm:"type:uuid;not null"`
	GroupID  uuid.UUID `gorm:"type:uuid;not null"`
	Goal     string    `gorm:"not null"`
	Mentor   *User     `gorm:"foreignKey:mentor_id"`
	Student  *User     `gorm:"foreignKey:user_id"`
}

func (_ *Pair) TableName() string {
	return "pairs"
}
func (_ *Pair) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{DoNothing: true})
	return nil
}

type File struct {
	Filename string
	Size     int64
	File     io.Reader
	Mimetype string
}

func (u *File) TableName() string {
	return "file"
}
