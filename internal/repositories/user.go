package repositories

import (
	"context"
	"go_test_task/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Exists(ctx context.Context, id uint64) (bool, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return userRepository{DB: db}
}

func (r userRepository) Exists(ctx context.Context, id uint64) (bool, error) {
	var amount int64
	return amount > 0, r.DB.
		Model(&models.User{}).
		Where("id = ?", id).
		Count(&amount).
		Error
}
