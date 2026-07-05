package service

import (
	"context"
	"errors"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/model"
	"backend-seapedia/internal/repository"
	"backend-seapedia/pkg/crypto"
	"backend-seapedia/pkg/jwt"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	SelectRole(ctx context.Context, userID uuid.UUID, role string) (*dto.LoginResponse, error)
	AddRole(ctx context.Context, userID uuid.UUID, role string) error
}

type authService struct {
	userRepo   repository.UserRepository
	walletRepo repository.WalletRepository
	cartRepo   repository.CartRepository
	jwtSecret  string
}

func NewAuthService(userRepo repository.UserRepository, walletRepo repository.WalletRepository, cartRepo repository.CartRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, walletRepo: walletRepo, cartRepo: cartRepo, jwtSecret: jwtSecret}
}

// Register: setiap akun baru otomatis mendapat role "buyer" (buyer wajib punya wallet+cart+address,
// dan ini role default paling wajar untuk akun baru). User bisa menambah role lain (seller/driver)
// lewat AddRole setelah login.
func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashed, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{ID: uuid.New(), Name: req.Name, Email: req.Email, Password: hashed}
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}
	if err := s.userRepo.AddRole(ctx, newUser.ID, model.RoleBuyer); err != nil {
		return nil, err
	}
	// siapkan resource dasar buyer (Level 3) supaya tidak error saat pertama kali dipakai
	if _, err := s.walletRepo.EnsureWallet(ctx, newUser.ID); err != nil {
		return nil, err
	}
	if _, err := s.cartRepo.EnsureCart(ctx, newUser.ID); err != nil {
		return nil, err
	}

	return &dto.UserResponse{ID: newUser.ID, Name: newUser.Name, Email: newUser.Email}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}
	if !crypto.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("email atau password salah")
	}

	roles, err := s.userRepo.GetRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, errors.New("akun ini belum memiliki role apapun, hubungi admin")
	}

	userResp := dto.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email}

	// Admin: selalu langsung aktif tanpa perlu pilih role (business rule: admin dipisah dari multi-role non-admin)
	if len(roles) == 1 && roles[0] == model.RoleAdmin {
		token, err := jwt.GenerateToken(user.ID, model.RoleAdmin, s.jwtSecret)
		if err != nil {
			return nil, err
		}
		return &dto.LoginResponse{Token: token, NeedRoleSelection: false, Roles: roles, ActiveRole: model.RoleAdmin, User: userResp}, nil
	}

	// Non-admin dengan >1 role: WAJIB pilih active role dulu, jangan langsung ke dashboard privat manapun.
	if len(roles) > 1 {
		temp, err := jwt.GenerateTempToken(user.ID, s.jwtSecret)
		if err != nil {
			return nil, err
		}
		return &dto.LoginResponse{Token: temp, NeedRoleSelection: true, Roles: roles, ActiveRole: "", User: userResp}, nil
	}

	// Hanya 1 role -> langsung aktifkan
	token, err := jwt.GenerateToken(user.ID, roles[0], s.jwtSecret)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{Token: token, NeedRoleSelection: false, Roles: roles, ActiveRole: roles[0], User: userResp}, nil
}

func (s *authService) SelectRole(ctx context.Context, userID uuid.UUID, role string) (*dto.LoginResponse, error) {
	has, err := s.userRepo.HasRole(ctx, userID, role)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("Anda tidak memiliki role tersebut")
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	roles, err := s.userRepo.GetRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	token, err := jwt.GenerateToken(userID, role, s.jwtSecret)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{
		Token: token, NeedRoleSelection: false, Roles: roles, ActiveRole: role,
		User: dto.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email},
	}, nil
}

// AddRole: user menambah role baru untuk akunnya sendiri (misal buyer ingin juga jadi seller/driver).
// Admin TIDAK bisa ditambahkan lewat sini (hanya lewat seeder/setup khusus) - didokumentasikan di README.
func (s *authService) AddRole(ctx context.Context, userID uuid.UUID, role string) error {
	if role != model.RoleSeller && role != model.RoleBuyer && role != model.RoleDriver {
		return errors.New("role tidak valid")
	}
	return s.userRepo.AddRole(ctx, userID, role)
}
