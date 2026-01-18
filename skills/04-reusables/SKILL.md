---
name: 复用表
description: Go Scaffold 的公共工具、类型和扩展规则
category: 核心
priority: 5
dependencies: [00-project-map]
---

# 复用表

## 职责范围

- 列出所有可复用的公共工具
- 列出所有可复用的公共类型
- 定义扩展规则和最佳实践
- 禁止重复实现

## 适用场景

- 实现新功能前查找可复用组件
- 添加新工具前检查是否已存在
- 扩展现有组件
- 代码审查

## 核心规则

### 1. 公共工具（pkg/utils）

#### 1.1 ID 生成器（Snowflake）

**位置**: `pkg/utils/snowflake.go`

**接口**:
```go
type IDGenerator interface {
    Generate() int64
}
```

**使用示例**:
```go
// 通过 App 容器获取
id := app.IDGenerator.Generate()

// 用于模型 ID
user.ID = app.IDGenerator.Generate()
```

**禁止**:
- ❌ 自己实现 ID 生成逻辑
- ❌ 使用数据库自增 ID（分布式环境下有问题）
- ❌ 使用 UUID（性能较差，不适合作为主键）

#### 1.2 IP 工具

**位置**: `pkg/utils/ip.go`

**功能**:
- 获取客户端真实 IP
- 处理代理和负载均衡

**使用示例**:
```go
clientIP := utils.GetClientIP(c.Request)
```

#### 1.3 国际化工具

**位置**: `pkg/utils/i18n.go`

**功能**:
- 从上下文获取语言
- 翻译消息

**使用示例**:
```go
lang := utils.GetLangFromContext(ctx)
msg := utils.Translate(lang, "error.not_found")
```

### 2. 公共类型（types/）

#### 2.1 常量定义（types/constants）

**应用常量**: `types/constants/app.go`
```go
const (
    AppName    = "go-scaffold"
    AppVersion = "0.1.2"
)
```

**缓存常量**: `types/constants/cache.go`
```go
const (
    CacheKeyUser    = "user:%d"
    CacheKeySession = "session:%s"
)
```

**执行器常量**: `types/constants/executor.go`
```go
const (
    PoolNameHTTP       = "http"
    PoolNameDatabase   = "database"
    PoolNameCache      = "cache"
    PoolNameLogger     = "logger"
    PoolNameBackground = "background"
)
```

**使用规则**:
- ✅ 使用预定义常量
- ❌ 硬编码字符串或数字

#### 2.2 错误定义（types/errors）

**错误码**: `types/errors/codes.go`
```go
const (
    CodeSuccess         = 0
    CodeInvalidParams   = 400
    CodeUnauthorized    = 401
    CodeForbidden       = 403
    CodeNotFound        = 404
    CodeInternalError   = 500
)
```

**错误类型**: `types/errors/error.go`
```go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

var (
    ErrInvalidParams = &AppError{Code: CodeInvalidParams, Message: "invalid parameters"}
    ErrUnauthorized  = &AppError{Code: CodeUnauthorized, Message: "unauthorized"}
    ErrNotFound      = &AppError{Code: CodeNotFound, Message: "resource not found"}
)
```

**使用规则**:
- ✅ 使用预定义错误
- ✅ 创建新错误时遵循相同格式
- ❌ 直接返回字符串错误

#### 2.3 响应结果（types/result）

**统一响应**: `types/result/result.go`
```go
type Result struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) *Result {
    return &Result{
        Code:    0,
        Message: "success",
        Data:    data,
    }
}

func Error(err error) *Result {
    // ...
}
```

**分页结果**: `types/result/pagination.go`
```go
type PaginationResult struct {
    List       interface{} `json:"list"`
    Total      int64       `json:"total"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    TotalPages int         `json:"total_pages"`
}
```

**使用规则**:
- ✅ 所有 API 响应使用统一格式
- ✅ 分页查询使用 PaginationResult
- ❌ 自定义响应格式

#### 2.4 公共接口（types/interfaces.go）

**加密接口**:
```go
type Crypto interface {
    HashPassword(password string) (string, error)
    ComparePassword(hashedPassword, password string) error
}
```

**使用规则**:
- ✅ 依赖接口而不是具体实现
- ✅ 新增接口时添加到此文件
- ❌ 在业务代码中定义基础接口

### 3. 公共组件（pkg/）

#### 3.1 数据库（pkg/database）

**接口**:
```go
type Database interface {
    GetDB() *gorm.DB
    Close() error
}
```

**使用示例**:
```go
// 通过 App 容器获取
db := app.DB.GetDB()

// 在 Repository 中使用
func (r *repo) FindByID(ctx context.Context, id int64) (*Model, error) {
    var model Model
    err := r.db.GetDB().WithContext(ctx).First(&model, id).Error
    return &model, err
}
```

**禁止**:
- ❌ 直接创建数据库连接
- ❌ 绕过 Database 接口

#### 3.2 缓存（pkg/cache）

**接口**:
```go
type Cache interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Delete(ctx context.Context, key string) error
    Close() error
}
```

**使用示例**:
```go
// 设置缓存
err := app.Cache.Set(ctx, "user:123", userData, 5*time.Minute)

// 获取缓存
data, err := app.Cache.Get(ctx, "user:123")
```

**禁止**:
- ❌ 直接创建 Redis 连接
- ❌ 绕过 Cache 接口

#### 3.3 日志（pkg/logger）

**接口**:
```go
type Logger interface {
    Debug(msg string, fields ...interface{})
    Info(msg string, fields ...interface{})
    Warn(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
    Sync() error
}
```

**使用示例**:
```go
// 记录日志
app.Logger.Info("user created", "id", user.ID, "username", user.Username)
app.Logger.Error("failed to create user", "error", err)
```

**禁止**:
- ❌ 使用 fmt.Println 或 log.Println
- ❌ 直接创建 Logger 实例

#### 3.4 JWT（pkg/jwt）

**接口**:
```go
type JWT interface {
    GenerateToken(claims map[string]interface{}) (string, error)
    ParseToken(token string) (map[string]interface{}, error)
}
```

**使用示例**:
```go
// 生成 token
token, err := app.JWT.GenerateToken(map[string]interface{}{
    "user_id": user.ID,
    "username": user.Username,
})

// 解析 token
claims, err := app.JWT.ParseToken(token)
```

**禁止**:
- ❌ 自己实现 JWT 生成和解析
- ❌ 硬编码 JWT 密钥

#### 3.5 RBAC（pkg/rbac）

**接口**:
```go
type RBAC interface {
    Enforce(sub, obj, act string) (bool, error)
    AddPolicy(sub, obj, act string) (bool, error)
    RemovePolicy(sub, obj, act string) (bool, error)
    Close()
}
```

**使用示例**:
```go
// 检查权限
ok, err := app.RBAC.Enforce("user:123", "/api/users", "GET")
if !ok {
    return ErrForbidden
}

// 添加权限
app.RBAC.AddPolicy("role:admin", "/api/users", "POST")
```

**禁止**:
- ❌ 自己实现权限检查逻辑
- ❌ 硬编码权限规则

#### 3.6 加密（pkg/crypto）

**接口**:
```go
type Crypto interface {
    HashPassword(password string) (string, error)
    ComparePassword(hashedPassword, password string) error
}
```

**使用示例**:
```go
// 加密密码
hashedPassword, err := app.Crypto.HashPassword(password)

// 验证密码
err := app.Crypto.ComparePassword(user.Password, password)
```

**禁止**:
- ❌ 使用 MD5 或 SHA1 加密密码
- ❌ 明文存储密码

#### 3.7 事务管理（pkg/dbtx）

**接口**:
```go
type Manager interface {
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
```

**使用示例**:
```go
err := app.DBTx.WithTransaction(ctx, func(ctx context.Context) error {
    if err := repo.Create(ctx, user); err != nil {
        return err
    }
    if err := repo.CreateProfile(ctx, profile); err != nil {
        return err
    }
    return nil
})
```

**禁止**:
- ❌ 手动管理事务（Begin, Commit, Rollback）
- ❌ 在事务外执行多个写操作

#### 3.8 协程池（pkg/executor）

**接口**:
```go
type Manager interface {
    Submit(poolName string, task func()) error
    SubmitWithTimeout(poolName string, task func(), timeout time.Duration) error
    Shutdown()
}
```

**使用示例**:
```go
// 提交任务到 HTTP 池
err := app.Executor.Submit(constants.PoolNameHTTP, func() {
    // 异步任务
})

// 提交任务到后台池（带超时）
err := app.Executor.SubmitWithTimeout(constants.PoolNameBackground, func() {
    // 异步任务
}, 30*time.Second)
```

**禁止**:
- ❌ 直接使用 go 关键字启动协程（除非有充分理由）
- ❌ 不限制协程数量

### 4. 扩展规则

#### 4.1 添加新工具

**步骤**:
1. 确认工具不存在于 `pkg/utils`
2. 在 `pkg/utils` 中创建新文件
3. 实现工具函数
4. 添加单元测试
5. 更新 `pkg/utils/README.md`
6. 更新本技能书

**命名规范**:
- 文件名：小写+下划线（如 `string_utils.go`）
- 函数名：大驼峰（如 `ConvertToSnakeCase`）

#### 4.2 添加新类型

**步骤**:
1. 确认类型不存在于 `types/`
2. 选择合适的子目录（constants/errors/result）
3. 添加类型定义
4. 添加注释说明
5. 更新本技能书

**命名规范**:
- 常量：大驼峰或全大写+下划线
- 类型：大驼峰
- 变量：小驼峰

#### 4.3 添加新组件

**步骤**:
1. 确认组件不存在于 `pkg/`
2. 在 `pkg/` 中创建新目录
3. 定义接口（`interface.go`）
4. 实现接口（`impl.go`）
5. 添加配置（`config.go`）
6. 添加常量（`constants.go`）
7. 添加文档（`README.md`）
8. 添加示例（`examples/`）
9. 添加单元测试
10. 集成到 App 容器
11. 更新本技能书

**目录结构**:
```
pkg/newcomponent/
├── interface.go      # 接口定义
├── impl.go           # 接口实现
├── config.go         # 配置结构
├── constants.go      # 常量定义
├── doc.go            # 包文档
├── README.md         # 使用说明
├── examples/         # 示例代码
│   └── basic/
│       └── main.go
└── *_test.go         # 单元测试
```

## AI 行为约束

### 必须遵守

1. **优先复用**: 实现功能前必须先查找可复用组件
2. **禁止重复**: 禁止重复实现已有功能
3. **遵循规范**: 添加新组件必须遵循扩展规则
4. **更新文档**: 添加新组件必须更新本技能书

### 禁止行为

1. **重复实现**: 禁止重复实现已有功能
2. **绕过接口**: 禁止绕过公共接口
3. **硬编码**: 禁止硬编码应该使用常量的值
4. **自定义格式**: 禁止自定义响应格式

### 违规处理

- 检测到重复实现 → **中断并提示使用现有组件**
- 检测到绕过接口 → **直接拒绝执行**
- 检测到硬编码 → **中断并要求使用常量**
- 检测到自定义格式 → **中断并要求使用统一格式**

## 更新规则

### 必须更新的情况

1. 添加新工具
2. 添加新类型
3. 添加新组件
4. 修改扩展规则

### 更新流程

1. 更新本文件
2. 更新相关 README
3. 通知相关技能书同步更新
4. 更新本文件的修改日期

---

**最后更新**: 2026-01-19
**维护者**: AI Assistant
