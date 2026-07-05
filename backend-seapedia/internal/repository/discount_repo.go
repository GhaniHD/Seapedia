package repository

import (
	"backend-seapedia/internal/model"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DiscountRepository interface {
	CreateVoucher(ctx context.Context, v *model.Voucher) error
	CreatePromo(ctx context.Context, p *model.Promo) error
	FindVoucherByCode(ctx context.Context, code string) (*model.Voucher, error)
	FindPromoByCode(ctx context.Context, code string) (*model.Promo, error)
	IncrementVoucherUsage(ctx context.Context, id uuid.UUID) error
	ListVouchers(ctx context.Context) ([]model.Voucher, error)
	ListPromos(ctx context.Context) ([]model.Promo, error)
	CountVouchers(ctx context.Context) (int64, error)
	CountPromos(ctx context.Context) (int64, error)
}

type discountRepository struct{ db *pgxpool.Pool }

func NewDiscountRepository(db *pgxpool.Pool) DiscountRepository { return &discountRepository{db: db} }

func (r *discountRepository) CreateVoucher(ctx context.Context, v *model.Voucher) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO vouchers (id, code, discount_type, discount_value, expiry_date, usage_limit) VALUES ($1,$2,$3,$4,$5,$6)`,
		v.ID, v.Code, v.DiscountType, v.DiscountValue, v.ExpiryDate, v.UsageLimit)
	return err
}

func (r *discountRepository) CreatePromo(ctx context.Context, p *model.Promo) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO promos (id, code, discount_type, discount_value, expiry_date) VALUES ($1,$2,$3,$4,$5)`,
		p.ID, p.Code, p.DiscountType, p.DiscountValue, p.ExpiryDate)
	return err
}

func (r *discountRepository) FindVoucherByCode(ctx context.Context, code string) (*model.Voucher, error) {
	var v model.Voucher
	err := r.db.QueryRow(ctx,
		`SELECT id, code, discount_type, discount_value, expiry_date, usage_limit, usage_count, created_at FROM vouchers WHERE code=$1`,
		code).Scan(&v.ID, &v.Code, &v.DiscountType, &v.DiscountValue, &v.ExpiryDate, &v.UsageLimit, &v.UsageCount, &v.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("voucher tidak ditemukan")
		}
		return nil, err
	}
	return &v, nil
}

func (r *discountRepository) FindPromoByCode(ctx context.Context, code string) (*model.Promo, error) {
	var p model.Promo
	err := r.db.QueryRow(ctx,
		`SELECT id, code, discount_type, discount_value, expiry_date, created_at FROM promos WHERE code=$1`,
		code).Scan(&p.ID, &p.Code, &p.DiscountType, &p.DiscountValue, &p.ExpiryDate, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("promo tidak ditemukan")
		}
		return nil, err
	}
	return &p, nil
}

func (r *discountRepository) IncrementVoucherUsage(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE vouchers SET usage_count = usage_count + 1 WHERE id = $1 AND usage_count < usage_limit`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("kuota voucher sudah habis")
	}
	return nil
}

func (r *discountRepository) ListVouchers(ctx context.Context) ([]model.Voucher, error) {
	rows, err := r.db.Query(ctx, `SELECT id, code, discount_type, discount_value, expiry_date, usage_limit, usage_count, created_at FROM vouchers ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Voucher
	for rows.Next() {
		var v model.Voucher
		if err := rows.Scan(&v.ID, &v.Code, &v.DiscountType, &v.DiscountValue, &v.ExpiryDate, &v.UsageLimit, &v.UsageCount, &v.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

func (r *discountRepository) ListPromos(ctx context.Context) ([]model.Promo, error) {
	rows, err := r.db.Query(ctx, `SELECT id, code, discount_type, discount_value, expiry_date, created_at FROM promos ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Promo
	for rows.Next() {
		var p model.Promo
		if err := rows.Scan(&p.ID, &p.Code, &p.DiscountType, &p.DiscountValue, &p.ExpiryDate, &p.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (r *discountRepository) CountVouchers(ctx context.Context) (int64, error) {
	var c int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM vouchers`).Scan(&c)
	return c, err
}
func (r *discountRepository) CountPromos(ctx context.Context) (int64, error) {
	var c int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM promos`).Scan(&c)
	return c, err
}

var _ = time.Now
