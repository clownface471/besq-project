package middleware

import (
	"net/http"
	"fmt"
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
			userRole := claims["role"].(string)

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

			c.Set("userRole", userRole)
			c.Set("userRole", claims["role"])
			c.Set("userID", claims["user_id"])
			c.Set("username", claims["username"]) // Pastikan saat login, username juga dimasukkan ke claims
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid", "details": err.Error()})
			c.Abort()
		}
	}
}

func ActionMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("userRole")

		// Jika mencoba POST, PUT, atau DELETE
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