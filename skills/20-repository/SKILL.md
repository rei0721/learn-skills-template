---
name: 数据访问规范
description: internal/repository 数据访问层的创建规范
category: 支撑
priority: 7
dependencies: [00-project-map, 02-design-specs, 13-database, 11-backend]
---

# 数据访问规范

## 职责范围

- 规范 Repository 接口与实现
- 统一数据库访问与事务使用方式
- 保持数据访问层与业务层隔离

## 适用场景

- 新增实体的数据访问
- 扩展已有仓库
- 重构查询实现

## 核心规则

### 1. 目录与文件结构

- 接口文件: `internal/repository/<domain>.go`
- 实现文件: `internal/repository/<domain>_impl.go`

### 2. 接口定义

```go
type XxxRepository interface {
    FindByID(ctx context.Context, id int64) (*models.DBXxx, error)
    Create(ctx context.Context, tx *gorm.DB, entity *models.DBXxx) error
    Update(ctx context.Context, tx *gorm.DB, entity *models.DBXxx) error
}
```

### 3. 实现规范

- 使用 gorm.DB 注入实现
- 查不到记录返回 nil, nil
- 写操作必须接收事务对象

```go
type xxxRepository struct {
    db *gorm.DB
}

func NewXxxRepository(db *gorm.DB) XxxRepository {
    return &xxxRepository{db: db}
}
```

## AI 行为约束

- 禁止在 Repository 中实现业务逻辑
- 禁止在 Repository 中返回未包装的业务错误
- 禁止在 Repository 中访问上层依赖

## 更新规则

- 新增或调整 Repository 的结构与接口
- 变更事务与查询约定

---

**最后更新**: 2026-01-19
**维护者**: AI Assistant
