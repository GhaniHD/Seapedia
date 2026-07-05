package service

import (
	"context"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/model"
	"backend-seapedia/internal/repository"

	"github.com/google/uuid"
)

type OverdueService interface {
	// SimulateNextDay memajukan virtual clock sistem sebanyak N hari, lalu langsung
	// menjalankan proses overdue handling terhadap order yang lewat deadline (Level 6).
	SimulateNextDay(ctx context.Context, days int) (*dto.SimulateNextDayResponse, error)
	ProcessOverdueOrders(ctx context.Context) ([]string, error)
}

type overdueService struct {
	clockRepo    repository.ClockRepository
	orderRepo    repository.OrderRepository
	productRepo  repository.ProductRepository
	walletRepo   repository.WalletRepository
}

func NewOverdueService(clockRepo repository.ClockRepository, orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository, walletRepo repository.WalletRepository) OverdueService {
	return &overdueService{clockRepo: clockRepo, orderRepo: orderRepo, productRepo: productRepo, walletRepo: walletRepo}
}

func (s *overdueService) SimulateNextDay(ctx context.Context, days int) (*dto.SimulateNextDayResponse, error) {
	if days <= 0 {
		days = 1
	}
	newNow, err := s.clockRepo.AdvanceDay(ctx, days)
	if err != nil {
		return nil, err
	}
	handled, err := s.ProcessOverdueOrders(ctx)
	if err != nil {
		return nil, err
	}
	return &dto.SimulateNextDayResponse{VirtualNow: newNow.Format("2006-01-02 15:04:05"), OverdueHandled: handled}, nil
}

// ProcessOverdueOrders: order "Sedang Dikirim" yang sudah lewat deadline_at (dihitung dari
// DeliverySLA per delivery_method) dipindah ke "Dikembalikan", lalu:
//   - saldo buyer di-refund penuh & dicatat di wallet_transactions (kalau belum refund)
//   - kalau income seller sudah tercatat (order Pesanan Selesai), dibalik lewat flag reversed
//     supaya tidak dihitung dobel di laporan income seller
//   - stok produk dikembalikan sesuai qty di order_items
// Semua guarded dengan flag refunded/seller_income_reversed/stock_restored supaya TIDAK terjadi
// double refund / double reversal / double stock restore untuk order yang sama (business rule Level 6).
func (s *overdueService) ProcessOverdueOrders(ctx context.Context) ([]string, error) {
	overdue, err := s.orderRepo.ListAllActiveOverdue(ctx)
	if err != nil {
		return nil, err
	}

	var handled []string
	for _, o := range overdue {
		if err := s.orderRepo.UpdateStatus(ctx, o.ID, model.StatusDikembalikan); err != nil {
			return handled, err
		}
		note := "Overdue: order tidak selesai sebelum deadline (" + o.DeliveryMethod + "), otomatis dikembalikan"
		if err := s.orderRepo.AddStatusHistory(ctx, &model.OrderStatusHistory{ID: uuid.New(), OrderID: o.ID, Status: model.StatusDikembalikan, Note: note}); err != nil {
			return handled, err
		}

		if !o.Refunded {
			if err := s.walletRepo.AddBalance(ctx, o.BuyerID, o.Total); err == nil {
				if w, werr := s.walletRepo.GetByUserID(ctx, o.BuyerID); werr == nil {
					_ = s.walletRepo.RecordTransaction(ctx, &model.WalletTransaction{
						ID: uuid.New(), WalletID: w.ID, Type: "refund", Amount: o.Total,
						Description: "Refund overdue order " + o.OrderNo, OrderID: &o.ID,
					})
				}
				_ = s.orderRepo.MarkRefunded(ctx, o.ID)
			}
		}

		if !o.SellerIncomeReversed {
			_ = s.orderRepo.MarkSellerIncomeReversed(ctx, o.ID)
		}

		if !o.StockRestored {
			items, ierr := s.orderRepo.ListItems(ctx, o.ID)
			if ierr == nil {
				for _, it := range items {
					_ = s.productRepo.IncreaseStock(ctx, it.ProductID, it.Quantity)
				}
			}
			_ = s.orderRepo.MarkStockRestored(ctx, o.ID)
		}

		handled = append(handled, o.OrderNo)
	}
	return handled, nil
}
