package repository

import (
	"backend-seapedia/internal/model"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository interface {
	EnsureWallet(ctx context.Context, userID uuid.UUID) (*model.Wallet, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Wallet, error)
	AddBalance(ctx context.Context, userID uuid.UUID, amount float64) error
	DeductBalance(ctx context.Context, userID uuid.UUID, amount float64) error
	RecordTransaction(ctx context.Context, tx *model.WalletTransaction) error
	ListTransactions(ctx context.Context, walletID uuid.UUID) ([]model.WalletTransaction, error)
}

type walletRepository struct{ db *pgxpool.Pool }

func NewWalletRepository(db *pgxpool.Pool) WalletRepository { return &walletRepository{db: db} }

func (r *walletRepository) EnsureWallet(ctx context.Context, userID uuid.UUID) (*model.Wallet, error) {
	w, err := r.GetByUserID(ctx, userID)
	if err == nil {
		return w, nil
	}
	id := uuid.New()
	_, err = r.db.Exec(ctx,
		`INSERT INTO wallets (id, user_id, balance) VALUES ($1,$2,0) ON CONFLICT (user_id) DO NOTHING`, id, userID)
	if err != nil {
		return nil, err
	}
	return r.GetByUserID(ctx, userID)
}

func (r *walletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Wallet, error) {
	var w model.Wallet
	err := r.db.QueryRow(ctx, `SELECT id, user_id, balance, updated_at FROM wallets WHERE user_id=$1`, userID).
		Scan(&w.ID, &w.UserID, &w.Balance, &w.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("wallet tidak ditemukan")
		}
		return nil, err
	}
	return &w, nil
}

func (r *walletRepository) AddBalance(ctx context.Context, userID uuid.UUID, amount float64) error {
	_, err := r.db.Exec(ctx, `UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE user_id = $2`, amount, userID)
	return err
}

func (r *walletRepository) DeductBalance(ctx context.Context, userID uuid.UUID, amount float64) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE wallets SET balance = balance - $1, updated_at = NOW() WHERE user_id = $2 AND balance >= $1`,
		amount, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("saldo wallet tidak cukup")
	}
	return nil
}

func (r *walletRepository) RecordTransaction(ctx context.Context, tx *model.WalletTransaction) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO wallet_transactions (id, wallet_id, type, amount, description, order_id) VALUES ($1,$2,$3,$4,$5,$6)`,
		tx.ID, tx.WalletID, tx.Type, tx.Amount, tx.Description, tx.OrderID)
	return err
}

func (r *walletRepository) ListTransactions(ctx context.Context, walletID uuid.UUID) ([]model.WalletTransaction, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, wallet_id, type, amount, description, order_id, created_at FROM wallet_transactions WHERE wallet_id=$1 ORDER BY created_at DESC`,
		walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.WalletTransaction
	for rows.Next() {
		var t model.WalletTransaction
		if err := rows.Scan(&t.ID, &t.WalletID, &t.Type, &t.Amount, &t.Description, &t.OrderID, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}
