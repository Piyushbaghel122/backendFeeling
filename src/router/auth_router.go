package router

import (
	"net/http"
	"go-api/src/controller"
)

// RegisterAuthRoutes adds the authentication routes to the HTTP mux
func RegisterAuthRoutes() {
	http.HandleFunc("/api/auth/signup", controller.SignUpController)
	http.HandleFunc("/api/auth/login", controller.LoginController)
	http.HandleFunc("/api/auth/logout", controller.LogoutController)
	http.HandleFunc("/api/auth/profile", controller.GetProfileController)
}
