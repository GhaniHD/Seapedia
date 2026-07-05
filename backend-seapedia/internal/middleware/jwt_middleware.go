package middleware

import (
	"net/http"
	"strings"

	"backend-seapedia/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak ditemukan"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ParseToken(tokenString, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak valid"})
			c.Abort()
			return
		}

		// titipkan data dari token ke context, supaya bisa dipakai handler berikutnya
		c.Set("user_id", claims["userID"])
		c.Set("role", claims["role"])

		c.Next()
	}
}
