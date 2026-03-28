# Vue3 OA 办公自动化系统

基于 Vue 3 + TypeScript + Naive UI 构建的企业级办公自动化系统。

## 技术栈

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
| vite-plugin-mock | 3.0.1 | Mock 数据插件 |
| ESLint | 8.57.0 | 代码质量检查 |
| Prettier | 3.2.5 | 代码格式化工具 |

## 功能模块

### 1. 工作台 (Dashboard)
- 数据统计展示
- 待办事项提醒
- 快捷操作入口

### 2. 审批中心
- **发起申请** - 支持多种审批类型
  - 请假申请（年假、病假、事假等）
  - 报销申请（差旅、办公、招待等）
  - 加班申请
  - 出差申请
  - 通用审批
- **我的申请** - 查看个人申请记录及状态
- **待我审批** - 处理待审批事项
- **已办审批** - 查看已处理的审批记录

### 3. 公告通知
- 公告列表（支持分类、置顶）
- 公告详情查看
- 未读提醒

### 4. 日程管理
- 日程日历视图
- 日程列表视图
- 日程增删改查
- 日程提醒

### 5. 系统管理 (需管理员权限)
- **用户管理** - 用户增删改查、状态管理
- **部门管理** - 组织架构管理
- **角色管理** - 角色权限配置
- **流程配置** - 审批流程设置

### 6. 个人中心
- 个人信息查看与修改
- 密码修改

---

## 开发进度

### 已完成

| 模块 | 状态 | 说明 |
|------|------|------|
| 项目架构 | ✅ 已完成 | Vite + Vue3 + TS 项目搭建 |
| 路由配置 | ✅ 已完成 | 完整的路由定义与导航守卫 |
| 类型定义 | ✅ 已完成 | User, Approval, Notice, Schedule, Common |
| 状态管理 | ✅ 已完成 | user, app, approval, notice stores |
| HTTP 封装 | ✅ 已完成 | Axios 请求拦截、响应处理、错误处理 |
| 工具函数 | ✅ 已完成 | storage, date, 通用工具 |
| Mock 数据 | ✅ 已完成 | 完整的模拟数据接口 |
| 布局组件 | ✅ 已完成 | OALayout 侧边栏 + 头部 + 内容区 |
| 样式变量 | ✅ 已完成 | SCSS 变量、重置样式 |

### 待开发

| 模块 | 状态 | 文件路径 |
|------|------|----------|
| 登录页面 | ⏳ 待开发 | `src/views/login/index.vue` |
| 工作台 | ⏳ 待开发 | `src/views/dashboard/index.vue` |
| 发起申请 | ⏳ 待开发 | `src/views/approval/Apply.vue` |
| 我的申请 | ⏳ 待开发 | `src/views/approval/MyApply.vue` |
| 待我审批 | ⏳ 待开发 | `src/views/approval/Pending.vue` |
| 已办审批 | ⏳ 待开发 | `src/views/approval/Done.vue` |
| 审批详情 | ⏳ 待开发 | `src/views/approval/Detail.vue` |
| 公告列表 | ⏳ 待开发 | `src/views/notice/List.vue` |
| 公告详情 | ⏳ 待开发 | `src/views/notice/Detail.vue` |
| 日程日历 | ⏳ 待开发 | `src/views/schedule/Calendar.vue` |
| 日程列表 | ⏳ 待开发 | `src/views/schedule/List.vue` |
| 用户管理 | ⏳ 待开发 | `src/views/system/User.vue` |
| 部门管理 | ⏳ 待开发 | `src/views/system/Dept.vue` |
| 角色管理 | ⏳ 待开发 | `src/views/system/Role.vue` |
| 流程配置 | ⏳ 待开发 | `src/views/system/Flow.vue` |
| 个人信息 | ⏳ 待开发 | `src/views/profile/Info.vue` |
| 修改密码 | ⏳ 待开发 | `src/views/profile/Password.vue` |
| 404 页面 | ⏳ 待开发 | `src/views/error/404.vue` |

---

## 项目结构

```
vue3-oa/
├── src/
│   ├── assets/             # 静态资源
│   ├── components/         # 公共组件
│   │   └── Layout/
│   │       └── OALayout.vue    # 主布局组件
│   ├── mock/               # Mock 数据
│   │   ├── user.ts
│   │   ├── approval.ts
│   │   ├── notice.ts
│   │   └── schedule.ts
│   ├── router/             # 路由配置
│   │   ├── index.ts
│   │   └── routes.ts
│   ├── stores/             # Pinia 状态管理
│   │   ├── user.ts
│   │   ├── app.ts
│   │   ├── approval.ts
│   │   └── notice.ts
│   ├── styles/             # 全局样式
│   │   ├── variables.scss
│   │   ├── reset.scss
│   │   └── index.scss
│   ├── types/              # TypeScript 类型定义
│   │   ├── user.d.ts
│   │   ├── approval.d.ts
│   │   ├── notice.d.ts
│   │   ├── schedule.d.ts
│   │   └── common.d.ts
│   ├── utils/              # 工具函数
│   │   ├── request.ts
│   │   ├── storage.ts
│   │   ├── date.ts
│   │   └── index.ts
│   ├── views/              # 页面组件 (待开发)
│   ├── App.vue
│   └── main.ts
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── .env
```

---

## 快速开始

### 安装依赖

```bash
npm install
```

### 启动开发服务器

```bash
npm run dev
```

### 构建生产版本

```bash
npm run build
```

### 代码检查

```bash
npm run lint
```

---

## 测试账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | 123456 | 管理员 |
| user | 123456 | 普通员工 |
| manager | 123456 | 部门经理 |

---

## Mock 接口列表

### 认证相关
- `POST /api/auth/login` - 登录
- `GET /api/auth/info` - 获取用户信息
- `PUT /api/auth/password` - 修改密码

### 用户管理
- `GET /api/user/list` - 获取用户列表
- `GET /api/dept/list` - 获取部门列表
- `GET /api/role/list` - 获取角色列表

### 审批管理
- `GET /api/approval/my` - 我的申请列表
- `GET /api/approval/pending` - 待审批列表
- `GET /api/approval/done` - 已办审批列表
- `GET /api/approval/:id` - 审批详情
- `POST /api/approval/create` - 发起申请
- `POST /api/approval/action` - 审批操作
- `POST /api/approval/withdraw` - 撤回申请
- `GET /api/approval/stats` - 审批统计

### 公告管理
- `GET /api/notice/list` - 公告列表
- `GET /api/notice/:id` - 公告详情
- `POST /api/notice/create` - 发布公告
- `POST /api/notice/:id/read` - 标记已读
- `GET /api/notice/unread-count` - 未读数量

### 日程管理
- `GET /api/schedule/list` - 日程列表
- `GET /api/schedule/:id` - 日程详情
- `POST /api/schedule/create` - 创建日程
- `PUT /api/schedule/:id` - 更新日程
- `DELETE /api/schedule/:id` - 删除日程
- `GET /api/schedule/week` - 本周日程

---

## License

MIT
