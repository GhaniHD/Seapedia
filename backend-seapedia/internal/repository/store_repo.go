package repository

import (
	"backend-seapedia/internal/model"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoreRepository interface {
	Create(ctx context.Context, s *model.Store) error
	Update(ctx context.Context, s *model.Store) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Store, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Store, error)
	NameExists(ctx context.Context, name string, excludeStoreID *uuid.UUID) (bool, error)
	List(ctx context.Context) ([]model.Store, error)
	Count(ctx context.Context) (int64, error)
}

type storeRepository struct{ db *pgxpool.Pool }

func NewStoreRepository(db *pgxpool.Pool) StoreRepository { return &storeRepository{db: db} }

func (r *storeRepository) Create(ctx context.Context, s *model.Store) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO stores (id, user_id, name, description) VALUES ($1,$2,$3,$4)`,
		s.ID, s.UserID, s.Name, s.Description)
	return err
}

func (r *storeRepository) Update(ctx context.Context, s *model.Store) error {
	_, err := r.db.Exec(ctx,
		`UPDATE stores SET name=$1, description=$2, updated_at=NOW() WHERE id=$3`,
		s.Name, s.Description, s.ID)
	return err
}

func (r *storeRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Store, error) {
	var s model.Store
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, name, description, created_at, updated_at FROM stores WHERE user_id=$1`,
		userID).Scan(&s.ID, &s.UserID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("store belum dibuat")
		}
		return nil, err
	}
	return &s, nil
}

func (r *storeRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Store, error) {
	var s model.Store
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, name, description, created_at, updated_at FROM stores WHERE id=$1`,
		id).Scan(&s.ID, &s.UserID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("store tidak ditemukan")
		}
		return nil, err
	}
	return &s, nil
}

func (r *storeRepository) NameExists(ctx context.Context, name string, excludeStoreID *uuid.UUID) (bool, error) {
	var exists bool
	if excludeStoreID != nil {
		err := r.db.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM stores WHERE LOWER(name)=LOWER($1) AND id != $2)`,
			name, *excludeStoreID).Scan(&exists)
		return exists, err
	}
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM stores WHERE LOWER(name)=LOWER($1))`, name).Scan(&exists)
	return exists, err
}

func (r *storeRepository) List(ctx context.Context) ([]model.Store, error) {
	rows, err := r.db.Query(ctx, `SELECT id, user_id, name, description, created_at, updated_at FROM stores ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Store
	for rows.Next() {
		var s model.Store
		if err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

func (r *storeRepository) Count(ctx context.Context) (int64, error) {
	var c int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM stores`).Scan(&c)
	return c, err
}
