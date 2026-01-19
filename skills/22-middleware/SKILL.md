---
name: 中间件规范
description: internal/middleware Gin 中间件的创建规范
category: 支撑
priority: 7
dependencies: [00-project-map, 11-backend, 12-api]
---

# 中间件规范

## 职责范围

- 规范 Gin 中间件的创建与组织
- 统一请求上下文与错误处理方式

## 适用场景

- 新增请求预处理逻辑
- 增加认证授权、日志、限流等能力

## 核心规则

1. 中间件必须放在 `internal/middleware/`
2. 必须返回 `gin.HandlerFunc`
3. 需要失败时统一返回标准响应格式
4. 必须使用 context 传递必要的请求信息

## AI 行为约束

- 禁止在中间件中实现业务逻辑
- 禁止绕过统一响应格式
- 禁止在中间件中访问 Repository

## 更新规则

- 新增中间件类型
- 调整中间件组织方式

---

**最后更新**: 2026-01-19
**维护者**: AI Assistant
