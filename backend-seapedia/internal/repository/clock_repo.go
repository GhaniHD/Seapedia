package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ClockRepository menyimpan "virtual now" sistem, dipakai untuk mensimulasikan
// pergantian hari (Level 6) tanpa harus menunggu waktu asli berjalan.
type ClockRepository interface {
	Now(ctx context.Context) (time.Time, error)
	AdvanceDay(ctx context.Context, days int) (time.Time, error)
}

type clockRepository struct{ db *pgxpool.Pool }

func NewClockRepository(db *pgxpool.Pool) ClockRepository { return &clockRepository{db: db} }

func (r *clockRepository) Now(ctx context.Context) (time.Time, error) {
	var t time.Time
	err := r.db.QueryRow(ctx, `SELECT virtual_now FROM system_clock WHERE id = 1`).Scan(&t)
	return t, err
}

func (r *clockRepository) AdvanceDay(ctx context.Context, days int) (time.Time, error) {
	var t time.Time
	err := r.db.QueryRow(ctx,
		`UPDATE system_clock SET virtual_now = virtual_now + ($1 || ' days')::interval WHERE id = 1 RETURNING virtual_now`,
		days).Scan(&t)
	return t, err
}
