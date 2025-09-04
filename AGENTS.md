# Repository Guidelines

## Project Structure & Module Organization
- `main.go`: Application entrypoint (Fx lifecycle, flags).
- `config/`: INI configs (`dev.ini`, `prod.ini`, `deploy.ini`) and `db.sql`.
- `internal/`: Business code
  - `router/`, `controller/`, `service/`, `model/`, `middleware/`, `provider/`, `request/`, `vo/`, `job/`, `config/`.
- `tools/`: Shared utilities (logger, jwt, http_client, conv, random, resp, cron, key_utils).
- `deploy/`: Packaging and runtime scripts (`deploy.sh`, `run.sh`).
- `log/`: Runtime/access logs (git-ignored). PID file: `tg-admin-service.pid`.

## Build, Test, and Development Commands
- Run dev: `go run main.go -mode=dev`
- Run prod locally: `go run main.go -mode=prod`
- Build binary: `go build -o app .` then run `./app -mode=prod`
- Package (Linux example): `cd deploy && ./deploy.sh tg-admin-service linux amd64 0`
- Start packaged app on server: `cd deploy && ./run.sh tg-admin-service`
- Logs: tail `log/access.log` and `tg-admin-service.log` (if configured by logger).

## Coding Style & Naming Conventions
- Formatting: `go fmt ./...` before pushing. Tabs/standard Go style.
- Packages: short, lower-case; files use snake_case when needed.
- Exported API: `PascalCase`; internal/private: `camelCase`.
- Errors: return `error` as last value; wrap with context; use `err`.
- HTTP handlers live under `router/` and `controller/`; business logic in `service/`; persistence in `model/`.

## Testing Guidelines
- Framework: Go `testing`. Place tests alongside code as `*_test.go`.
- Naming: `TestXxx(t *testing.T)`; table-driven where helpful.
- Run all: `go test ./... -v` (coverage: `go test ./... -cover`).
- Prefer testing services and providers (DB/redis can use fakes or containers if added later).

## Commit & Pull Request Guidelines
- Commits: concise, imperative; scope first if useful (e.g., `router: add task routes`, `fix: token refresh`). Chinese or English OK—be consistent.
- PRs: include purpose, major changes, how to run/test, config impacts, and any screenshots/curl examples for new endpoints. Link related issues.
- Keep diffs focused; update docs/comments when changing behavior.

## Security & Configuration Tips
- Never commit secrets. Use `config/dev.ini` locally; `-mode=prod` loads `config/prod.ini`.
- Validate inputs in controllers; rely on middleware for CORS/JWT. Review `whitelist` in `internal/router/router.go` when adding public routes.

# Response
- Reply in Chinese by default