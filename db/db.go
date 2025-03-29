package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/PabloPei/SmartSpend-backend/conf"
	_ "github.com/lib/pq"
)

func NewPostgresStorage(cfg conf.PostgreSqlConfig) (*sql.DB, error) {

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBAddress, cfg.DBPort, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Error connecting to the database:", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging the database:", err)
		return nil, err
	}

	return db, nil
}
