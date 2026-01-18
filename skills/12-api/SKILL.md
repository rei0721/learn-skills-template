---
name: API 规范
description: RESTful API 设计规范、路由命名和版本管理
category: 支撑
priority: 7
dependencies: [00-project-map, 02-design-specs, 04-reusables]
---

# API 规范

## 职责范围

- RESTful API 设计规范
- 路由命名规范
- 请求响应格式
- API 版本管理

## 适用场景

- 设计新 API
- 修改现有 API
- API 文档编写
- 前后端联调

## 核心规则

### 1. RESTful 设计规范

#### 资源命名

```
✅ 正确：使用复数名词
GET    /api/users          # 获取用户列表
GET    /api/users/:id      # 获取单个用户
POST   /api/users          # 创建用户
PUT    /api/users/:id      # 更新用户
DELETE /api/users/:id      # 删除用户

❌ 错误：使用动词
GET    /api/getUsers
POST   /api/createUser
```

#### HTTP 方法

| 方法   | 用途     | 幂等性 | 示例                    |
| ------ | -------- | ------ | ----------------------- |
| GET    | 查询     | 是     | GET /api/users          |
| POST   | 创建     | 否     | POST /api/users         |
| PUT    | 完整更新 | 是     | PUT /api/users/:id      |
| PATCH  | 部分更新 | 否     | PATCH /api/users/:id    |
| DELETE | 删除     | 是     | DELETE /api/users/:id   |

### 2. 统一响应格式

#### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 123,
    "username": "john"
  }
}
```

#### 错误响应

```json
{
  "code": 400,
  "message": "invalid parameters",
  "data": null
}
```

#### 分页响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "page_size": 10,
    "total_pages": 10
  }
}
```

### 3. 路由组织

```go
// internal/router/routes.go
func SetupRoutes(r *gin.Engine, app *app.App) {
    // API v1
    v1 := r.Group("/api/v1")
    {
        // 公开路由
        public := v1.Group("")
        {
            public.POST("/auth/login", authHandler.Login)
            public.POST("/auth/register", authHandler.Register)
        }
        
        // 需要认证的路由
        auth := v1.Group("")
        auth.Use(middleware.AuthMiddleware(app.JWT))
        {
            auth.GET("/users/me", userHandler.GetMe)
            auth.PUT("/users/me", userHandler.UpdateMe)
        }
        
        // 需要权限的路由
        admin := v1.Group("/admin")
        admin.Use(middleware.AuthMiddleware(app.JWT))
        admin.Use(middleware.RBACMiddleware(app.RBAC))
        {
            admin.GET("/users", userHandler.List)
            admin.POST("/users", userHandler.Create)
        }
    }
}
```

### 4. 状态码规范

| 状态码 | 含义                 | 使用场景               |
| ------ | -------------------- | ---------------------- |
| 200    | OK                   | 请求成功               |
| 201    | Created              | 资源创建成功           |
| 204    | No Content           | 删除成功               |
| 400    | Bad Request          | 参数错误               |
| 401    | Unauthorized         | 未认证                 |
| 403    | Forbidden            | 无权限                 |
| 404    | Not Found            | 资源不存在             |
| 500    | Internal Server Error| 服务器错误             |

## AI 行为约束

### 必须遵守

1. 使用 RESTful 风格
2. 使用统一响应格式
3. 正确使用 HTTP 状态码
4. 遵循路由命名规范

### 禁止行为

1. 禁止在 URL 中使用动词
2. 禁止自定义响应格式
3. 禁止滥用 POST 方法
4. 禁止返回敏感信息

---

**最后更新**: 2026-01-19
