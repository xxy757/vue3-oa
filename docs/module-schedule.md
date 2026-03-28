# 日程管理模块

## 1. 模块功能说明

日程管理模块是 OA 办公系统的核心功能之一，用于管理用户的日常日程安排。主要功能包括：

- **日程日历视图**：以日历形式展示日程，支持按日期查看和管理
- **日程列表视图**：以表格形式展示日程列表，支持筛选和分页
- **日程创建**：支持创建包含标题、描述、时间、优先级、提醒等信息的日程
- **日程编辑**：支持修改已有日程的详细信息
- **日程删除**：支持删除不需要的日程
- **日程详情**：支持查看日程的完整信息

## 2. 页面组件说明

### 2.1 日程日历页面 (`Calendar.vue`)

**路由路径**: `/schedule/calendar`

**功能特性**:
- 左侧显示 Naive UI 的 `n-calendar` 日历组件
- 日历上显示每天的日程数量标记（Badge 形式）
- 点击日期后，右侧显示该日期的所有日程列表
- 日程按优先级用不同颜色的左边框标识
- 支持新增日程（弹出模态框）
- 点击日程可查看详情，并支持编辑和删除

**布局结构**:
```
+------------------+------------------+
|                  |   2024年3月28日   |
|     日历组件      |   [新增日程按钮]  |
|                  |------------------|
|  (带日程数量标记)  |     日程列表      |
|                  |   - 日程项1       |
|                  |   - 日程项2       |
+------------------+------------------+
```

### 2.2 日程列表页面 (`List.vue`)

**路由路径**: `/schedule/list`

**功能特性**:
- 使用 Naive UI 的 `n-data-table` 表格组件展示日程列表
- 支持按日期范围筛选
- 支持按优先级筛选
- 支持按标题关键字搜索
- 分页展示日程数据
- 表格操作列支持查看、编辑、删除操作
- 支持新增日程

**表格列定义**:
| 列名 | 字段 | 说明 |
|-----|------|------|
| 标题 | title | 带颜色标识的日程标题 |
| 优先级 | priority | 使用标签显示 |
| 开始时间 | startTime | 日期+时间，全天只显示日期 |
| 结束时间 | endTime | 日期+时间，全天只显示日期 |
| 地点 | location | 可选字段 |
| 提醒 | remind | 提醒时间设置 |
| 操作 | actions | 查看/编辑/删除按钮 |

## 3. 日程数据结构

### 3.1 Schedule 类型定义

```typescript
interface Schedule {
  id: number                    // 日程ID
  title: string                 // 日程标题
  description: string           // 日程描述
  startDate: string             // 开始日期 (YYYY-MM-DD)
  endDate: string               // 结束日期 (YYYY-MM-DD)
  startTime: string             // 开始时间 (HH:mm)
  endTime: string               // 结束时间 (HH:mm)
  isAllDay: boolean             // 是否全天
  priority: SchedulePriority    // 优先级
  remind: ScheduleRemind        // 提醒设置
  location: string              // 地点
  creatorId: number             // 创建者ID
  creatorName: string           // 创建者姓名
  participantIds: number[]      // 参与者ID列表
  participants: Participant[]   // 参与者详情
  color: string                 // 日程颜色
  createTime: string            // 创建时间
  updateTime: string            // 更新时间
}
```

### 3.2 优先级类型

```typescript
type SchedulePriority = 'low' | 'medium' | 'high'

const SchedulePriorityLabels = {
  low: '低',
  medium: '中',
  high: '高'
}
```

优先级颜色对应：
- `low`: #18a058 (绿色)
- `medium`: #f0a020 (橙色)
- `high`: #d03050 (红色)

### 3.3 提醒类型

```typescript
type ScheduleRemind = 'none' | '5min' | '15min' | '30min' | '1hour' | '1day'

const ScheduleRemindLabels = {
  none: '不提醒',
  '5min': '5分钟前',
  '15min': '15分钟前',
  '30min': '30分钟前',
  '1hour': '1小时前',
  '1day': '1天前'
}
```

### 3.4 创建日程参数

```typescript
interface CreateScheduleParams {
  title: string                   // 必填：日程标题
  description?: string            // 可选：日程描述
  startDate: string               // 必填：开始日期
  endDate?: string                // 可选：结束日期（默认同开始日期）
  startTime?: string              // 可选：开始时间
  endTime?: string                // 可选：结束时间
  isAllDay?: boolean              // 可选：是否全天（默认 false）
  priority?: SchedulePriority     // 可选：优先级（默认 medium）
  remind?: ScheduleRemind         // 可选：提醒（默认 none）
  location?: string               // 可选：地点
  participantIds?: number[]       // 可选：参与者ID列表
  color?: string                  // 可选：日程颜色
}
```

## 4. 接口信息

### 4.1 获取日程列表

**接口**: `GET /api/schedule/list`

**参数**:
| 参数 | 类型 | 必填 | 说明 |
|-----|------|-----|------|
| startDate | string | 是 | 开始日期 |
| endDate | string | 是 | 结束日期 |
| creatorId | number | 否 | 创建者ID |

**响应**:
```json
{
  "code": 200,
  "message": "成功",
  "data": [
    {
      "id": 1,
      "title": "项目会议",
      "description": "讨论项目进度",
      ...
    }
  ]
}
```

### 4.2 获取日程详情

**接口**: `GET /api/schedule/:id`

**响应**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "id": 1,
    "title": "项目会议",
    ...
  }
}
```

### 4.3 创建日程

**接口**: `POST /api/schedule/create`

**请求体**: `CreateScheduleParams`

**响应**:
```json
{
  "code": 200,
  "message": "创建成功",
  "data": { /* 新创建的日程对象 */ }
}
```

### 4.4 更新日程

**接口**: `PUT /api/schedule/:id`

**请求体**: `Partial<CreateScheduleParams>`

**响应**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": { /* 更新后的日程对象 */ }
}
```

### 4.5 删除日程

**接口**: `DELETE /api/schedule/:id`

**响应**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

### 4.6 获取本周日程

**接口**: `GET /api/schedule/week`

**响应**:
```json
{
  "code": 200,
  "message": "成功",
  "data": [ /* 本周的日程列表 */ ]
}
```

## 5. 主要实现逻辑

### 5.1 Store 状态管理 (`stores/schedule.ts`)

使用 Pinia 进行状态管理，主要包含：

**状态**:
- `schedules`: 日程列表数据
- `weekSchedules`: 本周日程数据
- `currentSchedule`: 当前选中的日程
- `loading`: 加载状态

**Actions**:
- `getSchedules`: 获取日程列表
- `getWeekSchedules`: 获取本周日程
- `getScheduleDetail`: 获取日程详情
- `createSchedule`: 创建日程
- `updateSchedule`: 更新日程
- `deleteSchedule`: 删除日程

**Getters/辅助方法**:
- `getSchedulesByDate`: 根据日期获取日程
- `getScheduleCountByDate`: 获取日期的日程数量
- `getPriorityColor`: 获取优先级颜色
- `getPriorityLabel`: 获取优先级标签
- `getRemindLabel`: 获取提醒标签

### 5.2 日历页面实现要点

1. **日历组件**: 使用 Naive UI 的 `n-calendar` 组件
2. **日程标记**: 在日历单元格中通过 Badge 显示日程数量
3. **日期选择**: 点击日期后更新右侧日程列表
4. **日程表单**: 使用 `n-modal` 弹窗，包含完整的表单字段
5. **表单验证**: 使用 `n-form` 的验证规则，确保必填字段

### 5.3 列表页面实现要点

1. **数据表格**: 使用 `n-data-table` 组件
2. **筛选功能**: 支持日期范围、优先级、关键字筛选
3. **分页功能**: 使用 `n-pagination` 组件
4. **操作按钮**: 在表格操作列渲染查看、编辑、删除按钮
5. **删除确认**: 使用 `useDialog` 进行删除确认

### 5.4 表单处理

1. **日期时间处理**: 使用 `n-date-picker` 和 `n-time-picker` 组件
2. **全天切换**: 当选择全天时，清空时间字段
3. **颜色选择**: 使用自定义颜色选择器，提供预设颜色
4. **表单重置**: 新增和编辑共用表单，需要正确重置状态

## 6. 文件结构

```
src/
├── views/
│   └── schedule/
│       ├── Calendar.vue      # 日程日历页面
│       └── List.vue          # 日程列表页面
├── stores/
│   └── schedule.ts           # 日程 Store
├── types/
│   └── schedule.d.ts         # 日程类型定义
└── mock/
    └── schedule.ts           # 日程 Mock 数据
```

## 7. 路由配置

```typescript
{
  path: '/schedule',
  name: 'ScheduleLayout',
  component: () => import('@/components/Layout/OALayout.vue'),
  redirect: '/schedule/calendar',
  meta: {
    title: '日程管理',
    icon: 'calendar-outline',
    requiresAuth: true
  },
  children: [
    {
      path: 'calendar',
      name: 'ScheduleCalendar',
      component: () => import('@/views/schedule/Calendar.vue'),
      meta: {
        title: '日程日历',
        icon: 'calendar-outline',
        requiresAuth: true
      }
    },
    {
      path: 'list',
      name: 'ScheduleList',
      component: () => import('@/views/schedule/List.vue'),
      meta: {
        title: '日程列表',
        icon: 'list-outline',
        requiresAuth: true
      }
    }
  ]
}
```
