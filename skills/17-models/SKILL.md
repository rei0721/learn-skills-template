---
name: 模型规范
description: internal/models 的模型定义与 GORM 映射规范
category: 支撑
priority: 7
dependencies: [00-project-map, 13-database]
---

# 模型规范

## 职责范围

- 统一 internal/models 的模型定义方式
- 规范表名、字段、索引与关联关系
- 约束敏感字段的序列化规则

## 适用场景

- 新增或修改 internal/models 的数据模型
- 为模型添加索引、唯一约束或关联
- 设计数据库结构与模型映射

## 核心规则

### 1. 基础字段

- 业务模型必须嵌入 `BaseDBModel`
- 主键类型统一为 `int64`
- 使用 `gorm.DeletedAt` 启用软删除

### 2. 表名规范

- 每个模型必须实现 `TableName() string`
- 表名统一使用复数下划线命名
- 禁止依赖 GORM 默认推断

### 3. 字段与索引

- 所有字符串字段必须设置长度 `size:<N>`
- 需要唯一约束的字段使用 `uniqueIndex`
- 需要查询优化的字段使用 `index`
- 复合唯一索引使用命名索引

### 4. 关联关系

- 多对多使用 `many2many:<join_table>`
- 外键关系显式声明 `foreignKey` 与 `references`
- 关联字段默认不暴露到 JSON

### 5. 序列化安全

- 敏感字段必须 `json:"-"`
- API 响应必须通过 types 层 DTO 返回

## 示例

```go
type DBUser struct {
    BaseDBModel
    Username string  `gorm:"uniqueIndex;size:50;not null" json:"username"`
    Email    *string `gorm:"uniqueIndex;size:100" json:"email,omitempty"`
    Password string  `gorm:"size:255;not null" json:"-"`
    Status   int     `gorm:"default:1" json:"status"`
}

func (DBUser) TableName() string { return "users" }
```

## AI 行为约束

- 新增模型时必须嵌入 `BaseDBModel` 并实现 `TableName()`
- 变更字段或索引时必须同步数据库迁移
- 对外响应必须使用 DTO，禁止直接返回 DB 模型

## 更新规则

- 新增模型 → 更新 `13-database` 与 `16-updates`
- 修改字段/索引 → 提供迁移与回滚方案
- 规则变更 → 同步到技能书索引
