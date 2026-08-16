package main

import (
	"log"
	"net/http"
	"go-api/src"
	"go-api/src/config"
	// "go-api/src/services"
	"go-api/src/router"
)

func main() {
	config.ConnectDB()
	config.ConnectRedis()

	// Test the picture uploading using genAI and storing in pinecone
	// 	answer, err := services.QueryPinecone("What is the capital of France?")
	var err error = nil
	if err != nil {
		log.Println("Pinecone Error:", err)
	} else {
		var answer any = nil
		log.Println("Pinecone Answer:", answer)
	}
	// Register all application routes
	router.RegisterAIChatRoutes()
	router.RegisterAuthRoutes()

	http.HandleFunc("/", src.HomeHandler)
	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
