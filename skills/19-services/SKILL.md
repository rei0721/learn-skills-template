---
name: 业务服务规范
description: internal/service 新业务服务的创建规范
category: 支撑
priority: 7
dependencies: [00-project-map, 02-design-specs, 11-backend, 04-reusables]
---

# 业务服务规范

## 职责范围

- 定义业务服务在 internal/service 的组织方式
- 规范 Service 接口、实现与依赖注入
- 统一事务与错误处理方式

## 适用场景

- 新增业务能力
- 扩展已有服务
- 重构服务实现

## 核心规则

### 1. 目录与文件结构

- 服务目录: `internal/service/<domain>/`
- 接口文件: `<domain>.go`
- 实现文件: `<domain>_impl.go`
- 常量文件: `constants.go`

### 2. 接口定义

```go
type XxxService interface {
    Create(ctx context.Context, req *types.CreateXxxRequest) (*types.XxxResponse, error)
    GetByID(ctx context.Context, id int64) (*types.XxxResponse, error)
    List(ctx context.Context, page, pageSize int) (*types.PageResult, error)
    Update(ctx context.Context, id int64, req *types.UpdateXxxRequest) error
    Delete(ctx context.Context, id int64) error
}
```

### 3. 实现规范

- 采用 BaseService 统一管理依赖
- 只通过 Repository 访问数据
- 多写操作必须使用事务管理器

```go
type xxxService struct {
    service.BaseService[repository.XxxRepository]
}

func NewXxxService(repo repository.XxxRepository) XxxService {
    s := &xxxService{}
    s.BaseService.SetRepository(repo)
    return s
}
```

### 4. 错误处理

- 使用 types/errors 中的 BizError 作为对外错误
- 仅在 Service 内部包装底层错误

## AI 行为约束

- 禁止在 Service 中直接操作数据库
- 禁止在 Service 中返回原始数据库错误
- 禁止绕过事务管理器执行多写操作

## 更新规则

- 新增或调整 Service 的组织方式
- 变更依赖注入或事务规范

---

**最后更新**: 2026-01-19
**维护者**: AI Assistant
