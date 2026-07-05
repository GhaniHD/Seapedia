package service

import (
	"context"
	"errors"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/model"
	"backend-seapedia/internal/repository"
	"backend-seapedia/pkg/crypto"
	"backend-seapedia/pkg/jwt"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "buyer",
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
		Role:  newUser.Role,
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	if !crypto.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("email atau password salah")
	}

	token, err := jwt.GenerateToken(user.ID, user.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
