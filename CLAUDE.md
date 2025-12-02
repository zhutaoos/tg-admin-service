# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
# Run in development mode
go run main.go -mode=dev

# Run in production mode
go run main.go -mode=prod

# Initialize database tables on startup
go run main.go -mode=dev -initDb=true

# Build the application
go build -o tg-admin-service main.go

# Run tests
go test ./...

# Run a single package test
go test ./internal/queue/...
```

## Architecture Overview

This is a Telegram bot administration service built with a layered architecture using Uber FX for dependency injection.

**Request Flow:**
```
HTTP Request → Gin Router → Middleware (CORS, Auth, Response) → Controller → Service → Model/GORM → MySQL/Redis
```

**FX Module Structure (main.go):**
- `InfrastructureModule`: Database, Redis, queue system, services
- `ControllerModule`: HTTP controllers
- `RouterModule`: Route definitions

## Key Packages

| Package | Purpose |
|---------|---------|
| `internal/queue/` | Redis Stream-based message queue with Producer, Worker, RunnerManager, FailureTracker, Limiter |
| `internal/provider/telegram/` | Telegram API client abstraction |
| `internal/provider/botregistry/` | Bot selection and failover management |
| `internal/job/` | Async task handlers using hibiken/asynq |
| `internal/model/` | GORM models inheriting from `MysqlBaseModel` |
| `tools/logger/` | Structured logging with date-based rotation |

## Configuration

- Config files: `config/{dev,prod,deploy}.ini`
- Uses Viper for INI parsing
- Mode selection via `-mode` flag
- Global access: `config.Db()`, `config.Get[T]()`

## Queue System Design

The Redis Stream queue (`internal/queue/`) implements:
- **Per-chat message ordering**: One runner per chat ID
- **Dynamic worker scaling**: RunnerManager creates/destroys workers based on load
- **Bot failover**: FailureTracker with exponential backoff
- **Rate limiting**: Configurable per-bot rate limits
- **Backpressure**: Queue depth monitoring to prevent overload

## API Routes

- `/api/admin/*` - Admin authentication
- `/api/user/*`, `/api/group/*`, `/api/bot/*` - Entity management
- `/api/message/*`, `/api/task/*` - Message and task operations
- `/api/file/*` - File operations (public)
- `/api/index/health` - Health check

## Model Pattern

All models implement `BaseModel` interface with:
- `ToMap()` - Convert to map
- `InitWithMap()` - Initialize from map
- `CreateTable()` - Auto-migrate table
