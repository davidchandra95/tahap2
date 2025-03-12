package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"tahap2/internal/domain"
	"time"
)

type AuthService struct {
	userRepo domain.UserRepository
}

func NewAuthService(userRepo domain.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, user domain.User) (domain.User, error) {
	existUser, err := s.userRepo.GetUserByPhoneNumber(user.PhoneNumber)
	if err != nil {
		return domain.User{}, err
	}
	if existUser.PhoneNumber != "" {
		return domain.User{}, errors.New("phone number already registered")
	}
	hashedPin, err := hashPin(user.Pin)
	if err != nil {
		return domain.User{}, errors.New("failed to hash pin")
	}

	user.Pin = hashedPin
	err = s.userRepo.CreateUser(ctx, &user)
	return user, err
}

func (s *AuthService) Login(ctx context.Context, phoneNumber, pin string) (string, error) {
	user, err := s.userRepo.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return "", errors.New("phone number not found")
	}
	if err = checkPin(user.Pin, pin); err != nil {
		return "", errors.New("Phone number and PIN doesn't match")
	}

	return user.ID.String(), nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, firstname, lastname, address string) (domain.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return domain.User{}, err
	}

	user.FirstName = firstname
	user.LastName = lastname
	user.Address = address
	user.UpdatedAt = time.Now()
	err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return *user, nil
}

type RegisterParam struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	Address     string
	Pin         string
}

type User struct {
	ID          uuid.UUID
	FirstName   string
	LastName    string
	PhoneNumber string
	Address     string
	Pin         string
	CreatedAt   string
}

func hashPin(pin string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	return string(hashed), err
}

func checkPin(hashedPin, pin string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPin), []byte(pin))
}
