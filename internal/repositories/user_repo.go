package repositories

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"tahap2/internal/domain"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepo) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepo) GetUserByID(userID uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.DB.Where("id = ?", userID).First(&user).Error
	return &user, err
}

func (r *UserRepo) GetUserByPhoneNumber(phoneNumber string) (*domain.User, error) {
	var user domain.User
	err := r.DB.Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user *domain.User) error {
	return r.DB.Save(user).Error
}
