package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/model"
	"backend-seapedia/internal/repository"

	"github.com/google/uuid"
)

type DiscountService interface {
	CreateVoucher(ctx context.Context, req dto.CreateVoucherRequest) (*dto.VoucherResponse, error)
	CreatePromo(ctx context.Context, req dto.CreatePromoRequest) (*dto.PromoResponse, error)
	ListVouchers(ctx context.Context) ([]dto.VoucherResponse, error)
	ListPromos(ctx context.Context) ([]dto.PromoResponse, error)

	// Validate mengecek kode (voucher ATAU promo, dicari di kedua tabel) dan menghitung nilai
	// diskonnya untuk subtotal tertentu. Business rule: voucher & promo TIDAK bisa digabung -
	// satu kode dipakai per checkout, hasilnya jelas dibedakan lewat field "kind".
	Validate(ctx context.Context, code string, subtotal float64) (amount float64, kind string, voucherID, promoID *uuid.UUID, err error)
	CommitUsage(ctx context.Context, voucherID *uuid.UUID) error
}

type discountService struct{ repo repository.DiscountRepository }

func NewDiscountService(repo repository.DiscountRepository) DiscountService { return &discountService{repo: repo} }

func (s *discountService) CreateVoucher(ctx context.Context, req dto.CreateVoucherRequest) (*dto.VoucherResponse, error) {
	expiry, err := time.Parse(time.RFC3339, req.ExpiryDate)
	if err != nil {
		return nil, errors.New("format expiry_date harus RFC3339, contoh: 2026-12-31T23:59:59Z")
	}
	v := &model.Voucher{
		ID: uuid.New(), Code: strings.ToUpper(strings.TrimSpace(req.Code)), DiscountType: req.DiscountType,
		DiscountValue: req.DiscountValue, ExpiryDate: expiry, UsageLimit: req.UsageLimit,
	}
	if err := s.repo.CreateVoucher(ctx, v); err != nil {
		return nil, errors.New("gagal membuat voucher, kemungkinan kode sudah dipakai")
	}
	return toVoucherResp(v), nil
}

func (s *discountService) CreatePromo(ctx context.Context, req dto.CreatePromoRequest) (*dto.PromoResponse, error) {
	expiry, err := time.Parse(time.RFC3339, req.ExpiryDate)
	if err != nil {
		return nil, errors.New("format expiry_date harus RFC3339, contoh: 2026-12-31T23:59:59Z")
	}
	p := &model.Promo{
		ID: uuid.New(), Code: strings.ToUpper(strings.TrimSpace(req.Code)), DiscountType: req.DiscountType,
		DiscountValue: req.DiscountValue, ExpiryDate: expiry,
	}
	if err := s.repo.CreatePromo(ctx, p); err != nil {
		return nil, errors.New("gagal membuat promo, kemungkinan kode sudah dipakai")
	}
	return toPromoResp(p), nil
}

func (s *discountService) ListVouchers(ctx context.Context) ([]dto.VoucherResponse, error) {
	vs, err := s.repo.ListVouchers(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.VoucherResponse, 0, len(vs))
	for _, v := range vs {
		out = append(out, *toVoucherResp(&v))
	}
	return out, nil
}

func (s *discountService) ListPromos(ctx context.Context) ([]dto.PromoResponse, error) {
	ps, err := s.repo.ListPromos(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.PromoResponse, 0, len(ps))
	for _, p := range ps {
		out = append(out, *toPromoResp(&p))
	}
	return out, nil
}

func (s *discountService) Validate(ctx context.Context, code string, subtotal float64) (float64, string, *uuid.UUID, *uuid.UUID, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return 0, "", nil, nil, nil
	}

	if v, err := s.repo.FindVoucherByCode(ctx, code); err == nil {
		if time.Now().After(v.ExpiryDate) {
			return 0, "", nil, nil, errors.New("voucher sudah kedaluwarsa")
		}
		if v.UsageCount >= v.UsageLimit {
			return 0, "", nil, nil, errors.New("kuota voucher sudah habis")
		}
		return calcDiscount(v.DiscountType, v.DiscountValue, subtotal), "voucher", &v.ID, nil, nil
	}

	if p, err := s.repo.FindPromoByCode(ctx, code); err == nil {
		if time.Now().After(p.ExpiryDate) {
			return 0, "", nil, nil, errors.New("promo sudah kedaluwarsa")
		}
		return calcDiscount(p.DiscountType, p.DiscountValue, subtotal), "promo", nil, &p.ID, nil
	}

	return 0, "", nil, nil, errors.New("kode voucher/promo tidak ditemukan")
}

func (s *discountService) CommitUsage(ctx context.Context, voucherID *uuid.UUID) error {
	if voucherID == nil {
		return nil
	}
	return s.repo.IncrementVoucherUsage(ctx, *voucherID)
}

func calcDiscount(discountType string, value, subtotal float64) float64 {
	if discountType == model.DiscountPercent {
		amt := subtotal * (value / 100)
		if amt > subtotal {
			amt = subtotal
		}
		return amt
	}
	if value > subtotal {
		return subtotal
	}
	return value
}

func toVoucherResp(v *model.Voucher) *dto.VoucherResponse {
	return &dto.VoucherResponse{ID: v.ID.String(), Code: v.Code, DiscountType: v.DiscountType, DiscountValue: v.DiscountValue, ExpiryDate: v.ExpiryDate, UsageLimit: v.UsageLimit, UsageCount: v.UsageCount}
}
func toPromoResp(p *model.Promo) *dto.PromoResponse {
	return &dto.PromoResponse{ID: p.ID.String(), Code: p.Code, DiscountType: p.DiscountType, DiscountValue: p.DiscountValue, ExpiryDate: p.ExpiryDate}
}
