package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/PabloPei/TreeSense-Backend/conf"
	"github.com/PabloPei/TreeSense-Backend/internal/middlewares"
	"github.com/PabloPei/TreeSense-Backend/internal/users"
	"github.com/PabloPei/TreeSense-Backend/internal/roles"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(cfg conf.ApiServerConfig, db *sql.DB) *APIServer {
	return &APIServer{
		addr: fmt.Sprintf("%s:%s", cfg.PublicHost, cfg.Port),
		db:   db, 
	}
}

func (s *APIServer) Run() error {

	router := mux.NewRouter()

	// global middlewares
	router.Use(middlewares.LoggingMiddleware)
	router.Use(middlewares.RecoveryMiddleware)

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// role routes
	roleRepository := roles.NewSQLRepository(s.db)
	roleService := roles.NewService(roleRepository)
	roleHandler := roles.NewHandler(roleService)
	roleHandler.RegisterRoutes(subrouter)

	// user routes
	userRepository := users.NewSQLRepository(s.db)
	userService := users.NewService(userRepository, roleRepository)
	userHandler := users.NewHandler(userService)
	userHandler.RegisterRoutes(subrouter)

	

	log.Println("Server running on", s.addr)
	return http.ListenAndServe(s.addr, router)

}
