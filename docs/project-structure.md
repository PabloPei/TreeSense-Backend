# Project Structure 

Domain-Driven Design (DDD), the application is divided into domains or bounded contexts, where each domain owns its own layers, including models, repositories, and services.

/SmartSpend
│── /cmd # Entry point of the application (main.go)
│── /config # Configuration files (environment variables, etc.)
│── /internal # Internal backend code (not accessible from other modules)
│ │── /domain # domain 
│ │   |── /domain_service.go # Business logic
│ │   │── /domain_repository.go # Database access
│ │   │── /domain_handlers.go # HTTP controllers (handle requests)
│ │   │── /domain_test.go # Unit test for the domain
│ │── /models # domain 
│ │   │── /domain.go # Models and interfaces for the domain
│ │── /middlewares # Middlewares (authentication, logging, etc.)
│── /pkg # Reusable code (can be used by other projects)
│── /db # Configuration, access and migrations for DB
│ │── /migrations # SQL Scripts for migrations and change in db
│── /scripts # Useful scripts (e.g., initialize data)
│── /test # E2E and integration tests
│── .env # Environment variables (do not upload to git)
│── go.mod # Go module file
│── go.sum # Dependency checksums
│── Makefile # Automation commands
│── README.md # Project documentation
