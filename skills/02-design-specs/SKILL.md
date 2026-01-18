---
name: 设计规范
description: Go Scaffold 的架构风格、分层原则和设计模式
category: 核心
priority: 3
dependencies: [00-project-map, 01-project-foundation]
---

# 设计规范

## 职责范围

- 定义架构风格和设计模式
- 规范分层原则和依赖方向
- 定义接口设计规范
- 禁止的反模式

## 适用场景

- 设计新功能
- 重构现有代码
- 代码审查
- 架构决策

## 核心规则

### 1. 架构风格

#### 分层架构（Layered Architecture）

```
┌─────────────────────────────────────┐
│     Presentation Layer (表现层)     │
│  Handler, Middleware, Router        │
└────────────┬────────────────────────┘
             │ 依赖
             ↓
┌─────────────────────────────────────┐
│      Business Layer (业务层)        │
│  Service, Business Logic            │
└────────────┬────────────────────────┘
             │ 依赖
             ↓
┌─────────────────────────────────────┐
│   Data Access Layer (数据访问层)    │
│  Repository, Models                 │
└────────────┬────────────────────────┘
             │ 依赖
             ↓
┌─────────────────────────────────────┐
│  Infrastructure Layer (基础设施层)  │
│  Database, Cache, Logger, etc.      │
└─────────────────────────────────────┘
```

**核心原则**:
- 上层依赖下层，下层不依赖上层
- 每层只能访问相邻的下一层
- 通过接口隔离层与层之间的依赖

#### 依赖注入容器（Dependency Injection Container）

```go
// App 是依赖注入容器
type App struct {
    // 基础设施
    Config  *config.Config
    Logger  logger.Logger
    DB      database.Database
    Cache   cache.Cache
    
    // 业务组件
    Services    map[string]Service
    Repositories map[string]Repository
    
    // Web 组件
    Router     *gin.Engine
    HTTPServer httpserver.HTTPServer
}
```

**核心原则**:
- 所有组件通过容器创建和管理
- 组件之间通过接口依赖
- 支持延迟注入和热重载

### 2. 分层原则

#### 表现层（Presentation Layer）

**职责**:
- HTTP 请求处理
- 参数验证和绑定
- 响应格式化
- 错误处理

**规范**:
```go
// Handler 结构
type XxxHandler struct {
    service XxxService  // 依赖业务层
    logger  logger.Logger
}

// Handler 方法签名
func (h *XxxHandler) MethodName(c *gin.Context) {
    // 1. 参数绑定和验证
    var req XxxRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, result.Error(errors.InvalidParams))
        return
    }
    
    // 2. 调用业务层
    resp, err := h.service.MethodName(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, result.Error(err))
        return
    }
    
    // 3. 返回响应
    c.JSON(200, result.Success(resp))
}
```

**禁止**:
- ❌ 在 Handler 中实现业务逻辑
- ❌ 在 Handler 中直接访问数据库
- ❌ 在 Handler 中直接访问缓存

#### 业务层（Business Layer）

**职责**:
- 业务逻辑实现
- 业务规则验证
- 事务边界管理
- 服务间协调

**规范**:
```go
// Service 接口
type XxxService interface {
    MethodName(ctx context.Context, req XxxRequest) (*XxxResponse, error)
}

// Service 实现
type xxxServiceImpl struct {
    repo   XxxRepository  // 依赖数据访问层
    txMgr  dbtx.Manager
    logger logger.Logger
}

// Service 方法实现
func (s *xxxServiceImpl) MethodName(ctx context.Context, req XxxRequest) (*XxxResponse, error) {
    // 1. 业务规则验证
    if err := s.validateBusinessRules(req); err != nil {
        return nil, err
    }
    
    // 2. 使用事务（如果需要）
    var result *XxxResponse
    err := s.txMgr.WithTransaction(ctx, func(ctx context.Context) error {
        // 3. 调用数据访问层
        data, err := s.repo.FindByXxx(ctx, req.Xxx)
        if err != nil {
            return err
        }
        
        // 4. 业务逻辑处理
        result = s.processBusinessLogic(data)
        return nil
    })
    
    return result, err
}
```

**禁止**:
- ❌ 在 Service 中处理 HTTP 请求
- ❌ 在 Service 中直接操作数据库（使用 Repository）
- ❌ 在 Service 中硬编码配置

#### 数据访问层（Data Access Layer）

**职责**:
- 数据访问抽象
- 数据库操作封装
- 缓存策略实现
- 查询优化

**规范**:
```go
// Repository 接口
type XxxRepository interface {
    Create(ctx context.Context, entity *models.Xxx) error
    FindByID(ctx context.Context, id int64) (*models.Xxx, error)
    FindAll(ctx context.Context, page, pageSize int) ([]*models.Xxx, int64, error)
    Update(ctx context.Context, entity *models.Xxx) error
    Delete(ctx context.Context, id int64) error
}

// Repository 实现
type xxxRepositoryImpl struct {
    db     database.Database
    cache  cache.Cache
    logger logger.Logger
}

// Repository 方法实现
func (r *xxxRepositoryImpl) FindByID(ctx context.Context, id int64) (*models.Xxx, error) {
    // 1. 尝试从缓存获取
    if r.cache != nil {
        cacheKey := fmt.Sprintf("xxx:%d", id)
        if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
            var entity models.Xxx
            if err := json.Unmarshal([]byte(cached), &entity); err == nil {
                return &entity, nil
            }
        }
    }
    
    // 2. 从数据库查询
    var entity models.Xxx
    if err := r.db.GetDB().WithContext(ctx).First(&entity, id).Error; err != nil {
        return nil, err
    }
    
    // 3. 写入缓存
    if r.cache != nil {
        cacheKey := fmt.Sprintf("xxx:%d", id)
        if data, err := json.Marshal(entity); err == nil {
            _ = r.cache.Set(ctx, cacheKey, data, 5*time.Minute)
        }
    }
    
    return &entity, nil
}
```

**禁止**:
- ❌ 在 Repository 中实现业务逻辑
- ❌ 在 Repository 中处理事务（使用 Service 层）
- ❌ 在 Repository 中直接返回 GORM 错误（转换为业务错误）

### 3. 依赖方向

#### 依赖倒置原则（Dependency Inversion Principle）

```
高层模块不应该依赖低层模块，两者都应该依赖抽象
抽象不应该依赖细节，细节应该依赖抽象
```

**实现方式**:
```go
// 定义接口（抽象）
type UserRepository interface {
    FindByID(ctx context.Context, id int64) (*models.User, error)
}

// 高层模块依赖接口
type UserService struct {
    repo UserRepository  // 依赖抽象，不依赖具体实现
}

// 低层模块实现接口
type userRepositoryImpl struct {
    db database.Database
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id int64) (*models.User, error) {
    // 具体实现
}
```

#### 禁止的依赖方向

```
❌ Handler → Repository (跨层依赖)
❌ Service → Handler (反向依赖)
❌ Repository → Service (反向依赖)
❌ pkg → internal (反向依赖)
❌ types → pkg (反向依赖)
```

### 4. 接口设计规范

#### 接口定义

```go
// 接口命名：名词或动词+er
type Logger interface {
    Debug(msg string, fields ...interface{})
    Info(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
}

// 接口应该小而专注（Interface Segregation Principle）
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 而不是
type ReadWriter interface {
    Read(p []byte) (n int, err error)
    Write(p []byte) (n int, err error)
}
```

#### 接口实现

```go
// 实现命名：接口名+Impl 或 具体实现名
type LoggerImpl struct {
    // ...
}

// 或
type ZapLogger struct {
    // ...
}

// 确保实现了接口（编译时检查）
var _ Logger = (*ZapLogger)(nil)
```

### 5. 错误处理规范

#### 错误定义

```go
// 使用自定义错误类型
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func (e *AppError) Error() string {
    return e.Message
}

// 预定义错误
var (
    ErrNotFound      = &AppError{Code: 404, Message: "resource not found"}
    ErrUnauthorized  = &AppError{Code: 401, Message: "unauthorized"}
    ErrInvalidParams = &AppError{Code: 400, Message: "invalid parameters"}
)
```

#### 错误处理

```go
// 在边界处处理错误
func (h *Handler) Method(c *gin.Context) {
    result, err := h.service.Method(c.Request.Context())
    if err != nil {
        // 转换为 HTTP 错误
        if errors.Is(err, ErrNotFound) {
            c.JSON(404, result.Error(err))
            return
        }
        c.JSON(500, result.Error(err))
        return
    }
    c.JSON(200, result.Success(result))
}

// 在内部传递错误
func (s *Service) Method(ctx context.Context) (*Result, error) {
    data, err := s.repo.FindByID(ctx, id)
    if err != nil {
        // 包装错误，添加上下文
        return nil, fmt.Errorf("failed to find data: %w", err)
    }
    return data, nil
}
```

### 6. 禁止的反模式

#### 反模式 1: 上帝对象（God Object）

```go
// ❌ 错误：一个结构体包含所有功能
type App struct {
    // 包含所有业务逻辑
    CreateUser()
    UpdateUser()
    DeleteUser()
    CreateOrder()
    UpdateOrder()
    // ... 100+ 方法
}

// ✅ 正确：按职责分离
type UserService struct {
    CreateUser()
    UpdateUser()
    DeleteUser()
}

type OrderService struct {
    CreateOrder()
    UpdateOrder()
}
```

#### 反模式 2: 循环依赖（Circular Dependency）

```go
// ❌ 错误：A 依赖 B，B 依赖 A
type ServiceA struct {
    serviceB *ServiceB
}

type ServiceB struct {
    serviceA *ServiceA
}

// ✅ 正确：通过接口打破循环
type ServiceA struct {
    serviceB ServiceBInterface
}

type ServiceB struct {
    // 不依赖 ServiceA
}
```

#### 反模式 3: 全局变量（Global Variables）

```go
// ❌ 错误：使用全局变量
var globalDB *gorm.DB
var globalLogger *zap.Logger

// ✅ 正确：通过依赖注入
type Service struct {
    db     database.Database
    logger logger.Logger
}
```

#### 反模式 4: 硬编码配置（Hard-coded Configuration）

```go
// ❌ 错误：硬编码配置
func Connect() {
    db, _ := gorm.Open(mysql.Open("root:password@tcp(localhost:3306)/db"))
}

// ✅ 正确：从配置读取
func Connect(cfg *config.DatabaseConfig) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
    db, _ := gorm.Open(mysql.Open(dsn))
}
```

## AI 行为约束

### 必须遵守

1. **分层原则**: 严格遵守分层架构
2. **依赖方向**: 严格遵守依赖方向
3. **接口驱动**: 通过接口隔离依赖
4. **错误处理**: 使用统一的错误处理方式

### 禁止行为

1. **跨层依赖**: 禁止跨层直接依赖
2. **反向依赖**: 禁止下层依赖上层
3. **循环依赖**: 禁止循环依赖
4. **全局变量**: 禁止使用全局变量（除常量外）

### 违规处理

- 检测到跨层依赖 → **直接拒绝执行**
- 检测到反向依赖 → **直接拒绝执行**
- 检测到循环依赖 → **直接拒绝执行**
- 检测到反模式 → **中断并提示修正**

## 更新规则

### 必须更新的情况

1. 架构风格变化
2. 分层原则变化
3. 新增设计模式
4. 新增反模式

### 更新流程

1. 更新本文件
2. 更新相关示例代码
3. 通知相关技能书同步更新
4. 更新本文件的修改日期

---

**最后更新**: 2026-01-19
**维护者**: AI Assistant
