---
name: 更新
description: 记录项目任意变更，并同步到技能书与项目地图
category: 核心
priority: 10
dependencies: [00-project-map, 15-documentation]
---

# 更新

## 职责范围

- 记录本项目任意变更（功能、修复、重构、规范、依赖、配置、文档）
- 作为技能书与代码变更的同步入口
- 约束更新时必须补齐的“最低信息集”

## 适用场景

- 新增/修改/删除任何业务能力（含 API、Service、Repository）
- 影响认证、权限、事务、配置、数据库模型等高风险区域的改动
- 纯重构但影响对外行为、错误码、返回结构、路由等

## 核心规则

### 1. 最低信息集

每条更新必须包含：

- 变更类型（Added/Changed/Fixed/Refactor/Docs/Test）
- 变更范围（模块/路径）
- 行为变化（对外接口、请求/响应、默认值、权限）
- 风险点（安全/数据/兼容性）
- 验证方式（go test、lint、人工验证）

### 2. 同步顺序

当更新涉及新增模块/职责变化/依赖变化时，必须按顺序同步：

1. `skills/00-project-map/module-index.md`
2. `skills/00-project-map/risk-zones.md`（如果命中高风险）
3. 本文件的更新记录
4. 必要时同步到 docs（按 `15-documentation` 的触发条件）

## 更新记录

| 日期 | 类型 | 摘要 | 影响范围 | 验证 |
| ---- | ---- | ---- | -------- | ---- |
| 2026-01-19 | Added | 新增用户名+密码注册/登录（无邮箱手机号验证码），注册默认分配 RBAC 角色 | internal/service/auth, internal/handler/auth_handler.go, internal/router, types, internal/models | go test ./...; golangci-lint run ./... |
| 2026-01-19 | Docs | 启动命令统一改为 go run ./cmd/server <cmd> | docs/getting-started, docs/faq.md, docs/contributing, docs/development/configuration.md, docs/changelog.md | 人工检查文档 |
| 2026-01-19 | Fixed | 修复 RBAC 适配器自定义表名/前缀参数顺序导致初始化失败 | pkg/rbac/rbac_impl.go | go test ./...; golangci-lint run ./... |
| 2026-01-19 | Fixed | 修复 RBAC 适配器自定义表名使用错误参数导致初始化失败 | pkg/rbac/rbac_impl.go | go test ./...; golangci-lint run ./... |
| 2026-01-19 | Fixed | 统一 Casbin 与 Gorm Adapter 版本以修复 RBAC 初始化崩溃 | pkg/rbac/rbac_impl.go, pkg/rbac/README.md | go test ./...; golangci-lint run ./... |
| 2026-01-19 | Docs | 新增 internal/models 模型定义技能书 | skills/17-models/SKILL.md, skills/SKILLS.md, skills/README.md | 人工检查文档 |
| 2026-01-19 | Added | 新增技能补齐、业务服务、数据访问、可复用工具包、中间件、处理器、错误处理、配置集成与工作流库技能书 | skills/18-skill-gap/SKILL.md, skills/19-services/SKILL.md, skills/20-repository/SKILL.md, skills/21-pkg-reusables/SKILL.md, skills/22-middleware/SKILL.md, skills/23-handler/SKILL.md, skills/24-error-handling/SKILL.md, skills/25-config-integration/SKILL.md, skills/workflows/SKILL.md, skills/SKILLS.md, skills/README.md | 人工检查文档 |
| 2026-01-19 | Added | 新增项目目标规划技能书并设为必加载 | skills/27-project-goals/SKILL.md, skills/SKILLS.md, skills/README.md | 人工检查文档 |

