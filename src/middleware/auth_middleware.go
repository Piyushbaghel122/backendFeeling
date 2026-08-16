package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware protects routes by requiring a valid JWT token in the "auth_token" cookie
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Error(w, "Unauthorized - No token provided", http.StatusUnauthorized)
			return
		}

		// 2. Parse the token
		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// 3. Check if it's valid
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}

		// 4. Extract claims and put user_id in the request context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized - Invalid claims", http.StatusUnauthorized)
			return
		}

		userID := claims["user_id"]
		ctx := context.WithValue(r.Context(), "user_id", userID)

		// 5. Call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
