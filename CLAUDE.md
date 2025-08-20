# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**tg-admin-service** is a Go-based admin service API built with Gin + Gorm framework for rapid API development. It's a Chinese-language project focused on administrative functionality.

## Technology Stack

- **Framework**: Gin (HTTP routing), Gorm (ORM)
- **Database**: MySQL (via Gorm), Redis (for caching/queues)
- **Configuration**: Viper for config file management
- **Authentication**: JWT tokens for API authentication
- **Dependency Injection**: Uber FX for DI and lifecycle management
- **Logging**: Custom structured logging
- **Queue System**: Redis-based message and delay queues

## Architecture

### Layer Structure (Clean Architecture)
```
├── internal/
│   ├── config/           # Configuration management
│   ├── controller/       # HTTP handlers (presentation layer)
│   ├── converter/        # Data conversion utilities
│   ├── middleware/       # HTTP middleware (CORS, JWT, response formatting)
│   ├── model/           # Database models (domain layer)
│   ├── provider/        # Dependency injection providers
│   ├── query/           # Database query builders
│   ├── request/         # Request DTOs
│   ├── router/          # Route definitions
│   ├── service/         # Business logic layer
│   └── vo/             # View objects (response DTOs)
├── tools/               # Shared utilities
│   ├── jwt/            # JWT token utilities
│   ├── logger/         # Logging utilities
│   ├── resp/           # Response formatting
│   └── random/         # Random generation utilities
```

### Key Components

**Dependency Injection**: Uses Uber FX for DI with modules defined in `internal/provider/provider_set.go`

**Database Models**: Located in `internal/model/` with auto-migration support

**Middleware Stack**: 
- CORS handling
- JWT authentication (with whitelist for public endpoints)
- Response formatting
- Access logging

**Queue System**: Redis-based implementation with:
- Normal message queues (Redis Streams)
- Delay queues (Redis ZSET)
- Callback registration system

## Development Commands

### Running the Application
```bash
# Development mode
go run main.go -mode=dev

# Production mode  
go run main.go -mode=prod

# With database initialization
go run main.go -mode=dev -initDb=true
```

### Building
```bash
# Build for current platform
go build -o tg-admin-service

# Cross-platform builds (see deploy/deploy.sh)
cd deploy && ./deploy.sh tg-admin-service linux amd64 0
```

### Configuration
Configuration files are in `config/` directory:
- `dev.ini` - Development environment
- `prod.ini` - Production environment  
- `deploy.ini` - Deployment configuration

Key configuration sections:
- `[mysql]` - Database connection
- `[redis]` - Redis connection
- `[server]` - Server port and name

### Database Setup
The application uses Gorm with auto-migration. Models are auto-migrated on startup when `-initDb=true` is used.

## API Structure

### Authentication
- JWT-based authentication
- Whitelist endpoints bypass JWT: `/admin/login`, `/admin/initPwd`, `/api/index/health`

### Route Groups
- `/admin/*` - Admin management endpoints
- `/api/user/*` - User management endpoints  
- `/api/evaluate/*` - Evaluation-related endpoints
- `/api/index/health` - Health check endpoint

### Response Format
All API responses use standardized JSON format:
```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

## Testing

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./internal/service/...

# Run tests with coverage
go test -cover ./...
```

### Environment Setup
1. Ensure MySQL and Redis are running
2. Configure database credentials in `config/dev.ini`
3. Run with `-initDb=true` to initialize database schema

## Deployment

### Manual Deployment
```bash
# Build for production
cd deploy
./deploy.sh tg-admin-service linux amd64 0

# Deploy to server
tar -xzf tg-admin-service.tar.gz
nohup ./tg-admin-service -mode=prod > tg-admin-service.log 2>&1 &
```

### Docker (if needed)
Standard Go containerization patterns apply.

## Key Files to Know

- `main.go` - Application entry point with FX lifecycle
- `internal/provider/provider_set.go` - Dependency injection configuration
- `internal/router/router.go` - Main router setup and middleware
- `internal/config/config.go` - Configuration management
- `tools/logger/logger.go` - Logging utilities

## Common Development Tasks

### Adding New API Endpoints
1. Create request DTO in `internal/request/`
2. Add model in `internal/model/` (if needed)
3. Implement service in `internal/service/`
4. Create controller in `internal/controller/`
5. Add route in `internal/router/`
6. Register provider in `internal/provider/`

### Database Changes
1. Update model in `internal/model/`
2. Run with `-initDb=true` to auto-migrate
3. Update service layer as needed

### Queue Usage
See README.md for detailed queue implementation examples using Redis Streams and ZSET for delay queues.
- 默认用中文回复