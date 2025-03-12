package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FirstName   string    `gorm:"not null"`
	LastName    string    `gorm:"not null"`
	PhoneNumber string    `gorm:"unique;not null"`
	Address     string    `gorm:"not null"`
	Pin         string    `gorm:"not null"`
	Balance     int64     `gorm:"default:0;not null"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// UserRepository defines the methods for database operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByPhoneNumber(phoneNumber string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	GetUserByID(userID uuid.UUID) (*User, error)
}

// UserService defines the methods for business logic
type UserService interface {
	Register(ctx context.Context, user User) (User, error)
	Login(ctx context.Context, phoneNumber, pin string) (string, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, firstname, lastname, address string) (User, error)
}
