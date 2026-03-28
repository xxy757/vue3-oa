# 公告通知模块文档

## 1. 模块功能说明

公告通知模块是 OA 系统中用于发布和管理企业内部公告、通知、政策和紧急消息的功能模块。该模块支持：

- 按类型分类浏览公告（通知、公告、制度、紧急）
- 关键词搜索公告
- 置顶公告展示
- 未读公告高亮显示
- 公告详情阅读
- 阅读次数统计

## 2. 页面组件说明

### 2.1 公告列表页 (`src/views/notice/List.vue`)

公告列表页面是公告模块的主入口，主要功能包括：

**功能特性：**
- 使用 `n-data-table` 表格展示公告列表
- 支持按类型筛选（通知、公告、制度、紧急）
- 支持关键词搜索（标题）
- 置顶公告显示置顶图标
- 未读公告行背景高亮，标题加粗
- 分页功能，支持切换每页条数

**组件依赖：**
- `NCard` - 卡片容器
- `NSelect` - 类型筛选下拉
- `NInput` - 搜索输入框
- `NButton` - 按钮
- `NDataTable` - 数据表格
- `NTag` - 类型标签
- `NPagination` - 分页组件

**主要方法：**
| 方法名 | 说明 |
|--------|------|
| `fetchNoticeList` | 获取公告列表数据 |
| `handleSearch` | 执行搜索，重置到第一页 |
| `handleFilterChange` | 类型筛选变化处理 |
| `handlePageChange` | 分页变化处理 |
| `handlePageSizeChange` | 每页条数变化处理 |
| `handleViewDetail` | 跳转到公告详情 |

### 2.2 公告详情页 (`src/views/notice/Detail.vue`)

公告详情页面用于展示单条公告的完整内容。

**功能特性：**
- 显示公告标题、类型标签、置顶标识
- 显示发布者、发布时间、阅读次数
- 富文本内容渲染
- 返回列表按钮

**组件依赖：**
- `NCard` - 卡片容器
- `NSpin` - 加载状态
- `NTag` - 类型和置顶标签
- `NIcon` - 图标
- `NButton` - 按钮
- `NSpace` - 间距布局
- `NEmpty` - 空状态

**主要方法：**
| 方法名 | 说明 |
|--------|------|
| `fetchNoticeDetail` | 根据 ID 获取公告详情 |
| `handleGoBack` | 返回公告列表 |

## 3. 公告类型说明

### 3.1 类型枚举 (`NoticeType`)

| 类型值 | 中文名称 | 标签颜色 | 说明 |
|--------|----------|----------|------|
| `notice` | 通知 | info (蓝色) | 一般性通知消息 |
| `announcement` | 公告 | default (灰色) | 正式公告 |
| `policy` | 制度 | success (绿色) | 公司制度、规范 |
| `urgent` | 紧急 | error (红色) | 紧急重要消息 |

### 3.2 状态枚举 (`NoticeStatus`)

| 状态值 | 中文名称 | 说明 |
|--------|----------|------|
| `draft` | 草稿 | 未发布的草稿状态 |
| `published` | 已发布 | 已正式发布 |
| `archived` | 已归档 | 已归档的历史公告 |

### 3.3 数据结构 (`Notice`)

```typescript
interface Notice {
  id: number                 // 公告ID
  title: string              // 公告标题
  type: NoticeType           // 公告类型
  content: string            // 公告内容（富文本HTML）
  summary: string            // 摘要
  coverImage: string         // 封面图片URL
  publisherId: number        // 发布者ID
  publisherName: string      // 发布者名称
  publisherAvatar: string    // 发布者头像
  status: NoticeStatus       // 状态
  isTop: boolean             // 是否置顶
  readCount: number          // 阅读次数
  createTime: string         // 创建时间
  updateTime: string         // 更新时间
  isRead?: boolean           // 当前用户是否已读
}
```

## 4. 接口信息

### 4.1 获取公告列表

- **URL**: `GET /api/notice/list`
- **参数**:
  | 参数名 | 类型 | 必填 | 说明 |
  |--------|------|------|------|
  | page | number | 否 | 页码，默认 1 |
  | pageSize | number | 否 | 每页条数，默认 10 |
  | type | string | 否 | 类型筛选 |
  | keyword | string | 否 | 关键词搜索 |

- **响应**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "list": [/* Notice数组 */],
    "total": 100,
    "page": 1,
    "pageSize": 10,
    "totalPages": 10
  }
}
```

### 4.2 获取公告详情

- **URL**: `GET /api/notice/:id`
- **参数**: URL 路径参数 `id`
- **响应**: Notice 对象
- **说明**: 获取详情时会自动增加阅读次数并标记为已读

### 4.3 标记已读

- **URL**: `POST /api/notice/:id/read`
- **参数**: URL 路径参数 `id`
- **响应**: `{ code: 200, message: "成功", data: null }`

### 4.4 获取未读数量

- **URL**: `GET /api/notice/unread-count`
- **响应**: `{ code: 200, message: "成功", data: { count: 5 } }`

### 4.5 发布公告

- **URL**: `POST /api/notice/create`
- **参数**:
```json
{
  "title": "公告标题",
  "type": "notice",
  "content": "<p>公告内容</p>",
  "summary": "摘要",
  "coverImage": "图片URL",
  "isTop": false
}
```
- **响应**: 新创建的 Notice 对象

## 5. 主要实现逻辑

### 5.1 列表排序逻辑

公告列表的排序规则：
1. 置顶公告优先显示
2. 同为置顶或非置顶时，按发布时间倒序

```typescript
filteredNotices.sort((a, b) => {
  if (a.isTop && !b.isTop) return -1
  if (!a.isTop && b.isTop) return 1
  return new Date(b.createTime).getTime() - new Date(a.createTime).getTime()
})
```

### 5.2 未读高亮逻辑

在表格中通过 `row-class-name` 属性为未读行添加自定义样式类：

```typescript
const getRowClass = (row: Notice) => {
  return row.isRead ? '' : 'unread-row'
}
```

对应的 CSS：
```scss
:deep(.unread-row) {
  background-color: rgba(24, 144, 255, 0.05);
}
```

### 5.3 类型标签颜色映射

```typescript
const getTypeTagType = (type: NoticeType) => {
  const typeMap: Record<NoticeType, 'default' | 'info' | 'success' | 'warning' | 'error'> = {
    notice: 'info',      // 蓝色
    announcement: 'default', // 灰色
    policy: 'success',   // 绿色
    urgent: 'error'      // 红色
  }
  return typeMap[type]
}
```

### 5.4 Store 状态管理

`useNoticeStore` 提供以下状态和方法：

**状态：**
- `noticeList` - 公告列表数据
- `currentNotice` - 当前查看的公告详情
- `unreadCount` - 未读公告数量

**方法：**
- `getNoticeList(params)` - 获取公告列表
- `getNoticeDetail(id)` - 获取公告详情
- `createNotice(params)` - 发布公告
- `markAsRead(id)` - 标记已读
- `getUnreadCount()` - 获取未读数量
