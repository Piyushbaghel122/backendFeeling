package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"go-api/src/config"
	"go-api/src/modules"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUpController(w http.ResponseWriter, r *http.Request) {
	var user modules.User

	// In Go, we decode the JSON body directly into our struct (no need for destructuring)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Hash the password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	// Update the user struct with the hashed password and timestamps
	user.Password = string(hashPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Save the user to the database
	if err := config.DB.Create(&user).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create user"})
		return
	}

	// Create JWT token
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserId,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := tokenObj.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	// Set Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		MaxAge:   24 * 60 * 60, // 24 hours
		HttpOnly: true,
	})

	// Store token in Redis
	ctx := context.Background()
	userIdStr := fmt.Sprint(user.UserId)
	err = config.RedisClient.Set(ctx, userIdStr, tokenString, 24*time.Hour).Err()
	if err != nil {
		http.Error(w, "Failed to save token to Redis", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"user": map[string]interface{}{
			"username":   user.UserName,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
		},
		"token": tokenString,
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginController(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user modules.User
	// Fetch user from DB using the email
	result := config.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	
	// Compare password with DB hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} 

	// Create JWT token
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserId,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := tokenObj.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	// Set Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		MaxAge:   24 * 60 * 60, // 24 hours
		HttpOnly: true,
	})

	// Store token in Redis
	ctx := context.Background()
	userIdStr := fmt.Sprint(user.UserId)
	config.RedisClient.Set(ctx, userIdStr, tokenString, 24*time.Hour)

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   tokenString,
		"user": map[string]interface{}{
			"username":   user.UserName,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
		},
	})
}

func LogoutController(w http.ResponseWriter, r *http.Request) {
	// Clear the cookie by setting MaxAge to -1
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// TODO: Ideally, you should also delete the token from Redis here using the user_id from the context

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}

// GetProfileController gets the currently logged-in user using their JWT cookie or Authorization header
func GetProfileController(w http.ResponseWriter, r *http.Request) {
	var tokenString string

	// Extract the cookie
	cookie, err := r.Cookie("auth_token")
	if err == nil {
		tokenString = cookie.Value
	} else {
		// Fallback to Authorization header
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}
	}

	if tokenString == "" {
		http.Error(w, "Unauthorized - No token found", http.StatusUnauthorized)
		return
	}

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
		return
	}

	// Extract the user_id from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID := claims["user_id"]

	// Fetch full user details from DB using userID
	var user modules.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Profile fetched successfully",
		"user": map[string]interface{}{
			"id":         user.UserId,
			"username":   user.UserName,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
		},
	})
}
