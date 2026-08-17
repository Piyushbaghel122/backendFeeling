package router

import (
	"net/http"
	"go-api/src/controller"
)


func RegisterAuthRoutes() {
	http.HandleFunc("/api/auth/signup", controller.SignUpController)
	http.HandleFunc("/api/auth/login", controller.LoginController)
	http.HandleFunc("/api/auth/logout", controller.LogoutController)
	http.HandleFunc("/api/auth/profile", controller.GetProfileController)
	http.HandleFunc("/api/auth/upload-avatar", controller.UploadAvatarController)
}
