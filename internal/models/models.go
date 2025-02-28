package models

import (
	"io"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Surname   string    `gorm:"not null" json:"surname"`
	Email     string    `gorm:"not null" json:"email"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Password  []byte    `gorm:"not null" json:"-"`
	Age       int       `gorm:"not null" json:"age"`
}

func (u *User) TableName() string {
	return "users"
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
