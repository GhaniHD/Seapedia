package main

import (
	"log"

	"backend-seapedia/api"
	"backend-seapedia/db"
	"backend-seapedia/internal/config"
	"backend-seapedia/internal/handler"
	"backend-seapedia/internal/migration"
	"backend-seapedia/internal/repository"
	"backend-seapedia/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	pool, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := migration.RunMigrations(pool, "db/migrations"); err != nil {
		log.Fatal("gagal migrasi:", err)
	}

	db.SeedDemoData(pool)

	// ---------- Repository layer ----------
	userRepo := repository.NewUserRepository(pool)
	reviewRepo := repository.NewReviewRepository(pool)
	storeRepo := repository.NewStoreRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	walletRepo := repository.NewWalletRepository(pool)
	addressRepo := repository.NewAddressRepository(pool)
	cartRepo := repository.NewCartRepository(pool)
	discountRepo := repository.NewDiscountRepository(pool)
	orderRepo := repository.NewOrderRepository(pool)
	deliveryRepo := repository.NewDeliveryRepository(pool)
	clockRepo := repository.NewClockRepository(pool)

	// ---------- Service layer ----------
	authService := service.NewAuthService(userRepo, walletRepo, cartRepo, cfg.JWTSecret)
	userService := service.NewUserService(userRepo, walletRepo, storeRepo, orderRepo, deliveryRepo)
	reviewService := service.NewReviewService(reviewRepo)
	storeService := service.NewStoreService(storeRepo)
	productService := service.NewProductService(productRepo, storeRepo)
	walletService := service.NewWalletService(walletRepo, addressRepo)
	cartService := service.NewCartService(cartRepo, productRepo, storeRepo)
	discountService := service.NewDiscountService(discountRepo)
	checkoutService := service.NewCheckoutService(cartRepo, productRepo, addressRepo, walletRepo, orderRepo, deliveryRepo, discountService)
	orderService := service.NewOrderService(orderRepo, storeRepo, deliveryRepo)
	deliveryService := service.NewDeliveryService(deliveryRepo, orderRepo, storeRepo, addressRepo)
	adminService := service.NewAdminService(userRepo, storeRepo, productRepo, orderRepo, discountRepo, deliveryRepo)
	overdueService := service.NewOverdueService(clockRepo, orderRepo, productRepo, walletRepo)

	// ---------- Handler layer ----------
	h := &api.Handlers{
		Auth:     handler.NewAuthHandler(authService),
		User:     handler.NewUserHandler(userService),
		Review:   handler.NewReviewHandler(reviewService),
		Store:    handler.NewStoreHandler(storeService),
		Product:  handler.NewProductHandler(productService),
		Wallet:   handler.NewWalletHandler(walletService),
		Cart:     handler.NewCartHandler(cartService),
		Checkout: handler.NewCheckoutHandler(checkoutService),
		Order:    handler.NewOrderHandler(orderService),
		Discount: handler.NewDiscountHandler(discountService),
		Delivery: handler.NewDeliveryHandler(deliveryService),
		Admin:    handler.NewAdminHandler(adminService, overdueService),
	}

	r := api.SetupRoutes(h, cfg.JWTSecret)

	log.Println("Server jalan di port:", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
