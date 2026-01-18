---
name: 规则约束
description: Go Scaffold 的硬性规则、技术红线和违规定义
category: 核心
priority: 4
dependencies: [00-project-map, 01-project-foundation, 02-design-specs]
---

# 规则约束

## 职责范围

- 定义硬性规则（不可违反）
- 定义技术红线（触碰即拒绝）
- 定义违规行为和处理方式
- 提供规则检查清单

## 适用场景

- 代码审查
- 自动化检查
- AI 行为约束
- 团队规范

## 核心规则

### 1. 硬性规则（Hard Rules）

#### 规则 1: 依赖方向必须单向

**规则描述**:
- 上层可以依赖下层，下层不能依赖上层
- cmd → internal → pkg → types

**检查方式**:
```bash
# 检查 pkg 是否依赖 internal
grep -r "github.com/rei0721/go-scaffold/internal" pkg/
# 应该返回空结果

# 检查 types 是否依赖 pkg
grep -r "github.com/rei0721/go-scaffold/pkg" types/
# 应该返回空结果
```

**违规处理**: **直接拒绝执行**

#### 规则 2: 不允许使用全局变量

**规则描述**:
- 禁止使用全局变量（除常量外）
- 所有依赖通过依赖注入传递

**允许的全局变量**:
```go
// ✅ 允许：常量
const MaxRetries = 3

// ✅ 允许：错误定义
var ErrNotFound = errors.New("not found")

// ❌ 禁止：全局状态
var globalDB *gorm.DB
var globalLogger *zap.Logger
```

**违规处理**: **直接拒绝执行**

#### 规则 3: 所有错误必须处理

**规则描述**:
- 不允许忽略错误（除非有充分理由）
- 错误必须记录或返回

**正确示例**:
```go
// ✅ 正确：处理错误
data, err := repo.FindByID(ctx, id)
if err != nil {
    logger.Error("failed to find data", "error", err)
    return nil, err
}

// ✅ 正确：有充分理由忽略错误
_ = logger.Sync()  // Sync 在某些平台会返回无害的错误
```

**错误示例**:
```go
// ❌ 错误：忽略错误
data, _ := repo.FindByID(ctx, id)
```

**违规处理**: **中断并要求修正**

#### 规则 4: 数据库操作必须使用事务

**规则描述**:
- 涉及多个写操作的必须使用事务
- 使用 dbtx.Manager 管理事务

**正确示例**:
```go
// ✅ 正确：使用事务
err := s.txMgr.WithTransaction(ctx, func(ctx context.Context) error {
    if err := s.repo.Create(ctx, user); err != nil {
        return err
    }
    if err := s.repo.CreateProfile(ctx, profile); err != nil {
        return err
    }
    return nil
})
```

**错误示例**:
```go
// ❌ 错误：不使用事务
s.repo.Create(ctx, user)
s.repo.CreateProfile(ctx, profile)
```

**违规处理**: **中断并要求修正**

#### 规则 5: 密码必须加密存储

**规则描述**:
- 密码必须使用 bcrypt 加密
- 不允许明文存储密码
- 使用 pkg/crypto 包

**正确示例**:
```go
// ✅ 正确：加密密码
hashedPassword, err := crypto.HashPassword(password)
if err != nil {
    return err
}
user.Password = hashedPassword
```

**错误示例**:
```go
// ❌ 错误：明文存储
user.Password = password
```

**违规处理**: **直接拒绝执行**

#### 规则 6: JWT 密钥不允许硬编码

**规则描述**:
- JWT 密钥必须从环境变量或配置文件读取
- 不允许硬编码在代码中

**正确示例**:
```go
// ✅ 正确：从配置读取
jwtSecret := config.JWT.Secret
```

**错误示例**:
```go
// ❌ 错误：硬编码
jwtSecret := "my-secret-key"
```

**违规处理**: **直接拒绝执行**

#### 规则 7: 不允许 SQL 注入

**规则描述**:
- 使用 GORM 的参数化查询
- 不允许拼接 SQL 字符串

**正确示例**:
```go
// ✅ 正确：参数化查询
db.Where("username = ?", username).First(&user)
```

**错误示例**:
```go
// ❌ 错误：拼接 SQL
db.Raw("SELECT * FROM users WHERE username = '" + username + "'").Scan(&user)
```

**违规处理**: **直接拒绝执行**

#### 规则 8: 日志不允许输出敏感信息

**规则描述**:
- 不允许在日志中输出密码、token、密钥等敏感信息
- 敏感信息必须脱敏

**正确示例**:
```go
// ✅ 正确：脱敏
logger.Info("user login", "username", username)

// ✅ 正确：不输出密码
logger.Info("user created", "id", user.ID)
```

**错误示例**:
```go
// ❌ 错误：输出密码
logger.Info("user login", "username", username, "password", password)

// ❌ 错误：输出 token
logger.Info("token generated", "token", token)
```

**违规处理**: **直接拒绝执行**

#### 规则 9: 代码变更必须同步更新文档

**规则描述**:
- 触发文档更新条件的代码变更必须同步更新 docs 文档
- 参考 `15-documentation` 技能书的触发条件表

**触发条件**（P0 - 必须立即更新）:
- 新增/删除模块或包
- 修改配置文件结构或新增配置项
- 修改 API 接口
- 修改数据库表结构
- 修改依赖关系
- 修改启动流程或环境要求

**检查方式**:
```bash
# 提交前检查
- [ ] 确认是否触发文档更新条件
- [ ] 更新相关文档
- [ ] 验证文档内容准确性
```

**违规处理**: **直接拒绝执行**

### 2. 技术红线（Red Lines）

#### 红线 1: 不允许修改 BaseModel

**原因**:
- BaseModel 是所有模型的基础
- 修改会影响所有模型
- 可能导致数据库迁移问题

**定义**:
```go
// internal/models/db_base.go
type BaseModel struct {
    ID        int64          `gorm:"primaryKey;autoIncrement:false" json:"id"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

**违规处理**: **直接拒绝执行**

#### 红线 2: 不允许绕过 Repository 直接访问数据库

**原因**:
- Repository 提供统一的数据访问接口
- 绕过 Repository 会导致缓存失效
- 难以维护和测试

**正确示例**:
```go
// ✅ 正确：通过 Repository
user, err := s.userRepo.FindByID(ctx, id)
```

**错误示例**:
```go
// ❌ 错误：直接访问数据库
var user models.User
s.db.GetDB().First(&user, id)
```

**违规处理**: **直接拒绝执行**

#### 红线 3: 不允许在 Handler 中实现业务逻辑

**原因**:
- Handler 只负责 HTTP 协议相关的处理
- 业务逻辑应该在 Service 层
- 便于复用和测试

**正确示例**:
```go
// ✅ 正确：调用 Service
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, result.Error(err))
        return
    }
    
    user, err := h.service.CreateUser(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, result.Error(err))
        return
    }
    
    c.JSON(200, result.Success(user))
}
```

**错误示例**:
```go
// ❌ 错误：在 Handler 中实现业务逻辑
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    c.ShouldBindJSON(&req)
    
    // 业务逻辑不应该在这里
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
    user := &models.User{
        Username: req.Username,
        Password: string(hashedPassword),
    }
    h.db.Create(user)
    
    c.JSON(200, user)
}
```

**违规处理**: **直接拒绝执行**

#### 红线 4: 不允许在 pkg 中依赖 internal

**原因**:
- pkg 是公共组件，应该业务无关
- 依赖 internal 会导致循环依赖
- 降低 pkg 的复用性

**违规处理**: **直接拒绝执行**

#### 红线 5: 不允许修改 App 容器的初始化顺序

**原因**:
- 初始化顺序是精心设计的
- 修改可能导致依赖错误
- 可能导致应用无法启动

**正确顺序**:
```
1. Config
2. Logger
3. I18n
4. IDGenerator
5. Database
6. Cache
7. Crypto
8. JWT
9. RBAC
10. Executor
11. Repository
12. Service
13. Handler
14. Router
15. HTTPServer
```

**违规处理**: **直接拒绝执行**

### 3. 违规定义

#### 违规级别

| 级别 | 名称     | 处理方式         | 示例                     |
| ---- | -------- | ---------------- | ------------------------ |
| P0   | 严重违规 | 直接拒绝执行     | 反向依赖、SQL 注入       |
| P1   | 高危违规 | 中断并要求修正   | 忽略错误、不使用事务     |
| P2   | 中危违规 | 警告并建议修正   | 代码风格、命名不规范     |
| P3   | 低危违规 | 提示但允许继续   | 注释不完整、文档缺失     |

#### 违规检查清单

**代码提交前检查**:
- [ ] 依赖方向正确
- [ ] 没有全局变量
- [ ] 所有错误已处理
- [ ] 数据库操作使用事务
- [ ] 密码已加密
- [ ] JWT 密钥未硬编码
- [ ] 没有 SQL 注入风险
- [ ] 日志未输出敏感信息
- [ ] 没有绕过 Repository
- [ ] Handler 没有业务逻辑
- [ ] pkg 没有依赖 internal
- [ ] **docs 文档已同步更新（如果需要）**

**代码审查检查**:
- [ ] 遵守分层架构
- [ ] 遵守接口设计规范
- [ ] 遵守错误处理规范
- [ ] 遵守命名规范
- [ ] 有足够的单元测试
- [ ] 有足够的注释
- [ ] 更新了相关文档

## AI 行为约束

### 必须遵守

1. **硬性规则**: 严格遵守所有硬性规则
2. **技术红线**: 不触碰任何技术红线
3. **违规检查**: 执行前进行违规检查
4. **违规处理**: 按照违规级别处理

### 禁止行为

1. **忽略规则**: 禁止忽略任何硬性规则
2. **触碰红线**: 禁止触碰任何技术红线
3. **跳过检查**: 禁止跳过违规检查
4. **降低级别**: 禁止降低违规级别

### 违规处理流程

```
检测到违规
    ↓
判断违规级别
    ↓
P0: 直接拒绝执行 → 输出错误信息 → 停止
P1: 中断并要求修正 → 输出修正建议 → 等待用户确认
P2: 警告并建议修正 → 输出警告信息 → 继续执行
P3: 提示但允许继续 → 输出提示信息 → 继续执行
```

## 更新规则

### 必须更新的情况

1. 新增硬性规则
2. 新增技术红线
3. 修改违规定义
4. 修改违规级别

### 更新流程

1. 更新本文件
2. 更新违规检查清单
3. 通知相关技能书同步更新
4. 更新本文件的修改日期

---

**最后更新**: 2026-01-19
**维护者**: AI Assistant
