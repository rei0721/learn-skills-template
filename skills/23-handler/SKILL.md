---
name: 处理器规范
description: internal/handler HTTP 处理器的创建规范
category: 支撑
priority: 7
dependencies: [00-project-map, 11-backend, 12-api]
---

# 处理器规范

## 职责范围

- 规范 HTTP Handler 的职责边界
- 统一请求解析与响应输出

## 适用场景

- 新增 API 接口
- 变更路由与响应结构

## 核心规则

1. Handler 只负责参数解析与调用 Service
2. 统一使用 types/result 输出响应
3. 错误必须转为 BizError 并映射 HTTP 状态码
4. 不允许直接访问 Repository

## AI 行为约束

- 禁止在 Handler 中实现业务逻辑
- 禁止自定义响应格式
- 禁止暴露内部错误详情

## 更新规则

- 新增或调整 Handler 结构
- 变更统一响应或错误处理方式

---

**最后更新**: 2026-01-19
**维护者**: AI Assistant
