# Repository Guidelines

本指南面向贡献者，概述项目结构、开发调试方式、代码风格、测试与协作规范。阅读后即可快速上手。

## 项目结构与模块组织
- `main.go`：应用入口（Fx 生命周期、flags）。
- `config/`：INI 配置（`dev.ini`、`prod.ini`、`deploy.ini`）与 `db.sql`。
- `internal/`：业务代码目录（`router/`、`controller/`、`service/`、`model/`、`middleware/`、`provider/`、`request/`、`vo/`、`job/`、`config/`）。
- `tools/`：通用工具（`logger`、`jwt`、`http_client`、`conv`、`random`、`resp`、`cron`、`key_utils`）。
- `deploy/`：打包与运行脚本（`deploy.sh`、`run.sh`）。
- `log/`：运行/访问日志（已忽略），PID：`tg-admin-service.pid`。
- 测试：与源码同目录，命名 `*_test.go`。

## 构建、测试与本地开发
- 开发运行：`go run main.go -mode=dev`（加载 `config/dev.ini`）。
- 本地生产：`go run main.go -mode=prod`（加载 `config/prod.ini`）。
- 构建二进制：`go build -o app .`，运行：`./app -mode=prod`。
- 打包（Linux 示例）：`cd deploy && ./deploy.sh tg-admin-service linux amd64 0`。
- 启动打包产物：`cd deploy && ./run.sh tg-admin-service`。
- 查看日志：`tail -f log/access.log`、`log/tg-admin-service.log`（如启用）。

## 代码风格与命名约定
- 格式化：提交前执行 `go fmt ./...`（标准 Go 风格，Tab 缩进）。
- 包：短小、小写；文件名按需使用下划线（snake_case）。
- 导出使用 PascalCase；内部使用 camelCase。
- 错误：作为最后一个返回值，变量名 `err`，并包裹上下文信息。
- HTTP：路由/控制器在 `router/`、`controller/`；业务在 `service/`；持久化在 `model/`。

## 测试指南
- 框架：Go `testing`；优先表驱动测试。
- 命名：`TestXxx(t *testing.T)`；与被测文件同目录 `*_test.go`。
- 运行：`go test ./... -v`；覆盖率：`go test ./... -cover`。
- 建议优先覆盖 `service/` 与 `provider/`；外部依赖可用 fake/容器（若后续引入）。

## 提交与 Pull Request
- Commit：简洁祈使句，可加范围前缀（例：`router: add task routes`、`fix: token refresh`）；中/英文均可但需一致。
- PR 必须包含：目的、主要变更、运行/测试方法、配置影响、相关 issue；新增接口请附 curl 示例或截图。
- 变更应聚焦单一主题，涉及行为修改需同步更新文档/注释。

## 安全与配置提示
- 切勿提交任何密钥/凭据。开发使用 `config/dev.ini`；生产通过 `-mode=prod` 读取 `config/prod.ini`。
- 输入校验在控制器层；CORS/JWT 交由中间件处理。新增公开路由请更新 `internal/router/router.go` 的白名单。
- 变更部署脚本或日志路径请同步 `deploy/` 与说明文档。

