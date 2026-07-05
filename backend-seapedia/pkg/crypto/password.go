package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword mengubah password asli (plain text) menjadi hash
// yang aman disimpan di database. Dipanggil saat Register.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(bytes), nil
}

// CheckPassword membandingkan password yang diinput user saat login
// dengan hash yang tersimpan di database. Return true kalau cocok.
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
