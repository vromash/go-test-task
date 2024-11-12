package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
}

type Balance struct {
	gorm.Model
	UserID uint64
	Amount string
}
