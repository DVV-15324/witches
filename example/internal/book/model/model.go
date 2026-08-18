package model

import (
	"time"
)

type Book struct {
	ID        int      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (*Book) TableName() string {
	return "book"
}


