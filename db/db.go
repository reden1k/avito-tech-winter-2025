package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func waitForDB() {
	for {
		err := DB.Ping()
		if err == nil {
			log.Println("База данных готова!")
			return
		}
		log.Printf("Не удалось подключиться к базе данных: %v. Повторная попытка через 5 секунд...", err)
		time.Sleep(5 * time.Second)
	}
}

func InitDB() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	DB, err = sql.Open("postgres", dsn)

	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(time.Minute * 5)

	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	waitForDB()

	log.Println("Успешное подключение к БД")
}
