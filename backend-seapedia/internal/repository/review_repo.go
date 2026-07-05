package repository

import (
	"backend-seapedia/internal/model"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepository interface {
	Create(ctx context.Context, r *model.AppReview) error
	List(ctx context.Context) ([]model.AppReview, error)
}

type reviewRepository struct{ db *pgxpool.Pool }

func NewReviewRepository(db *pgxpool.Pool) ReviewRepository { return &reviewRepository{db: db} }

func (r *reviewRepository) Create(ctx context.Context, rv *model.AppReview) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO app_reviews (id, reviewer_name, rating, comment) VALUES ($1,$2,$3,$4)`,
		rv.ID, rv.ReviewerName, rv.Rating, rv.Comment)
	return err
}

func (r *reviewRepository) List(ctx context.Context) ([]model.AppReview, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, reviewer_name, rating, comment, created_at FROM app_reviews ORDER BY created_at DESC LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.AppReview
	for rows.Next() {
		var rv model.AppReview
		if err := rows.Scan(&rv.ID, &rv.ReviewerName, &rv.Rating, &rv.Comment, &rv.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, rv)
	}
	return out, nil
}

var _ = uuid.Nil
