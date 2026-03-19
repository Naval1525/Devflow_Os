package service

import (
	"context"
	"devflowos/api/internal/model"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already registered")
)

type AuthService struct {
	userRepo   UserRepository
	jwtSecret  []byte
	bcryptCost int
}

type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

func NewAuthService(userRepo UserRepository, jwtSecret string, bcryptCost int) *AuthService {
	if bcryptCost <= 0 {
		bcryptCost = bcrypt.DefaultCost
	}
	return &AuthService{
		userRepo:   userRepo,
		jwtSecret:  []byte(jwtSecret),
		bcryptCost: bcryptCost,
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *AuthService) Signup(ctx context.Context, email, password string) (*model.User, string, error) {
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, "", ErrEmailExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		return nil, "", err
	}
	u, err := s.userRepo.Create(ctx, email, string(hash))
	if err != nil {
		return nil, "", err
	}
	token, err := s.issueToken(u)
	if err != nil {
		return u, "", err
	}
	return u, token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}
	return s.issueToken(u)
}

func (s *AuthService) issueToken(u *model.User) (string, error) {
	claims := Claims{
		UserID: u.ID.String(),
		Email:  u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (uuid.UUID, error) {
	t, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil || !t.Valid {
		return uuid.Nil, errors.New("invalid token")
	}
	claims, ok := t.Claims.(*Claims)
	if !ok {
		return uuid.Nil, errors.New("invalid claims")
	}
	return uuid.MustParse(claims.UserID), nil
}
