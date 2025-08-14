# QWEN.md - Context for tg-admin-service

## Project Overview

This project, `tg-admin-service`, is a Go-based API service framework built using the Gin web framework and GORM for database interactions. It's designed for rapid API development and includes features like JWT authentication, Redis integration (for queues and caching), and a structured directory layout for maintainability.

Key Technologies:
- **Gin**: For HTTP routing and middleware.
- **GORM**: For ORM (Object-Relational Mapping) with MySQL.
- **Viper**: For configuration file parsing (`.ini` files).
- **JWT**: For secure user authentication.
- **Redis**: Used for implementing queues (Stream and ZSet for delayed queues) and potentially caching.
- **Uber FX**: For dependency injection.

## Directory Structure

Based on the `README.md` and directory listings, the structure is as follows:

```
├── config              // Project configuration files (dev.ini, prod.ini, deploy.ini)
├── deploy              // Deployment scripts and configurations
├── internal            // Core business logic
│   ├── config          // Internal configuration handling
│   ├── controller      // HTTP request handlers
│   ├── converter       // Data converters
│   ├── dto             // Data Transfer Objects
│   ├── error           // Custom error definitions
│   ├── middleware      // Gin middleware (e.g., CORS, JWT, Response)
│   ├── model           // Database models (GORM structs)
│   ├── provider        // Uber FX providers/modules
│   ├── query           // Database query logic (potentially GORM)
│   ├── request         // HTTP request validation structs
│   ├── router          // API route definitions
│   ├── service         // Business logic layer
│   └── vo              // Value Objects
├── log                 // Application runtime logs
├── main.go             // Application entry point
├── tools               // Utility functions and packages
│   ├── conv
│   ├── jwt            // JWT token generation and parsing
│   ├── key_utils
│   ├── logger         // Logging utilities
│   ├── queue          // Redis-based queue implementation
│   ├── random
│   ├── resp           // Standardized API response handling
│   └── utils.go
├── go.mod             // Go module dependencies
├── go.sum             // Go module checksums
├── README.md          // Project documentation
└── ...                // Other files like .gitignore, LICENSE
```

## Building and Running

### Running Locally

The primary way to run the application is via `go run`:

```bash
go run main.go -mode=dev
```

**Run Parameters:**
- `-mode=dev`: Runs the application using the `dev.ini` configuration file.
- `-mode=prod`: Runs the application using the `prod.ini` configuration file.
- `-initDb=true`: (Commented out in `main.go`) Initializes the database tables based on GORM model structs.

### Configuration

Configuration is managed via `.ini` files located in the `config/` directory.
Example `config/dev.ini`:
```ini
[mysql]
ip = rm-bp19d0k4o53v434g45o.mysql.rds.aliyuncs.com
port = 3306
username = tg
password = GiYLk9mbCu224nf1Vw
db_name = tg

[redis]
ip = r-bp1lj3mjfa1ph8me2hpd.redis.rds.aliyuncs.com
port = 6379
username = r-bp1lj3mjfa1ph8me2h
password = Bp1s0232xixf1zdin2
db = 0
max_total = 100

[server]
port = 8081
name = goingo
```

### Dependencies

Dependencies are managed using Go modules, defined in `go.mod`. Key dependencies include Gin, GORM, Viper, JWT, and Redis client.

## Development Conventions

- **Architecture**: Follows a layered architecture pattern (Controller -> Service -> Model/Query). Middleware handles cross-cutting concerns like authentication and logging.
- **Routing**: Routes are defined in `internal/router/` and grouped by functionality (e.g., Admin, User). The main router aggregates these.
- **Authentication**: Uses JWT for authentication, implemented as middleware (`internal/middleware`). Certain paths can be whitelisted from JWT checks.
- **Responses**: Standardized API responses are handled using structs and functions in the `tools/resp` package.
- **Configuration**: Uses Viper to load settings from `.ini` files based on the `-mode` flag.
- **Logging**: Custom logging is implemented in `tools/logger`.
- **Queues**: Redis-based queues (normal and delayed) are implemented in `tools/queue`.
- **Dependency Injection**: Uber FX is used for managing dependencies and application lifecycle in `main.go`.
