# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

OA SaaS 多租户办公自动化系统。前端 Vue 3 + TypeScript（Naive UI），后端 Go（Gin + GORM + MySQL）。界面语言为中文。

## 开发命令

### 前端（项目根目录）
```bash
npm install          # 安装依赖
npm run dev          # 开发服务器，端口 3000，/api/v1 代理到 localhost:8080
npm run build        # vue-tsc 类型检查 + vite 构建
npm run lint         # eslint 并自动修复
npx vue-tsc --noEmit # 仅类型检查，不生成文件
```

### 后端（backend/ 目录）
```bash
cd backend
go run ./cmd/server           # 启动服务（首次运行自动建表 + 种子数据）
make run                      # 同上，通过 Makefile
make build                    # 编译输出到 bin/server
make tidy                     # go mod tidy
```

后端依赖 MySQL 8.0，默认连接 localhost:3306，数据库 `oa_saas`，用户 `root`，密码 `12345678`。可通过 `backend/configs/config.yaml` 或 `config.local.yaml` 覆盖配置。环境变量 `OA_DB_PASSWORD` 覆盖数据库密码，`OA_CONFIG_PATH` 指定自定义配置路径。Redis 为可选依赖，未配置时自动降级为内存缓存。

### 测试账号
默认租户 slug: `demo`，所有密码: `123456`
- admin（管理员）、zhangsan（部门经理）、lisi（部门经理）、wangwu/zhaoliu/sunqi（普通员工）

## 架构

### 多租户数据隔离
- 共享数据库，12 张业务表均含 `tenant_id` 字段及索引
- `backend/internal/model/scopes.go` 中 `TenantScope(tenantID)` 为 GORM Scope，自动添加 `WHERE tenant_id = ?`
- 租户中间件 `backend/internal/middleware/tenant.go` 从 `X-Tenant-Slug` 请求头或子域名解析租户，将 `tenant_id` 注入 Gin 上下文（通过 `c.Get("tenant_id")` 获取）
- JWT Claims 同时携带 UserID 和 TenantID
- 前端 `src/utils/request.ts` 在每次请求中自动注入 `X-Tenant-Slug` 头

### 后端结构（Go）
- **Go 模块路径**: `oa-saas`（见 go.mod）
- **入口**: `backend/cmd/server/main.go` — 加载配置 → 初始化数据库（AutoMigrate + Seed）→ 初始化缓存（Redis 或内存降级）→ 启动 Gin
- **Handler**: `backend/internal/handler/` — 按模块分文件（auth、user、dept、role、approval、notice、schedule、flow、tenant、admin），通过构造函数接收 `*gorm.DB`。Handler 中通过 `c.GetUint("tenant_id")` 获取当前租户 ID，所有查询必须应用 `TenantScope`
- **Model**: `backend/internal/model/` — GORM 模型，共 14 张表，`seed.go` 含初始数据
- **中间件**: `backend/internal/middleware/` — CORS、Logger、Tenant（租户识别）、Auth（JWT 验证）
- **路由**: `backend/internal/router/router.go` — 所有路由前缀 `/api/v1`。四层分组：公开（无认证）、租户解析（仅 tenant 中间件）、认证（tenant + auth 中间件）、平台管理（仅 auth，`/api/v1/admin/*`）
- **配置**: `backend/internal/config/config.go` — 加载 `config.yaml`，合并 `config.local.yaml` 非零字段，再叠加 `OA_DB_PASSWORD` 环境变量
- **缓存**: `backend/internal/pkg/cache/` — 接口定义 + 内存实现 + Redis 实现
- **工具**: `backend/internal/pkg/utils/` — bcrypt 密码哈希（`password.go`）、统一 API 响应封装 `Success()`/`Error()`（`response.go`）
- **JWT**: `backend/internal/pkg/jwt/jwt.go` — `GenerateToken()` / `ParseToken()`，Claims 包含 UserID + TenantID

### 前端结构（Vue 3）
- **状态管理**: `src/stores/` — Pinia store，按领域划分（user、approval、notice、schedule、tenant、app）
- **类型定义**: `src/types/` — TypeScript 接口，与后端模型对应
- **HTTP 封装**: `src/utils/request.ts` — axios 封装，自动注入 Bearer Token 和 X-Tenant-Slug。响应自动解包：`data.code === 200` 时直接返回 `data.data`，否则抛出 Error。401 时自动清除认证信息并跳转登录页
- **路由**: `src/router/routes.ts` — `constantRoutes`（登录、注册、404）+ `asyncRoutes`（业务页面用 OALayout，超管页面用 AdminLayout）。路由通过 `meta.roles` 控制访问权限（`admin` = 企业管理员，`super_admin` = 平台超管）
- **布局**: `OALayout.vue`（主应用框架含侧边栏）、`AdminLayout.vue`（超管框架）
- **API 前缀**: 所有请求走 `/api/v1`，开发环境通过 vite.config.ts 代理到 `http://127.0.0.1:8080`

### 前后端 API 契约
- 所有接口返回统一格式 `{ code: number, message?: string, data?: T }`
- 成功时 `code = 200`，前端 `request.ts` 自动解包直接返回 `data` 字段
- 分页接口统一参数 `{ page, pageSize }`，返回 `{ list: T[], total: number }`
- 新增后端接口时，需在 `src/types/` 中定义对应 TypeScript 类型，在 `src/stores/` 中通过 `request` 对象调用

## 代码规范

- 代码不添加任何注释，除非用户明确要求
- TypeScript 禁止使用 `any` 类型（确实无法确定时需加注释说明）
- Naive UI 组件按需导入，不全局注册
- 图标从 `@vicons/ionicons5` 导入
- 组件样式使用 `<style lang="scss" scoped>`
- 组件 props 使用 `defineProps<T>()` 泛型语法
- SCSS 变量通过 vite.config.ts 全局注入（`additionalData`），**禁止**在组件中再次 `@import` 或 `@use` `variables.scss`
- 禁止在同一文件中重复导出同名变量/函数
- UI 设计遵循 Ant Design Pro 风格：主色 `#1677FF`，大量留白，轻阴影，卡片圆角 8px，禁用渐变，禁用药丸按钮

## Bug 修复流程

修复 Bug 前必须先全面诊断：运行 `vue-tsc --noEmit`、检查所有 `.ts` 文件语法完整性、验证 store 中 API URL 与后端路由匹配、验证返回结构与 TypeScript 类型定义一致。一次性列出所有问题再逐一修复，禁止看到表面报错就逐个修补。
