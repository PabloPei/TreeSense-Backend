package main

import (
	"log"

	"github.com/PabloPei/SmartSpend-backend/conf"
	"github.com/PabloPei/SmartSpend-backend/db"
	"github.com/PabloPei/SmartSpend-backend/internal/api"
)

func main() {

	// PSQL Connection //

	log.Println("Starting PostgreSQL connection...")

	db, err := db.NewPostgresStorage(conf.DatabaseConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to the database")

	// API Server //

	log.Println("Starting Api Server...")

	server := api.NewAPIServer(conf.ServerConfig, db)
	err = server.Run()

	log.Fatal("Server Crash:", err)

}
