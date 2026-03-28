// 审批相关类型

// 审批类型
export type ApprovalType = 'leave' | 'expense' | 'overtime' | 'travel' | 'general'

// 审批状态
export type ApprovalStatus = 'pending' | 'approved' | 'rejected' | 'withdrawn' | 'transferred'

// 审批流程节点
export interface ApprovalNode {
  id: number
  nodeId: number
  nodeName: string
  approverId: number
  approverName: string
  approverAvatar: string
  status: ApprovalStatus
  comment: string
  sort: number
  createTime: string
}

// 基础申请表单
export interface ApprovalBase {
  id: number
  type: ApprovalType
  title: string
  applicantId: number
  applicantName: string
  applicantDept: string
  status: ApprovalStatus
  currentStep: number
  totalStep: number
  createTime: string
  updateTime: string
  flowNodes: ApprovalNode[]
}

// 请假申请
export interface LeaveApproval extends ApprovalBase {
  type: 'leave'
  leaveType: 'annual' | 'sick' | 'personal' | 'maternity' | 'marriage' | 'bereavement'
  startDate: string
  endDate: string
  days: number
  reason: string
}

// 报销申请
export interface ExpenseApproval extends ApprovalBase {
  type: 'expense'
  expenseType: 'travel' | 'office' | 'entertainment' | 'other'
  amount: number
  description: string
  attachments: string[]
}

// 加班申请
export interface OvertimeApproval extends ApprovalBase {
  type: 'overtime'
  overtimeDate: string
  startTime: string
  endTime: string
  hours: number
  reason: string
}

// 出差申请
export interface TravelApproval extends ApprovalBase {
  type: 'travel'
  destination: string
  startDate: string
  endDate: string
  days: number
  reason: string
  budget: number
}

// 通用审批
export interface GeneralApproval extends ApprovalBase {
  type: 'general'
  content: string
  attachments: string[]
}

export type Approval = LeaveApproval | ExpenseApproval | OvertimeApproval | TravelApproval | GeneralApproval

// 发起申请表单
export interface CreateApprovalParams {
  type: ApprovalType
  title: string
  [key: string]: unknown
}

// 审批操作参数
export interface ApprovalActionParams {
  approvalId: number
  action: 'approve' | 'reject' | 'transfer' | 'withdraw'
  comment?: string
  transferTo?: number
}

// 审批类型标签
export const ApprovalTypeLabels: Record<ApprovalType, string> = {
  leave: '请假申请',
  expense: '报销申请',
  overtime: '加班申请',
  travel: '出差申请',
  general: '通用审批'
}

// 审批状态标签
export const ApprovalStatusLabels: Record<ApprovalStatus, string> = {
  pending: '待审批',
  approved: '已通过',
  rejected: '已驳回',
  withdrawn: '已撤回',
  transferred: '已转交'
}

// 请假类型标签
export const LeaveTypeLabels: Record<string, string> = {
  annual: '年假',
  sick: '病假',
  personal: '事假',
  maternity: '产假',
  marriage: '婚假',
  bereavement: '丧假'
}

// 报销类型标签
export const ExpenseTypeLabels: Record<string, string> = {
  travel: '差旅费',
  office: '办公费',
  entertainment: '招待费',
  other: '其他'
}
