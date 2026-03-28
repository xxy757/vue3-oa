# 登录模块文档

## 1. 模块功能说明

登录模块是 OA 办公系统的入口模块，提供用户身份认证功能。主要功能包括：

- 用户名密码登录
- 记住我功能（可选）
- 表单验证
- 登录状态管理
- 登录成功后自动跳转

## 2. 页面组件说明

### 2.1 登录页面 (`src/views/login/index.vue`)

登录页面使用 Naive UI 组件库构建，包含以下核心组件：

| 组件 | 说明 |
|------|------|
| `n-card` | 登录表单容器卡片 |
| `n-form` | 表单组件，支持验证 |
| `n-input` | 用户名和密码输入框 |
| `n-checkbox` | 记住我复选框 |
| `n-button` | 登录按钮 |
| `n-icon` | 图标组件 |

### 2.2 404 页面 (`src/views/error/404.vue`)

404 页面用于处理未匹配到的路由，包含：

- 404 错误码展示
- 错误提示信息
- 返回首页按钮

## 3. 接口信息

### 3.1 登录接口

- **URL**: `POST /api/auth/login`
- **请求参数**:

```typescript
interface LoginForm {
  username: string  // 用户名（必填）
  password: string  // 密码（必填）
  remember?: boolean // 记住我（可选）
}
```

- **响应数据**:

```typescript
interface LoginResult {
  token: string  // JWT Token
  user: User     // 用户信息
}
```

### 3.2 测试账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | 123456 | 管理员 |
| user | 123456 | 普通用户 |
| manager | 123456 | 经理 |

## 4. 主要实现逻辑

### 4.1 登录流程

1. 用户输入用户名和密码
2. 点击登录按钮触发 `handleLogin` 方法
3. 首先进行表单验证（用户名和密码必填）
4. 调用 `userStore.login()` 方法发送登录请求
5. 登录成功后，Token 和用户信息存储到 localStorage
6. 跳转到 `/dashboard` 页面

### 4.2 状态管理

登录状态通过 Pinia 的 `userStore` 进行管理：

- `token`: 存储 JWT Token
- `userInfo`: 存储用户信息
- `isLoggedIn`: 计算属性，判断是否已登录

### 4.3 路由守卫

路由守卫在 `src/router/index.ts` 中实现：

- 已登录用户访问 `/login` 时自动跳转到 `/dashboard`
- 未登录用户访问需要认证的页面时跳转到 `/login`
- `/login` 和 `/404` 在白名单中，无需认证

### 4.4 表单验证规则

```typescript
const rules: FormRules = {
  username: {
    required: true,
    message: '请输入用户名',
    trigger: ['blur', 'input']
  },
  password: {
    required: true,
    message: '请输入密码',
    trigger: ['blur', 'input']
  }
}
```

## 5. 相关文件

| 文件路径 | 说明 |
|----------|------|
| `src/views/login/index.vue` | 登录页面组件 |
| `src/views/error/404.vue` | 404 错误页面 |
| `src/stores/user.ts` | 用户状态管理 |
| `src/router/index.ts` | 路由配置和守卫 |
| `src/router/routes.ts` | 路由定义 |
| `src/types/user.d.ts` | 用户相关类型定义 |
| `src/utils/storage.ts` | Token 存储工具 |
| `src/utils/request.ts` | HTTP 请求封装 |
