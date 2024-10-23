package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type Database struct {
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
}

// getConnectionString - formats the database connection string
func getConnectionString(d Database) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.DB_HOST, d.DB_PORT, d.DB_USER, d.DB_PASSWORD, d.DB_NAME)
}

// Connection - open database connection
func Connection() error {
	config := Database{
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
	}
	connectionString := getConnectionString(config)
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("error open connection: %w", err)
	}
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connection to database: %w", err)
	}
	log.Println("Successful connection to database")
	return nil
}

// Close - closes the database connection
func Close() {
	if err := DB.Close(); err != nil {
		log.Printf("Ошибка при закрытии соединения с базой данных: %v", err)
	} else {
		log.Println("Соединение с базой данных успешно закрыто.")
	}
}
