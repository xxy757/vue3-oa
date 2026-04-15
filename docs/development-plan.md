# OA 系统 SaaS 化开发计划

> 版本：v2.0
> 创建日期：2026-03-31
> 最后更新：2026-04-03
> 关联文档：[SaaS 架构设计文档](./saas-design.md)

---

## 总览

| 阶段  | 名称             | 状态       | 产出                                              |
| ----- | ---------------- | ---------- | ------------------------------------------------- |
| **A** | 前端页面开发     | ✅ 已完成  | 所有 OA 页面开发完成，对接 Mock 可运行            |
| **B** | Go 后端搭建      | ✅ 已完成  | 多租户后端 API，前后端联调通过                    |
| **C** | 多租户 SaaS 改造 | ✅ 已完成  | 后端多租户 + 前端 SaaS 页面全部完成              |
| **D** | 计费系统         | ⏳ 未开始  | 支付对接、定时账单生成、过期检查                  |
| **E** | 部署上线         | ⏳ 未开始  | Docker 化、域名、监控、上线                       |

## 技术决策

| 决策项     | 选择                                | 实际情况                                      |
| ---------- | ----------------------------------- | --------------------------------------------- |
| 后端语言   | Go + Gin                            | ✅ Go 1.26.1 + Gin 1.12                       |
| 数据隔离   | 共享数据库 + tenant_id 字段         | ✅ 所有业务表均含 tenant_id                    |
| 支付渠道   | 支付宝优先                          | ⏳ 未实现                                     |
| 执行策略   | 先完成单租户 OA，再进行 SaaS 化改造 | ✅ 实际跳过单租户阶段，后端直接实现多租户架构  |
| 缓存方案   | Redis 可选                          | ✅ 内存缓存 + Redis 可插拔                    |

---

## 阶段 A：前端页面开发 — ✅ 已完成

> 前端页面在项目创建时已开发完毕，包含完整的 OA 业务页面

### 已完成页面清单

| 页面         | 文件路径                       | 状态 |
| ------------ | ------------------------------ | ---- |
| 登录页       | `views/login/index.vue`        | ✅   |
| 工作台       | `views/dashboard/index.vue`    | ✅   |
| 发起申请     | `views/approval/Apply.vue`     | ✅   |
| 我的申请     | `views/approval/MyApply.vue`   | ✅   |
| 待我审批     | `views/approval/Pending.vue`   | ✅   |
| 已办审批     | `views/approval/Done.vue`      | ✅   |
| 审批详情     | `views/approval/Detail.vue`    | ✅   |
| 公告列表     | `views/notice/List.vue`        | ✅   |
| 公告详情     | `views/notice/Detail.vue`      | ✅   |
| 日程日历     | `views/schedule/Calendar.vue`  | ✅   |
| 日程列表     | `views/schedule/List.vue`      | ✅   |
| 用户管理     | `views/system/User.vue`        | ✅   |
| 部门管理     | `views/system/Dept.vue`        | ✅   |
| 角色管理     | `views/system/Role.vue`        | ✅   |
| 流程配置     | `views/system/Flow.vue`        | ✅   |
| 个人信息     | `views/profile/Info.vue`       | ✅   |
| 修改密码     | `views/profile/Password.vue`   | ✅   |
| 404 页面     | `views/error/404.vue`          | ✅   |

### 前端 Store 清单

| Store         | 文件               | 功能                                          |
| ------------- | ------------------ | --------------------------------------------- |
| userStore     | `stores/user.ts`   | 登录/登出、用户信息、租户信息、权限判断       |
| approvalStore | `stores/approval.ts` | 审批 CRUD、审批操作、审批列表               |
| noticeStore   | `stores/notice.ts` | 公告列表、详情、创建、已读、未读数            |
| scheduleStore | `stores/schedule.ts` | 日程 CRUD、周日程、按日期筛选              |
| appStore      | `stores/app.ts`    | 侧边栏状态、面包屑、全局 loading              |

### 前端工具/配置

| 模块        | 文件              | 说明                                                    |
| ----------- | ----------------- | ------------------------------------------------------- |
| HTTP 封装   | `utils/request.ts` | Axios 封装，自动附加 Bearer Token + X-Tenant-Slug，401 自动跳转登录 |
| 存储        | `utils/storage.ts` | Token、用户信息、租户信息的 localStorage 封装            |
| 日期工具    | `utils/date.ts`   | getMonthDates / getWeekDates 等日历工具函数              |
| 路由        | `router/routes.ts` | constantRoutes + asyncRoutes，admin 角色可见系统管理     |
| 路由守卫    | `router/index.ts`  | Token 校验、白名单、未登录跳转                           |
| 布局组件    | `components/Layout/OALayout.vue` | 侧边栏菜单 + 顶部导航 + 面包屑 + 用户信息 |

---

## 阶段 B：Go 后端搭建 — ✅ 已完成

> 后端目录位于 `backend/`，采用多租户架构直接实现（未走单租户过渡路线）

### B-1：项目初始化 ✅

| 项目     | 详情                                                       |
| -------- | ---------------------------------------------------------- |
| Go 版本  | 1.26.1                                                     |
| 模块名   | `oa-saas`                                                  |
| Web 框架 | Gin 1.12                                                   |
| ORM      | GORM 1.31                                                  |
| 依赖     | gin, gorm, golang-jwt/jwt/v5, go-redis/redis/v9(可选), yaml.v3 |
| 配置管理 | `configs/config.yaml` + `configs/config.local.yaml` 覆盖 + 环境变量 |
| 构建工具 | Makefile（build/run/clean）                                |

### B-2：数据库 ✅

| 项目         | 详情                                                       |
| ------------ | ---------------------------------------------------------- |
| 建表方式     | GORM AutoMigrate（14 个模型自动建表）                      |
| 数据库       | MySQL 8.0，库名 `oa_saas`                                  |
| 多租户       | 所有业务表均含 `tenant_id` 字段 + 索引                     |
| 种子数据     | 4 个套餐 + 1 个示例租户(demo) + 3 个角色 + 3 个部门 + 6 个用户 + 5 个审批流程 |
| SQL 迁移文件 | `migrations/001_init.up.sql`（早期单租户版本，仅作参考）   |

**已建表清单（14 个）：**

| 模型文件           | 表名                    | 含 tenant_id |
| ------------------ | ----------------------- | ------------ |
| `model/tenant.go`  | `plans`                 | -            |
|                    | `tenants`               | -            |
|                    | `invoices`              | ✅           |
|                    | `tenant_logs`           | ✅           |
| `model/user.go`    | `users`                 | ✅           |
| `model/department.go` | `departments`        | ✅           |
| `model/role.go`    | `roles`                 | ✅           |
| `model/approval_flow.go` | `approval_flows`  | ✅           |
| `model/approval.go` | `approvals`            | ✅           |
|                    | `approval_nodes`        | ✅           |
| `model/notice.go`  | `notices`               | ✅           |
|                    | `notice_reads`          | ✅           |
| `model/schedule.go` | `schedules`            | ✅           |
|                    | `schedule_participants` | ✅           |

### B-3：认证模块 ✅

| 功能         | 实现                                                        |
| ------------ | ----------------------------------------------------------- |
| JWT 工具     | `pkg/jwt/jwt.go`：HS256 签名，Claims 含 UserID + TenantID |
| 密码加密     | `pkg/utils/password.go`：bcrypt Hash + Verify               |
| 登录接口     | `POST /api/v1/auth/login` → 返回 token + user + tenant     |
| 获取用户信息 | `GET /api/v1/auth/info`（需 JWT）                           |
| Auth 中间件  | 解析 Bearer Token → 注入 user_id + tenant_id 到 context    |
| CORS 中间件  | 允许前端跨域                                                |
| Logger 中间件 | 请求日志记录                                                |

### B-4：用户/部门/角色 CRUD ✅

| 功能       | API                          | 说明                                     |
| ---------- | ---------------------------- | ---------------------------------------- |
| 用户列表   | `GET /api/v1/user/list`      | 分页 + 关键词搜索                        |
| 创建用户   | `POST /api/v1/user`          | 含套餐用户数上限校验                     |
| 更新用户   | `PUT /api/v1/user/:id`       |                                          |
| 删除用户   | `DELETE /api/v1/user/:id`    |                                          |
| 切换状态   | `PUT /api/v1/user/:id/status` | 启用/禁用                               |
| 部门列表   | `GET /api/v1/dept/list`      | 树形结构                                 |
| 部门 CRUD  | POST/PUT/DELETE `/api/v1/dept` | 有子部门时拒绝删除                     |
| 角色列表   | `GET /api/v1/role/list`      |                                          |
| 角色 CRUD  | POST/PUT/DELETE `/api/v1/role` | 管理员角色不可删除                     |
| 审批流程   | CRUD `/api/v1/flows`         | 审批流程配置管理                         |

### B-5：审批模块 ✅

| 功能     | API                                   | 说明                                     |
| -------- | ------------------------------------- | ---------------------------------------- |
| 发起申请 | `POST /api/v1/approvals`              | 自动匹配审批流程，创建审批+节点          |
| 我的申请 | `GET /api/v1/approvals/my`            | 分页 + 类型/状态筛选                     |
| 待我审批 | `GET /api/v1/approvals/pending`       | 当前用户为审批人且状态 pending            |
| 已办审批 | `GET /api/v1/approvals/done`          |                                          |
| 审批统计 | `GET /api/v1/approvals/stats`         | 各状态数量 + 待我审批数                  |
| 审批详情 | `GET /api/v1/approvals/:id`           | 申请信息 + 流程节点列表                  |
| 审批操作 | `POST /api/v1/approvals/:id/action`   | approve/reject/transfer                  |
| 撤回     | `POST /api/v1/approvals/:id/withdraw` | 仅 pending 状态可撤回                    |

### B-6：公告模块 ✅

| 功能     | API                                | 说明                         |
| -------- | ---------------------------------- | ---------------------------- |
| 公告列表 | `GET /api/v1/notices`              | 分页 + 类型/关键词筛选       |
| 未读数   | `GET /api/v1/notices/unread-count` |                              |
| 公告详情 | `GET /api/v1/notices/:id`          | 自动标记已读                 |
| 发布公告 | `POST /api/v1/notices`             |                              |
| 标记已读 | `POST /api/v1/notices/:id/read`    |                              |

### B-7：日程模块 ✅

| 功能     | API                            | 说明               |
| -------- | ------------------------------ | ------------------ |
| 日程列表 | `GET /api/v1/schedules`        | 日期范围 + 参与人  |
| 周日程   | `GET /api/v1/schedules/week`   |                    |
| 日程详情 | `GET /api/v1/schedules/:id`    | 含参与人           |
| 创建日程 | `POST /api/v1/schedules`       | 含参与人           |
| 更新日程 | `PUT /api/v1/schedules/:id`    |                    |
| 删除日程 | `DELETE /api/v1/schedules/:id` |                    |

### B-8：前后端联调 ✅

| 任务         | 详情                                                        |
| ------------ | ----------------------------------------------------------- |
| 移除 Mock    | `vite.config.ts` 中移除 `viteMockServe` 插件               |
| 配置代理     | `vite.config.ts` 添加 `proxy: '/api/v1' → http://127.0.0.1:8080` |
| 修改 baseURL | `utils/request.ts` 中 baseURL 改为 `/api/v1`                |
| Token 注入   | 请求拦截器自动附加 `Authorization: Bearer {token}`          |
| 租户标识     | 请求拦截器自动附加 `X-Tenant-Slug` 请求头                   |
| 错误处理     | 401 自动跳转登录、403/404/500 错误提示                      |
| 构建验证     | `npm run build` 零错误通过                                  |

### 后端完整文件清单

```
backend/
├── cmd/server/main.go                 ✅ 入口（DB初始化 + Redis可插拔 + AutoMigrate + Seed + 启动服务）
├── configs/
│   ├── config.yaml                    ✅ 基础配置
│   └── config.local.yaml              ✅ 本地覆盖配置
├── internal/
│   ├── config/config.go               ✅ 配置加载（YAML + 环境变量覆盖）
│   ├── middleware/
│   │   ├── auth.go                    ✅ JWT 认证中间件（含 RateLimit）
│   │   ├── cors.go                    ✅ 跨域中间件
│   │   ├── logger.go                  ✅ 请求日志中间件
│   │   └── tenant.go                  ✅ 租户识别中间件（Header + 子域名）
│   ├── model/
│   │   ├── tenant.go                  ✅ Plan/Tenant/Invoice/TenantLog 模型
│   │   ├── user.go                    ✅ 用户模型（含 tenant_id）
│   │   ├── department.go              ✅ 部门模型（含 tenant_id）
│   │   ├── role.go                    ✅ 角色模型（JSON 数组类型 + tenant_id）
│   │   ├── approval_flow.go           ✅ 审批流程模型（JSON 节点 + tenant_id）
│   │   ├── approval.go                ✅ 审批 + 审批节点模型（含 tenant_id）
│   │   ├── notice.go                  ✅ 公告 + 公告阅读模型（含 tenant_id）
│   │   ├── schedule.go                ✅ 日程 + 参与人模型（含 tenant_id）
│   │   ├── scopes.go                  ✅ 租户 Context 存取 + TenantScope
│   │   └── seed.go                    ✅ 种子数据（套餐/租户/角色/部门/用户/流程）
│   ├── handler/
│   │   ├── auth_handler.go            ✅ 登录 + 获取用户信息
│   │   ├── user_handler.go            ✅ 用户 CRUD + 状态切换
│   │   ├── dept_handler.go            ✅ 部门树 + CRUD
│   │   ├── role_handler.go            ✅ 角色 CRUD
│   │   ├── approval_handler.go        ✅ 审批全流程
│   │   ├── notice_handler.go          ✅ 公告 CRUD + 未读 + 标记已读
│   │   ├── schedule_handler.go        ✅ 日程 CRUD + 周视图
│   │   ├── flow_handler.go            ✅ 审批流程 CRUD
│   │   └── tenant_handler.go          ✅ 租户注册/信息/升级/账单
│   ├── router/
│   │   └── router.go                  ✅ 路由注册（37 个 API 端点）
│   └── pkg/
│       ├── jwt/jwt.go                 ✅ JWT 生成与解析（含 TenantID）
│       ├── cache/
│       │   ├── cache.go               ✅ 缓存接口定义
│       │   ├── memory.go              ✅ 内存缓存实现（默认）
│       │   └── redis.go               ✅ Redis 缓存实现（可选）
│       └── utils/
│           ├── password.go            ✅ bcrypt 密码加密
│           └── response.go            ✅ 统一 API 响应
├── migrations/
│   ├── 001_init.up.sql                ✅ 早期建表 SQL（仅参考，已被 AutoMigrate 取代）
│   └── 001_init.down.sql              ✅ 回滚 SQL
├── Makefile                            ✅ 构建/运行/清理
├── go.mod                              ✅ Go 1.26.1
└── go.sum                              ✅ 依赖校验
```

### 后端 API 端点完整清单（37 个）

#### 公开端点（无需认证）

| 方法 | 路径                          | 说明                     |
| ---- | ----------------------------- | ------------------------ |
| POST | `/api/v1/tenant/register`     | 租户注册                 |
| GET  | `/api/v1/plans`               | 套餐列表                 |

#### 需租户识别 + 登录

| 方法 | 路径                          | 说明                     |
| ---- | ----------------------------- | ------------------------ |
| POST | `/api/v1/auth/login`          | 用户登录                 |
| GET  | `/api/v1/auth/info`           | 获取当前用户信息         |

#### 需租户识别 + JWT 认证

| 方法   | 路径                                    | 说明             |
| ------ | --------------------------------------- | ---------------- |
| GET    | `/api/v1/user/list`                     | 用户列表         |
| POST   | `/api/v1/user`                          | 创建用户         |
| PUT    | `/api/v1/user/:id`                      | 更新用户         |
| DELETE | `/api/v1/user/:id`                      | 删除用户         |
| PUT    | `/api/v1/user/:id/status`               | 切换用户状态     |
| GET    | `/api/v1/dept/list`                     | 部门列表（树形） |
| POST   | `/api/v1/dept`                          | 创建部门         |
| PUT    | `/api/v1/dept/:id`                      | 更新部门         |
| DELETE | `/api/v1/dept/:id`                      | 删除部门         |
| GET    | `/api/v1/role/list`                     | 角色列表         |
| POST   | `/api/v1/role`                          | 创建角色         |
| PUT    | `/api/v1/role/:id`                      | 更新角色         |
| DELETE | `/api/v1/role/:id`                      | 删除角色         |
| POST   | `/api/v1/approvals`                     | 发起审批         |
| GET    | `/api/v1/approvals/my`                  | 我的申请         |
| GET    | `/api/v1/approvals/pending`             | 待我审批         |
| GET    | `/api/v1/approvals/done`                | 已办审批         |
| GET    | `/api/v1/approvals/stats`               | 审批统计         |
| GET    | `/api/v1/approvals/:id`                 | 审批详情         |
| POST   | `/api/v1/approvals/:id/action`          | 审批操作         |
| POST   | `/api/v1/approvals/:id/withdraw`        | 撤回申请         |
| GET    | `/api/v1/notices`                       | 公告列表         |
| GET    | `/api/v1/notices/unread-count`          | 未读数           |
| GET    | `/api/v1/notices/:id`                   | 公告详情         |
| POST   | `/api/v1/notices`                       | 发布公告         |
| POST   | `/api/v1/notices/:id/read`              | 标记已读         |
| GET    | `/api/v1/schedules`                     | 日程列表         |
| GET    | `/api/v1/schedules/week`                | 周日程           |
| GET    | `/api/v1/schedules/:id`                 | 日程详情         |
| POST   | `/api/v1/schedules`                     | 创建日程         |
| PUT    | `/api/v1/schedules/:id`                 | 更新日程         |
| DELETE | `/api/v1/schedules/:id`                 | 删除日程         |
| GET    | `/api/v1/flows`                         | 流程列表         |
| POST   | `/api/v1/flows`                         | 创建流程         |
| PUT    | `/api/v1/flows/:id`                     | 更新流程         |
| DELETE | `/api/v1/flows/:id`                     | 删除流程         |
| GET    | `/api/v1/tenant/info`                   | 租户信息         |
| PUT    | `/api/v1/tenant/info`                   | 更新租户信息     |
| POST   | `/api/v1/tenant/plan/upgrade`           | 升级套餐         |
| GET    | `/api/v1/tenant/invoices`               | 账单列表         |

### 技术要点

- Go 版本：1.26.1，使用 `goproxy.cn` 国内代理
- Redis 可插拔：`config.yaml` 中 `redis.enabled` 控制，未配置时自动降级为内存缓存
- 缓存接口：`Cache` 接口 + `MemoryCache`（默认）+ `RedisCache`（可选）
- 认证：JWT HS256，Claims 含 UserID + TenantID，24h 过期
- 多租户：共享数据库 + tenant_id 隔离，TenantMiddleware 自动识别
- 数据库：GORM AutoMigrate 自动建表，所有业务模型含 tenant_id
- 种子数据：首次运行自动初始化，含 1 个示例租户（slug=demo，标准版，14天试用）

---

## 阶段 C：多租户 SaaS 改造 — ✅ 已完成

### 后端多租户 — ✅ 已完成

| 功能             | 状态 | 实现详情                                                   |
| ---------------- | ---- | ---------------------------------------------------------- |
| 租户模型         | ✅   | Plan / Tenant / Invoice / TenantLog 四表                   |
| 业务表 tenant_id | ✅   | 12 张业务表均含 tenant_id 字段 + 索引                      |
| 租户识别中间件   | ✅   | 支持 X-Tenant-Slug Header + 子域名解析                     |
| GORM TenantScope | ✅   | `model/scopes.go`：context 存取 + TenantScope 查询过滤     |
| 租户注册         | ✅   | `POST /api/v1/tenant/register`，自动创建租户+角色+管理员   |
| 租户管理         | ✅   | 信息查看/编辑、套餐升级、账单列表                          |
| 套餐体系         | ✅   | 4 个套餐等级 + 功能特性 JSON + 升级接口                    |
| JWT 含 tenant_id | ✅   | Claims 同时携带 UserID 和 TenantID                         |
| 种子数据         | ✅   | 含示例租户 + 4 个套餐                                      |

### 前端 SaaS 页面 — ✅ 已完成

| 任务             | 涉及文件                            | 状态 | 说明                                       |
| ---------------- | ----------------------------------- | ---- | ------------------------------------------ |
| 新增 types       | `types/tenant.ts`, `types/admin.ts` | ✅   | 租户/套餐/账单/仪表盘类型定义              |
| 新增 stores      | `stores/tenant.ts`                  | ✅   | 租户状态管理（信息/套餐/账单/权限检查）    |
| 新增 utils       | `utils/plan.ts`                     | ✅   | 套餐权限检查函数 + 格式化工具              |
| 企业注册页       | `views/auth/Register.vue`           | ✅   | 3 步向导：企业信息 → 选择套餐 → 完成注册  |
| 选择套餐页       | `views/auth/ChoosePlan.vue`         | ✅   | 套餐对比卡片 + 升级确认                    |
| 租户信息页       | `views/tenant/Info.vue`             | ✅   | 企业信息查看/编辑 + 用户数进度条           |
| 套餐管理页       | `views/tenant/Plan.vue`             | ✅   | 当前套餐 + 功能对比表 + 升级弹窗           |
| 账单管理页       | `views/tenant/Invoices.vue`         | ✅   | 账单列表 + 详情弹窗                        |
| 超管后台布局     | `components/Layout/AdminLayout.vue` | ✅   | 独立布局（仪表盘/租户/套餐）               |
| 超管仪表盘       | `views/admin/Dashboard.vue`         | ✅   | 租户数/活跃/用户/收入统计 + 分布图         |
| 超管租户管理     | `views/admin/Tenants.vue`           | ✅   | 租户列表 + 状态管理 + CRUD                 |
| 超管套餐管理     | `views/admin/Plans.vue`             | ✅   | 套餐配置 + 功能特性编辑                    |
| 路由改造         | `router/routes.ts`                  | ✅   | 新增 SaaS 路由 + 租户/超管路由守卫         |
| 路由守卫改造     | `router/index.ts`                   | ✅   | 白名单添加 /register                       |
| OALayout 菜单    | `components/Layout/OALayout.vue`    | ✅   | 新增企业管理菜单组                          |
| Mock 数据        | `mock/tenant.ts`                    | ✅   | 15 个 SaaS Mock 端点                       |
| Vite 配置        | `vite.config.ts`                    | ✅   | 添加 viteMockServe 插件                    |

> **注意**：前端 `stores/user.ts` 和 `utils/request.ts` 已预置多租户支持：
> - `login()` 返回结果中已包含 tenant 信息并持久化
> - `request.ts` 拦截器已自动附加 `X-Tenant-Slug` 请求头
> - `storage.ts` 已有 `setTenantSlug` / `getTenantSlug` 等工具函数

---

## 阶段 D：计费系统 — ⏳ 未开始

### D-1：计费服务

| 任务         | 详情                                         | 状态 |
| ------------ | -------------------------------------------- | ---- |
| 月度账单生成 | 定时任务，每月 1 日统计活跃用户数×单价       | ⏳   |
| 过期检查     | 检查 `plan_expire_at`，过期则 `suspended`    | ⏳   |
| 支付宝对接   | 创建支付→回调处理→更新账单状态               | ⏳   |

> **注**：后端已有 `invoices` 表和 `tenant_handler.go` 中的 `UpgradePlan`/`ListInvoices` 基础接口，定时任务和支付对接待实现。

### D-2：前端支付流程

| 任务           | 详情                                       | 状态 |
| -------------- | ------------------------------------------ | ---- |
| 账单列表页     | 状态标签 + 支付按钮                        | ⏳   |
| 支付跳转/结果  | 跳转支付宝→回调→展示结果                   | ⏳   |
| 续费提醒       | 套餐即将过期时头部显示警告条                | ⏳   |

---

## 阶段 E：部署上线 — ⏳ 未开始

| 任务            | 详情                                         | 状态 |
| --------------- | -------------------------------------------- | ---- |
| 后端 Dockerfile | 多阶段构建（golang:1.26-alpine → alpine）    | ⏳   |
| docker-compose  | Nginx + Go API + MySQL 8.0 + Redis 7（可选） | ⏳   |
| Nginx 配置      | 通配符子域名 → X-Tenant-Slug Header + SSL    | ⏳   |
| CI/CD           | GitHub Actions：自动 build + 部署            | ⏳   |
| 监控            | Prometheus + Grafana                         | ⏳   |

---

## 文件变更清单

### 前端（vue3-oa）— 当前实际文件

```
src/
├── views/
│   ├── admin/
│   │   ├── Dashboard.vue           ✅ 超管仪表盘（统计卡片+分布+最近租户）
│   │   ├── Plans.vue               ✅ 超管套餐管理（套餐配置+功能编辑）
│   │   └── Tenants.vue             ✅ 超管租户管理（租户列表+状态+CRUD）
│   ├── approval/
│   │   ├── Apply.vue              ✅ 发起申请（5种类型表单）
│   │   ├── Detail.vue             ✅ 审批详情（时间线 + 操作区）
│   │   ├── Done.vue               ✅ 已办审批列表
│   │   ├── MyApply.vue            ✅ 我的申请列表
│   │   └── Pending.vue            ✅ 待我审批列表
│   ├── auth/
│   │   ├── ChoosePlan.vue          ✅ 套餐选择页（套餐卡片+升级确认）
│   │   └── Register.vue            ✅ 企业注册页（3步向导）
│   ├── dashboard/
│   │   └── index.vue              ✅ 工作台（统计卡片+待办+快捷入口+图表）
│   ├── error/
│   │   └── 404.vue                ✅ 404 页面
│   ├── login/
│   │   └── index.vue              ✅ 登录页
│   ├── notice/
│   │   ├── Detail.vue             ✅ 公告详情
│   │   └── List.vue               ✅ 公告列表
│   ├── profile/
│   │   ├── Info.vue               ✅ 个人信息
│   │   └── Password.vue           ✅ 修改密码
│   ├── schedule/
│   │   ├── Calendar.vue           ✅ 日程日历视图
│   │   └── List.vue               ✅ 日程列表视图
│   ├── system/
│   │   ├── Dept.vue               ✅ 部门管理
│   │   ├── Flow.vue               ✅ 流程配置
│   │   ├── Role.vue               ✅ 角色管理
│   │   └── User.vue               ✅ 用户管理
│   └── tenant/
│       ├── Invoices.vue            ✅ 账单管理（列表+详情弹窗）
│       ├── Info.vue                ✅ 企业信息（查看/编辑+用户数进度）
│       └── Plan.vue                ✅ 套餐管理（当前套餐+对比+升级）
├── components/Layout/
│   ├── AdminLayout.vue             ✅ 超管后台布局（侧边栏+头部）
│   └── OALayout.vue               ✅ 主布局（侧边栏+头部+面包屑+企业管理菜单）
├── stores/
│   ├── app.ts                     ✅ 应用状态
│   ├── approval.ts                ✅ 审批 Store
│   ├── notice.ts                  ✅ 公告 Store
│   ├── schedule.ts                ✅ 日程 Store
│   ├── tenant.ts                   ✅ 租户 Store（信息/套餐/账单/权限）
│   └── user.ts                    ✅ 用户 Store（含租户信息）
├── router/
│   ├── index.ts                   ✅ 路由实例 + 守卫（含 /register 白名单）
│   └── routes.ts                  ✅ 路由配置（含 SaaS 路由）
├── types/
│   ├── admin.ts                    ✅ 超管类型定义（DashboardStats）
│   ├── approval.ts                ✅ 审批类型定义
│   ├── common.ts                  ✅ 通用类型（ApiResponse, PageResult）
│   ├── notice.ts                  ✅ 公告类型定义
│   ├── schedule.ts                ✅ 日程类型定义
│   ├── tenant.ts                   ✅ 租户类型定义（Plan/Tenant/Invoice等）
│   ├── user.d.ts                  ✅ 用户类型声明
│   └── user.ts                    ✅ 用户类型定义（含 TenantInfo）
├── utils/
│   ├── date.ts                    ✅ 日期工具
│   ├── index.ts                   ✅ 通用工具
│   ├── plan.ts                     ✅ 套餐工具（标签/状态映射/权限检查）
│   ├── request.ts                 ✅ HTTP 封装（含租户 Header）
│   └── storage.ts                 ✅ 本地存储（含租户 Slug）
├── styles/
│   ├── index.scss                 ✅ 全局样式
│   ├── reset.scss                 ✅ 重置样式
│   └── variables.scss             ✅ SCSS 变量
├── mock/
│   ├── approval.ts                📦 审批 Mock（联调后停用）
│   ├── notice.ts                  📦 公告 Mock（联调后停用）
│   ├── schedule.ts                📦 日程 Mock（联调后停用）
│   ├── tenant.ts                   📦 SaaS Mock（15个端点，开发阶段使用）
│   └── user.ts                    📦 用户 Mock（联调后停用）
├── App.vue                        ✅ 根组件
├── main.ts                        ✅ 入口
└── env.d.ts                       ✅ 环境类型声明
```

### 后端（backend/）

```
backend/
├── cmd/server/main.go             ✅ 入口文件
├── configs/
│   ├── config.yaml                ✅ 基础配置
│   └── config.local.yaml          ✅ 本地覆盖
├── internal/
│   ├── config/config.go           ✅ 配置加载
│   ├── middleware/
│   │   ├── auth.go                ✅ JWT + RateLimit
│   │   ├── cors.go                ✅ 跨域
│   │   ├── logger.go              ✅ 日志
│   │   └── tenant.go              ✅ 租户识别
│   ├── model/
│   │   ├── tenant.go              ✅ 租户体系模型
│   │   ├── user.go                ✅ 用户模型
│   │   ├── department.go          ✅ 部门模型
│   │   ├── role.go                ✅ 角色模型
│   │   ├── approval_flow.go       ✅ 审批流程模型
│   │   ├── approval.go            ✅ 审批模型
│   │   ├── notice.go              ✅ 公告模型
│   │   ├── schedule.go            ✅ 日程模型
│   │   ├── scopes.go              ✅ 租户 Scope
│   │   └── seed.go                ✅ 种子数据
│   ├── handler/
│   │   ├── auth_handler.go        ✅ 认证处理器
│   │   ├── user_handler.go        ✅ 用户处理器
│   │   ├── dept_handler.go        ✅ 部门处理器
│   │   ├── role_handler.go        ✅ 角色处理器
│   │   ├── approval_handler.go    ✅ 审批处理器
│   │   ├── notice_handler.go      ✅ 公告处理器
│   │   ├── schedule_handler.go    ✅ 日程处理器
│   │   ├── flow_handler.go        ✅ 流程处理器
│   │   └── tenant_handler.go      ✅ 租户处理器
│   ├── router/
│   │   └── router.go              ✅ 路由注册
│   └── pkg/
│       ├── jwt/jwt.go             ✅ JWT 工具
│       ├── cache/
│       │   ├── cache.go           ✅ 缓存接口
│       │   ├── memory.go          ✅ 内存缓存
│       │   └── redis.go           ✅ Redis 缓存
│       └── utils/
│           ├── password.go        ✅ 密码加密
│           └── response.go        ✅ 统一响应
├── migrations/
│   ├── 001_init.up.sql            ✅ 建表 SQL（参考）
│   └── 001_init.down.sql          ✅ 回滚 SQL
├── Makefile                        ✅ 构建/运行
├── go.mod                          ✅ 模块定义
└── go.sum                          ✅ 依赖校验
```

---

## 变更记录

| 日期       | 版本 | 变更内容                                                                                                    |
| ---------- | ---- | ----------------------------------------------------------------------------------------------------------- |
| 2026-03-31 | v1.0 | 初始版本，5 阶段开发计划                                                                                    |
| 2026-03-31 | v1.1 | Redis 可插拔改造：Redis 改为可选组件，未配置时自动降级为内存缓存                                            |
| 2026-03-31 | v1.2 | 阶段 A 全部完成（前端页面已有），修复全部 TS 编译错误，`npm run build` 通过                                 |
| 2026-03-31 | v1.3 | B-1 完成：Go 项目初始化；B-3 完成：认证模块；`go build` 通过                                                |
| 2026-03-31 | v1.4 | B-2~B-7 完成：全部后端模块开发完成；B-8 完成：前后端联调（移除 Mock、代理配置、API 路径对齐、统一响应格式） |
| 2026-04-03 | v2.0 | 文档系统性修订：反映实际实现状态（后端直接采用多租户架构）；更新 API 端点清单（37 个）；完善文件清单；阶段 B 标记为已完成；阶段 C 拆分为后端已完成/前端待开发 |
| 2026-04-03 | v2.1 | 阶段 C 前端 SaaS 页面全部完成：14 个新文件（types/stores/utils/8个视图/1个布局/1个Mock）+ 6 个修改文件；代码审查通过（5 严重+12 警告已修复）；`vue-tsc` + `npm run build` 零错误 |

---

## 开发进度

### 阶段 A：前端页面开发 — ✅ 已完成

### 阶段 B：Go 后端搭建 — ✅ 已完成

| 子阶段 | 名称                | 状态    | 说明                                                      |
| ------ | ------------------- | ------- | --------------------------------------------------------- |
| B-1    | 项目初始化          | ✅ 完成 | Go 1.26.1 + Gin + GORM + 目录结构                         |
| B-2    | 数据库建表+种子数据 | ✅ 完成 | 14 个模型 AutoMigrate + 种子数据                          |
| B-3    | 认证模块            | ✅ 完成 | JWT + bcrypt + Auth/CORS/Logger 中间件                    |
| B-4    | 用户/部门/角色 CRUD | ✅ 完成 | 含状态切换 + 套餐用户数校验                               |
| B-5    | 审批模块            | ✅ 完成 | 审批全流程 + 自动匹配审批流程                             |
| B-6    | 公告模块            | ✅ 完成 | 公告 CRUD + 未读统计 + 自动标记已读                       |
| B-7    | 日程模块            | ✅ 完成 | 日程 CRUD + 周视图 + 参与人                               |
| B-8    | 前后端联调          | ✅ 完成 | 移除 Mock 插件 + 配置 proxy + baseURL 改为 /api/v1        |

### 阶段 C：多租户 SaaS 改造 — ✅ 已完成

| 子阶段   | 名称             | 状态       | 说明                                              |
| -------- | ---------------- | ---------- | ------------------------------------------------- |
| C-后端   | 多租户后端       | ✅ 完成    | 租户模型/中间件/Scope/注册/管理接口全部实现        |
| C-前端   | SaaS 前端页面    | ✅ 完成    | 注册页/套餐页/租户管理/超管后台 + Mock 全部完成    |

### 阶段 D：计费系统 — ⏳ 未开始

### 阶段 E：部署上线 — ⏳ 未开始
