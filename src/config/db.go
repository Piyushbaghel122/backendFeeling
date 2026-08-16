package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it (ignoring if using system env vars).")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Set defaults if environment variables are missing
	if dbUser == "" {
		dbUser = "root"
	}
	if dbPassword == "" {
		dbPassword = "password123"
	}
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbName == "" {
		dbName = "freelancing"
	}

	// Data Source Name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	log.Println("Connected to MySQL database successfully!")
}

