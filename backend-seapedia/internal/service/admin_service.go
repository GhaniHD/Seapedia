package service

import (
	"context"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/repository"
)

type AdminService interface {
	Dashboard(ctx context.Context) (*dto.AdminDashboardResponse, error)
}

type adminService struct {
	userRepo     repository.UserRepository
	storeRepo    repository.StoreRepository
	productRepo  repository.ProductRepository
	orderRepo    repository.OrderRepository
	discountRepo repository.DiscountRepository
	deliveryRepo repository.DeliveryRepository
}

func NewAdminService(userRepo repository.UserRepository, storeRepo repository.StoreRepository, productRepo repository.ProductRepository,
	orderRepo repository.OrderRepository, discountRepo repository.DiscountRepository, deliveryRepo repository.DeliveryRepository) AdminService {
	return &adminService{userRepo, storeRepo, productRepo, orderRepo, discountRepo, deliveryRepo}
}

// Dashboard mengumpulkan data monitoring lintas resource untuk Admin (Level 6):
// users, stores, products, orders, vouchers/promos, delivery jobs, overdue orders.
func (s *adminService) Dashboard(ctx context.Context) (*dto.AdminDashboardResponse, error) {
	resp := &dto.AdminDashboardResponse{}
	var err error
	if resp.TotalUsers, err = s.userRepo.CountUsers(ctx); err != nil {
		return nil, err
	}
	if resp.TotalStores, err = s.storeRepo.Count(ctx); err != nil {
		return nil, err
	}
	if resp.TotalProducts, err = s.productRepo.Count(ctx); err != nil {
		return nil, err
	}
	if resp.TotalOrders, err = s.orderRepo.CountAll(ctx); err != nil {
		return nil, err
	}
	if resp.TotalVouchers, err = s.discountRepo.CountVouchers(ctx); err != nil {
		return nil, err
	}
	if resp.TotalPromos, err = s.discountRepo.CountPromos(ctx); err != nil {
		return nil, err
	}
	if resp.TotalDeliveryJobs, err = s.deliveryRepo.Count(ctx); err != nil {
		return nil, err
	}
	if resp.OverdueOrdersCount, err = s.orderRepo.CountOverdue(ctx); err != nil {
		return nil, err
	}
	return resp, nil
}
