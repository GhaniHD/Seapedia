package api

import (
	"backend-seapedia/internal/handler"
	"backend-seapedia/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, jwtSecret string) *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/api/v1")
	{
		// endpoint publik, tidak butuh login
		apiV1.POST("/register", authHandler.Register)
		apiV1.POST("/login", authHandler.Login)

		// endpoint yang wajib login
		protected := apiV1.Group("/")
		protected.Use(middleware.JWTMiddleware(jwtSecret))
		{
			protected.GET("/profile", userHandler.GetProfile)
		}
	}

	return r
}