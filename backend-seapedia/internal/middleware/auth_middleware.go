package middleware

import (
	"net/http"
	"strings"

	"backend-seapedia/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AuthRequired mem-parse token dan mewajibkan active_role sudah dipilih
// (menolak temp-token pending). Semua endpoint privat (Seller/Buyer/Driver/Admin) wajib pakai ini.
func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := parseAuthHeader(c, secret)
		if !ok {
			return
		}
		activeRole, _ := claims["activeRole"].(string)
		if activeRole == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "silakan pilih role aktif terlebih dahulu (POST /select-role)"})
			c.Abort()
			return
		}
		c.Set("user_id", claims["userID"])
		c.Set("active_role", activeRole)
		c.Next()
	}
}

// AuthPending dipakai khusus endpoint /select-role: menerima temp-token (active_role kosong) maupun token biasa.
func AuthPending(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := parseAuthHeader(c, secret)
		if !ok {
			return
		}
		c.Set("user_id", claims["userID"])
		c.Next()
	}
}

func parseAuthHeader(c *gin.Context, secret string) (map[string]interface{}, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak ditemukan"})
		c.Abort()
		return nil, false
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := jwt.ParseToken(tokenString, secret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return nil, false
	}
	return claims, true
}

// RequireRole membatasi endpoint hanya untuk active_role tertentu.
// PENTING: ini mengecek active_role dari TOKEN (server-side), bukan cuma dari tampilan UI (Level 7).
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool)
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		activeRole := c.GetString("active_role")
		if !allowed[activeRole] {
			c.JSON(http.StatusForbidden, gin.H{"error": "role aktif Anda tidak memiliki akses ke resource ini"})
			c.Abort()
			return
		}
		c.Next()
	}
}
