package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/PabloPei/TreeSense-Backend/conf"
	"github.com/PabloPei/TreeSense-Backend/internal/audit"
	"github.com/PabloPei/TreeSense-Backend/internal/middlewares"
	"github.com/PabloPei/TreeSense-Backend/internal/permission"
	"github.com/PabloPei/TreeSense-Backend/internal/roles"
	"github.com/PabloPei/TreeSense-Backend/internal/trees"
	"github.com/PabloPei/TreeSense-Backend/internal/users"
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
	router.Methods(http.MethodOptions).Handler(middlewares.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))) //TODO: revisar este cambio. Hace que las rutas con OPTIONS devuelvan 204 en vez de 404 y permitan que se ejecute la ruta que de verdad quiero.
	// El navegador envia una preflight request (OPTIONS...) como parte del CORS. Como no encontraba la ruta, bloqueaba la request principal.
	// Esto registra una ruta OPTIONS universal para que las OPTIONS.... den 204 y luego se ejecute la request principal.

	// Global middlewares
	router.Use(middlewares.CORSMiddleware) //TODO: Chequear si dejar en prod
	router.Use(middlewares.LoggingMiddleware)
	router.Use(middlewares.RecoveryMiddleware)
	api := router.PathPrefix("/api/v1").Subrouter()

	// Repositories
	userRepository := users.NewSQLRepository(s.db)
	roleRepository := roles.NewSQLRepository(s.db)
	treeRepository := trees.NewSQLRepository(s.db)
	auditRepository := audit.NewSQLRepository(s.db)
	permissionRepository := permission.NewSQLRepository(s.db)

	// Services
	userService := users.NewService(userRepository)
	roleService := roles.NewService(roleRepository, userRepository)
	treeService := trees.NewService(treeRepository)
	permissionService := permission.NewService(permissionRepository, userRepository)
	auditService := audit.NewService(auditRepository)

	// Middlewares
	authMiddleware := middlewares.NewAuthMiddleware(permissionService, userService, auditService)
	auditMiddleware := middlewares.NewAuditMiddleware()

	/// Subrouters

	// without audit

	treeRouter := api.PathPrefix("/tree").Subrouter()
	treeHandler := trees.NewHandler(treeService)
	treeHandler.RegisterRoutes(treeRouter, authMiddleware)

	// with audit
	userRouter := api.PathPrefix("/user").Subrouter()
	userHandler := users.NewHandler(userService)
	userHandler.RegisterRoutes(userRouter, authMiddleware)
	userRouter.Use(auditMiddleware)

	roleRouter := api.PathPrefix("/role").Subrouter()
	roleHandler := roles.NewHandler(roleService)
	roleHandler.RegisterRoutes(roleRouter, authMiddleware)
	roleRouter.Use(auditMiddleware)

	permissionRouter := api.PathPrefix("/permission").Subrouter()
	permissionRouter.Use(auditMiddleware)
	permissionHandler := permission.NewHandler(permissionService)
	permissionHandler.RegisterRoutes(permissionRouter, authMiddleware)

	log.Println("Server running on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
