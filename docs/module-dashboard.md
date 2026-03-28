# 工作台(Dashboard)模块文档

## 1. 模块功能说明

工作台模块是 OA 系统的首页/入口页面，为用户提供以下功能：

- **欢迎区域**：显示当前用户名、问候语和当前日期
- **统计卡片**：展示待审批数量、未读公告数、今日日程数、我的申请数量
- **待办事项**：显示最近的待审批列表（最多5条），支持点击跳转到详情
- **快捷入口**：提供发起申请、公告列表、日程管理、我的申请的快速访问按钮
- **图表展示**：使用 ECharts 展示审批统计趋势和审批类型分布

## 2. 页面布局结构

```
+------------------------------------------+
|              欢迎区域                      |
+------------------------------------------+
|  待审批  |  未读公告  |  今日日程  | 我的申请 |
+------------------------------------------+
|              待办审批列表        | 快捷入口 |
|                                  +--------+
|                                  |本周日程 |
+----------------------------------+--------+
|          审批统计      |    审批类型分布   |
+------------------------------------------+
```

### 响应式设计

- **大屏幕 (l)**: 统计卡片4列，主内容区左侧待办/右侧快捷，图表区2列
- **中等屏幕 (m)**: 统计卡片2列
- **小屏幕**: 统计卡片1列，主内容区堆叠显示

## 3. 使用的接口信息

### 3.1 审批统计接口

**GET** `/api/approval/stats`

返回数据结构：
```typescript
{
  code: 200,
  message: '成功',
  data: {
    myPending: number,    // 我的待处理申请数
    myApproved: number,   // 我的已通过申请数
    myRejected: number,   // 我的已驳回申请数
    todoApproval: number  // 待我审批数量
  }
}
```

### 3.2 待审批列表接口

**GET** `/api/approval/pending`

请求参数：
```typescript
{
  page: number,      // 页码
  pageSize: number   // 每页数量
}
```

返回数据结构：
```typescript
{
  code: 200,
  message: '成功',
  data: {
    list: Approval[],    // 审批列表
    total: number,       // 总数
    page: number,        // 当前页
    pageSize: number,    // 每页数量
    totalPages: number   // 总页数
  }
}
```

### 3.3 未读公告数接口

**GET** `/api/notice/unread-count`

返回数据结构：
```typescript
{
  code: 200,
  message: '成功',
  data: {
    count: number  // 未读公告数
  }
}
```

### 3.4 本周日程接口

**GET** `/api/schedule/week`

返回数据结构：
```typescript
{
  code: 200,
  message: '成功',
  data: Schedule[]  // 本周日程列表
}
```

## 4. 主要组件和数据流

### 4.1 使用的 Store

| Store | 用途 |
|-------|------|
| `useUserStore` | 获取当前用户信息（用户名、头像等） |
| `useApprovalStore` | 获取审批相关数据和统计 |
| `useNoticeStore` | 获取未读公告数量 |
| `useScheduleStore` | 获取本周日程数据 |

### 4.2 数据流

```
组件挂载 (onMounted)
    |
    +--> fetchStats() --> 并行请求
    |       |
    |       +--> approvalStore.getPendingApprovals()
    |       +--> noticeStore.getUnreadCount()
    |       +--> scheduleStore.getWeekSchedules()
    |       +--> approvalStore.getMyApprovals()
    |
    +--> fetchPendingList() --> approvalStore.getPendingApprovals()
    |
    +--> initApprovalChart() --> 初始化审批趋势图
    |
    +--> initTypeChart() --> 初始化类型分布图
```

### 4.3 主要组件

| 组件 | 说明 |
|------|------|
| `n-card` | 卡片容器，用于各个区域的外层包裹 |
| `n-grid/ngi` | 栅格布局，实现响应式设计 |
| `n-statistic` | 统计数值展示 |
| `n-number-animation` | 数字动画效果 |
| `n-list/n-list-item` | 列表组件，用于待办事项 |
| `n-timeline` | 时间轴组件，用于本周日程 |
| `n-tag` | 标签组件，用于显示审批类型 |
| `echarts` | 图表库，用于可视化展示 |

### 4.4 类型定义

使用的类型主要来自：
- `@/types/approval` - 审批相关类型（Approval, ApprovalType, ApprovalStatus 等）
- `@/types/schedule` - 日程相关类型（Schedule 等）
- `@/types/notice` - 公告相关类型

## 5. 功能交互

### 5.1 点击事件

| 元素 | 点击行为 |
|------|----------|
| 待审批统计卡片 | 无跳转（展示用） |
| 未读公告统计卡片 | 跳转到 `/notice/list` |
| 今日日程统计卡片 | 跳转到 `/schedule/calendar` |
| 我的申请统计卡片 | 跳转到 `/approval/my-apply` |
| 待办审批项 | 跳转到 `/approval/detail/:id` |
| 快捷入口按钮 | 跳转到对应页面 |
| 查看更多按钮 | 跳转到对应列表页面 |

### 5.2 图表说明

**审批统计图（折线图）**：
- X轴：周一到周日
- Y轴：审批数量
- 系列：已通过、已驳回、待审批
- 交互：支持 hover 显示详情

**审批类型分布图（环形图）**：
- 数据：请假申请、报销申请、加班申请、出差申请、通用审批
- 交互：hover 高亮，显示百分比

## 6. 样式说明

### 6.1 统计卡片配色

| 卡片 | 图标背景色 | 图标颜色 |
|------|-----------|---------|
| 待审批 | rgba(240, 160, 32, 0.1) | #f0a020 |
| 未读公告 | rgba(24, 160, 88, 0.1) | #18a058 |
| 今日日程 | rgba(32, 128, 240, 0.1) | #2080f0 |
| 我的申请 | rgba(208, 48, 80, 0.1) | #d03050 |

### 6.2 响应式断点

遵循 Naive UI 的响应式断点：
- `l` (large): >= 1280px
- `m` (medium): >= 768px
- `s` (small): < 768px
