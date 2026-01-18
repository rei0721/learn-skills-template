---
name: 数据库规范
description: 数据库表设计、模型定义和迁移规范
category: 支撑
priority: 8
dependencies: [00-project-map, 02-design-specs, 03-rules]
---

# 数据库规范

## 职责范围

- 数据库表设计规范
- 模型定义规范
- 数据库迁移规范
- 字段不可变性规则

## 适用场景

- 设计新表
- 修改表结构
- 定义模型
- 数据库迁移

## 核心规则

### 1. 表设计规范

#### 表命名

```
✅ 正确：小写+下划线，复数形式
users
user_profiles
order_items

❌ 错误：驼峰或单数
User
userProfile
order
```

#### 字段命名

```
✅ 正确：小写+下划线
user_id
created_at
is_active

❌ 错误：驼峰
userId
createdAt
isActive
```

### 2. 模型定义规范

#### BaseModel

```go
// internal/models/db_base.go
type BaseModel struct {
    ID        int64          `gorm:"primaryKey;autoIncrement:false" json:"id"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

**规则**:
- 所有模型必须嵌入 BaseModel
- ID 使用 Snowflake 生成，不使用自增
- 支持软删除

#### 模型示例

```go
// internal/models/db_user.go
type User struct {
    BaseModel
    Username string `gorm:"uniqueIndex;not null;size:50" json:"username"`
    Email    string `gorm:"uniqueIndex;not null;size:100" json:"email"`
    Password string `gorm:"not null;size:255" json:"-"`
    Status   int    `gorm:"default:1;index" json:"status"`
}

func (User) TableName() string {
    return "users"
}
```

### 3. 字段约束

#### 必须约束

```go
// NOT NULL
Username string `gorm:"not null"`

// UNIQUE
Email string `gorm:"uniqueIndex"`

// DEFAULT
Status int `gorm:"default:1"`

// SIZE
Username string `gorm:"size:50"`
```

#### 索引

```go
// 单字段索引
Status int `gorm:"index"`

// 唯一索引
Email string `gorm:"uniqueIndex"`

// 复合索引
type User struct {
    BaseModel
    Username string `gorm:"index:idx_username_status"`
    Status   int    `gorm:"index:idx_username_status"`
}
```

### 4. 关联关系

#### 一对一

```go
type User struct {
    BaseModel
    Profile Profile `gorm:"foreignKey:UserID"`
}

type Profile struct {
    BaseModel
    UserID int64 `gorm:"uniqueIndex"`
    Bio    string
}
```

#### 一对多

```go
type User struct {
    BaseModel
    Orders []Order `gorm:"foreignKey:UserID"`
}

type Order struct {
    BaseModel
    UserID int64 `gorm:"index"`
}
```

### 5. 数据库迁移

#### 自动迁移

```go
// internal/app/app_initdb.go
func (a *App) AutoMigrate() error {
    return a.DB.GetDB().AutoMigrate(
        &models.User{},
        &models.Profile{},
        &models.Order{},
    )
}
```

#### 字段不可变性

**不可修改的字段**:
- ID（主键）
- CreatedAt（创建时间）
- 外键字段（除非级联更新）

**可修改的字段**:
- UpdatedAt（自动更新）
- 业务字段

## AI 行为约束

### 必须遵守

1. 所有模型必须嵌入 BaseModel
2. 表名使用复数形式
3. 字段名使用下划线命名
4. 添加必要的约束和索引

### 禁止行为

1. 禁止修改 BaseModel
2. 禁止使用数据库自增 ID
3. 禁止删除外键约束
4. 禁止在生产环境直接修改表结构

---

**最后更新**: 2026-01-19
