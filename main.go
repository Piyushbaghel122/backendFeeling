package main

import (
	"log"
	"net/http"
	"os"
	"go-api/src"
	"go-api/src/config"
	"go-api/src/middleware"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", src.HomeHandler)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	log.Println("Server running on port", port)
	
	// Wrap the default ServeMux with our CORS middleware
	handler := middleware.CORSMiddleware(http.DefaultServeMux)
	
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
