package repository

import (
	"backend-seapedia/internal/model"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeliveryRepository interface {
	Create(ctx context.Context, d *model.Delivery) error
	FindByOrderID(ctx context.Context, orderID uuid.UUID) (*model.Delivery, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Delivery, error)
	ListAvailable(ctx context.Context) ([]model.Delivery, error)
	TakeJob(ctx context.Context, id, driverID uuid.UUID) error
	CompleteJob(ctx context.Context, id uuid.UUID) error
	ListByDriver(ctx context.Context, driverID uuid.UUID) ([]model.Delivery, error)
	Count(ctx context.Context) (int64, error)
}

type deliveryRepository struct{ db *pgxpool.Pool }

func NewDeliveryRepository(db *pgxpool.Pool) DeliveryRepository { return &deliveryRepository{db: db} }

const deliveryCols = `id, order_id, driver_id, status, fee, driver_earning, taken_at, completed_at, created_at`

func scanDelivery(row pgx.Row) (*model.Delivery, error) {
	var d model.Delivery
	err := row.Scan(&d.ID, &d.OrderID, &d.DriverID, &d.Status, &d.Fee, &d.DriverEarning, &d.TakenAt, &d.CompletedAt, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *deliveryRepository) Create(ctx context.Context, d *model.Delivery) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO deliveries (id, order_id, status, fee, driver_earning) VALUES ($1,$2,$3,$4,$5)`,
		d.ID, d.OrderID, d.Status, d.Fee, d.DriverEarning)
	return err
}

func (r *deliveryRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*model.Delivery, error) {
	d, err := scanDelivery(r.db.QueryRow(ctx, `SELECT `+deliveryCols+` FROM deliveries WHERE order_id=$1`, orderID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("delivery job tidak ditemukan")
		}
		return nil, err
	}
	return d, nil
}

func (r *deliveryRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Delivery, error) {
	d, err := scanDelivery(r.db.QueryRow(ctx, `SELECT `+deliveryCols+` FROM deliveries WHERE id=$1`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("delivery job tidak ditemukan")
		}
		return nil, err
	}
	return d, nil
}

func (r *deliveryRepository) ListAvailable(ctx context.Context) ([]model.Delivery, error) {
	rows, err := r.db.Query(ctx, `SELECT `+deliveryCols+` FROM deliveries WHERE status='available' ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Delivery
	for rows.Next() {
		var d model.Delivery
		if err := rows.Scan(&d.ID, &d.OrderID, &d.DriverID, &d.Status, &d.Fee, &d.DriverEarning, &d.TakenAt, &d.CompletedAt, &d.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

// TakeJob: hanya berhasil kalau status masih 'available' -> mencegah 2 driver ambil job yang sama (atomic).
func (r *deliveryRepository) TakeJob(ctx context.Context, id, driverID uuid.UUID) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE deliveries SET status='taken', driver_id=$1, taken_at=NOW() WHERE id=$2 AND status='available'`,
		driverID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("job ini sudah diambil driver lain")
	}
	return nil
}

func (r *deliveryRepository) CompleteJob(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE deliveries SET status='completed', completed_at=NOW() WHERE id=$1 AND status='taken'`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("job belum diambil / sudah selesai")
	}
	return nil
}

func (r *deliveryRepository) ListByDriver(ctx context.Context, driverID uuid.UUID) ([]model.Delivery, error) {
	rows, err := r.db.Query(ctx, `SELECT `+deliveryCols+` FROM deliveries WHERE driver_id=$1 ORDER BY created_at DESC`, driverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Delivery
	for rows.Next() {
		var d model.Delivery
		if err := rows.Scan(&d.ID, &d.OrderID, &d.DriverID, &d.Status, &d.Fee, &d.DriverEarning, &d.TakenAt, &d.CompletedAt, &d.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

func (r *deliveryRepository) Count(ctx context.Context) (int64, error) {
	var c int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM deliveries`).Scan(&c)
	return c, err
}
