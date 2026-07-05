package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AdminRepository berisi query agregat lintas tabel untuk dashboard monitoring Admin (Level 6).
type AdminRepository interface {
	CountUsers(ctx context.Context) (int64, error)
}

type adminRepository struct{ db *pgxpool.Pool }

func NewAdminRepository(db *pgxpool.Pool) AdminRepository { return &adminRepository{db: db} }

func (r *adminRepository) CountUsers(ctx context.Context) (int64, error) {
	var c int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&c)
	return c, err
}
