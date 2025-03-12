package handlers

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
	"tahap2/internal/domain"
	"tahap2/internal/middlewares"
	"time"
)

var (
	refreshTokenSecret = []byte("refreshsecret") // testing purpose
	accessTokenTTL     = time.Minute * 15
	refreshTokenTTL    = time.Hour * 24 * 7
)

type AuthHandler struct {
	authService domain.UserService
}

func NewAuthHandler(authService domain.UserService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterParam
	if err := c.Bind(&req); err != nil {
		return err
	}

	err := req.validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid body request"})
	}
	newUser, err := h.authService.Register(context.Background(), domain.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Pin:         req.Pin,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"status": "success",
		"result": toUserResponse(newUser),
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginParam
	if err := c.Bind(&req); err != nil {
		return err
	}
	userID, err := h.authService.Login(context.Background(), req.PhoneNumber, req.PIN)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": err.Error()})
	}

	accessToken, err := generateToken(userID, middlewares.AccessTokenSecret, accessTokenTTL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to generate token"})
	}
	refreshToken, err := generateToken(userID, refreshTokenSecret, refreshTokenTTL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to generate token"})
	}

	result := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"result": result,
	})
}

func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	userID := c.Get(middlewares.UserIDKey).(uuid.UUID)
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Address   string `json:"address"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}
	user, err := h.authService.UpdateProfile(context.Background(), userID, req.FirstName, req.LastName, req.Address)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": err.Error()})
	}

	result := UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Address:   user.Address,
		UpdatedAt: user.UpdatedAt.Format(time.DateTime),
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"result": result,
	})
}

// Generate JWT token
func generateToken(userID string, secret []byte, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (r RegisterParam) validate() error {
	// Check required fields
	if r.FirstName == "" {
		return errors.New("first_name is required")
	}
	if r.LastName == "" {
		return errors.New("last_name is required")
	}
	if r.PhoneNumber == "" {
		return errors.New("phone_number is required")
	}
	if r.Address == "" {
		return errors.New("address is required")
	}
	if r.Pin == "" {
		return errors.New("pin is required")
	}

	// Validate phone number length (10-15 characters)
	if len(r.PhoneNumber) < 10 || len(r.PhoneNumber) > 15 {
		return errors.New("phone_number must be between 10 and 15 characters")
	}

	// Validate address length (min 5 characters)
	if len(r.Address) < 5 {
		return errors.New("address must be at least 5 characters")
	}

	// Validate PIN (exactly 6 digits)
	pinRegex := regexp.MustCompile(`^\d{6}$`)
	if !pinRegex.MatchString(r.Pin) {
		return errors.New("pin must be exactly 6 numeric digits")
	}

	return nil
}

func toUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
		CreatedAt:   user.CreatedAt.Format(time.DateTime),
	}
}

type LoginParam struct {
	PhoneNumber string `json:"phone_number"`
	PIN         string `json:"pin"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterParam struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Pin         string `json:"pin"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Address     string    `json:"address"`
	CreatedAt   string    `json:"created_at,omitempty"`
	UpdatedAt   string    `json:"updated_at,omitempty"`
}
