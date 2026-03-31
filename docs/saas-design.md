# 企业 OA 办公系统 SaaS 化设计文档

> 版本：v1.0  
> 更新日期：2026-03-30  
> 作者：技术团队

---

## 目录

1. [项目概述](#1-项目概述)
2. [整体架构设计](#2-整体架构设计)
3. [多租户设计](#3-多租户设计)
4. [数据库设计](#4-数据库设计)
5. [API 接口设计](#5-api-接口设计)
6. [前端改造方案](#6-前端改造方案)
7. [计费系统设计](#7-计费系统设计)
8. [部署方案](#8-部署方案)
9. [开发里程碑](#9-开发里程碑)

---

## 1. 项目概述

### 1.1 产品定位

将现有的单租户 OA 办公系统改造为多租户 SaaS 平台，为中小企业提供一站式办公自动化服务。

### 1.2 目标用户

| 用户群体             | 特征                   | 核心需求                   |
| -------------------- | ---------------------- | -------------------------- |
| 小微企业（10-50人）  | 成本敏感、IT能力弱     | 快速上手、低月费、免运维   |
| 中型企业（50-500人） | 有一定IT预算、流程规范 | 功能完整、可定制、数据安全 |
| 大型企业分支机构     | 需独立部署或混合云     | 私有化部署、深度定制       |

### 1.3 核心价值

- **零部署**：注册即用，无需购买服务器、安装软件
- **按需付费**：按用户数计费，用多少付多少
- **持续迭代**：功能自动更新，无需手动升级
- **数据隔离**：多租户架构，数据安全隔离

### 1.4 功能模块

| 模块     | 说明                         | 套餐要求     |
| -------- | ---------------------------- | ------------ |
| 工作台   | 数据统计、待办提醒、快捷入口 | 所有套餐     |
| 审批中心 | 请假/报销/加班/出差/通用审批 | 标准版及以上 |
| 公告通知 | 公告发布、阅读统计、未读提醒 | 所有套餐     |
| 日程管理 | 日历视图、日程提醒、共享日程 | 标准版及以上 |
| 系统管理 | 用户/部门/角色/流程配置      | 所有套餐     |
| 租户管理 | 企业信息、套餐管理、账单查看 | 所有套餐     |

---

## 2. 整体架构设计

### 2.1 系统架构图

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              客户端层                                    │
├─────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                  │
│  │  企业用户前端  │  │  超管后台    │  │  注册/登录页  │                  │
│  │  (Vue3 SPA)  │  │  (Vue3 SPA)  │  │  (Vue3 SPA)  │                  │
│  └──────────────┘  └──────────────┘  └──────────────┘                  │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                              接入层                                      │
├─────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │                      Nginx (反向代理 + 负载均衡)                    │  │
│  │   - 子域名路由: *.oa-saas.com → 租户识别                          │  │
│  │   - SSL 终止                                                      │  │
│  │   - 静态资源服务                                                   │  │
│  └──────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                              服务层 (Go)                                 │
├─────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐      │
│  │  API Gateway│ │  Auth Svc   │ │ Tenant Svc  │ │ Billing Svc │      │
│  │  (Gin/Fiber)│ │  (JWT+租户) │ │ (租户管理)   │ │  (计费服务)  │      │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘      │
│                                                                         │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐      │
│  │ Approval Svc│ │ Notice Svc  │ │ Schedule Svc│ │  User Svc   │      │
│  │  (审批服务)  │ │ (公告服务)   │ │ (日程服务)   │ │ (用户服务)   │      │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘      │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐  │
│  │                    中间件层                                       │  │
│  │  - 租户识别中间件 (TenantMiddleware)                              │  │
│  │  - 权限校验中间件 (PermissionMiddleware)                          │  │
│  │  - 套餐限制中间件 (PlanLimitMiddleware)                           │  │
│  │  - 请求日志中间件 (LoggingMiddleware)                             │  │
│  │  - 限流中间件 (RateLimitMiddleware)                               │  │
│  └─────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                              数据层                                      │
├─────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐      │
│  │   MySQL     │ │  Redis(可选)  │ │ 阿里云 OSS  │ │ Elasticsearch│      │
│  │  (主数据库)  │ │ (缓存/会话,可选) │ │ (文件存储)   │ │  (日志/搜索) │      ││  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘      │
└─────────────────────────────────────────────────────────────────────────┘
```

### 2.2 技术选型

| 层级       | 技术栈                    | 版本  | 说明                               |
| ---------- | ------------------------- | ----- | ---------------------------------- |
| **前端**   | Vue 3 + TypeScript + Vite | 3.4+  | 现有技术栈                         |
|            | Naive UI                  | 2.38+ | 组件库                             |
|            | Pinia                     | 2.1+  | 状态管理                           |
|            | Vue Router                | 4.3+  | 路由管理                           |
| **后端**   | Go                        | 1.21+ | 主语言                             |
|            | Gin / Fiber               | 最新  | Web 框架（推荐 Gin）               |
|            | GORM                      | 1.25+ | ORM 框架                           |
|            | go-jwt/jwt                | v5    | JWT 认证                           |
|            | go-redis/redis            | v9    | Redis 客户端（可选）               |
|            | go-sql-driver/mysql       | 1.7+  | MySQL 驱动                         |
| **数据库** | MySQL                     | 8.0+  | 主数据库                           |
|            | Redis                     | 7.0+  | 缓存、会话、分布式锁（推荐，可选） |
| **存储**   | 阿里云 OSS / 腾讯云 COS   | -     | 文件存储                           |
| **部署**   | Docker                    | 24+   | 容器化                             |
|            | Docker Compose            | 2.20+ | 本地编排                           |
|            | Nginx                     | 1.24+ | 反向代理                           |
| **监控**   | Prometheus + Grafana      | -     | 指标监控                           |
|            | ELK Stack                 | -     | 日志收集                           |

### 2.3 目录结构（后端）

```
oa-saas-backend/
├── cmd/
│   └── server/
│       └── main.go              # 入口文件
├── configs/
│   ├── config.yaml              # 配置文件
│   └── config.example.yaml      # 配置示例
├── internal/
│   ├── config/                  # 配置加载
│   ├── middleware/              # 中间件
│   │   ├── tenant.go            # 租户识别
│   │   ├── auth.go              # JWT 认证
│   │   ├── permission.go        # 权限校验
│   │   ├── plan_limit.go        # 套餐限制
│   │   ├── ratelimit.go         # 限流
│   │   └── logger.go            # 日志
│   ├── model/                   # 数据模型
│   │   ├── tenant.go
│   │   ├── user.go
│   │   ├── approval.go
│   │   ├── notice.go
│   │   ├── schedule.go
│   │   └── billing.go
│   ├── repository/              # 数据访问层
│   │   ├── tenant_repo.go
│   │   ├── user_repo.go
│   │   └── ...
│   ├── service/                 # 业务逻辑层
│   │   ├── tenant_service.go
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── approval_service.go
│   │   ├── notice_service.go
│   │   ├── schedule_service.go
│   │   └── billing_service.go
│   ├── handler/                 # HTTP 处理器
│   │   ├── tenant_handler.go
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── approval_handler.go
│   │   ├── notice_handler.go
│   │   ├── schedule_handler.go
│   │   └── billing_handler.go
│   ├── dto/                     # 数据传输对象
│   │   ├── request/
│   │   └── response/
│   └── pkg/                     # 公共工具
│       ├── jwt/
│       ├── cache/
│       │   ├── cache.go          # 缓存接口定义（Cache Interface）
│       │   ├── memory.go         # 内存缓存实现（默认，无 Redis 时使用）
│       │   └── redis.go          # Redis 缓存实现（可选）
│       ├── oss/
│       ├── email/
│       ├── sms/
│       └── utils/
├── migrations/                  # 数据库迁移
│   ├── 001_init.up.sql
│   └── 001_init.down.sql
├── scripts/                     # 脚本
│   ├── build.sh
│   └── deploy.sh
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
└── go.sum
```

---

## 3. 多租户设计

### 3.1 租户模型

```
┌─────────────────────────────────────────────────────────────┐
│                        Tenant (租户)                         │
├─────────────────────────────────────────────────────────────┤
│  id              - 租户ID (UUID)                             │
│  name            - 企业名称                                   │
│  slug            - 企业标识 (用于子域名)                       │
│  logo            - 企业 Logo                                  │
│  contact_name    - 联系人                                     │
│  contact_phone   - 联系电话                                   │
│  contact_email   - 联系邮箱                                   │
│  plan_id         - 当前套餐ID                                 │
│  plan_expire_at  - 套餐过期时间                               │
│  max_users       - 最大用户数                                 │
│  status          - 状态 (trial/active/suspended/cancelled)   │
│  created_at      - 创建时间                                   │
│  updated_at      - 更新时间                                   │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 租户识别策略

采用 **子域名识别** 方案：

| 子域名                 | 说明                    |
| ---------------------- | ----------------------- |
| `company1.oa-saas.com` | company1 租户的 OA 系统 |
| `company2.oa-saas.com` | company2 租户的 OA 系统 |
| `admin.oa-saas.com`    | 超级管理员后台          |
| `www.oa-saas.com`      | 官网/注册页             |

**识别流程**：

```
1. 用户访问 company1.oa-saas.com
2. Nginx 提取子域名 "company1"
3. 请求转发到后端，Header 携带 X-Tenant-Slug: company1
4. 租户中间件根据 slug 查询租户信息
5. 将 tenant_id 注入到请求上下文
6. 后续所有数据库操作自动携带 tenant_id 条件
```

### 3.3 数据隔离实现

#### 3.3.1 数据库表设计原则

所有租户相关表必须包含 `tenant_id` 字段：

```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL COMMENT '租户ID',
    username VARCHAR(50) NOT NULL,
    -- ... 其他字段
    INDEX idx_tenant_id (tenant_id)
);
```

#### 3.3.2 Go 中间件实现

```go
// internal/middleware/tenant.go

package middleware

import (
    "strings"
    "github.com/gin-gonic/gin"
    "oa-saas/internal/service"
)

type TenantMiddleware struct {
    tenantService *service.TenantService
}

func NewTenantMiddleware(ts *service.TenantService) *TenantMiddleware {
    return &TenantMiddleware{tenantService: ts}
}

func (m *TenantMiddleware) Handle() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从 Header 获取租户标识
        tenantSlug := c.GetHeader("X-Tenant-Slug")

        // 2. 如果 Header 没有，从子域名提取
        if tenantSlug == "" {
            host := c.Request.Host
            parts := strings.Split(host, ".")
            if len(parts) >= 1 && parts[0] != "www" && parts[0] != "admin" {
                tenantSlug = parts[0]
            }
        }

        // 3. 查询租户信息
        tenant, err := m.tenantService.GetBySlug(c, tenantSlug)
        if err != nil {
            c.JSON(404, gin.H{"code": 404, "message": "租户不存在"})
            c.Abort()
            return
        }

        // 4. 检查租户状态
        if tenant.Status == "suspended" {
            c.JSON(403, gin.H{"code": 403, "message": "租户已暂停服务"})
            c.Abort()
            return
        }

        // 5. 注入租户信息到上下文
        c.Set("tenant_id", tenant.ID)
        c.Set("tenant", tenant)

        c.Next()
    }
}
```

#### 3.3.3 GORM 自动注入 tenant_id

```go
// internal/model/scopes.go

package model

import (
    "gorm.io/gorm"
)

// TenantScope 自动添加 tenant_id 条件
func TenantScope(tenantID uint64) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("tenant_id = ?", tenantID)
    }
}

// CreateHook 创建时自动填充 tenant_id
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if tenantID, ok := tx.Statement.Context.Value("tenant_id").(uint64); ok {
        u.TenantID = tenantID
    }
    return nil
}
```

**使用示例**：

```go
// 查询时自动过滤
func (r *UserRepo) GetByID(ctx context.Context, tenantID, userID uint64) (*model.User, error) {
    var user model.User
    err := r.db.WithContext(ctx).
        Scopes(model.TenantScope(tenantID)).
        First(&user, userID).Error
    return &user, err
}

// 创建时自动注入
func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}
```

### 3.4 租户生命周期管理

```
┌─────────┐    注册     ┌─────────┐    付费     ┌─────────┐
│  无     │ ─────────> │  试用   │ ─────────> │  正式   │
└─────────┘            └─────────┘            └─────────┘
                            │                      │
                            │ 欠费                  │ 欠费
                            ▼                      ▼
                       ┌─────────┐           ┌─────────┐
                       │  暂停   │           │  暂停   │
                       └─────────┘           └─────────┘
                            │                      │
                            │ 续费                 │ 续费
                            ▼                      ▼
                       ┌─────────┐           ┌─────────┐
                       │  试用   │           │  正式   │
                       └─────────┘           └─────────┘
```

| 状态        | 说明              | 功能限制             |
| ----------- | ----------------- | -------------------- |
| `trial`     | 试用期（14天）    | 最多 5 个用户        |
| `active`    | 正式付费          | 按套餐限制           |
| `suspended` | 暂停（欠费/违规） | 无法登录使用         |
| `cancelled` | 已注销            | 数据保留 30 天后删除 |

---

## 4. 数据库设计

### 4.1 ER 关系图

```
┌──────────────┐       ┌──────────────┐       ┌──────────────┐
│    Tenant    │       │     Plan     │       │   Invoice    │
├──────────────┤       ├──────────────┤       ├──────────────┤
│ id           │──┐    │ id           │◄──────│ tenant_id    │
│ name         │  │    │ name         │       │ plan_id      │
│ slug         │  │    │ price        │       │ amount       │
│ plan_id      │◄─┼────│ max_users    │       │ status       │
│ ...          │  │    │ features     │       │ ...          │
└──────────────┘  │    └──────────────┘       └──────────────┘
                  │
                  │    ┌──────────────┐       ┌──────────────┐
                  │    │     User     │       │    Role      │
                  │    ├──────────────┤       ├──────────────┤
                  └───►│ tenant_id    │──────►│ tenant_id    │
                       │ username     │       │ name         │
                       │ role_id      │       │ code         │
                       │ ...          │       │ permissions  │
                       └──────────────┘       └──────────────┘
                             │
                  ┌──────────┼──────────┐
                  ▼          ▼          ▼
           ┌──────────┐ ┌──────────┐ ┌──────────┐
           │ Approval │ │  Notice  │ │ Schedule │
           ├──────────┤ ├──────────┤ ├──────────┤
           │tenant_id │ │tenant_id │ │tenant_id │
           │ ...      │ │ ...      │ │ ...      │
           └──────────┘ └──────────┘ └──────────┘
```

### 4.2 核心表结构

#### 4.2.1 租户相关表

```sql
-- 套餐表
CREATE TABLE plans (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL COMMENT '套餐名称',
    code VARCHAR(30) NOT NULL COMMENT '套餐代码',
    price DECIMAL(10,2) NOT NULL COMMENT '每用户每月价格',
    min_users INT UNSIGNED NOT NULL DEFAULT 1 COMMENT '最少用户数',
    max_users INT UNSIGNED NOT NULL DEFAULT 100 COMMENT '最大用户数',
    features JSON COMMENT '功能权限配置',
    is_active TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='套餐表';

-- 租户表
CREATE TABLE tenants (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '企业名称',
    slug VARCHAR(50) NOT NULL COMMENT '企业标识(子域名)',
    logo VARCHAR(255) DEFAULT '' COMMENT '企业Logo',
    contact_name VARCHAR(50) NOT NULL COMMENT '联系人',
    contact_phone VARCHAR(20) NOT NULL COMMENT '联系电话',
    contact_email VARCHAR(100) NOT NULL COMMENT '联系邮箱',
    plan_id BIGINT UNSIGNED NOT NULL COMMENT '当前套餐ID',
    plan_start_at DATETIME COMMENT '套餐开始时间',
    plan_expire_at DATETIME COMMENT '套餐过期时间',
    current_users INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '当前用户数',
    max_users INT UNSIGNED NOT NULL DEFAULT 5 COMMENT '最大用户数',
    status VARCHAR(20) NOT NULL DEFAULT 'trial' COMMENT '状态: trial/active/suspended/cancelled',
    trial_ends_at DATETIME COMMENT '试用结束时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    UNIQUE KEY uk_slug (slug),
    KEY idx_status (status),
    KEY idx_plan_id (plan_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户表';

-- 租户账单表
CREATE TABLE invoices (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL COMMENT '租户ID',
    plan_id BIGINT UNSIGNED NOT NULL COMMENT '套餐ID',
    invoice_no VARCHAR(50) NOT NULL COMMENT '账单号',
    period_start DATE NOT NULL COMMENT '账单周期开始',
    period_end DATE NOT NULL COMMENT '账单周期结束',
    user_count INT UNSIGNED NOT NULL COMMENT '计费用户数',
    amount DECIMAL(10,2) NOT NULL COMMENT '账单金额',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending/paid/overdue/cancelled',
    paid_at DATETIME COMMENT '支付时间',
    payment_method VARCHAR(20) COMMENT '支付方式',
    payment_transaction_id VARCHAR(100) COMMENT '支付流水号',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_invoice_no (invoice_no),
    KEY idx_tenant_id (tenant_id),
    KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户账单表';

-- 租户操作日志表
CREATE TABLE tenant_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL COMMENT '租户ID',
    operator_id BIGINT UNSIGNED COMMENT '操作人ID',
    action VARCHAR(50) NOT NULL COMMENT '操作类型',
    target_type VARCHAR(50) COMMENT '操作对象类型',
    target_id BIGINT UNSIGNED COMMENT '操作对象ID',
    content JSON COMMENT '操作详情',
    ip VARCHAR(45) COMMENT 'IP地址',
    user_agent VARCHAR(500) COMMENT 'User-Agent',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_tenant_id (tenant_id),
    KEY idx_action (action),
    KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户操作日志表';
```

#### 4.2.2 用户与权限表

```sql
-- 部门表
CREATE TABLE departments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL COMMENT '租户ID',
    parent_id BIGINT UNSIGNED DEFAULT NULL COMMENT '父部门ID',
    name VARCHAR(50) NOT NULL COMMENT '部门名称',
    sort INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序',
    leader_id BIGINT UNSIGNED COMMENT '部门负责人ID',
    leader_name VARCHAR(50) COMMENT '部门负责人姓名',
    phone VARCHAR(20) COMMENT '联系电话',
    email VARCHAR(100) COMMENT '部门邮箱',
    status TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态: 0禁用 1启用',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_tenant_id (tenant_id),
    KEY idx_parent_id (parent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='部门表';

-- 角色表
CREATE TABLE roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL COMMENT '租户ID',
    name VARCHAR(50) NOT NULL COMMENT '角色名称',
    code VARCHAR(30) NOT NULL COMMENT '角色代码',
    description VARCHAR(255) COMMENT '角色描述',
    permissions JSON COMMENT '权限列表',
    is_system TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否系统内置角色',
    status TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态: 0禁用 1启用',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_tenant_id (tenant_id),
    UNIQUE KEY uk_tenant_code (tenant_id, code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 用户表
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL COMMENT '租户ID',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码(加密)',
    nickname VARCHAR(50) NOT NULL COMMENT '昵称',
    email VARCHAR(100) COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '手机号',
    avatar VARCHAR(255) DEFAULT '' COMMENT '头像',
    dept_id BIGINT UNSIGNED COMMENT '部门ID',
    dept_name VARCHAR(50) COMMENT '部门名称',
    role_id BIGINT UNSIGNED COMMENT '角色ID',
    role_name VARCHAR(50) COMMENT '角色名称',
    status TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态: 0禁用 1启用',
    last_login_at DATETIME COMMENT '最后登录时间',
    last_login_ip VARCHAR(45) COMMENT '最后登录IP',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    KEY idx_tenant_id (tenant_id),
    UNIQUE KEY uk_tenant_username (tenant_id, username),
    KEY idx_dept_id (dept_id),
    KEY idx_role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

#### 4.2.3 审批表

```sql
-- 审批流程配置表
CREATE TABLE approval_flows (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL COMMENT '租户ID',
    name VARCHAR(100) NOT NULL COMMENT '流程名称',
    type VARCHAR(30) NOT NULL COMMENT '审批类型: leave/expense/overtime/travel/general',
    nodes JSON NOT NULL COMMENT '流程节点配置',
    is_active TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_tenant_id (tenant_id),
    KEY idx_type (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='审批流程配置表';

-- 审批申请表
CREATE TABLE approvals (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL COMMENT '租户ID',
    flow_id BIGINT UNSIGNED COMMENT '流程ID',
    type VARCHAR(30) NOT NULL COMMENT '审批类型',
    title VARCHAR(200) NOT NULL COMMENT '申请标题',
    applicant_id BIGINT UNSIGNED NOT NULL COMMENT '申请人ID',
    applicant_name VARCHAR(50) NOT NULL COMMENT '申请人姓名',
    applicant_dept VARCHAR(50) COMMENT '申请人部门',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending/approved/rejected/withdrawn/transferred',
    current_step INT UNSIGNED NOT NULL DEFAULT 1 COMMENT '当前步骤',
    total_step INT UNSIGNED NOT NULL DEFAULT 1 COMMENT '总步骤数',
    form_data JSON COMMENT '表单数据',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_tenant_id (tenant_id),
    KEY idx_applicant_id (applicant_id),
    KEY idx_status (status),
    KEY idx_type (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='审批申请表';

-- 审批节点记录表
CREATE TABLE approval_nodes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL COMMENT '租户ID',
    approval_id BIGINT UNSIGNED NOT NULL COMMENT '审批ID',
    node_name VARCHAR(100) NOT NULL COMMENT '节点名称',
    step INT UNSIGNED NOT NULL COMMENT '步骤序号',
    approver_id BIGINT UNSIGNED NOT NULL COMMENT '审批人ID',
    approver_name VARCHAR(50) NOT NULL COMMENT '审批人姓名',
    approver_avatar VARCHAR(255) COMMENT '审批人头像',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending/approved/rejected/transferred',
    comment TEXT COMMENT '审批意见',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_tenant_id (tenant_id),
    KEY idx_approval_id (approval_id),
    KEY idx_approver_id (approver_id),
    KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='审批节点记录表';
```

#### 4.2.4 公告表

```sql
-- 公告表
CREATE TABLE notices (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL COMMENT '租户ID',
    title VARCHAR(200) NOT NULL COMMENT '标题',
    type VARCHAR(30) NOT NULL COMMENT '类型: notice/announcement/policy/urgent',
    content TEXT NOT NULL COMMENT '内容',
    summary VARCHAR(500) COMMENT '摘要',
    cover_image VARCHAR(255) COMMENT '封面图',
    publisher_id BIGINT UNSIGNED NOT NULL COMMENT '发布人ID',
    publisher_name VARCHAR(50) NOT NULL COMMENT '发布人姓名',
    publisher_avatar VARCHAR(255) COMMENT '发布人头像',
    status VARCHAR(20) NOT NULL DEFAULT 'published' COMMENT '状态: draft/published/archived',
    is_top TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否置顶',
    read_count INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '阅读次数',
    published_at DATETIME COMMENT '发布时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_tenant_id (tenant_id),
    KEY idx_type (type),
    KEY idx_status (status),
    KEY idx_published_at (published_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='公告表';

-- 公告阅读记录表
CREATE TABLE notice_reads (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL COMMENT '租户ID',
    notice_id BIGINT UNSIGNED NOT NULL COMMENT '公告ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    read_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '阅读时间',
    KEY idx_tenant_id (tenant_id),
    UNIQUE KEY uk_notice_user (notice_id, user_id),
    KEY idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='公告阅读记录表';
```

#### 4.2.5 日程表

```sql
-- 日程表
CREATE TABLE schedules (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL COMMENT '租户ID',
    title VARCHAR(200) NOT NULL COMMENT '标题',
    description TEXT COMMENT '描述',
    start_date DATE NOT NULL COMMENT '开始日期',
    end_date DATE NOT NULL COMMENT '结束日期',
    start_time TIME COMMENT '开始时间',
    end_time TIME COMMENT '结束时间',
    is_all_day TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否全天',
    priority VARCHAR(20) NOT NULL DEFAULT 'medium' COMMENT '优先级: low/medium/high',
    remind VARCHAR(20) NOT NULL DEFAULT 'none' COMMENT '提醒: none/5min/15min/30min/1hour/1day',
    location VARCHAR(200) COMMENT '地点',
    creator_id BIGINT UNSIGNED NOT NULL COMMENT '创建人ID',
    creator_name VARCHAR(50) NOT NULL COMMENT '创建人姓名',
    color VARCHAR(20) DEFAULT '#2080f0' COMMENT '颜色',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_tenant_id (tenant_id),
    KEY idx_start_date (start_date),
    KEY idx_creator_id (creator_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='日程表';

-- 日程参与者表
CREATE TABLE schedule_participants (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL COMMENT '租户ID',
    schedule_id BIGINT UNSIGNED NOT NULL COMMENT '日程ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    user_name VARCHAR(50) NOT NULL COMMENT '用户姓名',
    user_avatar VARCHAR(255) COMMENT '用户头像',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_tenant_id (tenant_id),
    KEY idx_schedule_id (schedule_id),
    UNIQUE KEY uk_schedule_user (schedule_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='日程参与者表';
```

### 4.3 初始数据

```sql
-- 初始套餐数据
INSERT INTO plans (name, code, price, min_users, max_users, features, is_active) VALUES
('免费版', 'free', 0.00, 1, 5,
 '{"approval": false, "schedule": false, "notice": true, "storage": 100}', 1),
('标准版', 'standard', 29.00, 5, 50,
 '{"approval": true, "schedule": true, "notice": true, "storage": 1024}', 1),
('专业版', 'professional', 59.00, 10, 200,
 '{"approval": true, "schedule": true, "notice": true, "storage": 5120, "api": true}', 1),
('企业版', 'enterprise', 99.00, 50, 1000,
 '{"approval": true, "schedule": true, "notice": true, "storage": 20480, "api": true, "sso": true}', 1);

-- 超级管理员 (不属于任何租户)
INSERT INTO users (tenant_id, username, password, nickname, status) VALUES
(0, 'superadmin', '$2a$10$...', '超级管理员', 1);
```

---

## 5. API 接口设计

### 5.1 接口规范

#### 5.1.1 基础 URL

| 环境     | URL                            |
| -------- | ------------------------------ |
| 生产环境 | `https://api.oa-saas.com`      |
| 测试环境 | `https://api-test.oa-saas.com` |
| 本地开发 | `http://localhost:8080`        |

#### 5.1.2 请求头

```
Content-Type: application/json
Authorization: Bearer <token>
X-Tenant-Slug: <tenant_slug>    # 租户标识(可选，服务端从子域名提取)
```

#### 5.1.3 统一响应格式

**成功响应**：

```json
{
    "code": 200,
    "message": "success",
    "data": { ... }
}
```

**分页响应**：

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "list": [ ... ],
        "total": 100,
        "page": 1,
        "page_size": 10,
        "total_pages": 10
    }
}
```

**错误响应**：

```json
{
  "code": 400,
  "message": "参数错误",
  "data": null
}
```

#### 5.1.4 错误码定义

| 错误码 | 说明                |
| ------ | ------------------- |
| 200    | 成功                |
| 400    | 参数错误            |
| 401    | 未登录或 Token 过期 |
| 403    | 无权限              |
| 404    | 资源不存在          |
| 409    | 资源冲突            |
| 422    | 业务逻辑错误        |
| 429    | 请求过于频繁        |
| 500    | 服务器内部错误      |
| 1001   | 租户不存在          |
| 1002   | 租户已暂停          |
| 1003   | 套餐用户数已满      |
| 1004   | 功能未开通          |
| 2001   | 用户名已存在        |
| 2002   | 密码错误            |

### 5.2 认证接口

#### 5.2.1 租户注册

```
POST /api/v1/tenant/register
```

**请求参数**：

```json
{
  "name": "示例科技有限公司",
  "slug": "example-tech",
  "contact_name": "张三",
  "contact_phone": "13800138000",
  "contact_email": "admin@example.com",
  "plan_id": 2
}
```

**响应**：

```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "tenant_id": 1,
    "name": "示例科技有限公司",
    "slug": "example-tech",
    "plan": {
      "id": 2,
      "name": "标准版"
    },
    "trial_ends_at": "2026-04-13T00:00:00Z",
    "admin_user": {
      "id": 1,
      "username": "admin",
      "temp_password": "Abc123456"
    }
  }
}
```

#### 5.2.2 登录

```
POST /api/v1/auth/login
```

**请求参数**：

```json
{
  "username": "admin",
  "password": "123456"
}
```

**响应**：

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": "2026-03-31T00:00:00Z",
    "user": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员",
      "avatar": "https://...",
      "dept_name": "技术部",
      "role_name": "管理员",
      "permissions": ["*"]
    },
    "tenant": {
      "id": 1,
      "name": "示例科技有限公司",
      "plan": {
        "id": 2,
        "name": "标准版",
        "features": {
          "approval": true,
          "schedule": true
        }
      }
    }
  }
}
```

#### 5.2.3 获取当前用户信息

```
GET /api/v1/auth/info
```

#### 5.2.4 修改密码

```
PUT /api/v1/auth/password
```

**请求参数**：

```json
{
  "old_password": "123456",
  "new_password": "newpassword"
}
```

### 5.3 租户管理接口

#### 5.3.1 获取租户信息

```
GET /api/v1/tenant/info
```

#### 5.3.2 更新租户信息

```
PUT /api/v1/tenant/info
```

**请求参数**：

```json
{
  "name": "新公司名称",
  "logo": "https://...",
  "contact_name": "李四",
  "contact_phone": "13900139000"
}
```

#### 5.3.3 获取可用套餐列表

```
GET /api/v1/plans
```

#### 5.3.4 升级/续费套餐

```
POST /api/v1/tenant/plan/upgrade
```

**请求参数**：

```json
{
  "plan_id": 3,
  "user_count": 20,
  "period": 12
}
```

#### 5.3.5 获取账单列表

```
GET /api/v1/tenant/invoices?page=1&page_size=10
```

### 5.4 用户管理接口

#### 5.4.1 获取用户列表

```
GET /api/v1/users?page=1&page_size=10&keyword=&dept_id=&status=
```

**响应**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "zhangsan",
        "nickname": "张三",
        "email": "zhangsan@example.com",
        "phone": "13800138001",
        "avatar": "https://...",
        "dept_id": 1,
        "dept_name": "技术部",
        "role_id": 2,
        "role_name": "员工",
        "status": 1,
        "created_at": "2026-01-01 00:00:00"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10,
    "total_pages": 5
  }
}
```

#### 5.4.2 创建用户

```
POST /api/v1/users
```

**请求参数**：

```json
{
  "username": "lisi",
  "password": "123456",
  "nickname": "李四",
  "email": "lisi@example.com",
  "phone": "13800138002",
  "dept_id": 1,
  "role_id": 2,
  "status": 1
}
```

#### 5.4.3 更新用户

```
PUT /api/v1/users/:id
```

#### 5.4.4 删除用户

```
DELETE /api/v1/users/:id
```

#### 5.4.5 批量导入用户

```
POST /api/v1/users/import
Content-Type: multipart/form-data

file: users.xlsx
```

### 5.5 部门管理接口

```
GET    /api/v1/departments       # 获取部门树
POST   /api/v1/departments       # 创建部门
PUT    /api/v1/departments/:id   # 更新部门
DELETE /api/v1/departments/:id   # 删除部门
```

### 5.6 角色管理接口

```
GET    /api/v1/roles             # 获取角色列表
POST   /api/v1/roles             # 创建角色
PUT    /api/v1/roles/:id         # 更新角色
DELETE /api/v1/roles/:id         # 删除角色
```

### 5.7 审批管理接口

#### 5.7.1 获取我的申请列表

```
GET /api/v1/approvals/my?page=1&page_size=10&type=&status=
```

#### 5.7.2 获取待审批列表

```
GET /api/v1/approvals/pending?page=1&page_size=10
```

#### 5.7.3 获取已办审批列表

```
GET /api/v1/approvals/done?page=1&page_size=10
```

#### 5.7.4 获取审批详情

```
GET /api/v1/approvals/:id
```

#### 5.7.5 发起申请

```
POST /api/v1/approvals
```

**请求参数（请假）**：

```json
{
  "type": "leave",
  "title": "年假申请",
  "form_data": {
    "leave_type": "annual",
    "start_date": "2026-04-01",
    "end_date": "2026-04-03",
    "days": 3,
    "reason": "回老家探亲"
  }
}
```

#### 5.7.6 审批操作

```
POST /api/v1/approvals/:id/action
```

**请求参数**：

```json
{
  "action": "approve",
  "comment": "同意"
}
```

**action 可选值**：

- `approve` - 通过
- `reject` - 驳回
- `transfer` - 转交

#### 5.7.7 撤回申请

```
POST /api/v1/approvals/:id/withdraw
```

### 5.8 公告管理接口

```
GET    /api/v1/notices                   # 获取公告列表
GET    /api/v1/notices/:id               # 获取公告详情
POST   /api/v1/notices                   # 发布公告
PUT    /api/v1/notices/:id               # 更新公告
DELETE /api/v1/notices/:id               # 删除公告
POST   /api/v1/notices/:id/read          # 标记已读
GET    /api/v1/notices/unread-count      # 获取未读数量
```

### 5.9 日程管理接口

```
GET    /api/v1/schedules                 # 获取日程列表
GET    /api/v1/schedules/:id             # 获取日程详情
POST   /api/v1/schedules                 # 创建日程
PUT    /api/v1/schedules/:id             # 更新日程
DELETE /api/v1/schedules/:id             # 删除日程
GET    /api/v1/schedules/week            # 获取本周日程
```

---

## 6. 前端改造方案

### 6.1 项目结构调整

```
vue3-oa/
├── src/
│   ├── views/
│   │   ├── auth/                  # 认证相关页面 (新增)
│   │   │   ├── Login.vue          # 登录页 (改造)
│   │   │   ├── Register.vue       # 企业注册页 (新增)
│   │   │   ├── ChoosePlan.vue     # 选择套餐页 (新增)
│   │   │   └── ForgotPassword.vue # 忘记密码页 (新增)
│   │   ├── dashboard/             # 工作台
│   │   ├── approval/              # 审批中心
│   │   ├── notice/                # 公告通知
│   │   ├── schedule/              # 日程管理
│   │   ├── system/                # 系统管理
│   │   ├── profile/               # 个人中心
│   │   ├── tenant/                # 租户管理 (新增)
│   │   │   ├── Info.vue           # 企业信息
│   │   │   ├── Plan.vue           # 套餐管理
│   │   │   └── Invoices.vue       # 账单管理
│   │   ├── admin/                 # 超管后台 (新增)
│   │   │   ├── Dashboard.vue      # 超管仪表盘
│   │   │   ├── Tenants.vue        # 租户管理
│   │   │   ├── Plans.vue          # 套餐管理
│   │   │   └── Statistics.vue     # 统计分析
│   │   └── error/
│   ├── stores/
│   │   ├── user.ts
│   │   ├── app.ts
│   │   ├── tenant.ts              # 租户状态 (新增)
│   │   ├── plan.ts                # 套餐状态 (新增)
│   │   ├── approval.ts
│   │   ├── notice.ts
│   │   └── schedule.ts
│   ├── router/
│   │   ├── index.ts
│   │   ├── routes.ts              # 路由配置 (改造)
│   │   ├── tenant.routes.ts       # 租户路由 (新增)
│   │   └── admin.routes.ts        # 超管路由 (新增)
│   ├── utils/
│   │   ├── request.ts             # HTTP 封装 (改造)
│   │   ├── tenant.ts              # 租户工具 (新增)
│   │   └── plan.ts                # 套餐权限检查 (新增)
│   └── ...
```

### 6.2 新增页面清单

| 页面           | 路径               | 说明                                   |
| -------------- | ------------------ | -------------------------------------- |
| 企业注册       | `/register`        | 填写企业信息、选择套餐、创建管理员账号 |
| 选择套餐       | `/choose-plan`     | 套餐对比、选择、计算费用               |
| 企业信息       | `/tenant/info`     | 查看/修改企业基本信息                  |
| 套餐管理       | `/tenant/plan`     | 查看当前套餐、升级续费                 |
| 账单管理       | `/tenant/invoices` | 查看账单列表、账单详情、支付           |
| 超管仪表盘     | `/admin/dashboard` | 租户统计、收入统计、系统监控           |
| 租户管理       | `/admin/tenants`   | 租户列表、状态管理、数据查看           |
| 套餐管理(超管) | `/admin/plans`     | 套餐配置、价格调整                     |

### 6.3 路由改造

```typescript
// src/router/routes.ts

import type { RouteRecordRaw } from 'vue-router'

// 公开路由 (无需登录)
export const publicRoutes: RouteRecordRaw[] = [
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue'),
    meta: { title: '企业注册', requiresAuth: false }
  },
  {
    path: '/choose-plan',
    name: 'ChoosePlan',
    component: () => import('@/views/auth/ChoosePlan.vue'),
    meta: { title: '选择套餐', requiresAuth: false }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { title: '登录', requiresAuth: false }
  }
]

// 租户路由 (需要登录)
export const tenantRoutes: RouteRecordRaw[] = [
  {
    path: '/tenant',
    name: 'TenantLayout',
    component: () => import('@/components/Layout/OALayout.vue'),
    redirect: '/tenant/info',
    meta: { title: '企业管理', icon: 'business-outline', requiresAuth: true },
    children: [
      {
        path: 'info',
        name: 'TenantInfo',
        component: () => import('@/views/tenant/Info.vue'),
        meta: { title: '企业信息', requiresAuth: true }
      },
      {
        path: 'plan',
        name: 'TenantPlan',
        component: () => import('@/views/tenant/Plan.vue'),
        meta: { title: '套餐管理', requiresAuth: true }
      },
      {
        path: 'invoices',
        name: 'TenantInvoices',
        component: () => import('@/views/tenant/Invoices.vue'),
        meta: { title: '账单管理', requiresAuth: true }
      }
    ]
  }
]

// 超管路由 (独立部署或特定域名)
export const adminRoutes: RouteRecordRaw[] = [
  {
    path: '/admin',
    name: 'AdminLayout',
    component: () => import('@/components/Layout/AdminLayout.vue'),
    redirect: '/admin/dashboard',
    meta: { title: '超管后台', requiresAuth: true, requiresSuperAdmin: true },
    children: [
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/Dashboard.vue'),
        meta: { title: '仪表盘', requiresAuth: true, requiresSuperAdmin: true }
      },
      {
        path: 'tenants',
        name: 'AdminTenants',
        component: () => import('@/views/admin/Tenants.vue'),
        meta: { title: '租户管理', requiresAuth: true, requiresSuperAdmin: true }
      },
      {
        path: 'plans',
        name: 'AdminPlans',
        component: () => import('@/views/admin/Plans.vue'),
        meta: { title: '套餐管理', requiresAuth: true, requiresSuperAdmin: true }
      }
    ]
  }
]
```

### 6.4 套餐权限控制

```typescript
// src/utils/plan.ts

import { useTenantStore } from '@/stores/tenant'

export interface PlanFeatures {
  approval: boolean
  schedule: boolean
  notice: boolean
  storage: number
  api?: boolean
  sso?: boolean
}

export function usePlanPermission() {
  const tenantStore = useTenantStore()

  function hasFeature(feature: keyof PlanFeatures): boolean {
    const features = tenantStore.planFeatures
    if (!features) return false
    return features[feature] === true
  }

  function checkUserLimit(): { exceeded: boolean; current: number; max: number } {
    return {
      exceeded: tenantStore.currentUsers >= tenantStore.maxUsers,
      current: tenantStore.currentUsers,
      max: tenantStore.maxUsers
    }
  }

  return {
    hasFeature,
    checkUserLimit
  }
}
```

**在组件中使用**：

```vue
<template>
  <div>
    <!-- 套餐限制提示 -->
    <n-alert v-if="!hasFeature('approval')" type="warning">
      当前套餐不支持审批功能，<a @click="upgradePlan">立即升级</a>
    </n-alert>

    <!-- 用户数限制 -->
    <n-button v-if="userLimit.exceeded" disabled>
      用户数已达上限 ({{ userLimit.current }}/{{ userLimit.max }})
    </n-button>
  </div>
</template>

<script setup lang="ts">
  import { usePlanPermission } from '@/utils/plan'

  const { hasFeature, checkUserLimit } = usePlanPermission()
  const userLimit = checkUserLimit()
</script>
```

### 6.5 租户状态管理

```typescript
// src/stores/tenant.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { request } from '@/utils/request'

export interface Tenant {
  id: number
  name: string
  slug: string
  logo: string
  planId: number
  planName: string
  planFeatures: Record<string, boolean | number>
  currentUsers: number
  maxUsers: number
  status: string
  planExpireAt: string
}

export const useTenantStore = defineStore('tenant', () => {
  const tenant = ref<Tenant | null>(null)

  const tenantName = computed(() => tenant.value?.name || '')
  const tenantLogo = computed(() => tenant.value?.logo || '')
  const planFeatures = computed(() => tenant.value?.planFeatures || {})
  const currentUsers = computed(() => tenant.value?.currentUsers || 0)
  const maxUsers = computed(() => tenant.value?.maxUsers || 0)
  const isExpired = computed(() => {
    if (!tenant.value?.planExpireAt) return false
    return new Date(tenant.value.planExpireAt) < new Date()
  })

  async function fetchTenantInfo() {
    const data = await request.get('/tenant/info')
    tenant.value = data
    return data
  }

  async function updateTenantInfo(params: Partial<Tenant>) {
    const data = await request.put('/tenant/info', params)
    tenant.value = { ...tenant.value, ...data }
    return data
  }

  return {
    tenant,
    tenantName,
    tenantLogo,
    planFeatures,
    currentUsers,
    maxUsers,
    isExpired,
    fetchTenantInfo,
    updateTenantInfo
  }
})
```

---

## 7. 计费系统设计

### 7.1 套餐定义

| 套餐       | 价格      | 用户范围  | 存储  | 功能               |
| ---------- | --------- | --------- | ----- | ------------------ |
| **免费版** | ¥0        | 1-5人     | 100MB | 公告、系统管理     |
| **标准版** | ¥29/人/月 | 5-50人    | 1GB   | 全功能             |
| **专业版** | ¥59/人/月 | 10-200人  | 5GB   | 全功能 + API       |
| **企业版** | ¥99/人/月 | 50-1000人 | 20GB  | 全功能 + API + SSO |

### 7.2 计费规则

```
月费用 = 用户数 × 单价

示例：
- 标准版 20 人 = 20 × ¥29 = ¥580/月
- 专业版 50 人 = 50 × ¥59 = ¥2950/月
```

**计费周期**：

- 按月计费，可预付多月（享折扣）
- 3个月：95折
- 6个月：9折
- 12个月：85折

### 7.3 用户数统计规则

```
计费用户数 = 启用状态的用户数量

不计入：
- 已禁用的用户
- 已删除的用户
- 待激活的用户
```

### 7.4 账单生成流程

```
┌─────────────┐
│  定时任务    │ 每月1日00:00执行
│ (Cron Job)  │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 遍历租户    │ 查询所有 active 状态的租户
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 统计用户数  │ COUNT(*) WHERE tenant_id = ? AND status = 1
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 计算金额    │ 用户数 × 套餐单价 × 周期折扣
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 生成账单    │ INSERT INTO invoices
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 发送通知    │ 邮件/短信提醒付款
└─────────────┘
```

### 7.5 Go 计费服务实现

```go
// internal/service/billing_service.go

package service

import (
    "context"
    "time"
    "oa-saas/internal/model"
    "oa-saas/internal/repository"
)

type BillingService struct {
    tenantRepo  *repository.TenantRepo
    userRepo    *repository.UserRepo
    invoiceRepo *repository.InvoiceRepo
    planRepo    *repository.PlanRepo
}

func NewBillingService(
    tr *repository.TenantRepo,
    ur *repository.UserRepo,
    ir *repository.InvoiceRepo,
    pr *repository.PlanRepo,
) *BillingService {
    return &BillingService{
        tenantRepo:  tr,
        userRepo:    ur,
        invoiceRepo: ir,
        planRepo:    pr,
    }
}

// GenerateMonthlyInvoices 生成月度账单
func (s *BillingService) GenerateMonthlyInvoices(ctx context.Context) error {
    // 获取所有活跃租户
    tenants, err := s.tenantRepo.GetActiveTenants(ctx)
    if err != nil {
        return err
    }

    now := time.Now()
    periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
    periodEnd := periodStart.AddDate(0, 1, -1)

    for _, tenant := range tenants {
        // 统计用户数
        userCount, err := s.userRepo.CountActiveByTenantID(ctx, tenant.ID)
        if err != nil {
            continue
        }

        // 获取套餐信息
        plan, err := s.planRepo.GetByID(ctx, tenant.PlanID)
        if err != nil {
            continue
        }

        // 计算金额
        amount := float64(userCount) * plan.Price

        // 创建账单
        invoice := &model.Invoice{
            TenantID:    tenant.ID,
            PlanID:      plan.ID,
            InvoiceNo:   generateInvoiceNo(),
            PeriodStart: periodStart,
            PeriodEnd:   periodEnd,
            UserCount:   userCount,
            Amount:      amount,
            Status:      "pending",
        }

        s.invoiceRepo.Create(ctx, invoice)
    }

    return nil
}

// CheckExpiredTenants 检查过期租户
func (s *BillingService) CheckExpiredTenants(ctx context.Context) error {
    tenants, err := s.tenantRepo.GetExpiredTenants(ctx)
    if err != nil {
        return err
    }

    for _, tenant := range tenants {
        // 暂停租户
        tenant.Status = "suspended"
        s.tenantRepo.Update(ctx, tenant)

        // 发送通知
        // ...
    }

    return nil
}
```

### 7.6 支付对接

支持的支付方式：

- 支付宝（企业支付）
- 微信支付
- 银行转账（线下）

**支付宝对接示例**：

```go
// internal/pkg/payment/alipay.go

package payment

import (
    "github.com/smartwalle/alipay/v3"
)

type AlipayClient struct {
    client *alipay.Client
}

func NewAlipayClient(appID, privateKey string, isProduction bool) (*AlipayClient, error) {
    var client, err = alipay.New(appID, privateKey, isProduction)
    if err != nil {
        return nil, err
    }

    return &AlipayClient{client: client}, nil
}

func (c *AlipayClient) CreatePayment(orderNo, subject, amount string) (string, error) {
    p := alipay.TradePagePay{}
    p.OutTradeNo = orderNo
    p.TotalAmount = amount
    p.Subject = subject
    p.ProductCode = "FAST_INSTANT_TRADE_PAY"

    url, err := c.client.TradePagePay(p)
    if err != nil {
        return "", err
    }

    return url.String(), nil
}
```

---

## 8. 部署方案

### 8.1 部署架构

```
┌─────────────────────────────────────────────────────────────────┐
│                         云服务商 (阿里云/腾讯云)                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐        │
│  │   CDN       │    │  SLB/CLB    │    │    WAF      │        │
│  │  (静态资源)  │───►│ (负载均衡)   │───►│  (防火墙)   │        │
│  └─────────────┘    └─────────────┘    └─────────────┘        │
│                                                 │               │
│                                                 ▼               │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │                      VPC (虚拟私有云)                     │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │  │
│  │  │   ECS #1    │  │   ECS #2    │  │   ECS #3    │     │  │
│  │  │  (Go API)   │  │  (Go API)   │  │  (Go API)   │     │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘     │  │
│  │         │                │                │              │  │
│  │         └────────────────┼────────────────┘              │  │
│  │                          ▼                               │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │  │
│  │  │   RDS       │  │Redis(可选)  │  │    OSS      │     │  │
│  │  │  (MySQL)    │  │ (缓存,可选) │  │  (文件)     │     │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘     │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 8.2 Docker 配置

#### 8.2.1 后端 Dockerfile

```dockerfile
# Dockerfile

FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

CMD ["./main"]
```

#### 8.2.2 docker-compose.yml

```yaml
version: '3.8'

services:
  nginx:
    image: nginx:1.24-alpine
    ports:
      - '80:80'
      - '443:443'
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
      - frontend-dist:/usr/share/nginx/html
    depends_on:
      - api
    restart: always

  api:
    build: .
    ports:
      - '8080:8080'
```

    environment:
      - GIN_MODE=release
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=oa_saas
      - REDIS_ENABLED=${REDIS_ENABLED:-false}
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      - mysql
      # Redis 可插拔：不配置则自动降级到内存缓存，使用 sync.Map 包 实现分布式锁
    restart: always

mysql:
image: mysql:8.0
ports: - '3306:3306'
environment: - MYSQL_ROOT_PASSWORD=${DB_PASSWORD} - MYSQL_DATABASE=oa_saas - TZ=Asia/Shanghai
volumes: - mysql-data:/var/lib/mysql - ./migrations:/docker-entrypoint-initdb.d
command: --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
restart: always

# 可选组件 - 不启动时系统使用内存缓存

redis:
image: redis:7-alpine
ports: - '6379:6379'
volumes: - redis-data:/data
restart: always

volumes:
mysql-data:
redis-data:
frontend-dist:

````

### 8.3 Nginx 配置

```nginx
# nginx.conf

worker_processes auto;

events {
    worker_connections 1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    sendfile        on;
    keepalive_timeout  65;

    # Gzip
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;

    # 上游服务
    upstream api_server {
        server api:8080;
    }

    # 主域名 - 官网/注册
    server {
        listen 80;
        server_name www.oa-saas.com oa-saas.com;
        return 301 https://www.oa-saas.com$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name www.oa-saas.com oa-saas.com;

        ssl_certificate /etc/nginx/ssl/oa-saas.com.pem;
        ssl_certificate_key /etc/nginx/ssl/oa-saas.com.key;

        root /usr/share/nginx/html;
        index index.html;

        location / {
            try_files $uri $uri/ /index.html;
        }

        location /api {
            proxy_pass http://api_server;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Tenant-Slug "";
        }
    }

    # 租户子域名
    server {
        listen 80;
        server_name ~^(?<tenant>.+)\.oa-saas\.com$;
        return 301 https://$host$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name ~^(?<tenant>.+)\.oa-saas\.com$;

        ssl_certificate /etc/nginx/ssl/oa-saas.com.pem;
        ssl_certificate_key /etc/nginx/ssl/oa-saas.com.key;

        root /usr/share/nginx/html;
        index index.html;

        location / {
            try_files $uri $uri/ /index.html;
        }

        location /api {
            proxy_pass http://api_server;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Tenant-Slug $tenant;
        }
    }

    # 超管后台
    server {
        listen 443 ssl http2;
        server_name admin.oa-saas.com;

        ssl_certificate /etc/nginx/ssl/oa-saas.com.pem;
        ssl_certificate_key /etc/nginx/ssl/oa-saas.com.key;

        root /usr/share/nginx/html/admin;
        index index.html;

        location / {
            try_files $uri $uri/ /index.html;
        }

        location /api {
            proxy_pass http://api_server;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
    }
}
````

### 8.4 CI/CD 配置

```yaml
# .github/workflows/deploy.yml

name: Deploy

on:
  push:
    branches: [main]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Build Backend
        run: |
          cd backend
          go mod download
          go build -o main ./cmd/server

      - name: Set up Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Build Frontend
        run: |
          cd frontend
          npm ci
          npm run build

      - name: Deploy to Server
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          script: |
            cd /opt/oa-saas
            docker-compose pull
            docker-compose up -d
```

---

## 9. 开发里程碑

### 9.1 阶段划分

| 阶段      | 名称       | 周期  | 目标                           |
| --------- | ---------- | ----- | ------------------------------ |
| **阶段1** | 单租户完善 | 2周   | 完成所有页面开发，对接真实后端 |
| **阶段2** | 多租户改造 | 2周   | 后端多租户支持，租户注册/管理  |
| **阶段3** | 计费系统   | 1.5周 | 套餐管理、账单生成、支付对接   |
| **阶段4** | 部署上线   | 1周   | 部署、域名、监控、上线         |

### 9.2 阶段1：单租户完善（2周）

#### 第1周：页面开发

| 任务         | 预计工时 | 验收标准                             |
| ------------ | -------- | ------------------------------------ |
| 工作台页面   | 1天      | 统计卡片、待办列表、快捷入口正常显示 |
| 发起申请页面 | 1天      | 5种申请类型表单正常提交              |
| 我的申请页面 | 0.5天    | 列表展示、筛选、分页正常             |
| 待我审批页面 | 0.5天    | 列表展示、审批操作正常               |
| 已办审批页面 | 0.5天    | 列表展示正常                         |
| 审批详情页面 | 1天      | 详情展示、审批流程、操作按钮正常     |

#### 第2周：后端对接

| 任务               | 预计工时 | 验收标准                         |
| ------------------ | -------- | -------------------------------- |
| 搭建 Go 后端框架   | 1天      | Gin 项目结构、配置加载、日志正常 |
| 用户认证模块       | 1天      | 登录、JWT、用户信息接口正常      |
| 用户/部门/角色管理 | 1天      | CRUD 接口正常                    |
| 审批模块           | 1天      | 所有审批接口正常                 |
| 公告模块           | 0.5天    | 所有公告接口正常                 |
| 日程模块           | 0.5天    | 所有日程接口正常                 |
| 前端对接调试       | 1天      | 所有页面数据正常展示和操作       |

### 9.3 阶段2：多租户改造（2周）

#### 第1周：后端改造

| 任务                 | 预计工时 | 验收标准                   |
| -------------------- | -------- | -------------------------- |
| 数据库表加 tenant_id | 0.5天    | 所有表添加字段、索引       |
| 租户中间件           | 1天      | 子域名识别、上下文注入正常 |
| GORM Scope 封装      | 1天      | 自动过滤 tenant_id         |
| 租户注册接口         | 1天      | 注册流程正常、试用租户创建 |
| 租户管理接口         | 0.5天    | 租户信息 CRUD 正常         |
| 套餐接口             | 0.5天    | 套餐列表、升级续费正常     |
| 超管后台接口         | 1天      | 租户管理、统计接口正常     |

#### 第2周：前端改造

| 任务         | 预计工时 | 验收标准                 |
| ------------ | -------- | ------------------------ |
| 企业注册页面 | 1天      | 注册流程正常             |
| 租户管理页面 | 1天      | 企业信息、套餐管理正常   |
| 套餐权限控制 | 1天      | 功能限制、用户数限制正常 |
| 超管后台页面 | 2天      | 仪表盘、租户管理正常     |

### 9.4 阶段3：计费系统（1.5周）

| 任务         | 预计工时 | 验收标准               |
| ------------ | -------- | ---------------------- |
| 账单生成服务 | 1天      | 定时任务、账单生成正常 |
| 账单管理接口 | 0.5天    | 账单列表、详情接口正常 |
| 账单管理页面 | 0.5天    | 账单列表、详情展示正常 |
| 支付宝对接   | 1天      | 支付、回调正常         |
| 微信支付对接 | 1天      | 支付、回调正常         |
| 过期检查服务 | 0.5天    | 过期租户暂停正常       |
| 支付流程测试 | 0.5天    | 完整支付流程正常       |

### 9.5 阶段4：部署上线（1周）

| 任务           | 预计工时 | 验收标准                         |
| -------------- | -------- | -------------------------------- |
| 服务器环境搭建 | 1天      | Docker、MySQL 正常（Redis 可选） |
| 域名配置       | 0.5天    | 域名解析、SSL 证书正常           |
| Nginx 配置     | 0.5天    | 子域名路由、反向代理正常         |
| CI/CD 配置     | 1天      | 自动构建、部署正常               |
| 监控配置       | 1天      | Prometheus、Grafana 正常         |
| 压力测试       | 1天      | 并发 1000 正常响应               |
| 上线检查清单   | 0.5天    | 所有检查项通过                   |
| 正式上线       | 0.5天    | 系统对外可用                     |

### 9.6 上线检查清单

- [ ] 所有接口功能正常
- [ ] 前端页面无报错
- [ ] 租户注册流程正常
- [ ] 支付流程正常
- [ ] 套餐限制生效
- [ ] 数据隔离验证通过
- [ ] SSL 证书有效
- [ ] 监控告警正常
- [ ] 日志收集正常
- [ ] 数据库备份配置
- [ ] 文档更新完成
- [ ] 客服渠道就绪

---

## 附录

### A. 技术文档链接

| 文档     | 链接                     |
| -------- | ------------------------ |
| Vue 3    | https://vuejs.org/       |
| Naive UI | https://www.naiveui.com/ |
| Gin      | https://gin-gonic.com/   |
| GORM     | https://gorm.io/         |
| Docker   | https://docs.docker.com/ |

### B. 代码仓库

| 仓库 | 地址                              |
| ---- | --------------------------------- |
| 前端 | https://github.com/xxy757/vue3-oa |
| 后端 | (待创建)                          |

### C. 变更记录

| 日期       | 版本 | 变更内容 |
| ---------- | ---- | -------- |
| 2026-03-30 | v1.0 | 初始版本 |
