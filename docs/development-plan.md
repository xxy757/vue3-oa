# OA 系统 SaaS 化开发计划

> 版本：v1.0
> 创建日期：2026-03-31
> 关联文档：[SaaS 架构设计文档](./saas-design.md)

---

## 总览

| 阶段  | 名称             | 周期   | 产出                                   |
| ----- | ---------------- | ------ | -------------------------------------- |
| **A** | 前端页面开发     | 2 周   | 所有 OA 页面开发完成，对接 Mock 可运行 |
| **B** | Go 后端搭建      | 2 周   | 单租户后端 API，前后端联调通过         |
| **C** | 多租户 SaaS 改造 | 2 周   | 多租户支持、租户注册/管理              |
| **D** | 计费系统         | 1.5 周 | 套餐、账单、支付宝对接                 |
| **E** | 部署上线         | 1 周   | Docker 化、域名、监控、上线            |

## 技术决策

| 决策项   | 选择                                |
| -------- | ----------------------------------- |
| 后端语言 | Go + Gin                            |
| 数据隔离 | 共享数据库 + tenant_id 字段         |
| 支付渠道 | 支付宝优先                          |
| 执行策略 | 先完成单租户 OA，再进行 SaaS 化改造 |

---

## 阶段 A：前端页面开发（2 周）

> 前提：所有页面目前都是空壳，stores/types/mock/request 已就绪

### A-1：登录页 + 布局完善（第 1 天）

| 任务         | 涉及文件                         | 详情                                                                  |
| ------------ | -------------------------------- | --------------------------------------------------------------------- |
| 登录页       | `views/login/index.vue`          | 用户名/密码表单 + 记住我 + 调用 `userStore.login()` + 错误提示 + 跳转 |
| 布局优化     | `components/Layout/OALayout.vue` | 确认侧边栏菜单渲染、头部用户信息、面包屑、退出登录                    |
| 404 页面     | `views/error/404.vue`            | 简单的 404 提示 + 返回首页按钮                                        |
| 路由守卫验证 | `router/index.ts`                | 确认 token 校验、白名单、权限路由过滤正常                             |

### A-2：工作台（第 2 天）

| 任务     | 涉及文件                    | 详情                                                             |
| -------- | --------------------------- | ---------------------------------------------------------------- |
| 统计卡片 | `views/dashboard/index.vue` | 待审批数、我的申请数、未读公告数、本周日程数（用 `n-statistic`） |
| 待办列表 | 同上                        | 最近 5 条待我审批 + 我的申请中 pending 的                        |
| 快捷入口 | 同上                        | 发起申请、发布公告、创建日程的快捷按钮                           |
| 通知轮播 | 同上                        | 未读公告标题滚动展示                                             |

### A-3：审批中心 - 发起申请（第 3 天）

| 任务         | 涉及文件                   | 详情                                                                          |
| ------------ | -------------------------- | ----------------------------------------------------------------------------- |
| 申请类型选择 | `views/approval/Apply.vue` | 5 种类型 Tab（请假/报销/加班/出差/通用），参考 `types/approval.ts` 的类型定义 |
| 请假表单     | 同上                       | 请假类型下拉 + 起止日期 + 天数计算 + 原因                                     |
| 报销表单     | 同上                       | 报销类型 + 金额 + 说明 + 附件上传                                             |
| 加班表单     | 同上                       | 日期 + 起止时间 + 小时数计算 + 原因                                           |
| 出差表单     | 同上                       | 目的地 + 起止日期 + 天数 + 原因 + 预算                                        |
| 通用审批表单 | 同上                       | 内容 + 附件                                                                   |
| 提交逻辑     | 同上                       | 调用 `approvalStore.createApproval()`，成功后跳转我的申请                     |

### A-4：审批中心 - 列表页（第 4 天）

| 任务     | 涉及文件                     | 详情                                                  |
| -------- | ---------------------------- | ----------------------------------------------------- |
| 我的申请 | `views/approval/MyApply.vue` | `n-data-table` 列表 + 状态/类型筛选 + 分页 + 跳转详情 |
| 待我审批 | `views/approval/Pending.vue` | 同上结构 + 审批操作按钮（通过/驳回/转交）             |
| 已办审批 | `views/approval/Done.vue`    | 同上结构，只读展示                                    |

### A-5：审批详情（第 5 天）

| 任务           | 涉及文件                    | 详情                                                        |
| -------------- | --------------------------- | ----------------------------------------------------------- |
| 申请信息展示   | `views/approval/Detail.vue` | 根据 `type` 动态展示不同表单字段                            |
| 审批流程时间线 | 同上                        | `n-timeline` 展示各节点：节点名称、审批人、状态、意见、时间 |
| 审批操作区     | 同上                        | 通过/驳回/转交的对话框 + 意见输入（仅待我审批的显示）       |
| 撤回按钮       | 同上                        | 仅 pending 状态且是自己的申请可撤回                         |

### A-6：公告通知（第 6 天）

| 任务     | 涉及文件                       | 详情                                                              |
| -------- | ------------------------------ | ----------------------------------------------------------------- |
| 公告列表 | `views/notice/List.vue`        | `n-data-table` + 类型筛选 + 关键词搜索 + 置顶标识 + 已读/未读标记 |
| 公告详情 | 新增 `views/notice/Detail.vue` | 富文本内容展示 + 阅读数 + 发布者信息 + 自动标记已读               |
| 发布公告 | `views/notice/List.vue` 内弹窗 | 标题 + 类型 + 内容（`n-input type=textarea`）+ 摘要 + 封面图      |
| 未读徽章 | 侧边栏/头部                    | 调用 `noticeStore.getUnreadCount()` 显示红点                      |

### A-7：日程管理（第 7-8 天）

| 任务           | 涉及文件                      | 详情                                                                   |
| -------------- | ----------------------------- | ---------------------------------------------------------------------- |
| 日历视图       | `views/schedule/Calendar.vue` | 月历网格 + 每日日程色块 + 点击日期弹出日程列表 + 新建日程弹窗          |
| 日程列表       | `views/schedule/List.vue`     | 按日期分组展示 + 筛选 + 分页                                           |
| 日程 CRUD 弹窗 | Calendar/List 中复用组件      | 标题 + 起止日期/时间 + 全天开关 + 优先级 + 提醒 + 地点 + 参与人 + 颜色 |
| 日历工具函数   | `utils/date.ts` 已有          | `getMonthDates()`/`getWeekDates()` 直接使用                            |

### A-8：系统管理（第 9-10 天）

| 任务     | 涉及文件                | 详情                                                                                                                        |
| -------- | ----------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| 用户管理 | `views/system/User.vue` | `n-data-table` + 关键词搜索 + 状态筛选 + 新增/编辑弹窗（用户名/昵称/邮箱/手机/部门选择/角色选择/状态）+ 删除确认 + 批量导入 |
| 部门管理 | `views/system/Dept.vue` | 树形表格 `n-tree` + 新增/编辑弹窗（名称/上级部门/负责人/电话/邮箱/排序/状态）+ 删除                                         |
| 角色管理 | `views/system/Role.vue` | 列表 + 新增/编辑弹窗（名称/代码/描述/权限勾选 `n-checkbox-group`）+ 删除                                                    |
| 流程配置 | `views/system/Flow.vue` | 审批类型列表 + 流程节点编辑器（可视化拖拽或表单式配置审批节点和审批人）                                                     |

### 阶段 A 交付物

- [ ] 所有 OA 页面开发完成
- [ ] 对接 Mock 数据可完整运行
- [ ] 登录→工作台→审批→公告→日程→系统管理全流程可用

---

## 阶段 B：Go 后端搭建（2 周）

> 新建独立后端仓库 `oa-saas-backend/`，按 `docs/saas-design.md` §2.3 的目录结构

### B-1：项目初始化（第 1 天） ✅ 已完成

| 任务     | 详情                                                                                         |
| -------- | -------------------------------------------------------------------------------------------- |
| 创建项目 | `oa-saas-backend/`，`go mod init oa-saas`                                                    |
| 安装依赖 | `gin`, `gorm`, `go-jwt/jwt/v5`, `go-redis/redis/v9`(可选), `go-sql-driver/mysql`, `godotenv` |
| 配置管理 | `configs/config.yaml`：数据库、Redis（可选）、JWT 密钥、端口等                               |
| 入口文件 | `cmd/server/main.go`：加载配置、初始化 DB、Redis（可选）、注册路由、启动服务                 |
| 目录结构 | 按设计文档创建 `internal/{config,middleware,model,repository,service,handler,dto,pkg}`       |

### B-2：数据库建表（第 2 天） ✅ 已完成

| 任务       | 详情                                                                                           |
| ---------- | ---------------------------------------------------------------------------------------------- |
| 迁移工具   | 使用 GORM AutoMigrate 或 `golang-migrate`                                                      |
| 用户表     | `users`（id, username, password, nickname, email, phone, avatar, dept_id, role_id, status...） |
| 部门表     | `departments`（id, parent_id, name, sort, leader_id, leader_name, phone, email, status）       |
| 角色表     | `roles`（id, name, code, description, permissions JSON, status）                               |
| 审批相关表 | `approval_flows`, `approvals`, `approval_nodes`（按设计文档 §4.2.3）                           |
| 公告相关表 | `notices`, `notice_reads`（按设计文档 §4.2.4）                                                 |
| 日程相关表 | `schedules`, `schedule_participants`（按设计文档 §4.2.5）                                      |
| 种子数据   | 默认管理员 admin/123456、默认部门、默认角色                                                    |

### B-3：认证模块（第 3 天） ✅ 已完成

| 任务         | 详情                                                            |
| ------------ | --------------------------------------------------------------- |
| JWT 工具     | `pkg/jwt/`：GenerateToken, ParseToken, Claims 结构              |
| 密码加密     | `pkg/utils/password.go`：bcrypt Hash + Verify                   |
| 登录接口     | `POST /api/v1/auth/login`：校验密码→签发 JWT→返回 token+user    |
| 获取用户信息 | `GET /api/v1/auth/info`：从 token 解析用户 ID→查询返回          |
| 修改密码     | `PUT /api/v1/auth/password`：校验旧密码→更新新密码              |
| Auth 中间件  | `middleware/auth.go`：解析 Bearer Token→注入 user_id 到 context |
| CORS 中间件  | 允许前端跨域                                                    |

### B-4：用户/部门/角色 CRUD（第 4-5 天） ✅ 已完成

| 任务      | API                        | 详情                                   |
| --------- | -------------------------- | -------------------------------------- |
| 用户列表  | `GET /api/v1/users`        | 分页 + 关键词搜索 + 部门/角色/状态筛选 |
| 创建用户  | `POST /api/v1/users`       | 唯一性校验 + 密码加密                  |
| 更新用户  | `PUT /api/v1/users/:id`    |                                        |
| 删除用户  | `DELETE /api/v1/users/:id` | 软删除                                 |
| 部门树    | `GET /api/v1/departments`  | 递归返回树形结构                       |
| 部门 CRUD | POST/PUT/DELETE            |                                        |
| 角色列表  | `GET /api/v1/roles`        |                                        |
| 角色 CRUD | POST/PUT/DELETE            |                                        |

### B-5：审批模块（第 6-7 天） ✅ 已完成

| 任务     | API                                   | 详情                                                               |
| -------- | ------------------------------------- | ------------------------------------------------------------------ |
| 发起申请 | `POST /api/v1/approvals`              | 根据 type 校验表单→查询审批流程配置→创建审批+节点                  |
| 我的申请 | `GET /api/v1/approvals/my`            | 分页 + 类型/状态筛选                                               |
| 待我审批 | `GET /api/v1/approvals/pending`       | 查询当前用户为审批人且状态 pending 的                              |
| 已办审批 | `GET /api/v1/approvals/done`          |                                                                    |
| 审批详情 | `GET /api/v1/approvals/:id`           | 申请信息 + 流程节点列表                                            |
| 审批操作 | `POST /api/v1/approvals/:id/action`   | approve/reject/transfer→更新节点状态→判断是否全部通过→更新审批状态 |
| 撤回     | `POST /api/v1/approvals/:id/withdraw` | 仅 pending 状态可撤回                                              |
| 审批统计 | `GET /api/v1/approvals/stats`         | 各状态计数                                                         |

### B-6：公告模块（第 8 天）

| 任务     | API                                | 详情                              |
| -------- | ---------------------------------- | --------------------------------- |
| 公告列表 | `GET /api/v1/notices`              | 分页 + 类型/关键词筛选 + 置顶排序 |
| 公告详情 | `GET /api/v1/notices/:id`          |                                   |
| 发布公告 | `POST /api/v1/notices`             |                                   |
| 标记已读 | `POST /api/v1/notices/:id/read`    | 写入 `notice_reads`               |
| 未读数量 | `GET /api/v1/notices/unread-count` |                                   |

### B-7：日程模块（第 9 天）

| 任务     | API                            | 详情           |
| -------- | ------------------------------ | -------------- |
| 日程列表 | `GET /api/v1/schedules`        | 按日期范围查询 |
| 日程详情 | `GET /api/v1/schedules/:id`    |                |
| 创建日程 | `POST /api/v1/schedules`       | 含参与人       |
| 更新日程 | `PUT /api/v1/schedules/:id`    |                |
| 删除日程 | `DELETE /api/v1/schedules/:id` |                |
| 本周日程 | `GET /api/v1/schedules/week`   |                |

### B-8：前后端联调（第 10 天） ✅ 已完成

| 任务         | 详情                                                              |
| ------------ | ----------------------------------------------------------------- |
| 前端去 Mock  | 移除 `vite-plugin-mock`，`vite.config.ts` 中配置 proxy 到 Go 后端 |
| 修改 baseURL | `.env` 中 `VITE_API_BASE_URL=/api/v1`                             |
| 接口路径对齐 | 确保 stores 中调用路径与后端路由一致                              |
| 错误处理     | 统一 401 跳转登录、403 权限提示、500 错误提示                     |
| 联调测试     | 全流程：登录→工作台→审批→公告→日程→系统管理→个人中心              |

### 阶段 B 交付物

- [ ] Go 后端完整 API 服务
- [ ] 前后端联调通过
- [ ] 单租户 OA 系统完整可用

---

## 阶段 C：多租户 SaaS 改造（2 周）

### C-1：数据库改造（第 1 天）

| 任务               | 详情                                                                                                                          |
| ------------------ | ----------------------------------------------------------------------------------------------------------------------------- |
| 新建租户表         | `plans`, `tenants`, `invoices`, `tenant_logs`（按设计文档 §4.2.1）                                                            |
| 业务表加 tenant_id | users, departments, roles, approvals, approval_nodes, approval_flows, notices, notice_reads, schedules, schedule_participants |
| 创建索引           | 所有表 `tenant_id` 加索引                                                                                                     |
| 种子数据           | 4 个套餐（免费版/标准版/专业版/企业版）、1 个默认租户 + 迁移现有数据                                                          |

### C-2：租户中间件 + GORM Scope（第 2-3 天）

| 任务             | 详情                                                                                                         |
| ---------------- | ------------------------------------------------------------------------------------------------------------ |
| 租户识别中间件   | `middleware/tenant.go`：从 `X-Tenant-Slug` Header 或子域名提取 slug→查询租户→注入 context（按设计文档 §3.2） |
| GORM TenantScope | `model/scopes.go`：自动 WHERE tenant_id=? + BeforeCreate Hook 自动填充（按 §3.3.3）                          |
| 权限中间件       | `middleware/permission.go`：角色权限校验                                                                     |
| 限流中间件       | `middleware/ratelimit.go`：可插拔限流：Redis 可用时基于 Redis，否则基于内存                                  |
| 日志中间件       | `middleware/logger.go`：请求日志 + 租户 ID                                                                   |

### C-3：租户注册 + 管理接口（第 4-5 天）

| 任务     | API                                | 详情                                                |
| -------- | ---------------------------------- | --------------------------------------------------- |
| 租户注册 | `POST /api/v1/tenant/register`     | 创建租户→创建管理员→关联套餐→14 天试用（按 §5.2.1） |
| 租户信息 | `GET /api/v1/tenant/info`          |                                                     |
| 更新租户 | `PUT /api/v1/tenant/info`          |                                                     |
| 套餐列表 | `GET /api/v1/plans`                |                                                     |
| 升级续费 | `POST /api/v1/tenant/plan/upgrade` |                                                     |

### C-4：前端 SaaS 页面（第 6-8 天）

| 任务         | 涉及文件                             | 详情                                           |
| ------------ | ------------------------------------ | ---------------------------------------------- |
| 新增 stores  | `stores/tenant.ts`, `stores/plan.ts` | 租户状态 + 套餐权限检查（按设计文档 §6.4-6.5） |
| 新增 utils   | `utils/tenant.ts`, `utils/plan.ts`   | 子域名解析、功能权限检查函数                   |
| 企业注册页   | `views/auth/Register.vue`            | 企业信息 + 选择套餐 + 创建管理员（按 §6.2）    |
| 选择套餐页   | `views/auth/ChoosePlan.vue`          | 套餐对比卡片 + 费用计算器                      |
| 租户信息页   | `views/tenant/Info.vue`              | 企业信息查看/编辑                              |
| 套餐管理页   | `views/tenant/Plan.vue`              | 当前套餐 + 升级续费                            |
| 账单管理页   | `views/tenant/Invoices.vue`          | 账单列表 + 详情                                |
| 超管后台布局 | `components/Layout/AdminLayout.vue`  | 独立布局                                       |
| 超管仪表盘   | `views/admin/Dashboard.vue`          | 租户数、收入、活跃度统计                       |
| 超管租户管理 | `views/admin/Tenants.vue`            | 租户列表 + 状态管理                            |
| 超管套餐管理 | `views/admin/Plans.vue`              | 套餐配置                                       |

### C-5：路由改造 + 权限控制（第 9-10 天）

| 任务         | 详情                                                                |
| ------------ | ------------------------------------------------------------------- |
| 路由拆分     | 新增 `router/tenant.routes.ts`, `router/admin.routes.ts`（按 §6.3） |
| 登录改造     | 响应中增加 tenant 信息；应用初始化时获取租户+套餐                   |
| 套餐限制 UI  | 菜单/按钮根据 `hasFeature()` 显隐；用户数超限提示升级               |
| 导航守卫增强 | 超管路由需 `requiresSuperAdmin`；租户过期跳转续费页                 |

### 阶段 C 交付物

- [ ] 多租户数据隔离
- [ ] 租户自助注册
- [ ] 套餐权限控制
- [ ] 超管后台

---

## 阶段 D：计费系统（1.5 周）

### D-1：计费服务（第 1-3 天）

| 任务         | 详情                                                                    |
| ------------ | ----------------------------------------------------------------------- |
| 月度账单生成 | `billing_service.go`：定时任务，每月 1 日统计活跃用户数×单价（按 §7.5） |
| 过期检查     | 定时任务，检查 `plan_expire_at`，过期则 `suspended`                     |
| 账单管理 API | 列表、详情、状态查询                                                    |
| 支付宝对接   | `pkg/payment/alipay.go`：创建支付→回调处理→更新账单状态（按 §7.6）      |
| 支付流程     | 选择账单→跳转支付宝→回调→标记已付→更新租户状态                          |

### D-2：前端支付流程（第 4-5 天）

| 任务           | 详情                                       |
| -------------- | ------------------------------------------ |
| 账单列表页完善 | 状态标签（待支付/已支付/已逾期）+ 支付按钮 |
| 支付跳转       | 点击支付→后端返回支付宝 URL→前端跳转       |
| 支付结果页     | 回调展示支付成功/失败                      |
| 续费提醒       | 套餐即将过期时，头部显示警告条             |

### 阶段 D 交付物

- [ ] 完整计费闭环
- [ ] 支付宝支付
- [ ] 过期自动暂停

---

## 阶段 E：部署上线（1 周）

| 任务            | 详情                                         |
| --------------- | -------------------------------------------- |
| 后端 Dockerfile | 多阶段构建（golang:1.21-alpine → alpine）    |
| docker-compose  | Nginx + Go API + MySQL 8.0 + Redis 7（可选） |
| Nginx 配置      | 通配符子域名 → X-Tenant-Slug Header + SSL    |
| CI/CD           | GitHub Actions：自动 build + 部署            |
| 监控            | Prometheus + Grafana                         |
| 上线检查        | 按 §9.6 清单逐项验证                         |

---

## 文件变更清单

### 前端（vue3-oa）

**新增文件（约 20 个）：**

```
├── views/auth/Register.vue
├── views/auth/ChoosePlan.vue
├── views/auth/ForgotPassword.vue
├── views/tenant/Info.vue
├── views/tenant/Plan.vue
├── views/tenant/Invoices.vue
├── views/admin/Dashboard.vue
├── views/admin/Tenants.vue
├── views/admin/Plans.vue
├── views/notice/Detail.vue
├── components/Layout/AdminLayout.vue
├── stores/tenant.ts
├── stores/plan.ts
├── router/tenant.routes.ts
├── router/admin.routes.ts
├── utils/tenant.ts
├── utils/plan.ts
```

**改造文件（约 15 个）：**

```
├── views/login/index.vue          → 重写
├── views/dashboard/index.vue      → 重写
├── views/approval/*.vue           → 重写（5个）
├── views/notice/List.vue          → 重写
├── views/schedule/*.vue           → 重写（2个）
├── views/system/*.vue             → 重写（4个）
├── router/routes.ts               → 增加 SaaS 路由
├── router/index.ts                → 增加租户守卫
├── stores/user.ts                 → 增加租户信息
├── utils/request.ts               → 增加 tenant header
```

### 后端（新建 oa-saas-backend）

```
约 40+ 文件，按 docs/saas-design.md §2.3 目录结构
├── cmd/server/main.go
├── configs/config.yaml
├── internal/
│   ├── config/
│   ├── middleware/       (5 个中间件)
│   ├── model/            (8 个模型)
│   ├── repository/       (8 个 repo)
│   ├── service/          (7 个 service)
│   ├── handler/          (7 个 handler)
│   ├── dto/              (request/response)
│   └── pkg/              (jwt, cache, oss, email, payment, utils)
├── migrations/
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
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

## 开发进度

### 阶段 A：前端页面开发 — ✅ 已完成

所有前端页面在项目创建时已开发完毕，包括：

- 登录页、工作台、审批中心（5页）、公告（2页）、日程（2页）、系统管理（4页）、个人中心（2页）、404页
- 修复了全部 TypeScript 编译错误，`npm run build` 零错误通过

### 阶段 B：Go 后端搭建 — 🔄 进行中

| 子阶段 | 名称                | 状态    | 说明                                                      |
| ------ | ------------------- | ------- | --------------------------------------------------------- |
| B-1    | 项目初始化          | ✅ 完成 | 目录结构、依赖安装、配置文件、入口文件、Makefile          |
| B-2    | 数据库建表+种子数据 | ✅ 完成 | 8 个模型 + SQL 迁移文件 + 种子数据                        |
| B-3    | 认证模块            | ✅ 完成 | JWT 工具、登录接口、获取用户信息、Auth/CORS/Logger 中间件 |
| B-4    | 用户/部门/角色 CRUD | ✅ 完成 | user/dept/role/flow handler + 全部 CRUD 路由              |
| B-5    | 审批模块            | ✅ 完成 | 审批创建/列表/详情/操作/撤回/统计 + 审批流程 CRUD         |
| B-6    | 公告模块            | ✅ 完成 | 公告列表/详情/创建/未读数/标记已读                        |
| B-7    | 日程模块            | ✅ 完成 | 日程 CRUD + 本周日程 + 参与人                             |
| B-8    | 前后端联调          | ✅ 完成 | 移除 Mock、配置代理、对齐 API 路径、统一响应格式          |

#### 已创建的后端文件清单

```
backend/
├── cmd/server/main.go             ✅ 入口文件（DB初始化 + Redis可插拔 + 启动服务）
├── configs/config.yaml             ✅ 配置文件
├── internal/
│   ├── config/config.go            ✅ 配置加载（YAML）
│   ├── middleware/
│   │   ├── cors.go                 ✅ 跨域中间件
│   │   ├── auth.go                 ✅ JWT 认证中间件
│   │   └── logger.go               ✅ 请求日志中间件
│   ├── model/
│   │   ├── user.go                 ✅ 用户模型
│   │   ├── department.go           ✅ 部门模型
│   │   ├── role.go                 ✅ 角色模型（含 JSON 数组类型）
│   │   ├── approval_flow.go        ✅ 审批流程模型（含 JSON 节点）
│   │   ├── approval.go             ✅ 审批+审批节点模型
│   │   ├── notice.go               ✅ 公告+公告阅读模型
│   │   └── schedule.go             ✅ 日程+参与人模型
│   ├── handler/
│   │   ├── auth_handler.go         ✅ 登录 + 获取用户信息
│   │   ├── user_handler.go         ✅ 用户 CRUD + 状态切换
│   │   ├── dept_handler.go         ✅ 部门树 + CRUD
│   │   ├── role_handler.go         ✅ 角色 CRUD
│   │   ├── approval_handler.go     ✅ 审批全流程
│   │   ├── notice_handler.go       ✅ 公告 CRUD + 未读 + 标记已读
│   │   ├── schedule_handler.go     ✅ 日程 CRUD + 周视图
│   │   └── flow_handler.go         ✅ 审批流程 CRUD
│   ├── router/
│   │   └── router.go               ✅ 路由注册
│   └── pkg/
│       ├── jwt/jwt.go              ✅ JWT 生成与解析
│       ├── cache/
│       │   ├── cache.go            ✅ 缓存接口定义
│       │   ├── memory.go           ✅ 内存缓存实现
│       │   └── redis.go            ✅ Redis 缓存实现（可选）
│       └── utils/
│           ├── password.go         ✅ bcrypt 密码加密
│           └── response.go         ✅ 统一 API 响应
├── migrations/
│   ├── 001_init.up.sql             ✅ 建表 SQL
│   └── 001_init.down.sql           ✅ 回滚 SQL
├── Makefile                         ✅ 构建/运行/清理
├── go.mod                           ✅ 模块定义
└── go.sum                           ✅ 依赖校验
```

#### 技术要点记录

- Go 版本：1.26.1
- 使用 `goproxy.cn` 国内代理解决依赖下载问题
- Redis 可插拔：通过 `config.yaml` 中 `redis.enabled` 控制是否启用
- 缓存接口抽象：`Cache` 接口 + `MemoryCache` 默认实现 + `RedisCache` 可选实现
- 认证流程：JWT HS256 签名，Bearer Token 方式
- 数据库：GORM AutoMigrate 自动建表
