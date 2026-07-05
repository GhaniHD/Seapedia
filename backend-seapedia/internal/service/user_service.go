package service

import (
	"context"
	"errors"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/repository"

	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(ctx context.Context, id uuid.UUID, activeRole string) (*dto.ProfileResponse, error)
}

type userService struct {
	userRepo   repository.UserRepository
	walletRepo repository.WalletRepository
	storeRepo  repository.StoreRepository
	orderRepo  repository.OrderRepository
	deliveryRepo repository.DeliveryRepository
}

func NewUserService(userRepo repository.UserRepository, walletRepo repository.WalletRepository,
	storeRepo repository.StoreRepository, orderRepo repository.OrderRepository, deliveryRepo repository.DeliveryRepository) UserService {
	return &userService{userRepo: userRepo, walletRepo: walletRepo, storeRepo: storeRepo, orderRepo: orderRepo, deliveryRepo: deliveryRepo}
}

func (s *userService) GetProfile(ctx context.Context, id uuid.UUID, activeRole string) (*dto.ProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	roles, err := s.userRepo.GetRoles(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := &dto.ProfileResponse{ID: user.ID, Name: user.Name, Email: user.Email, Roles: roles, ActiveRole: activeRole}

	// Entry point ringkasan saldo lintas role (business rule Level 1), diisi nyata sesuai data yang sudah ada.
	if w, err := s.walletRepo.GetByUserID(ctx, id); err == nil {
		resp.WalletBalance = &w.Balance
	}
	if store, err := s.storeRepo.FindByUserID(ctx, id); err == nil {
		income, _, _, err := s.orderRepo.SellerIncomeSum(ctx, store.ID)
		if err == nil {
			resp.StoreIncome = &income
		}
	}
	// driver earning ringkas: sum driver_earning dari delivery yang completed
	deliveries, err := s.deliveryRepo.ListByDriver(ctx, id)
	if err == nil {
		var total float64
		for _, d := range deliveries {
			if d.Status == "completed" {
				total += d.DriverEarning
			}
		}
		resp.DriverEarning = &total
	}

	return resp, nil
}
