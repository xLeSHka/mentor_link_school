package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" `
	FirstName  string    `gorm:"not null"`
	SecondName string    `gorm:"not null"`
	Email      string    `gorm:"not null"`
	AvatarURL  *string
	Password   []byte `gorm:"not null"`
	BIO        *string
}

func (_ *User) TableName() string {
	return "users"
}

type Group struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" `
	Name      string    `gorm:"not null"`
	AvatarURL *string
}

func (_ *Group) TableName() string {
	return "groups"
}

type GetMentorRequest struct {
	UserID  uuid.UUID `gorm:"type:uuid;not null"`
	GroupID uuid.UUID `gorm:"type:uuid;not null"`
	Goal    string    `gorm:"not null"`
	BIO     *string
}

func (_ *GetMentorRequest) TableName() string {
	return "get_mentor_requests"
}
func (_ *GetMentorRequest) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{DoNothing: true})
	return nil
}

type CreateMentorRequest struct {
	UserID  uuid.UUID `gorm:"type:uuid;not null"`
	GroupID uuid.UUID `gorm:"type:uuid;not null"`
	Goal    string    `gorm:"not null"`
	BIO     *string
}

func (_ *CreateMentorRequest) TableName() string {
	return "create_mentor_requests"
}
func (_ *CreateMentorRequest) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{DoNothing: true})
	return nil
}

type Role struct {
	UserID  uuid.UUID `gorm:"type:uuid;not null"`
	GroupID uuid.UUID `gorm:"type:uuid;not null"`
	Role    string    `gorm:"not null"`
}

func (_ *Role) TableName() string {
	return "roles"
}
func (_ *Role) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{DoNothing: true})
	return nil
}

type Mentor struct {
	UserID  uuid.UUID `gorm:"type:uuid;not null"`
	GroupID uuid.UUID `gorm:"type:uuid;not null"`
}

func (_ *Mentor) TableName() string {
	return "mentors"
}
func (_ *Mentor) BeforeCreate(tx *gorm.DB) (err error) {
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
