package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// GenerateToken membuat token akhir yang sudah punya active_role.
// Dipakai setelah user login (kalau cuma 1 role) atau setelah select-role.
func GenerateToken(userID uuid.UUID, activeRole string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"userID":     userID.String(),
		"activeRole": activeRole,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateTempToken dipakai saat user punya >1 role non-admin dan belum
// memilih role aktif. Token ini TIDAK bisa dipakai untuk mengakses
// endpoint privat manapun (activeRole kosong) sampai memilih role lewat /select-role.
func GenerateTempToken(userID uuid.UUID, secret string) (string, error) {
	claims := jwt.MapClaims{
		"userID":     userID.String(),
		"activeRole": "",
		"pending":    true,
		"exp":        time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing token tidak valid")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token tidak valid atau sudah kedaluwarsa")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("gagal membaca token")
	}
	return claims, nil
}
