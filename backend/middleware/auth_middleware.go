package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthAndRoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token dibutuhkan"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte("secret_key_besq_2026"), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Safety check: Pastikan role ada dan berbentuk string
			roleClaim, ok := claims["role"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Format Token Salah: Role tidak ditemukan"})
				c.Abort()
				return
			}

			userRole := roleClaim

			// Cek apakah role ada di daftar yang diizinkan
			authorized := false
			for _, role := range allowedRoles {
				if role == userRole {
					authorized = true
					break
				}
			}

			if !authorized {
				c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak punya akses ke halaman ini!"})
				c.Abort()
				return
			}

			// Set data ke Context
			c.Set("userID", claims["user_id"])
			c.Set("username", claims["username"])
			c.Set("userRole", claims["role"])
			c.Set("nama", claims["nama"])

			// --- PERBAIKAN UTAMA DI SINI ---
			// Cek apakah username ada di dalam token
			if username, ok := claims["username"].(string); ok {
				c.Set("username", username)
			} else {
				// Fallback jika token lama tidak punya username
				// Ini mencegah panic "interface conversion: interface {} is nil"
				c.Set("username", "UNKNOWN_USER")
			}
			// -------------------------------

			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid", "details": err.Error()})
			c.Abort()
		}
	}
}

// ... (ActionMiddleware biarkan tetap sama) ...
func ActionMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("userRole")

		if c.Request.Method != "GET" {
			if userRole != requiredRole && userRole != "ADMIN" {
				c.JSON(http.StatusForbidden, gin.H{
					"error": fmt.Sprintf("Hanya %s yang boleh melakukan aksi ini!", requiredRole),
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}