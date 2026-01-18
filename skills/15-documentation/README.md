# 文档维护规范

> 确保 docs 项目文档与代码保持同步

## 📚 文档列表

- **[SKILL.md](./SKILL.md)** - 完整的文档维护规范（必读）
- **[quick-reference.md](./quick-reference.md)** - 快速参考卡片（常用）

## 🎯 核心原则

> **代码变更必须同步更新文档！**

这不是建议，而是**硬性规则**：

1. ✅ 代码提交前检查是否触发文档更新
2. ✅ 同一个 PR 中包含代码和文档变更
3. ✅ 文档内容必须与代码实现一致
4. ❌ 不允许"稍后更新文档"

## 🚀 快速开始

### 1. 确认是否需要更新文档

查看 [quick-reference.md](./quick-reference.md) 中的触发条件表：

- 🔴 P0（必须立即更新）：新增模块、修改配置、新增 API、修改表结构等
- 🟠 P1（应该尽快更新）：新增功能、修改功能、性能优化、Bug 修复等
- 🟡 P2（建议更新）：代码重构（如果影响架构）

### 2. 找到需要更新的文档

根据变更类型，在 [quick-reference.md](./quick-reference.md) 中找到对应的文档：

```
新增模块 → project-structure.md + architecture.md
修改配置 → configuration.md
新增 API → api/ 目录
修改表结构 → database.md
...
```

### 3. 更新文档

按照 [SKILL.md](./SKILL.md) 中的具体场景示例更新文档：

- 场景 1: 新增模块
- 场景 2: 修改配置
- 场景 3: 新增 API
- 场景 4: 修改数据库表
- 场景 5: 修改依赖
- 场景 6: 重构架构

### 4. 验证文档质量

使用 [SKILL.md](./SKILL.md) 中的检查清单：

- [ ] 文档内容与代码一致
- [ ] 包含必要的示例
- [ ] 格式规范
- [ ] 链接有效

## 📖 文档结构

```
docs/
├── overview/           # 项目概览（架构、技术栈）
├── getting-started/    # 快速开始（安装、配置）
├── development/        # 开发指南（结构、配置）
├── api/                # API 文档（需要时创建）
├── contributing/       # 贡献指南
├── changelog.md        # 变更日志（重要！）
└── faq.md              # 常见问题
```

## ⚠️ 常见错误

### 错误 1: 忘记更新 changelog.md

```markdown
❌ 错误：新增功能但未记录到 changelog.md

✅ 正确：
## [0.1.3] - 2026-01-20

### Added
- 新增用户管理功能
- 新增 POST /api/v1/users 接口
```

### 错误 2: 配置示例与实际不一致

```yaml
❌ 错误：文档中的配置示例与 configs/config.yaml 不一致

✅ 正确：从实际配置文件复制示例
```

### 错误 3: 新增 API 但未写文档

```markdown
❌ 错误：新增 API 接口但未创建 API 文档

✅ 正确：在 docs/api/ 目录下创建接口文档
```

## 🔍 AI 行为约束

AI 在执行代码变更时**必须**：

1. **检查触发条件**: 每次代码变更前检查是否触发文档更新
2. **同步更新**: 代码和文档在同一次操作中更新
3. **验证一致性**: 确保文档内容与代码实现一致
4. **拒绝违规**: 如果检测到应该更新文档但未更新，直接拒绝执行

## 📝 示例

### 示例 1: 新增配置项

**代码变更**:
```yaml
# configs/config.yaml
new_feature:
  enabled: true
  timeout: 30
```

**文档更新**:
```markdown
# docs/development/configuration.md

## new_feature 配置

### enabled
- **类型**: boolean
- **默认值**: true
- **说明**: 是否启用新功能

### timeout
- **类型**: int
- **默认值**: 30
- **说明**: 超时时间（秒）

**示例**:
\```yaml
new_feature:
  enabled: true
  timeout: 30
\```
```

### 示例 2: 新增 API

**代码变更**:
```go
// internal/handler/user_handler.go
func (h *UserHandler) CreateUser(c *gin.Context) {
    // ...
}
```

**文档更新**:
```markdown
# docs/api/users.md

## POST /api/v1/users

创建新用户

**请求参数**:
\```json
{
  "username": "john",
  "email": "john@example.com"
}
\```

**响应示例**:
\```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 123,
    "username": "john"
  }
}
\```
```

## 🔗 相关链接

- [完整规范](./SKILL.md)
- [快速参考](./quick-reference.md)
- [项目文档](../../docs/)
- [变更日志](../../docs/changelog.md)

---

**维护者**: AI Assistant
**最后更新**: 2026-01-19
