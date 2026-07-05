package repository

import (
	"backend-seapedia/internal/model"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository mendefinisikan operasi yang dapat dilakukan
// terhadap data user di database, seperti mencari, membuat,
// dan memperbarui data user.
type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
}

// userRepository menyimpan koneksi database PostgreSQL
// yang akan digunakan untuk menjalankan query.
type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository membuat repository user baru dan
// menghubungkannya dengan koneksi database agar siap digunakan.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	query := `SELECT id,name,email,password,role FROM users WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	query := `SELECT id,name,email,password,role FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("email tidak ditemukan")
		}
		return nil, err
	}

	return &user, err
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (id,name,email,password,role) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query,
		user.ID, user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`

	_, err := r.db.Exec(ctx, query,
		user.Name, user.Email, user.ID)
	if err != nil {
		return err
	}
	return nil
}
