package service

import (
	"context"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/model"
	"backend-seapedia/internal/repository"
	"backend-seapedia/pkg/utils"

	"github.com/google/uuid"
)

type ReviewService interface {
	Create(ctx context.Context, req dto.CreateReviewRequest) (*dto.ReviewResponse, error)
	List(ctx context.Context) ([]dto.ReviewResponse, error)
}

type reviewService struct{ repo repository.ReviewRepository }

func NewReviewService(repo repository.ReviewRepository) ReviewService { return &reviewService{repo: repo} }

func (s *reviewService) Create(ctx context.Context, req dto.CreateReviewRequest) (*dto.ReviewResponse, error) {
	rv := &model.AppReview{
		ID:           uuid.New(),
		ReviewerName: utils.SanitizeText(req.ReviewerName),
		Rating:       req.Rating,
		// Sanitasi XSS: comment di-escape sebelum disimpan (Level 7 minimum requirement,
		// tapi diterapkan sejak Level 1 karena field ini publik).
		Comment: utils.SanitizeText(req.Comment),
	}
	if err := s.repo.Create(ctx, rv); err != nil {
		return nil, err
	}
	return toReviewResponse(rv), nil
}

func (s *reviewService) List(ctx context.Context) ([]dto.ReviewResponse, error) {
	reviews, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ReviewResponse, 0, len(reviews))
	for _, r := range reviews {
		out = append(out, *toReviewResponse(&r))
	}
	return out, nil
}

func toReviewResponse(r *model.AppReview) *dto.ReviewResponse {
	return &dto.ReviewResponse{ID: r.ID, ReviewerName: r.ReviewerName, Rating: r.Rating, Comment: r.Comment, CreatedAt: r.CreatedAt}
}
