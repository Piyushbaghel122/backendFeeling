package controller

import (
	"encoding/json"
	"net/http"
	"go-api/src/services"
)

type ChatRequest struct {
	Messages []services.Message `json:"messages"`
}

func HandleAgentChat(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Generate response using the Agent
	response, err := services.GenerateAgentResponse(req.Messages)
	if err != nil {
		http.Error(w, "Agent failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"response": response,
	})
}

type TitleRequest struct {
	Message string `json:"message"`
}

func HandleChatTitle(w http.ResponseWriter, r *http.Request) {
	var req TitleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
    
	title, err := services.GenerateChatTitle(req.Message)
	if err != nil {
		http.Error(w, "Failed to generate title: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"title": title,
	})
}
