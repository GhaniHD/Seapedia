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

type WalletService interface {
	Topup(ctx context.Context, userID uuid.UUID, req dto.TopupRequest) (*dto.WalletResponse, error)
	GetBalance(ctx context.Context, userID uuid.UUID) (*dto.WalletResponse, error)
	ListTransactions(ctx context.Context, userID uuid.UUID) ([]dto.WalletTransactionResponse, error)

	AddAddress(ctx context.Context, userID uuid.UUID, req dto.UpsertAddressRequest) (*dto.AddressResponse, error)
	ListAddresses(ctx context.Context, userID uuid.UUID) ([]dto.AddressResponse, error)
}

type walletService struct {
	walletRepo  repository.WalletRepository
	addressRepo repository.AddressRepository
}

func NewWalletService(walletRepo repository.WalletRepository, addressRepo repository.AddressRepository) WalletService {
	return &walletService{walletRepo: walletRepo, addressRepo: addressRepo}
}

// Topup: dummy top-up sesuai spek Level 3 (langsung berhasil, tidak ada payment gateway asli).
func (s *walletService) Topup(ctx context.Context, userID uuid.UUID, req dto.TopupRequest) (*dto.WalletResponse, error) {
	if req.Amount <= 0 {
		return nil, errors.New("jumlah top up harus lebih dari 0")
	}
	w, err := s.walletRepo.EnsureWallet(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.walletRepo.AddBalance(ctx, userID, req.Amount); err != nil {
		return nil, err
	}
	if err := s.walletRepo.RecordTransaction(ctx, &model.WalletTransaction{
		ID: uuid.New(), WalletID: w.ID, Type: "topup", Amount: req.Amount, Description: "Dummy top up",
	}); err != nil {
		return nil, err
	}
	w2, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.WalletResponse{Balance: w2.Balance}, nil
}

func (s *walletService) GetBalance(ctx context.Context, userID uuid.UUID) (*dto.WalletResponse, error) {
	w, err := s.walletRepo.EnsureWallet(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.WalletResponse{Balance: w.Balance}, nil
}

func (s *walletService) ListTransactions(ctx context.Context, userID uuid.UUID) ([]dto.WalletTransactionResponse, error) {
	w, err := s.walletRepo.EnsureWallet(ctx, userID)
	if err != nil {
		return nil, err
	}
	txs, err := s.walletRepo.ListTransactions(ctx, w.ID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.WalletTransactionResponse, 0, len(txs))
	for _, t := range txs {
		out = append(out, dto.WalletTransactionResponse{ID: t.ID, Type: t.Type, Amount: t.Amount, Description: t.Description, CreatedAt: t.CreatedAt})
	}
	return out, nil
}

func (s *walletService) AddAddress(ctx context.Context, userID uuid.UUID, req dto.UpsertAddressRequest) (*dto.AddressResponse, error) {
	if req.IsDefault {
		if err := s.addressRepo.ClearDefault(ctx, userID); err != nil {
			return nil, err
		}
	}
	a := &model.Address{
		ID: uuid.New(), UserID: userID,
		Label: utils.SanitizeText(req.Label), Detail: utils.SanitizeText(req.Detail), IsDefault: req.IsDefault,
	}
	if err := s.addressRepo.Create(ctx, a); err != nil {
		return nil, err
	}
	return &dto.AddressResponse{ID: a.ID, Label: a.Label, Detail: a.Detail, IsDefault: a.IsDefault}, nil
}

func (s *walletService) ListAddresses(ctx context.Context, userID uuid.UUID) ([]dto.AddressResponse, error) {
	addrs, err := s.addressRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.AddressResponse, 0, len(addrs))
	for _, a := range addrs {
		out = append(out, dto.AddressResponse{ID: a.ID, Label: a.Label, Detail: a.Detail, IsDefault: a.IsDefault})
	}
	return out, nil
}
