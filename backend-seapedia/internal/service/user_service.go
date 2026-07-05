package service

import (
	"context"
	"errors"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/repository"

	"github.com/google/uuid"
)

type UserService interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
