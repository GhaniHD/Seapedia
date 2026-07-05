package main

import (
	"log"

	"backend-seapedia/api"
	"backend-seapedia/internal/config"
	"backend-seapedia/internal/handler"
	"backend-seapedia/internal/migration"
	"backend-seapedia/internal/repository"
	"backend-seapedia/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	dsn := config.BuildDSN(cfg)
	if err := migration.RunMigrations(dsn); err != nil {
		log.Fatal("gagal migrasi:", err)
	}

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Dependency Injection: dari layer paling dalam ke luar
	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	userService := service.NewUserService(userRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	r := api.SetupRoutes(authHandler, userHandler, cfg.JWTSecret)

	log.Println("Server jalan di port:", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
