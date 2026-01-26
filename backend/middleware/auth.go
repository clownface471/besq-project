package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT secret key - should match the one in controllers/auth.go
// In production, this should be in environment variables
var jwtSecret = []byte("your-secret-key-change-this-in-production")

// Claims represents the JWT claims structure (should match controllers/auth.go)
type Claims struct {
	Sub  string `json:"sub"`  // User ID
	Role string `json:"role"` // User Role
	jwt.RegisteredClaims
}

// AuthMiddleware validates the Authorization header (Bearer token)
// Parses the JWT and sets the user role into the context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		// Check if it's a Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format. Use 'Bearer <token>'",
			})
			return
		}

		tokenString := tokenParts[1]

		// Parse and validate JWT token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		// Check for parsing errors
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		// Check if token is valid
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token is not valid",
			})
			return
		}

		// Set user information in context
		c.Set("role", claims.Role)
		c.Set("userId", claims.Sub)

		// Continue to next handler
		c.Next()
	}
}

// RequireRoles middleware allows access only to users with specific roles
func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve role from context
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Role not found in context. Authentication required.",
			})
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Invalid role format.",
			})
			return
		}

		// Check if user's role is in the allowed roles
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		// If role is not in allowed roles, deny access
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Akses ditolak: Role anda tidak memiliki izin",
		})
	}
}