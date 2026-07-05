package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// GenerateToken membuat token JWT baru, membawa informasi
// user_id dan role di dalamnya, berlaku selama 24 jam.
func GenerateToken(userID uuid.UUID, role string, secret string) (string, error) {
	claimns := jwt.MapClaims{
		"userID": userID.String(),
		"role":   role,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimns)
	return token.SignedString([]byte(secret))
}

// ParseToken membaca dan memverifikasi token, mengembalikan
// isi claims-nya (user_id, role) kalau valid.
func ParseToken(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token tidak valid!")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Gagal Membaca Token")
	}
	return claims, nil
}
