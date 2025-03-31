package main

import (
	"log"
	"os"

	"github.com/PabloPei/TreeSense-Backend/conf"
	"github.com/PabloPei/TreeSense-Backend/db"
	"github.com/PabloPei/TreeSense-Backend/internal/api"
)

func main() {

	log.SetOutput(os.Stdout)

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
