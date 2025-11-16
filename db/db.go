package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the connection string from environment variable
	connStr := os.Getenv("POSTGRESQL")
	if connStr == "" {
		log.Fatal("POSTGRESQL environment variable not set")
	}

	// Connect to the database
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Verify the connection
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
}
