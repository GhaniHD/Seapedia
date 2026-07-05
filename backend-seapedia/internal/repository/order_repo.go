package repository

import (
	"backend-seapedia/internal/model"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	Create(ctx context.Context, o *model.Order) error
	CreateItems(ctx context.Context, items []model.OrderItem) error
	AddStatusHistory(ctx context.Context, h *model.OrderStatusHistory) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	SetDriver(ctx context.Context, orderID, driverID uuid.UUID) error
	MarkRefunded(ctx context.Context, id uuid.UUID) error
	MarkSellerIncomeReversed(ctx context.Context, id uuid.UUID) error
	MarkStockRestored(ctx context.Context, id uuid.UUID) error
	ListByBuyer(ctx context.Context, buyerID uuid.UUID) ([]model.Order, error)
	ListByStore(ctx context.Context, storeID uuid.UUID) ([]model.Order, error)
	ListItems(ctx context.Context, orderID uuid.UUID) ([]model.OrderItem, error)
	ListStatusHistory(ctx context.Context, orderID uuid.UUID) ([]model.OrderStatusHistory, error)
	CountAll(ctx context.Context) (int64, error)
	CountOverdue(ctx context.Context) (int64, error)
	ListAllActiveOverdue(ctx context.Context) ([]model.Order, error)
	SellerIncomeSum(ctx context.Context, storeID uuid.UUID) (float64, float64, int, error)
	BuyerSpendingSum(ctx context.Context, buyerID uuid.UUID) (float64, int, error)
}

type orderRepository struct{ db *pgxpool.Pool }

func NewOrderRepository(db *pgxpool.Pool) OrderRepository { return &orderRepository{db: db} }

func (r *orderRepository) Create(ctx context.Context, o *model.Order) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO orders (id, order_no, buyer_id, store_id, address_id, delivery_method, subtotal, discount_amount,
			delivery_fee, tax_amount, total, status, voucher_id, promo_id, deadline_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		o.ID, o.OrderNo, o.BuyerID, o.StoreID, o.AddressID, o.DeliveryMethod, o.Subtotal, o.DiscountAmount,
		o.DeliveryFee, o.TaxAmount, o.Total, o.Status, o.VoucherID, o.PromoID, o.DeadlineAt)
	return err
}

func (r *orderRepository) CreateItems(ctx context.Context, items []model.OrderItem) error {
	for _, it := range items {
		_, err := r.db.Exec(ctx,
			`INSERT INTO order_items (id, order_id, product_id, product_name, price, quantity) VALUES ($1,$2,$3,$4,$5,$6)`,
			it.ID, it.OrderID, it.ProductID, it.ProductName, it.Price, it.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *orderRepository) AddStatusHistory(ctx context.Context, h *model.OrderStatusHistory) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO order_status_history (id, order_id, status, note) VALUES ($1,$2,$3,$4)`,
		h.ID, h.OrderID, h.Status, h.Note)
	return err
}

func (r *orderRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	var o model.Order
	err := r.db.QueryRow(ctx, `
		SELECT id, order_no, buyer_id, store_id, address_id, delivery_method, subtotal, discount_amount,
			delivery_fee, tax_amount, total, status, voucher_id, promo_id, driver_id, refunded,
			seller_income_reversed, stock_restored, deadline_at, created_at, updated_at
		FROM orders WHERE id=$1`, id).Scan(
		&o.ID, &o.OrderNo, &o.BuyerID, &o.StoreID, &o.AddressID, &o.DeliveryMethod, &o.Subtotal, &o.DiscountAmount,
		&o.DeliveryFee, &o.TaxAmount, &o.Total, &o.Status, &o.VoucherID, &o.PromoID, &o.DriverID, &o.Refunded,
		&o.SellerIncomeReversed, &o.StockRestored, &o.DeadlineAt, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("order tidak ditemukan")
		}
		return nil, err
	}
	return &o, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}

func (r *orderRepository) SetDriver(ctx context.Context, orderID, driverID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE orders SET driver_id = $1, updated_at = NOW() WHERE id = $2`, driverID, orderID)
	return err
}

func (r *orderRepository) MarkRefunded(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE orders SET refunded = TRUE WHERE id = $1`, id)
	return err
}
func (r *orderRepository) MarkSellerIncomeReversed(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE orders SET seller_income_reversed = TRUE WHERE id = $1`, id)
	return err
}
func (r *orderRepository) MarkStockRestored(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE orders SET stock_restored = TRUE WHERE id = $1`, id)
	return err
}

func scanOrders(rows pgx.Rows) ([]model.Order, error) {
	var out []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(
			&o.ID, &o.OrderNo, &o.BuyerID, &o.StoreID, &o.AddressID, &o.DeliveryMethod, &o.Subtotal, &o.DiscountAmount,
			&o.DeliveryFee, &o.TaxAmount, &o.Total, &o.Status, &o.VoucherID, &o.PromoID, &o.DriverID, &o.Refunded,
			&o.SellerIncomeReversed, &o.StockRestored, &o.DeadlineAt, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, nil
}

const orderCols = `id, order_no, buyer_id, store_id, address_id, delivery_method, subtotal, discount_amount,
	delivery_fee, tax_amount, total, status, voucher_id, promo_id, driver_id, refunded,
	seller_income_reversed, stock_restored, deadline_at, created_at, updated_at`

func (r *orderRepository) ListByBuyer(ctx context.Context, buyerID uuid.UUID) ([]model.Order, error) {
	rows, err := r.db.Query(ctx, `SELECT `+orderCols+` FROM orders WHERE buyer_id=$1 ORDER BY created_at DESC`, buyerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOrders(rows)
}

func (r *orderRepository) ListByStore(ctx context.Context, storeID uuid.UUID) ([]model.Order, error) {
	rows, err := r.db.Query(ctx, `SELECT `+orderCols+` FROM orders WHERE store_id=$1 ORDER BY created_at DESC`, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOrders(rows)
}

func (r *orderRepository) ListItems(ctx context.Context, orderID uuid.UUID) ([]model.OrderItem, error) {
	rows, err := r.db.Query(ctx, `SELECT id, order_id, product_id, product_name, price, quantity FROM order_items WHERE order_id=$1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.OrderItem
	for rows.Next() {
		var it model.OrderItem
		if err := rows.Scan(&it.ID, &it.OrderID, &it.ProductID, &it.ProductName, &it.Price, &it.Quantity); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, nil
}

func (r *orderRepository) ListStatusHistory(ctx context.Context, orderID uuid.UUID) ([]model.OrderStatusHistory, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, order_id, status, note, created_at FROM order_status_history WHERE order_id=$1 ORDER BY created_at ASC`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.OrderStatusHistory
	for rows.Next() {
		var h model.OrderStatusHistory
		if err := rows.Scan(&h.ID, &h.OrderID, &h.Status, &h.Note, &h.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, nil
}

// dipakai overdue service: order yang masih "Sedang Dikirim" dan sudah lewat deadline_at pada virtual clock
func (r *orderRepository) ListAllActiveOverdue(ctx context.Context) ([]model.Order, error) {
	rows, err := r.db.Query(ctx, `SELECT `+orderCols+` FROM orders o
		WHERE o.status = 'Sedang Dikirim' AND o.deadline_at IS NOT NULL
		AND o.deadline_at < (SELECT virtual_now FROM system_clock WHERE id = 1)`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOrders(rows)
}

func (r *orderRepository) CountAll(ctx context.Context) (int64, error) {
	var c int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders`).Scan(&c)
	return c, err
}

func (r *orderRepository) CountOverdue(ctx context.Context) (int64, error) {
	var c int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders o
		WHERE o.status = 'Sedang Dikirim' AND o.deadline_at IS NOT NULL
		AND o.deadline_at < (SELECT virtual_now FROM system_clock WHERE id = 1)`).Scan(&c)
	return c, err
}

// SellerIncomeSum: total income (order Pesanan Selesai, tidak refund) dan total yang direversal
func (r *orderRepository) SellerIncomeSum(ctx context.Context, storeID uuid.UUID) (income float64, reversed float64, count int, err error) {
	err = r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(CASE WHEN status = 'Pesanan Selesai' AND seller_income_reversed = FALSE THEN total ELSE 0 END), 0),
		       COALESCE(SUM(CASE WHEN seller_income_reversed = TRUE THEN total ELSE 0 END), 0),
		       COUNT(*)
		FROM orders WHERE store_id = $1`, storeID).Scan(&income, &reversed, &count)
	return
}

func (r *orderRepository) BuyerSpendingSum(ctx context.Context, buyerID uuid.UUID) (total float64, count int, err error) {
	err = r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(CASE WHEN refunded = FALSE THEN total ELSE 0 END), 0), COUNT(*)
		FROM orders WHERE buyer_id = $1`, buyerID).Scan(&total, &count)
	return
}
