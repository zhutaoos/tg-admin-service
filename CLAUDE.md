# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
tg-admin-service is a Go-based admin API service built with Gin + Gorm framework, focused on administrative functionality with Chinese language support.

## Technology Stack
- **Framework**: Gin (HTTP routing), Gorm (ORM)
- **Database**: MySQL (via Gorm), Redis (for caching/queues)
- **Configuration**: Viper for config file management
- **Authentication**: JWT tokens for API authentication
- **Dependency Injection**: Uber FX for DI and lifecycle management
- **Logging**: Custom structured logging
- **Queue System**: Redis-based message and delay queues with Asynq

## Architecture
Clean Architecture pattern with the following layers:
- `internal/controller/` - HTTP handlers
- `internal/service/` - Business logic
- `internal/model/` - Database models
- `internal/job/` - Background job processing (Asynq-based)
- `internal/middleware/` - HTTP middleware
- `internal/router/` - Route definitions
- `internal/config/` - Configuration management
- `tools/` - Shared utilities

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
# Usage: ./deploy.sh <project_name> <goos> <goarch> <cgo_enabled>
```

### Testing
```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./tools/cron

# Run tests with verbose output
go test -v ./...
```

### Testing Cron Tasks
The codebase includes built-in cron task testing via `JobService.ComprehensiveTest()` which:
- Validates cron expressions
- Tests immediate task execution
- Adds test cron tasks (1-minute intervals)
- Displays scheduler entries from Redis

### Database Setup
Models auto-migrate on startup when `-initDb=true` is used. Configuration files are in `config/` directory:
- `dev.ini` - Development environment
- `prod.ini` - Production environment
- `deploy.ini` - Deployment configuration

### Key Configuration Sections
- `[mysql]` - Database connection
- `[redis]` - Redis connection (for queues and caching)
- `[server]` - Server port and name

## Queue System Details
- **Task Types**: bot_msg, coupon:open, coupon:expire
- **Handler Registration**: Via `job.NewBotMsgHandler()` in provider.go
- **Cron Format**: 5-field format (`分 时 日 月 周`)
- **Timezone**: Asia/Shanghai
- **Validation**: Built-in cron expression validation

## Common API Endpoints
- `/api/task/create` - Create cron tasks
- `/api/admin/login` - Admin authentication  
- `/api/index/health` - Health check
- All endpoints return standardized JSON format with Chinese messages

## Important Notes
- **Module name**: The go.mod declares module as `app` (not the directory name)
- **Timezone**: All scheduling and logging uses Asia/Shanghai timezone
- **FX Dependency Injection**: Uses Uber FX for dependency injection and lifecycle management
- **Database**: Auto-migration occurs on startup with `-initDb=true` flag
- **Background Jobs**: Uses Asynq for Redis-based task queues and scheduling