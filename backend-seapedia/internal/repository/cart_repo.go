package repository

import (
	"backend-seapedia/internal/model"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepository interface {
	EnsureCart(ctx context.Context, userID uuid.UUID) (*model.Cart, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Cart, error)
	SetStore(ctx context.Context, cartID uuid.UUID, storeID *uuid.UUID) error
	UpsertItem(ctx context.Context, cartID, productID uuid.UUID, qty int) error
	UpdateItemQty(ctx context.Context, cartID, productID uuid.UUID, qty int) error
	RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error
	ListItems(ctx context.Context, cartID uuid.UUID) ([]model.CartItem, error)
	ClearItems(ctx context.Context, cartID uuid.UUID) error
}

type cartRepository struct{ db *pgxpool.Pool }

func NewCartRepository(db *pgxpool.Pool) CartRepository { return &cartRepository{db: db} }

func (r *cartRepository) EnsureCart(ctx context.Context, userID uuid.UUID) (*model.Cart, error) {
	c, err := r.GetByUserID(ctx, userID)
	if err == nil {
		return c, nil
	}
	id := uuid.New()
	_, err = r.db.Exec(ctx, `INSERT INTO carts (id, user_id) VALUES ($1,$2) ON CONFLICT (user_id) DO NOTHING`, id, userID)
	if err != nil {
		return nil, err
	}
	return r.GetByUserID(ctx, userID)
}

func (r *cartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Cart, error) {
	var c model.Cart
	err := r.db.QueryRow(ctx, `SELECT id, user_id, store_id FROM carts WHERE user_id=$1`, userID).
		Scan(&c.ID, &c.UserID, &c.StoreID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("cart tidak ditemukan")
		}
		return nil, err
	}
	return &c, nil
}

func (r *cartRepository) SetStore(ctx context.Context, cartID uuid.UUID, storeID *uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE carts SET store_id = $1, updated_at = NOW() WHERE id = $2`, storeID, cartID)
	return err
}

func (r *cartRepository) UpsertItem(ctx context.Context, cartID, productID uuid.UUID, qty int) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO cart_items (id, cart_id, product_id, quantity) VALUES ($1,$2,$3,$4)
		 ON CONFLICT (cart_id, product_id) DO UPDATE SET quantity = cart_items.quantity + $4, updated_at = NOW()`,
		uuid.New(), cartID, productID, qty)
	return err
}

func (r *cartRepository) UpdateItemQty(ctx context.Context, cartID, productID uuid.UUID, qty int) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE cart_items SET quantity = $1, updated_at = NOW() WHERE cart_id = $2 AND product_id = $3`,
		qty, cartID, productID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("item cart tidak ditemukan")
	}
	return nil
}

func (r *cartRepository) RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2`, cartID, productID)
	return err
}

func (r *cartRepository) ListItems(ctx context.Context, cartID uuid.UUID) ([]model.CartItem, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, cart_id, product_id, quantity, created_at FROM cart_items WHERE cart_id=$1`, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.CartItem
	for rows.Next() {
		var it model.CartItem
		if err := rows.Scan(&it.ID, &it.CartID, &it.ProductID, &it.Quantity, &it.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, nil
}

func (r *cartRepository) ClearItems(ctx context.Context, cartID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cart_items WHERE cart_id = $1`, cartID)
	return err
}
