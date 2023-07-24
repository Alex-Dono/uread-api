package model

import "github.com/google/uuid"

type Book struct {
	ID       uuid.UUID `gorm:"not null;type:uuid;primary_key"`
	Title    string    `gorm:"not null;type:varchar(50)"`
	Author   string    `gorm:"not null;type:varchar(50)"`
	FilePath string    `gorm:"not null;type:varchar(50)"`
}
