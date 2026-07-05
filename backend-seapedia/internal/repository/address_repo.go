package repository

import (
	"backend-seapedia/internal/model"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddressRepository interface {
	Create(ctx context.Context, a *model.Address) error
	ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Address, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Address, error)
	ClearDefault(ctx context.Context, userID uuid.UUID) error
}

type addressRepository struct{ db *pgxpool.Pool }

func NewAddressRepository(db *pgxpool.Pool) AddressRepository { return &addressRepository{db: db} }

func (r *addressRepository) Create(ctx context.Context, a *model.Address) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO addresses (id, user_id, label, detail, is_default) VALUES ($1,$2,$3,$4,$5)`,
		a.ID, a.UserID, a.Label, a.Detail, a.IsDefault)
	return err
}

func (r *addressRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Address, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, label, detail, is_default, created_at FROM addresses WHERE user_id=$1 ORDER BY created_at DESC`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Address
	for rows.Next() {
		var a model.Address
		if err := rows.Scan(&a.ID, &a.UserID, &a.Label, &a.Detail, &a.IsDefault, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func (r *addressRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Address, error) {
	var a model.Address
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, label, detail, is_default, created_at FROM addresses WHERE id=$1`, id).
		Scan(&a.ID, &a.UserID, &a.Label, &a.Detail, &a.IsDefault, &a.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("alamat tidak ditemukan")
		}
		return nil, err
	}
	return &a, nil
}

func (r *addressRepository) ClearDefault(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE addresses SET is_default = FALSE WHERE user_id = $1`, userID)
	return err
}
