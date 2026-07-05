package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/model"
	"backend-seapedia/internal/repository"
	"backend-seapedia/pkg/utils"

	"github.com/google/uuid"
)

type CheckoutService interface {
	Checkout(ctx context.Context, buyerID uuid.UUID, req dto.CheckoutRequest) (*dto.CheckoutSummaryResponse, error)
}

type checkoutService struct {
	cartRepo     repository.CartRepository
	productRepo  repository.ProductRepository
	addressRepo  repository.AddressRepository
	walletRepo   repository.WalletRepository
	orderRepo    repository.OrderRepository
	deliveryRepo repository.DeliveryRepository
	discountSvc  DiscountService
}

func NewCheckoutService(cartRepo repository.CartRepository, productRepo repository.ProductRepository,
	addressRepo repository.AddressRepository, walletRepo repository.WalletRepository,
	orderRepo repository.OrderRepository, deliveryRepo repository.DeliveryRepository, discountSvc DiscountService) CheckoutService {
	return &checkoutService{cartRepo, productRepo, addressRepo, walletRepo, orderRepo, deliveryRepo, discountSvc}
}

// Checkout mengimplementasikan aturan bisnis inti:
//   total = subtotal - discount + delivery_fee + ppn(12% dari (subtotal - discount))
// PPN dihitung dari (subtotal - discount), BUKAN dari delivery_fee. Ini dipilih supaya PPN
// murni mengenakan nilai barang, sesuai praktik umum e-commerce. Didokumentasikan di README.
func (s *checkoutService) Checkout(ctx context.Context, buyerID uuid.UUID, req dto.CheckoutRequest) (*dto.CheckoutSummaryResponse, error) {
	addressID, err := uuid.Parse(req.AddressID)
	if err != nil {
		return nil, errors.New("address_id tidak valid")
	}
	address, err := s.addressRepo.FindByID(ctx, addressID)
	if err != nil {
		return nil, err
	}
	if address.UserID != buyerID {
		return nil, errors.New("alamat ini bukan milik Anda")
	}

	cart, err := s.cartRepo.EnsureCart(ctx, buyerID)
	if err != nil {
		return nil, err
	}
	if cart.StoreID == nil {
		return nil, errors.New("cart Anda masih kosong")
	}
	items, err := s.cartRepo.ListItems(ctx, cart.ID)
	if err != nil || len(items) == 0 {
		return nil, errors.New("cart Anda masih kosong")
	}

	// hitung subtotal & siapkan order items dari harga produk TERKINI (bukan dari cache cart)
	var subtotal float64
	type line struct {
		productID uuid.UUID
		name      string
		price     float64
		qty       int
	}
	var lines []line
	for _, it := range items {
		p, err := s.productRepo.FindByID(ctx, it.ProductID)
		if err != nil {
			return nil, err
		}
		if p.Stock < it.Quantity {
			return nil, errors.New("stok produk '" + p.Name + "' tidak cukup")
		}
		subtotal += p.Price * float64(it.Quantity)
		lines = append(lines, line{productID: p.ID, name: p.Name, price: p.Price, qty: it.Quantity})
	}

	var discountAmount float64
	var discountKind string
	var voucherID, promoID *uuid.UUID
	if code := strings.TrimSpace(req.DiscountCode); code != "" {
		discountAmount, discountKind, voucherID, promoID, err = s.discountSvc.Validate(ctx, code, subtotal)
		if err != nil {
			return nil, err
		}
	}

	deliveryFee := utils.DeliveryFee(req.DeliveryMethod)
	taxBase := subtotal - discountAmount
	taxAmount := taxBase * utils.TaxRate
	total := taxBase + deliveryFee + taxAmount

	// Buyer tidak bisa checkout kalau saldo wallet tidak cukup
	wallet, err := s.walletRepo.EnsureWallet(ctx, buyerID)
	if err != nil {
		return nil, err
	}
	if wallet.Balance < total {
		return nil, errors.New("saldo wallet Anda tidak cukup untuk checkout ini")
	}

	now := time.Now()
	deadline := now.Add(utils.DeliverySLA(req.DeliveryMethod))
	order := &model.Order{
		ID: uuid.New(), OrderNo: utils.GenerateOrderNo(now, strings.ToUpper(uuid.New().String()[:6])),
		BuyerID: buyerID, StoreID: *cart.StoreID, AddressID: addressID, DeliveryMethod: req.DeliveryMethod,
		Subtotal: subtotal, DiscountAmount: discountAmount, DeliveryFee: deliveryFee, TaxAmount: taxAmount, Total: total,
		Status: model.StatusSedangDikemas, VoucherID: voucherID, PromoID: promoID, DeadlineAt: &deadline,
	}
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	var orderItems []model.OrderItem
	for _, l := range lines {
		orderItems = append(orderItems, model.OrderItem{ID: uuid.New(), OrderID: order.ID, ProductID: l.productID, ProductName: l.name, Price: l.price, Quantity: l.qty})
		if err := s.productRepo.DecreaseStock(ctx, l.productID, l.qty); err != nil {
			return nil, err
		}
	}
	if err := s.orderRepo.CreateItems(ctx, orderItems); err != nil {
		return nil, err
	}
	if err := s.orderRepo.AddStatusHistory(ctx, &model.OrderStatusHistory{ID: uuid.New(), OrderID: order.ID, Status: model.StatusSedangDikemas, Note: "Order dibuat"}); err != nil {
		return nil, err
	}

	// charge wallet
	if err := s.walletRepo.DeductBalance(ctx, buyerID, total); err != nil {
		return nil, err
	}
	if err := s.walletRepo.RecordTransaction(ctx, &model.WalletTransaction{
		ID: uuid.New(), WalletID: wallet.ID, Type: "checkout", Amount: -total,
		Description: "Checkout order " + order.OrderNo, OrderID: &order.ID,
	}); err != nil {
		return nil, err
	}

	if err := s.discountSvc.CommitUsage(ctx, voucherID); err != nil {
		return nil, err
	}

	// siapkan cart baru (kosong) untuk belanja selanjutnya
	if err := s.cartRepo.ClearItems(ctx, cart.ID); err != nil {
		return nil, err
	}
	_ = s.cartRepo.SetStore(ctx, cart.ID, nil)

	return &dto.CheckoutSummaryResponse{
		OrderID: order.ID, OrderNo: order.OrderNo, Subtotal: subtotal, DiscountAmount: discountAmount, DiscountKind: discountKind,
		DeliveryFee: deliveryFee, TaxAmount: taxAmount, TaxRate: utils.TaxRate, Total: total, Status: order.Status, CreatedAt: now,
	}, nil
}
