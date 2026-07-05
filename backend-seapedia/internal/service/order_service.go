package service

import (
	"context"
	"errors"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/model"
	"backend-seapedia/internal/repository"

	"github.com/google/uuid"
)

type OrderService interface {
	ListBuyerOrders(ctx context.Context, buyerID uuid.UUID) ([]dto.OrderResponse, error)
	GetBuyerOrderDetail(ctx context.Context, buyerID, orderID uuid.UUID) (*dto.OrderResponse, error)
	ListSellerOrders(ctx context.Context, sellerUserID uuid.UUID) ([]dto.OrderResponse, error)
	ProcessOrder(ctx context.Context, sellerUserID, orderID uuid.UUID) error
	BuyerSpendingReport(ctx context.Context, buyerID uuid.UUID) (*dto.SpendingReportResponse, error)
	SellerIncomeReport(ctx context.Context, sellerUserID uuid.UUID) (*dto.IncomeReportResponse, error)
}

type orderService struct {
	orderRepo    repository.OrderRepository
	storeRepo    repository.StoreRepository
	deliveryRepo repository.DeliveryRepository
}

func NewOrderService(orderRepo repository.OrderRepository, storeRepo repository.StoreRepository, deliveryRepo repository.DeliveryRepository) OrderService {
	return &orderService{orderRepo: orderRepo, storeRepo: storeRepo, deliveryRepo: deliveryRepo}
}

func (s *orderService) ListBuyerOrders(ctx context.Context, buyerID uuid.UUID) ([]dto.OrderResponse, error) {
	orders, err := s.orderRepo.ListByBuyer(ctx, buyerID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.OrderResponse, 0, len(orders))
	for _, o := range orders {
		out = append(out, *s.toOrderResp(ctx, &o, false))
	}
	return out, nil
}

func (s *orderService) GetBuyerOrderDetail(ctx context.Context, buyerID, orderID uuid.UUID) (*dto.OrderResponse, error) {
	o, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if o.BuyerID != buyerID {
		return nil, errors.New("order ini bukan milik Anda")
	}
	return s.toOrderResp(ctx, o, true), nil
}

func (s *orderService) ListSellerOrders(ctx context.Context, sellerUserID uuid.UUID) ([]dto.OrderResponse, error) {
	store, err := s.storeRepo.FindByUserID(ctx, sellerUserID)
	if err != nil {
		return nil, err
	}
	orders, err := s.orderRepo.ListByStore(ctx, store.ID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.OrderResponse, 0, len(orders))
	for _, o := range orders {
		out = append(out, *s.toOrderResp(ctx, &o, false))
	}
	return out, nil
}

// ProcessOrder: Seller memindahkan order dari "Sedang Dikemas" -> "Menunggu Pengirim",
// lalu membuat delivery job yang baru terlihat oleh Driver (business rule Level 4 & 5:
// order tidak boleh available buat Driver sebelum diproses seller).
func (s *orderService) ProcessOrder(ctx context.Context, sellerUserID, orderID uuid.UUID) error {
	store, err := s.storeRepo.FindByUserID(ctx, sellerUserID)
	if err != nil {
		return err
	}
	o, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	if o.StoreID != store.ID {
		return errors.New("Anda hanya bisa memproses order milik toko Anda sendiri")
	}
	if o.Status != model.StatusSedangDikemas {
		return errors.New("order ini sudah diproses sebelumnya")
	}

	if err := s.orderRepo.UpdateStatus(ctx, orderID, model.StatusMenungguPengirim); err != nil {
		return err
	}
	if err := s.orderRepo.AddStatusHistory(ctx, &model.OrderStatusHistory{ID: uuid.New(), OrderID: orderID, Status: model.StatusMenungguPengirim, Note: "Diproses oleh seller"}); err != nil {
		return err
	}

	driverEarning := o.DeliveryFee * 0.8 // lihat pkg/utils.DriverEarningRate, didokumentasikan di README
	return s.deliveryRepo.Create(ctx, &model.Delivery{
		ID: uuid.New(), OrderID: orderID, Status: model.DeliveryJobAvailable, Fee: o.DeliveryFee, DriverEarning: driverEarning,
	})
}

func (s *orderService) BuyerSpendingReport(ctx context.Context, buyerID uuid.UUID) (*dto.SpendingReportResponse, error) {
	total, count, err := s.orderRepo.BuyerSpendingSum(ctx, buyerID)
	if err != nil {
		return nil, err
	}
	return &dto.SpendingReportResponse{TotalOrders: count, TotalSpending: total}, nil
}

func (s *orderService) SellerIncomeReport(ctx context.Context, sellerUserID uuid.UUID) (*dto.IncomeReportResponse, error) {
	store, err := s.storeRepo.FindByUserID(ctx, sellerUserID)
	if err != nil {
		return nil, err
	}
	income, reversed, count, err := s.orderRepo.SellerIncomeSum(ctx, store.ID)
	if err != nil {
		return nil, err
	}
	return &dto.IncomeReportResponse{TotalOrders: count, TotalIncome: income, TotalReversed: reversed}, nil
}

func (s *orderService) toOrderResp(ctx context.Context, o *model.Order, withDetail bool) *dto.OrderResponse {
	resp := &dto.OrderResponse{
		ID: o.ID, OrderNo: o.OrderNo, DeliveryMethod: o.DeliveryMethod, Subtotal: o.Subtotal, DiscountAmount: o.DiscountAmount,
		DeliveryFee: o.DeliveryFee, TaxAmount: o.TaxAmount, Total: o.Total, Status: o.Status, DeadlineAt: o.DeadlineAt, CreatedAt: o.CreatedAt,
	}
	if st, err := s.storeRepo.FindByID(ctx, o.StoreID); err == nil {
		resp.StoreName = st.Name
	}
	if withDetail {
		items, _ := s.orderRepo.ListItems(ctx, o.ID)
		for _, it := range items {
			resp.Items = append(resp.Items, dto.OrderItemResponse{ProductName: it.ProductName, Price: it.Price, Quantity: it.Quantity})
		}
		history, _ := s.orderRepo.ListStatusHistory(ctx, o.ID)
		for _, h := range history {
			resp.StatusHistory = append(resp.StatusHistory, dto.StatusHistoryResponse{Status: h.Status, Note: h.Note, CreatedAt: h.CreatedAt})
		}
	}
	return resp
}
