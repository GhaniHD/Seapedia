package utils

import "html"

// SanitizeText mencegah XSS pada konten yang diinput publik (review aplikasi, dsb).
// Semua tag/karakter HTML di-escape sebelum disimpan, sehingga kalau ditampilkan
// mentah oleh frontend pun tidak akan pernah dieksekusi sebagai script.
// (Level 7 - Secure Inputs, Queries, and Public Comments)
func SanitizeText(input string) string {
	return html.EscapeString(input)
}
