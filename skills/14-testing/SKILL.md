---
name: 测试规范
description: 单元测试、集成测试和测试覆盖率规范
category: 支撑
priority: 9
dependencies: [00-project-map, 02-design-specs]
---

# 测试规范

## 职责范围

- 单元测试规范
- 集成测试规范
- 测试覆盖率要求
- Mock 使用规范

## 适用场景

- 编写测试
- 代码审查
- CI/CD 流程
- 质量保证

## 核心规则

### 1. 单元测试规范

#### 测试文件命名

```
源文件: user_service.go
测试文件: user_service_test.go
```

#### 测试函数命名

```go
func TestUserService_Create(t *testing.T) {
    // 测试创建用户
}

func TestUserService_Create_InvalidParams(t *testing.T) {
    // 测试无效参数
}
```

#### 测试结构

```go
func TestUserService_Create(t *testing.T) {
    // 1. Setup
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo, nil, nil)
    
    // 2. Execute
    user, err := service.Create(context.Background(), &CreateUserRequest{
        Username: "test",
        Email:    "test@example.com",
    })
    
    // 3. Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test", user.Username)
}
```

### 2. 集成测试规范

#### 测试标签

```go
// +build integration

package service_test

func TestUserService_Integration(t *testing.T) {
    // 集成测试
}
```

#### 运行集成测试

```bash
# 运行所有测试
go test ./...

# 只运行单元测试
go test -short ./...

# 只运行集成测试
go test -tags=integration ./...
```

### 3. 测试覆盖率

#### 目标

- 核心业务逻辑：>= 80%
- 工具函数：>= 70%
- Handler 层：>= 60%

#### 生成覆盖率报告

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 4. Mock 使用

#### 接口 Mock

```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}
```

## AI 行为约束

### 必须遵守

1. 新功能必须编写测试
2. 测试覆盖率达标
3. 测试必须可重复运行
4. 测试不依赖外部服务

### 禁止行为

1. 禁止跳过测试
2. 禁止提交失败的测试
3. 禁止在测试中使用 time.Sleep
4. 禁止测试依赖执行顺序

---

**最后更新**: 2026-01-19
