package router

import (
	"net/http"
	"go-api/src/controller"
)

// RegisterAIChatRoutes adds the AI chat routes to the HTTP mux
func RegisterAIChatRoutes() {
	http.HandleFunc("/api/chat/agent", controller.HandleAgentChat)
	http.HandleFunc("/api/chat/title", controller.HandleChatTitle)
}
