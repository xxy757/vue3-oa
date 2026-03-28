# 审批中心模块文档

## 1. 模块功能说明

审批中心是企业OA办公系统的核心模块之一，提供完整的申请和审批流程管理功能。主要功能包括：

- **发起申请**：支持多种申请类型（请假、报销、加班、出差、通用审批）
- **我的申请**：查看个人提交的申请记录，支持筛选和撤回
- **待我审批**：查看待处理的审批任务，支持批准/驳回操作
- **已办审批**：查看已处理的审批记录
- **审批详情**：查看申请详细信息和审批流程进度

## 2. 各页面组件说明

### 2.1 发起申请页面 (Apply.vue)

**文件路径**: `src/views/approval/Apply.vue`

**功能说明**:
- 选择申请类型（5种类型可选）
- 根据不同类型动态显示表单字段
- 表单验证和提交

**主要组件**:
- `NForm`: 表单容器
- `NSelect`: 申请类型、请假类型、报销类型选择
- `NDatePicker`: 日期选择
- `NTimePicker`: 时间选择
- `NInputNumber`: 数字输入（金额、天数、小时）
- `NUpload`: 文件上传

### 2.2 我的申请页面 (MyApply.vue)

**文件路径**: `src/views/approval/MyApply.vue`

**功能说明**:
- 表格展示个人申请列表
- 按状态、类型筛选
- 分页功能
- 查看详情、撤回操作

**主要组件**:
- `NDataTable`: 数据表格
- `NSelect`: 筛选下拉
- `NPagination`: 分页
- `NPopconfirm`: 撤回确认

### 2.3 待我审批页面 (Pending.vue)

**文件路径**: `src/views/approval/Pending.vue`

**功能说明**:
- 表格展示待审批列表
- 批准/驳回操作（弹窗确认）
- 查看详情
- 分页功能

**主要组件**:
- `NDataTable`: 数据表格
- `NModal`: 审批操作弹窗
- `NPagination`: 分页

### 2.4 已办审批页面 (Done.vue)

**文件路径**: `src/views/approval/Done.vue`

**功能说明**:
- 表格展示已处理的审批记录
- 查看详情
- 分页功能

**主要组件**:
- `NDataTable`: 数据表格
- `NPagination`: 分页

### 2.5 审批详情页面 (Detail.vue)

**文件路径**: `src/views/approval/Detail.vue`

**功能说明**:
- 显示申请基本信息（类型特定字段）
- 显示审批流程进度（时间线）
- 审批操作（批准/驳回/转交）
- 撤回功能（仅申请人可用）

**主要组件**:
- `NDescriptions`: 信息描述列表
- `NTimeline`: 审批流程时间线
- `NModal`: 审批/转交弹窗
- `NPopconfirm`: 撤回确认

## 3. 申请类型和表单字段

### 3.1 请假申请 (leave)

| 字段名 | 类型 | 说明 | 必填 |
|--------|------|------|------|
| title | string | 申请标题 | 是 |
| leaveType | enum | 请假类型：annual(年假)、sick(病假)、personal(事假)、maternity(产假)、marriage(婚假)、bereavement(丧假) | 是 |
| startDate | date | 开始日期 | 是 |
| endDate | date | 结束日期 | 是 |
| days | number | 请假天数（自动计算） | 是 |
| reason | string | 请假原因 | 是 |

### 3.2 报销申请 (expense)

| 字段名 | 类型 | 说明 | 必填 |
|--------|------|------|------|
| title | string | 申请标题 | 是 |
| expenseType | enum | 报销类型：travel(差旅费)、office(办公费)、entertainment(招待费)、other(其他) | 是 |
| amount | number | 报销金额 | 是 |
| description | string | 报销说明 | 是 |
| attachments | array | 附件列表 | 否 |

### 3.3 加班申请 (overtime)

| 字段名 | 类型 | 说明 | 必填 |
|--------|------|------|------|
| title | string | 申请标题 | 是 |
| overtimeDate | date | 加班日期 | 是 |
| startTime | time | 开始时间 | 是 |
| endTime | time | 结束时间 | 是 |
| hours | number | 加班时长（自动计算） | 是 |
| reason | string | 加班原因 | 是 |

### 3.4 出差申请 (travel)

| 字段名 | 类型 | 说明 | 必填 |
|--------|------|------|------|
| title | string | 申请标题 | 是 |
| destination | string | 目的地 | 是 |
| startDate | date | 开始日期 | 是 |
| endDate | date | 结束日期 | 是 |
| days | number | 出差天数（自动计算） | 是 |
| reason | string | 出差原因 | 是 |
| budget | number | 预算金额 | 是 |

### 3.5 通用审批 (general)

| 字段名 | 类型 | 说明 | 必填 |
|--------|------|------|------|
| title | string | 申请标题 | 是 |
| content | string | 申请内容 | 是 |
| attachments | array | 附件列表 | 否 |

## 4. 接口信息

### 4.1 获取我的申请列表

- **URL**: `GET /api/approval/my`
- **参数**:
  - `page`: 页码
  - `pageSize`: 每页数量
  - `status`: 审批状态（可选）
  - `type`: 申请类型（可选）
- **返回**: `PageResult<Approval>`

### 4.2 获取待审批列表

- **URL**: `GET /api/approval/pending`
- **参数**:
  - `page`: 页码
  - `pageSize`: 每页数量
- **返回**: `PageResult<Approval>`

### 4.3 获取已办审批列表

- **URL**: `GET /api/approval/done`
- **参数**:
  - `page`: 页码
  - `pageSize`: 每页数量
- **返回**: `PageResult<Approval>`

### 4.4 获取审批详情

- **URL**: `GET /api/approval/:id`
- **参数**: 无（ID在URL路径中）
- **返回**: `Approval`

### 4.5 发起申请

- **URL**: `POST /api/approval/create`
- **参数**: `CreateApprovalParams`
  - `type`: 申请类型
  - `title`: 申请标题
  - 其他类型特定字段
- **返回**: `Approval`

### 4.6 审批操作

- **URL**: `POST /api/approval/action`
- **参数**: `ApprovalActionParams`
  - `approvalId`: 审批ID
  - `action`: 操作类型（approve/reject/transfer）
  - `comment`: 审批意见（可选）
  - `transferTo`: 转交人员ID（转交时必填）
- **返回**: 无

### 4.7 撤回申请

- **URL**: `POST /api/approval/withdraw`
- **参数**:
  - `id`: 申请ID
- **返回**: 无

### 4.8 获取审批统计

- **URL**: `GET /api/approval/stats`
- **参数**: 无
- **返回**:
  ```json
  {
    "myPending": 2,
    "myApproved": 5,
    "myRejected": 1,
    "todoApproval": 3
  }
  ```

## 5. 审批流程说明

### 5.1 流程状态

| 状态 | 说明 |
|------|------|
| pending | 待审批 |
| approved | 已通过 |
| rejected | 已驳回 |
| withdrawn | 已撤回 |
| transferred | 已转交 |

### 5.2 流程节点

每个审批申请包含多个审批节点（`flowNodes`），每个节点包含：

- `nodeId`: 节点ID
- `nodeName`: 节点名称（如：部门经理审批、人事审批）
- `approverId`: 审批人ID
- `approverName`: 审批人姓名
- `approverAvatar`: 审批人头像
- `status`: 节点状态
- `comment`: 审批意见
- `sort`: 节点顺序

### 5.3 审批规则

1. **申请提交**: 用户提交申请后，状态变为`pending`，流程进入第一个审批节点
2. **逐级审批**: 审批按节点顺序依次进行，当前节点完成后进入下一节点
3. **批准操作**: 审批人批准后，流程进入下一节点；所有节点通过后，申请状态变为`approved`
4. **驳回操作**: 审批人驳回后，流程终止，申请状态变为`rejected`
5. **转交操作**: 审批人可将审批权限转交给其他人，状态变为`transferred`
6. **撤回操作**: 申请人在审批未完成前可撤回申请，状态变为`withdrawn`

### 5.4 权限控制

- **发起申请**: 所有登录用户
- **我的申请**: 查看本人提交的申请
- **待我审批**: 查看待本人审批的任务
- **已办审批**: 查看本人已处理的审批
- **撤回**: 仅申请人本人，且申请状态为`pending`

## 6. Store 使用说明

### 6.1 引入 Store

```typescript
import { useApprovalStore } from '@/stores/approval'

const approvalStore = useApprovalStore()
```

### 6.2 可用方法

```typescript
// 获取我的申请列表
const result = await approvalStore.getMyApprovals({ page: 1, pageSize: 10 })

// 获取待审批列表
const result = await approvalStore.getPendingApprovals({ page: 1, pageSize: 10 })

// 获取已办审批列表
const result = await approvalStore.getDoneApprovals({ page: 1, pageSize: 10 })

// 获取审批详情
const detail = await approvalStore.getApprovalDetail(id)

// 发起申请
const newApproval = await approvalStore.createApproval(params)

// 审批操作
await approvalStore.approvalAction({ approvalId, action, comment })

// 撤回申请
await approvalStore.withdrawApproval(id)
```

### 6.3 可用状态

```typescript
// 我的申请列表
approvalStore.myApprovals

// 待审批列表
approvalStore.pendingApprovals

// 已办审批列表
approvalStore.doneApprovals

// 当前审批详情
approvalStore.currentApproval

// 未读待审批数量
approvalStore.unreadCount
```

## 7. 类型定义

类型定义文件位于 `src/types/approval.d.ts`，主要类型包括：

- `ApprovalType`: 申请类型
- `ApprovalStatus`: 审批状态
- `ApprovalBase`: 基础申请信息
- `LeaveApproval`: 请假申请
- `ExpenseApproval`: 报销申请
- `OvertimeApproval`: 加班申请
- `TravelApproval`: 出差申请
- `GeneralApproval`: 通用审批
- `Approval`: 所有申请类型的联合类型
- `ApprovalNode`: 审批流程节点
- `CreateApprovalParams`: 发起申请参数
- `ApprovalActionParams`: 审批操作参数
