package service

import (
	"context"
	"errors"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/repository"

	"github.com/google/uuid"
)

type CartService interface {
	AddItem(ctx context.Context, userID uuid.UUID, req dto.AddCartItemRequest) (*dto.CartResponse, error)
	UpdateItem(ctx context.Context, userID, productID uuid.UUID, req dto.UpdateCartItemRequest) (*dto.CartResponse, error)
	RemoveItem(ctx context.Context, userID, productID uuid.UUID) (*dto.CartResponse, error)
	GetCart(ctx context.Context, userID uuid.UUID) (*dto.CartResponse, error)
	ClearCart(ctx context.Context, userID uuid.UUID) error
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	storeRepo   repository.StoreRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository, storeRepo repository.StoreRepository) CartService {
	return &cartService{cartRepo: cartRepo, productRepo: productRepo, storeRepo: storeRepo}
}

// AddItem menerapkan business rule single-store checkout:
// satu cart cuma boleh isi produk dari SATU toko. Kalau buyer coba nambah produk
// dari toko lain, ditolak dengan pesan jelas suruh clear cart dulu.
func (s *cartService) AddItem(ctx context.Context, userID uuid.UUID, req dto.AddCartItemRequest) (*dto.CartResponse, error) {
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return nil, errors.New("product_id tidak valid")
	}
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	cart, err := s.cartRepo.EnsureCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cart.StoreID != nil && *cart.StoreID != product.StoreID {
		return nil, errors.New("cart Anda sudah berisi produk dari toko lain. Kosongkan cart terlebih dahulu untuk belanja di toko ini (single-store checkout)")
	}
	if cart.StoreID == nil {
		if err := s.cartRepo.SetStore(ctx, cart.ID, &product.StoreID); err != nil {
			return nil, err
		}
	}

	if err := s.cartRepo.UpsertItem(ctx, cart.ID, productID, req.Quantity); err != nil {
		return nil, err
	}
	return s.GetCart(ctx, userID)
}

func (s *cartService) UpdateItem(ctx context.Context, userID, productID uuid.UUID, req dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	cart, err := s.cartRepo.EnsureCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.cartRepo.UpdateItemQty(ctx, cart.ID, productID, req.Quantity); err != nil {
		return nil, err
	}
	return s.GetCart(ctx, userID)
}

func (s *cartService) RemoveItem(ctx context.Context, userID, productID uuid.UUID) (*dto.CartResponse, error) {
	cart, err := s.cartRepo.EnsureCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.cartRepo.RemoveItem(ctx, cart.ID, productID); err != nil {
		return nil, err
	}
	items, err := s.cartRepo.ListItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		_ = s.cartRepo.SetStore(ctx, cart.ID, nil) // cart kosong -> boleh ganti toko lagi
	}
	return s.GetCart(ctx, userID)
}

func (s *cartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	cart, err := s.cartRepo.EnsureCart(ctx, userID)
	if err != nil {
		return err
	}
	if err := s.cartRepo.ClearItems(ctx, cart.ID); err != nil {
		return err
	}
	return s.cartRepo.SetStore(ctx, cart.ID, nil)
}

func (s *cartService) GetCart(ctx context.Context, userID uuid.UUID) (*dto.CartResponse, error) {
	cart, err := s.cartRepo.EnsureCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	items, err := s.cartRepo.ListItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	resp := &dto.CartResponse{StoreID: cart.StoreID, Items: []dto.CartItemResponse{}}
	if cart.StoreID != nil {
		if st, err := s.storeRepo.FindByID(ctx, *cart.StoreID); err == nil {
			resp.StoreName = st.Name
		}
	}

	var subtotal float64
	for _, it := range items {
		p, err := s.productRepo.FindByID(ctx, it.ProductID)
		if err != nil {
			continue
		}
		lineTotal := p.Price * float64(it.Quantity)
		subtotal += lineTotal
		resp.Items = append(resp.Items, dto.CartItemResponse{
			ID: it.ID, ProductID: it.ProductID, Name: p.Name, Price: p.Price, Quantity: it.Quantity, Subtotal: lineTotal,
		})
	}
	resp.Subtotal = subtotal
	return resp, nil
}
