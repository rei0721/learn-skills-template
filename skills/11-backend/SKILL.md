---
name: 后端规范
description: Go Scaffold 后端开发的 Service、Repository 和中间件规范
category: 支撑
priority: 6
dependencies: [00-project-map, 02-design-specs, 03-rules, 04-reusables]
---

# 后端规范

## 职责范围

- Service 层开发规范
- Repository 层开发规范
- 中间件开发规范
- 业务逻辑组织方式

## 适用场景

- 开发新的业务功能
- 实现数据访问层
- 添加中间件
- 代码审查

## 核心规则

### 1. Service 层规范

#### 接口定义

```go
// internal/service/xxx/xxx.go
type XxxService interface {
    Create(ctx context.Context, req *CreateXxxRequest) (*XxxResponse, error)
    GetByID(ctx context.Context, id int64) (*XxxResponse, error)
    List(ctx context.Context, page, pageSize int) (*PaginationResult, error)
    Update(ctx context.Context, id int64, req *UpdateXxxRequest) error
    Delete(ctx context.Context, id int64) error
}
```

#### 实现规范

```go
// internal/service/xxx/xxx_impl.go
type xxxServiceImpl struct {
    repo       XxxRepository
    txMgr      dbtx.Manager
    logger     logger.Logger
    cache      cache.Cache
    idGen      utils.IDGenerator
}

func NewXxxService(repo XxxRepository, txMgr dbtx.Manager, logger logger.Logger) XxxService {
    return &xxxServiceImpl{
        repo:   repo,
        txMgr:  txMgr,
        logger: logger,
    }
}
```

#### 方法实现模板

```go
func (s *xxxServiceImpl) Create(ctx context.Context, req *CreateXxxRequest) (*XxxResponse, error) {
    // 1. 参数验证
    if err := s.validateRequest(req); err != nil {
        return nil, errors.ErrInvalidParams
    }
    
    // 2. 业务规则检查
    if err := s.checkBusinessRules(ctx, req); err != nil {
        return nil, err
    }
    
    // 3. 构建实体
    entity := &models.Xxx{
        ID:   s.idGen.Generate(),
        Name: req.Name,
        // ...
    }
    
    // 4. 使用事务（如果需要）
    err := s.txMgr.WithTransaction(ctx, func(ctx context.Context) error {
        if err := s.repo.Create(ctx, entity); err != nil {
            return err
        }
        // 其他数据库操作
        return nil
    })
    
    if err != nil {
        s.logger.Error("failed to create xxx", "error", err)
        return nil, err
    }
    
    // 5. 清除缓存（如果需要）
    if s.cache != nil {
        _ = s.cache.Delete(ctx, fmt.Sprintf("xxx:%d", entity.ID))
    }
    
    // 6. 返回响应
    return s.toResponse(entity), nil
}
```

### 2. Repository 层规范

#### 接口定义

```go
// internal/repository/xxx.go
type XxxRepository interface {
    Create(ctx context.Context, entity *models.Xxx) error
    FindByID(ctx context.Context, id int64) (*models.Xxx, error)
    FindAll(ctx context.Context, page, pageSize int) ([]*models.Xxx, int64, error)
    Update(ctx context.Context, entity *models.Xxx) error
    Delete(ctx context.Context, id int64) error
}
```

#### 实现规范

```go
// internal/repository/xxx_impl.go
type xxxRepositoryImpl struct {
    db     database.Database
    cache  cache.Cache
    logger logger.Logger
}

func NewXxxRepository(db database.Database, cache cache.Cache, logger logger.Logger) XxxRepository {
    return &xxxRepositoryImpl{
        db:     db,
        cache:  cache,
        logger: logger,
    }
}
```

#### 缓存策略

```go
func (r *xxxRepositoryImpl) FindByID(ctx context.Context, id int64) (*models.Xxx, error) {
    // 1. 尝试从缓存获取
    cacheKey := fmt.Sprintf("xxx:%d", id)
    if r.cache != nil {
        if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
            var entity models.Xxx
            if err := json.Unmarshal([]byte(cached), &entity); err == nil {
                return &entity, nil
            }
        }
    }
    
    // 2. 从数据库查询
    var entity models.Xxx
    err := r.db.GetDB().WithContext(ctx).First(&entity, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.ErrNotFound
        }
        return nil, err
    }
    
    // 3. 写入缓存
    if r.cache != nil {
        if data, err := json.Marshal(entity); err == nil {
            _ = r.cache.Set(ctx, cacheKey, data, 5*time.Minute)
        }
    }
    
    return &entity, nil
}
```

### 3. 中间件规范

#### 认证中间件

```go
// internal/middleware/auth.go
func AuthMiddleware(jwt jwt.JWT) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, result.Error(errors.ErrUnauthorized))
            c.Abort()
            return
        }
        
        claims, err := jwt.ParseToken(token)
        if err != nil {
            c.JSON(401, result.Error(errors.ErrUnauthorized))
            c.Abort()
            return
        }
        
        c.Set("user_id", claims["user_id"])
        c.Set("username", claims["username"])
        c.Next()
    }
}
```

#### 权限中间件

```go
// internal/middleware/rbac.go
func RBACMiddleware(rbac rbac.RBAC) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        path := c.Request.URL.Path
        method := c.Request.Method
        
        ok, err := rbac.Enforce(userID, path, method)
        if err != nil || !ok {
            c.JSON(403, result.Error(errors.ErrForbidden))
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## AI 行为约束

### 必须遵守

1. Service 层必须通过 Repository 访问数据
2. 多个写操作必须使用事务
3. 敏感操作必须添加日志
4. 缓存策略必须一致

### 禁止行为

1. 禁止在 Service 中直接操作数据库
2. 禁止在 Repository 中实现业务逻辑
3. 禁止绕过中间件
4. 禁止硬编码业务规则

---

**最后更新**: 2026-01-19
