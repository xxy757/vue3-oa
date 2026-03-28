# 系统管理和个人中心模块文档

## 1. 模块功能说明

本模块包含系统管理和个人中心两大功能模块，共6个页面。

### 1.1 系统管理模块
系统管理模块提供系统的基础配置功能，仅限管理员（admin）角色访问。

- **用户管理**：管理系统用户，包括用户新增、编辑、启用/禁用、删除等功能
- **部门管理**：管理组织架构，支持树形结构的部门管理
- **角色管理**：管理系统角色，配置角色权限
- **流程配置**：配置审批流程，设置流程节点和审批人

### 1.2 个人中心模块
个人中心模块提供用户个人信息管理功能，所有已登录用户均可访问。

- **个人信息**：查看和编辑个人基本信息
- **修改密码**：修改登录密码，包含密码强度提示

## 2. 页面组件说明

### 2.1 用户管理 (User.vue)

**文件路径**: `src/views/system/User.vue`

**功能特性**:
- 用户列表展示（用户名、昵称、部门、角色、状态等）
- 关键词搜索功能
- 分页功能
- 新增用户弹窗表单
- 编辑用户弹窗表单
- 启用/禁用用户状态切换
- 删除用户确认

**组件依赖**:
- `NDataTable` - 数据表格
- `NModal` - 弹窗
- `NForm` - 表单
- `NTreeSelect` - 部门树选择器
- `NSelect` - 角色选择器

**核心方法**:
| 方法名 | 说明 |
|--------|------|
| `fetchUserList()` | 获取用户列表 |
| `fetchDeptList()` | 获取部门列表 |
| `fetchRoleList()` | 获取角色列表 |
| `handleAdd()` | 打开新增用户弹窗 |
| `handleEdit(row)` | 打开编辑用户弹窗 |
| `handleSubmit()` | 提交用户表单 |
| `handleToggleStatus(row)` | 切换用户状态 |
| `handleDelete(row)` | 删除用户 |

---

### 2.2 部门管理 (Dept.vue)

**文件路径**: `src/views/system/Dept.vue`

**功能特性**:
- 左侧部门树形结构展示
- 右侧部门详情展示
- 新增部门（支持添加顶级部门和子部门）
- 编辑部门信息
- 删除部门（有子部门时不可删除）

**组件依赖**:
- `NTree` - 树形组件
- `NDescriptions` - 描述列表
- `NTreeSelect` - 上级部门选择器

**核心方法**:
| 方法名 | 说明 |
|--------|------|
| `fetchDeptList()` | 获取部门列表 |
| `buildTree()` | 构建树形数据 |
| `handleAdd(parent)` | 打开新增部门弹窗 |
| `handleEdit(dept)` | 打开编辑部门弹窗 |
| `canDelete(dept)` | 判断部门是否可删除 |

---

### 2.3 角色管理 (Role.vue)

**文件路径**: `src/views/system/Role.vue`

**功能特性**:
- 角色列表展示
- 新增角色
- 编辑角色
- 删除角色（admin角色不可删除）
- 权限配置（多选复选框）

**组件依赖**:
- `NDataTable` - 数据表格
- `NCheckboxGroup` - 权限多选组

**权限选项**:
| 权限编码 | 权限名称 |
|----------|----------|
| `*` | 所有权限 |
| `approval:apply` | 发起申请 |
| `approval:approve` | 审批申请 |
| `notice:view` | 查看公告 |
| `notice:create` | 发布公告 |
| `notice:manage` | 管理公告 |
| `schedule:view` | 查看日程 |
| `schedule:manage` | 管理日程 |
| `system:user` | 用户管理 |
| `system:dept` | 部门管理 |
| `system:role` | 角色管理 |
| `system:flow` | 流程配置 |

---

### 2.4 流程配置 (Flow.vue)

**文件路径**: `src/views/system/Flow.vue`

**功能特性**:
- 流程列表展示
- 新增审批流程
- 编辑流程节点
- 流程可视化预览
- 删除流程

**流程节点类型**:
| 类型 | 说明 |
|------|------|
| `submit` | 提交节点 |
| `approval` | 审批节点 |
| `notify` | 通知节点 |
| `condition` | 条件分支 |

**组件依赖**:
- `NSelect` - 节点类型选择器
- `NModal` - 弹窗

**核心数据结构**:
```typescript
interface FlowNode {
  name: string
  type: 'submit' | 'approval' | 'notify' | 'condition'
  approver: number[]
}

interface FlowConfig {
  id: number
  name: string
  code: string
  description: string
  nodes: FlowNode[]
  status: 0 | 1
  createTime: string
}
```

---

### 2.5 个人信息 (Info.vue)

**文件路径**: `src/views/profile/Info.vue`

**功能特性**:
- 左侧个人信息卡片（头像、昵称、角色、部门）
- 右侧信息编辑表单
- 头像URL编辑
- 账号安全提示

**组件依赖**:
- `NAvatar` - 头像组件
- `NDescriptions` - 描述列表
- `NForm` - 表单
- `NList` - 安全设置列表

**核心方法**:
| 方法名 | 说明 |
|--------|------|
| `initFormData()` | 初始化表单数据 |
| `handleSubmit()` | 保存个人信息 |
| `handleReset()` | 重置表单 |

---

### 2.6 修改密码 (Password.vue)

**文件路径**: `src/views/profile/Password.vue`

**功能特性**:
- 旧密码输入
- 新密码输入
- 确认新密码
- 密码强度实时计算和提示
- 密码要求说明
- 安全提示信息

**密码强度规则**:
- 长度 >= 6: +20分
- 长度 >= 10: +10分
- 包含小写字母: +15分
- 包含大写字母: +15分
- 包含数字: +20分
- 包含特殊字符: +20分

**强度等级**:
- 0-40%: 弱 (红色)
- 40-70%: 中 (黄色)
- 70-100%: 强 (绿色)

---

## 3. 数据结构说明

### 3.1 用户 (User)

```typescript
interface User {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  deptId: number
  deptName: string
  roleId: number
  roleName: string
  status: 0 | 1  // 0: 禁用, 1: 启用
  createTime: string
}
```

### 3.2 部门 (Dept)

```typescript
interface Dept {
  id: number
  name: string
  parentId: number | null
  children?: Dept[]
  sort: number
  leader: string
  phone: string
  email: string
  status: 0 | 1
}
```

### 3.3 角色 (Role)

```typescript
interface Role {
  id: number
  name: string
  code: string
  description: string
  permissions: string[]
  status: 0 | 1
  createTime: string
}
```

### 3.4 流程节点 (FlowNode)

```typescript
interface FlowNode {
  name: string
  type: 'submit' | 'approval' | 'notify' | 'condition'
  approver: number[]
}
```

### 3.5 流程配置 (FlowConfig)

```typescript
interface FlowConfig {
  id: number
  name: string
  code: string
  description: string
  nodes: FlowNode[]
  status: 0 | 1
  createTime: string
}
```

---

## 4. 接口信息

### 4.1 用户相关接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/user/list` | GET | 获取用户列表（支持分页和搜索） |
| `/api/user` | POST | 新增用户 |
| `/api/user/:id` | PUT | 编辑用户 |
| `/api/user/:id` | DELETE | 删除用户 |
| `/api/user/:id/status` | PUT | 切换用户状态 |

**用户列表请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | number | 否 | 页码，默认1 |
| pageSize | number | 否 | 每页数量，默认10 |
| keyword | string | 否 | 搜索关键词 |

### 4.2 部门相关接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/dept/list` | GET | 获取部门列表（树形结构） |
| `/api/dept` | POST | 新增部门 |
| `/api/dept/:id` | PUT | 编辑部门 |
| `/api/dept/:id` | DELETE | 删除部门 |

### 4.3 角色相关接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/role/list` | GET | 获取角色列表 |
| `/api/role` | POST | 新增角色 |
| `/api/role/:id` | PUT | 编辑角色 |
| `/api/role/:id` | DELETE | 删除角色 |

### 4.4 认证相关接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/auth/info` | GET | 获取当前用户信息 |
| `/api/auth/info` | PUT | 更新当前用户信息 |
| `/api/auth/password` | PUT | 修改密码 |

**修改密码请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| oldPassword | string | 是 | 旧密码 |
| newPassword | string | 是 | 新密码 |

---

## 5. 权限控制说明

### 5.1 路由权限

在路由配置中，通过 `meta.roles` 字段控制访问权限：

```typescript
{
  path: 'user',
  name: 'SystemUser',
  component: () => import('@/views/system/User.vue'),
  meta: {
    title: '用户管理',
    requiresAuth: true,
    roles: ['admin']  // 仅管理员可访问
  }
}
```

### 5.2 权限矩阵

| 页面 | 路径 | admin | manager | employee |
|------|------|-------|---------|----------|
| 用户管理 | /system/user | Y | N | N |
| 部门管理 | /system/dept | Y | N | N |
| 角色管理 | /system/role | Y | N | N |
| 流程配置 | /system/flow | Y | N | N |
| 个人信息 | /profile/info | Y | Y | Y |
| 修改密码 | /profile/password | Y | Y | Y |

### 5.3 权限验证流程

1. 用户登录后获取用户信息和角色
2. 路由守卫检查目标路由的 `meta.roles`
3. 如果用户角色不在允许的角色列表中，跳转到403页面
4. 个人中心页面所有登录用户均可访问

### 5.4 Mock 数据中的角色定义

| 角色ID | 角色名称 | 角色编码 | 权限 |
|--------|----------|----------|------|
| 1 | 管理员 | admin | * (所有权限) |
| 2 | 员工 | employee | approval:apply, notice:view, schedule:view |
| 3 | 部门经理 | manager | approval:apply, approval:approve, notice:view, notice:create, schedule:view |

---

## 6. 使用说明

### 6.1 开发环境

1. 确保 Mock 服务已启用（vite-plugin-mock）
2. 访问 `/system/*` 路由需要使用 admin 账号登录
3. 默认管理员账号：用户名 `admin`，密码 `123456`

### 6.2 页面跳转

- 系统管理菜单：侧边栏 -> 系统管理
- 个人中心：点击右上角用户头像 -> 个人信息 / 修改密码

### 6.3 注意事项

1. 用户管理中，用户名创建后不可修改
2. 角色管理中，admin 角色不可删除
3. 部门管理中，有子部门的部门不可删除
4. 修改密码后需要重新登录
5. 流程配置目前使用前端模拟数据，实际项目需要对接后端接口
