package service

import (
	"context"
	"errors"
	"strings"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/model"
	"backend-seapedia/internal/repository"
	"backend-seapedia/pkg/utils"

	"github.com/google/uuid"
)

type StoreService interface {
	UpsertMyStore(ctx context.Context, userID uuid.UUID, req dto.UpsertStoreRequest) (*dto.StoreResponse, error)
	GetMyStore(ctx context.Context, userID uuid.UUID) (*dto.StoreResponse, error)
	GetPublicStore(ctx context.Context, storeID uuid.UUID) (*dto.StoreResponse, error)
	ListStores(ctx context.Context) ([]dto.StoreResponse, error)
}

type storeService struct{ repo repository.StoreRepository }

func NewStoreService(repo repository.StoreRepository) StoreService { return &storeService{repo: repo} }

func (s *storeService) UpsertMyStore(ctx context.Context, userID uuid.UUID, req dto.UpsertStoreRequest) (*dto.StoreResponse, error) {
	name := strings.TrimSpace(req.Name)
	existing, err := s.repo.FindByUserID(ctx, userID)

	var excludeID *uuid.UUID
	if err == nil {
		excludeID = &existing.ID
	}
	taken, tErr := s.repo.NameExists(ctx, name, excludeID)
	if tErr != nil {
		return nil, tErr
	}
	if taken {
		return nil, errors.New("nama toko sudah dipakai, silakan pilih nama lain")
	}

	desc := utils.SanitizeText(req.Description)

	if err == nil {
		existing.Name = name
		existing.Description = desc
		if err := s.repo.Update(ctx, existing); err != nil {
			return nil, err
		}
		return &dto.StoreResponse{ID: existing.ID, Name: existing.Name, Description: existing.Description, CreatedAt: existing.CreatedAt}, nil
	}

	store := &model.Store{ID: uuid.New(), UserID: userID, Name: name, Description: desc}
	if err := s.repo.Create(ctx, store); err != nil {
		return nil, errors.New("gagal membuat toko, kemungkinan nama sudah dipakai")
	}
	return &dto.StoreResponse{ID: store.ID, Name: store.Name, Description: store.Description}, nil
}

func (s *storeService) GetMyStore(ctx context.Context, userID uuid.UUID) (*dto.StoreResponse, error) {
	store, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.StoreResponse{ID: store.ID, Name: store.Name, Description: store.Description, CreatedAt: store.CreatedAt}, nil
}

func (s *storeService) GetPublicStore(ctx context.Context, storeID uuid.UUID) (*dto.StoreResponse, error) {
	store, err := s.repo.FindByID(ctx, storeID)
	if err != nil {
		return nil, err
	}
	return &dto.StoreResponse{ID: store.ID, Name: store.Name, Description: store.Description, CreatedAt: store.CreatedAt}, nil
}

func (s *storeService) ListStores(ctx context.Context) ([]dto.StoreResponse, error) {
	stores, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.StoreResponse, 0, len(stores))
	for _, st := range stores {
		out = append(out, dto.StoreResponse{ID: st.ID, Name: st.Name, Description: st.Description, CreatedAt: st.CreatedAt})
	}
	return out, nil
}
