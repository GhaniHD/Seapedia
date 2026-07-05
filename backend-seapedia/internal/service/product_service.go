package service

import (
	"context"
	"errors"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/model"
	"backend-seapedia/internal/repository"
	"backend-seapedia/pkg/utils"

	"github.com/google/uuid"
)

type ProductService interface {
	Create(ctx context.Context, userID uuid.UUID, req dto.UpsertProductRequest) (*dto.ProductResponse, error)
	Update(ctx context.Context, userID, productID uuid.UUID, req dto.UpsertProductRequest) (*dto.ProductResponse, error)
	Delete(ctx context.Context, userID, productID uuid.UUID) error
	ListMine(ctx context.Context, userID uuid.UUID) ([]dto.ProductResponse, error)
	ListPublic(ctx context.Context) ([]dto.ProductResponse, error)
	GetPublicDetail(ctx context.Context, productID uuid.UUID) (*dto.ProductResponse, error)
}

type productService struct {
	repo      repository.ProductRepository
	storeRepo repository.StoreRepository
}

func NewProductService(repo repository.ProductRepository, storeRepo repository.StoreRepository) ProductService {
	return &productService{repo: repo, storeRepo: storeRepo}
}

// ownStoreOrErr memastikan Seller punya toko sendiri (ownership check - Level 2 & 7)
func (s *productService) ownStoreOrErr(ctx context.Context, userID uuid.UUID) (*model.Store, error) {
	store, err := s.storeRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("Anda belum membuat toko, buat toko terlebih dahulu")
	}
	return store, nil
}

func (s *productService) Create(ctx context.Context, userID uuid.UUID, req dto.UpsertProductRequest) (*dto.ProductResponse, error) {
	store, err := s.ownStoreOrErr(ctx, userID)
	if err != nil {
		return nil, err
	}
	p := &model.Product{
		ID: uuid.New(), StoreID: store.ID,
		Name: utils.SanitizeText(req.Name), Description: utils.SanitizeText(req.Description),
		Price: req.Price, Stock: req.Stock,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return toProductResponse(p, store.Name), nil
}

func (s *productService) Update(ctx context.Context, userID, productID uuid.UUID, req dto.UpsertProductRequest) (*dto.ProductResponse, error) {
	store, err := s.ownStoreOrErr(ctx, userID)
	if err != nil {
		return nil, err
	}
	p, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if p.StoreID != store.ID {
		return nil, errors.New("Anda hanya bisa mengubah produk milik toko Anda sendiri")
	}
	p.Name = utils.SanitizeText(req.Name)
	p.Description = utils.SanitizeText(req.Description)
	p.Price = req.Price
	p.Stock = req.Stock
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return toProductResponse(p, store.Name), nil
}

func (s *productService) Delete(ctx context.Context, userID, productID uuid.UUID) error {
	store, err := s.ownStoreOrErr(ctx, userID)
	if err != nil {
		return err
	}
	p, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return err
	}
	if p.StoreID != store.ID {
		return errors.New("Anda hanya bisa menghapus produk milik toko Anda sendiri")
	}
	return s.repo.Delete(ctx, productID)
}

func (s *productService) ListMine(ctx context.Context, userID uuid.UUID) ([]dto.ProductResponse, error) {
	store, err := s.ownStoreOrErr(ctx, userID)
	if err != nil {
		return nil, err
	}
	products, err := s.repo.ListByStore(ctx, store.ID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ProductResponse, 0, len(products))
	for _, p := range products {
		out = append(out, *toProductResponse(&p, store.Name))
	}
	return out, nil
}

func (s *productService) ListPublic(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := s.repo.ListPublic(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ProductResponse, 0, len(products))
	for _, p := range products {
		storeName := ""
		if st, err := s.storeRepo.FindByID(ctx, p.StoreID); err == nil {
			storeName = st.Name
		}
		out = append(out, *toProductResponse(&p, storeName))
	}
	return out, nil
}

func (s *productService) GetPublicDetail(ctx context.Context, productID uuid.UUID) (*dto.ProductResponse, error) {
	p, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	storeName := ""
	if st, err := s.storeRepo.FindByID(ctx, p.StoreID); err == nil {
		storeName = st.Name
	}
	return toProductResponse(p, storeName), nil
}

func toProductResponse(p *model.Product, storeName string) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID: p.ID, StoreID: p.StoreID, StoreName: storeName, Name: p.Name, Description: p.Description,
		Price: p.Price, Stock: p.Stock, CreatedAt: p.CreatedAt,
	}
}
