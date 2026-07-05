package repository

import (
	"backend-seapedia/internal/model"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	Create(ctx context.Context, p *model.Product) error
	Update(ctx context.Context, p *model.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Product, error)
	ListByStore(ctx context.Context, storeID uuid.UUID) ([]model.Product, error)
	ListPublic(ctx context.Context) ([]model.Product, error)
	DecreaseStock(ctx context.Context, id uuid.UUID, qty int) error
	IncreaseStock(ctx context.Context, id uuid.UUID, qty int) error
	Count(ctx context.Context) (int64, error)
}

type productRepository struct{ db *pgxpool.Pool }

func NewProductRepository(db *pgxpool.Pool) ProductRepository { return &productRepository{db: db} }

func (r *productRepository) Create(ctx context.Context, p *model.Product) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO products (id, store_id, name, description, price, stock) VALUES ($1,$2,$3,$4,$5,$6)`,
		p.ID, p.StoreID, p.Name, p.Description, p.Price, p.Stock)
	return err
}

func (r *productRepository) Update(ctx context.Context, p *model.Product) error {
	_, err := r.db.Exec(ctx,
		`UPDATE products SET name=$1, description=$2, price=$3, stock=$4, updated_at=NOW() WHERE id=$5`,
		p.Name, p.Description, p.Price, p.Stock, p.ID)
	return err
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)
	return err
}

func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	var p model.Product
	err := r.db.QueryRow(ctx,
		`SELECT id, store_id, name, description, price, stock, created_at, updated_at FROM products WHERE id=$1`,
		id).Scan(&p.ID, &p.StoreID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) ListByStore(ctx context.Context, storeID uuid.UUID) ([]model.Product, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, store_id, name, description, price, stock, created_at, updated_at FROM products WHERE store_id=$1 ORDER BY created_at DESC`,
		storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProducts(rows)
}

func (r *productRepository) ListPublic(ctx context.Context) ([]model.Product, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, store_id, name, description, price, stock, created_at, updated_at FROM products ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProducts(rows)
}

func scanProducts(rows pgx.Rows) ([]model.Product, error) {
	var out []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.StoreID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (r *productRepository) DecreaseStock(ctx context.Context, id uuid.UUID, qty int) error {
	// WHERE stock >= qty mencegah stock jadi negatif (dicek di level SQL, aman dari race condition)
	tag, err := r.db.Exec(ctx,
		`UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2 AND stock >= $1`,
		qty, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("stok produk tidak cukup")
	}
	return nil
}

func (r *productRepository) IncreaseStock(ctx context.Context, id uuid.UUID, qty int) error {
	_, err := r.db.Exec(ctx, `UPDATE products SET stock = stock + $1, updated_at = NOW() WHERE id = $2`, qty, id)
	return err
}

func (r *productRepository) Count(ctx context.Context) (int64, error) {
	var c int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM products`).Scan(&c)
	return c, err
}
