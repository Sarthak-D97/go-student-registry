package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Sarthak-D97/go_stuAPI/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthorizeJWT validates the token from the Authorization header
// It expects the service.JWTService to be passed from main.go
func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header found"})
			return
		}

		// 1. Extract the token string
		tokenString := authHeader
		if strings.HasPrefix(tokenString, BEARER_SCHEMA) {
			tokenString = tokenString[len(BEARER_SCHEMA):]
		}

		// Clean up any extra whitespace
		tokenString = strings.TrimSpace(tokenString)

		// 2. Validate the token using the injected service
		token, err := jwtService.ValidateToken(tokenString)

		// 3. Check for validity
		if err != nil || token == nil || !token.Valid {
			// Log the actual error to the console for debugging
			fmt.Printf("Token Validation Failed: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// 4. Token is valid - Extract Claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Log for debugging (remove in production)
			fmt.Println("User Authenticated:", claims["username"])

			// IMPORTANT: Save claims to context so Controllers can use them
			// Usage in Controller: claims, _ := c.Get("claims")
			c.Set("claims", claims)

			// If you want to store specific values:
			// c.Set("userID", claims["user_id"])
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Next()
	}
}
