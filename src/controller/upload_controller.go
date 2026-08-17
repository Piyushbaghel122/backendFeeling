package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go-api/src/config"
	"go-api/src/modules"

	"github.com/golang-jwt/jwt/v5"
)

// UploadAvatarController handles image uploads for user profiles
func UploadAvatarController(w http.ResponseWriter, r *http.Request) {
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

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID := claims["user_id"]

	// Parse multipart form data (limit to 10MB)
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Get the file from form data
	file, handler, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Error retrieving file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create uploads directory if it doesn't exist
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	// Create a unique filename
	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("avatar_%v_%d%s", userID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Save the file
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	// Create the URL to be served to the frontend
	avatarURL := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)

	// Update the database with the new AvatarURL
	var user modules.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.AvatarURL = avatarURL
	if err := config.DB.Save(&user).Error; err != nil {
		http.Error(w, "Error updating database", http.StatusInternalServerError)
		return
	}

	// Return the new URL
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":    "Avatar uploaded successfully",
		"avatar_url": avatarURL,
	})
}
