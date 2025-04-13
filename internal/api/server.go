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
	"github.com/PabloPei/TreeSense-Backend/internal/trees"
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

	// Global middlewares
	router.Use(middlewares.LoggingMiddleware)
	router.Use(middlewares.RecoveryMiddleware)
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// Repositories
	userRepository := users.NewSQLRepository(s.db)
	roleRepository := roles.NewSQLRepository(s.db)
	treeRepository := trees.NewSQLRepository(s.db)

	// Servicies
	userService := users.NewService(userRepository)
	roleService := roles.NewService(roleRepository, userRepository)
	treeService := trees.NewService(treeRepository)

	// Local Middlewares
	authMiddleware := middlewares.NewAuthMiddleware(roleService, userService)

	// Routes
	userHandler := users.NewHandler(userService)
	userHandler.RegisterRoutes(subrouter, authMiddleware)
	roleHandler := roles.NewHandler(roleService)
	roleHandler.RegisterRoutes(subrouter, authMiddleware)
	treeHandler := trees.NewHandler(treeService)
	treeHandler.RegisterRoutes(subrouter, authMiddleware)

	log.Println("Server running on", s.addr)
	return http.ListenAndServe(s.addr, router)

}
