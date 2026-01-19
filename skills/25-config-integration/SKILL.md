---
name: 配置集成流程
description: 配置项从定义到应用集成的完整流程
category: 支撑
priority: 7
dependencies: [00-project-map, 01-project-foundation, 15-documentation]
---

# 配置集成流程

## 职责范围

- 规范新增配置项的定义、验证与加载
- 统一配置变更的文档同步要求

## 适用场景

- 新增配置项
- 调整配置结构
- 引入新的基础设施配置

## 核心规则

1. 在 `internal/config` 中新增配置结构体与验证逻辑
2. 在 `internal/config/config.go` 中注册到主配置结构
3. 更新 `configs/config.example.yaml` 与本地配置样例
4. 配置接入应用容器并注入到使用方
5. 配置变更必须同步 docs 文档

## AI 行为约束

- 禁止在代码中硬编码配置值
- 禁止绕过配置验证直接使用外部输入
- 禁止更新配置但不更新文档

## 更新规则

- 新增配置集成步骤
- 变更配置验证策略

---

**最后更新**: 2026-01-19
**维护者**: AI Assistant
