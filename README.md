# OA SaaS 办公自动化系统

基于 Vue 3 + Go + 多租户架构构建的企业级 SaaS 办公自动化系统。

## 技术栈

### 前端

| 技术 | 版本 | 说明 |
|------|------|------|
| Vue | 3.4.21 | 渐进式 JavaScript 框架 |
| TypeScript | 5.4.2 | JavaScript 的超集，提供类型支持 |
| Vue Router | 4.3.0 | Vue.js 官方路由管理器 |
| Pinia | 2.1.7 | Vue.js 官方状态管理库 |
| Naive UI | 2.38.1 | Vue 3 组件库 |
| Axios | 1.6.8 | HTTP 请求库 |
| ECharts | 5.5.0 | 数据可视化图表库 |
| vue-echarts | 6.6.9 | ECharts 的 Vue 封装 |
| Day.js | 1.11.10 | 轻量级日期处理库 |
| Vite | 5.2.0 | 下一代前端构建工具 |
| Sass | 1.72.0 | CSS 预处理器 |
| ESLint | 8.57.0 | 代码质量检查 |
| Prettier | 3.2.5 | 代码格式化工具 |

### 后端

| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.26.1 | 后端开发语言 |
| Gin | 1.12 | HTTP Web 框架 |
| GORM | 1.31 | ORM 框架 |
| MySQL | 8.0 | 关系型数据库 |
| JWT (golang-jwt/v5) | 5.x | 身份认证 |
| Redis (go-redis/v9) | 9.x | 可选缓存 |
| gopkg.in/yaml.v3 | 3.x | 配置文件解析 |

---

## 功能模块

### 1. 工作台 (Dashboard)
- 数据统计展示（审批/公告/日程/日程统计卡片）
- 待办事项提醒
- 快捷操作入口
- 审批趋势图表

### 2. 审批中心
- **发起申请** - 支持 5 种审批类型
  - 请假申请（年假、病假、事假等）
  - 报销申请（差旅、办公、招待等）
  - 加班申请
  - 出差申请
  - 通用审批
- **我的申请** - 查看个人申请记录及状态
- **待我审批** - 处理待审批事项（审批/驳回/转交）
- **已办审批** - 查看已处理的审批记录
- **审批详情** - 时间线展示审批流程节点
- **审批统计** - 各状态数量统计

### 3. 公告通知
- 公告列表（支持分类筛选、关键词搜索、已读状态）
- 公告详情查看（自动标记已读）
- 未读数量提醒

### 4. 日程管理
- 日程日历视图
- 日程列表视图
- 日程增删改查（支持参与人、优先级、全天事件）
- 本周日程快捷查看

### 5. 系统管理 (需管理员权限)
- **用户管理** - 用户增删改查、状态启用/禁用、套餐用户数限制
- **部门管理** - 树形组织架构管理
- **角色管理** - 角色权限配置（JSON 权限数组）
- **流程配置** - 审批流程可视化设置

### 6. 个人中心
- 个人信息查看与修改
- 密码修改

### 7. 多租户 SaaS
- **企业注册** - 3 步向导（企业信息 → 选择套餐 → 完成注册）
- **套餐管理** - 4 个套餐等级（免费版/基础版/标准版/专业版）
- **企业管理** - 企业信息查看/编辑、用户数进度条
- **账单管理** - 账单列表、套餐升级
- **超管后台** - 仪表盘统计、租户管理、套餐配置

---

## 开发进度

| 阶段 | 名称 | 状态 | 说明 |
|------|------|------|------|
| A | 前端页面开发 | ✅ 已完成 | 26 个页面全部完成 |
| B | Go 后端搭建 | ✅ 已完成 | 39 个 API 接口，前后端联调通过 |
| C | 多租户 SaaS 改造 | ✅ 已完成 | 租户隔离 + SaaS 页面 + 超管后台 |
| D | 计费系统 | ⏳ 未开始 | 支付对接、定时账单、过期检查 |
| E | 部署上线 | ⏳ 未开始 | Docker 化、域名、监控 |

### 已完成页面清单

| 页面 | 文件路径 | 代码行数 | 状态 |
|------|----------|----------|------|
| 登录页 | `src/views/login/index.vue` | 180 | ✅ |
| 工作台 | `src/views/dashboard/index.vue` | 623 | ✅ |
| 发起申请 | `src/views/approval/Apply.vue` | 501 | ✅ |
| 我的申请 | `src/views/approval/MyApply.vue` | 234 | ✅ |
| 待我审批 | `src/views/approval/Pending.vue` | 220 | ✅ |
| 已办审批 | `src/views/approval/Done.vue` | 149 | ✅ |
| 审批详情 | `src/views/approval/Detail.vue` | 362 | ✅ |
| 公告列表 | `src/views/notice/List.vue` | 250 | ✅ |
| 公告详情 | `src/views/notice/Detail.vue` | 177 | ✅ |
| 日程日历 | `src/views/schedule/Calendar.vue` | 573 | ✅ |
| 日程列表 | `src/views/schedule/List.vue` | 680 | ✅ |
| 用户管理 | `src/views/system/User.vue` | 429 | ✅ |
| 部门管理 | `src/views/system/Dept.vue` | 347 | ✅ |
| 角色管理 | `src/views/system/Role.vue` | 297 | ✅ |
| 流程配置 | `src/views/system/Flow.vue` | 563 | ✅ |
| 个人信息 | `src/views/profile/Info.vue` | 207 | ✅ |
| 修改密码 | `src/views/profile/Password.vue` | 206 | ✅ |
| 404 页面 | `src/views/error/404.vue` | 56 | ✅ |
| 企业注册 | `src/views/auth/Register.vue` | 328 | ✅ |
| 选择套餐 | `src/views/auth/ChoosePlan.vue` | 240 | ✅ |
| 企业信息 | `src/views/tenant/Info.vue` | 166 | ✅ |
| 套餐管理 | `src/views/tenant/Plan.vue` | 254 | ✅ |
| 账单管理 | `src/views/tenant/Invoices.vue` | 133 | ✅ |
| 超管仪表盘 | `src/views/admin/Dashboard.vue` | 251 | ✅ |
| 超管租户管理 | `src/views/admin/Tenants.vue` | 352 | ✅ |
| 超管套餐管理 | `src/views/admin/Plans.vue` | 266 | ✅ |

---

## 项目结构

```
vue3-oa/
├── backend/                        # Go 后端
│   ├── cmd/server/main.go          # 入口（DB初始化 + AutoMigrate + Seed + 启动服务）
│   ├── configs/
│   │   ├── config.yaml             # 基础配置
│   │   └── config.local.yaml       # 本地覆盖配置
│   ├── internal/
│   │   ├── config/config.go        # 配置加载（YAML + 环境变量覆盖）
│   │   ├── middleware/
│   │   │   ├── auth.go             # JWT 认证中间件
│   │   │   ├── cors.go             # 跨域中间件
│   │   │   ├── logger.go           # 请求日志中间件
│   │   │   └── tenant.go           # 租户识别中间件（Header + 子域名）
│   │   ├── model/                  # 数据模型（14 个表）
│   │   │   ├── tenant.go           # Plan/Tenant/Invoice/TenantLog
│   │   │   ├── user.go             # 用户
│   │   │   ├── department.go       # 部门
│   │   │   ├── role.go             # 角色
│   │   │   ├── approval_flow.go    # 审批流程
│   │   │   ├── approval.go         # 审批 + 审批节点
│   │   │   ├── notice.go           # 公告 + 已读
│   │   │   ├── schedule.go         # 日程 + 参与人
│   │   │   ├── scopes.go           # 租户 Scope
│   │   │   └── seed.go             # 种子数据
│   │   ├── handler/                # API 处理器（39 个接口）
│   │   │   ├── auth_handler.go     # 认证（登录/用户信息）
│   │   │   ├── user_handler.go     # 用户 CRUD
│   │   │   ├── dept_handler.go     # 部门 CRUD
│   │   │   ├── role_handler.go     # 角色 CRUD
│   │   │   ├── approval_handler.go # 审批全流程
│   │   │   ├── notice_handler.go   # 公告 CRUD + 已读
│   │   │   ├── schedule_handler.go # 日程 CRUD + 周视图
│   │   │   ├── flow_handler.go     # 流程配置 CRUD
│   │   │   └── tenant_handler.go   # 租户注册/管理/升级/账单
│   │   ├── router/router.go        # 路由注册
│   │   └── pkg/
│   │       ├── jwt/jwt.go          # JWT 生成与解析
│   │       ├── cache/              # 缓存接口 + 内存/Redis 实现
│   │       └── utils/              # 密码加密 + 统一响应
│   ├── migrations/                 # SQL 迁移文件（参考用）
│   ├── Makefile                    # 构建/运行/清理
│   ├── go.mod
│   └── go.sum
├── src/                            # Vue 3 前端
│   ├── assets/                     # 静态资源
│   ├── components/Layout/
│   │   ├── OALayout.vue            # 主布局（侧边栏 + 头部 + 面包屑）
│   │   └── AdminLayout.vue         # 超管后台布局
│   ├── mock/                       # Mock 数据（开发阶段使用）
│   │   ├── user.ts
│   │   ├── approval.ts
│   │   ├── notice.ts
│   │   ├── schedule.ts
│   │   └── tenant.ts
│   ├── router/
│   │   ├── index.ts                # 路由实例 + 导航守卫
│   │   └── routes.ts               # 路由配置（含 SaaS 路由）
│   ├── stores/                     # Pinia 状态管理
│   │   ├── user.ts                 # 用户 + 租户信息
│   │   ├── app.ts                  # 应用全局状态
│   │   ├── approval.ts             # 审批
│   │   ├── notice.ts               # 公告
│   │   ├── schedule.ts             # 日程
│   │   └── tenant.ts               # 租户/套餐/账单
│   ├── styles/
│   │   ├── variables.scss          # SCSS 变量
│   │   ├── reset.scss              # 重置样式
│   │   └── index.scss              # 全局样式
│   ├── types/                      # TypeScript 类型定义
│   │   ├── user.ts                 # 用户类型（含 TenantInfo）
│   │   ├── approval.ts             # 审批类型
│   │   ├── notice.ts               # 公告类型
│   │   ├── schedule.ts             # 日程类型
│   │   ├── common.ts               # 通用类型（ApiResponse, PageResult）
│   │   ├── tenant.ts               # 租户/套餐/账单类型
│   │   └── admin.ts                # 超管类型
│   ├── utils/
│   │   ├── request.ts              # HTTP 封装（Bearer Token + X-Tenant-Slug）
│   │   ├── storage.ts              # 本地存储（Token + 租户信息）
│   │   ├── date.ts                 # 日期工具
│   │   ├── plan.ts                 # 套餐工具（权限检查/格式化）
│   │   └── index.ts                # 通用工具
│   ├── views/                      # 页面组件（26 个，全部完成）
│   ├── App.vue
│   └── main.ts
├── docs/
│   ├── development-plan.md         # 开发计划（详细进度跟踪）
│   └── saas-design.md              # SaaS 架构设计文档
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── .env
```

---

## 快速开始

### 环境要求

- Node.js >= 16
- Go >= 1.22
- MySQL 8.0
- Redis 7（可选）

### 前端

```bash
npm install
npm run dev
```

### 后端

```bash
cd backend

# 配置数据库连接
# 编辑 configs/config.yaml 或创建 configs/config.local.yaml 覆盖

# 运行（首次运行自动建表 + 种子数据）
go run cmd/server/main.go

# 或使用 Makefile
make run
```

### 构建生产版本

```bash
# 前端
npm run build

# 后端
cd backend
make build
```

### 代码检查

```bash
npm run lint
```

---

## 测试账号

种子数据自动创建，默认租户 Slug: `demo`

| 用户名 | 密码 | 角色 | 昵称 |
|--------|------|------|------|
| admin | 123456 | 管理员 | 管理员 |
| zhangsan | 123456 | 部门经理 | 张三 |
| lisi | 123456 | 部门经理 | 李四 |
| wangwu | 123456 | 普通员工 | 王五 |
| zhaoliu | 123456 | 普通员工 | 赵六 |
| sunqi | 123456 | 普通员工 | 孙七 |

---

## 后端 API 接口列表

> 基础路径：`/api/v1`

### 公开端点（无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/tenant/register` | 租户注册 |
| GET | `/plans` | 套餐列表 |

### 认证相关

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/auth/login` | 用户登录 |
| GET | `/auth/info` | 获取当前用户信息 |

### 用户管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/user/list` | 用户列表（分页 + 搜索） |
| POST | `/user` | 创建用户 |
| PUT | `/user/:id` | 更新用户 |
| DELETE | `/user/:id` | 删除用户 |
| PUT | `/user/:id/status` | 启用/禁用 |

### 部门管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/dept/list` | 部门列表（树形） |
| POST | `/dept` | 创建部门 |
| PUT | `/dept/:id` | 更新部门 |
| DELETE | `/dept/:id` | 删除部门 |

### 角色管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/role/list` | 角色列表 |
| POST | `/role` | 创建角色 |
| PUT | `/role/:id` | 更新角色 |
| DELETE | `/role/:id` | 删除角色 |

### 审批管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/approvals` | 发起审批 |
| GET | `/approvals/my` | 我的申请 |
| GET | `/approvals/pending` | 待我审批 |
| GET | `/approvals/done` | 已办审批 |
| GET | `/approvals/stats` | 审批统计 |
| GET | `/approvals/:id` | 审批详情 |
| POST | `/approvals/:id/action` | 审批操作（approve/reject/transfer） |
| POST | `/approvals/:id/withdraw` | 撤回申请 |

### 公告管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/notices` | 公告列表 |
| GET | `/notices/unread-count` | 未读数量 |
| GET | `/notices/:id` | 公告详情（自动标记已读） |
| POST | `/notices` | 发布公告 |
| POST | `/notices/:id/read` | 标记已读 |

### 日程管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/schedules` | 日程列表 |
| GET | `/schedules/week` | 本周日程 |
| GET | `/schedules/:id` | 日程详情 |
| POST | `/schedules` | 创建日程 |
| PUT | `/schedules/:id` | 更新日程 |
| DELETE | `/schedules/:id` | 删除日程 |

### 流程配置

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/flows` | 流程列表 |
| POST | `/flows` | 创建流程 |
| PUT | `/flows/:id` | 更新流程 |
| DELETE | `/flows/:id` | 删除流程 |

### 租户管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/tenant/info` | 租户信息 |
| PUT | `/tenant/info` | 更新租户信息 |
| POST | `/tenant/plan/upgrade` | 升级套餐 |
| GET | `/tenant/invoices` | 账单列表 |

---

## 多租户架构

- **数据隔离**：共享数据库 + `tenant_id` 字段隔离（12 张业务表均含 tenant_id + 索引）
- **租户识别**：支持 `X-Tenant-Slug` 请求头 + 子域名解析
- **GORM TenantScope**：自动过滤当前租户数据
- **JWT 含 TenantID**：Claims 同时携带 UserID 和 TenantID
- **套餐体系**：4 个套餐等级（免费版/基础版/标准版/专业版），含用户数限制
- **缓存方案**：内存缓存（默认）+ Redis（可选，配置 `redis.enabled` 控制）

---

## License

MIT
