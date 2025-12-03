# Repository Guidelines

## Project Structure & Module Organization
- `main.go` wires Fx modules and CLI flags `-mode`, `-initDb`, then starts the Gin server and job lifecycle.
- `internal/` holds layered code: `router` (routes/middleware), `controller` (handlers returning `tools/resp` helpers), `service` (business logic), `model` (GORM entities), `provider` (Telegram, bot registry, Redis, HTTP clients), `queue` (Redis Stream runners/workers), `job` (asynq tasks), `middleware`, and `dto/request/query/vo` for transport structs, plus shared `config` and `error`.
- Config files and SQL seed live in `config/*.ini` and `config/db.sql`; uploaded artifacts are under `file/`; shared utilities sit in `tools/` (logger, jwt, cron, resp, etc.).

## Build, Test, and Development Commands
- `go run main.go -mode=dev` boots the API with dev config; append `-initDb=true` to auto-migrate tables on startup.
- `go run main.go -mode=prod` runs with production settings; keep `config/prod.ini` in sync with deployed env.
- `go build -o tg-admin-service main.go` produces the deployable binary.
- `go test ./...` runs the full suite; use `go test ./internal/queue/...` for queue-focused checks.

## Coding Style & Naming Conventions
- Go 1.23+; always `gofmt -w` (tabs) and `goimports` before committing.
- Package names are lowercase and singular; exported types/functions use PascalCase, locals camelCase; errors start with lower-case messages.
- Keep handlers thin: router → controller → service → model; prefer dependency injection (Fx) over globals; use `tools/logger` for structured logs and `tools/resp` helpers for consistent API responses.
- Place DTOs/VOs near transport layers; keep database concerns inside `model` and cross-cutting logic in `internal/middleware`.

## Testing Guidelines
- Standard Go `testing` framework; current example lives in `tools/cron/cron_utils_test.go`.
- Name tests `*_test.go` beside the code; prefer table-driven tests for validation/helpers; use fakes over real Telegram/Redis unless explicitly gated.
- Run `go test ./...` before pushing; add coverage for queue scheduling, provider failover paths, and middleware changes when touching them.

## Commit & Pull Request Guidelines
- History favors short, scope-limited messages; keep an imperative one-liner (English or Chinese), e.g., `fix: handle redis reconnect`.
- Reference the issue/feature branch when relevant and avoid bundling unrelated changes.
- PRs should describe behavior changes, configs touched, and test evidence; include API samples or screenshots for handler updates.

## Configuration & Security Tips
- Do not commit secrets; keep `config/dev.ini` for local use and supply prod values via env/secret manager.
- Redis/MySQL URLs and Telegram tokens come from config; ensure `-mode` matches the intended config file.
- Generated artifacts like `tg-admin-service.pid` should stay out of commits; clean them before packaging releases.
