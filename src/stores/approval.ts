import { defineStore } from 'pinia'
import { ref } from 'vue'
import type {
  Approval,
  ApprovalType,
  ApprovalStatus,
  CreateApprovalParams,
  ApprovalActionParams
} from '@/types/approval'
import type { PageResult } from '@/types/common'
import { request } from '@/utils/request'

export const useApprovalStore = defineStore('approval', () => {
  const myApprovals = ref<Approval[]>([])
  const pendingApprovals = ref<Approval[]>([])
  const doneApprovals = ref<Approval[]>([])
  const currentApproval = ref<Approval | null>(null)
  const unreadCount = ref(0)

  // 获取我的申请列表
  async function getMyApprovals(params: {
    page: number
    pageSize: number
    status?: ApprovalStatus
    type?: ApprovalType
  }): Promise<PageResult<Approval>> {
    const result: PageResult<Approval> = await request.get('/approvals/my', { params })
    myApprovals.value = result.list
    return result
  }

  async function getPendingApprovals(params: {
    page: number
    pageSize: number
  }): Promise<PageResult<Approval>> {
    const result: PageResult<Approval> = await request.get('/approvals/pending', { params })
    pendingApprovals.value = result.list
    unreadCount.value = result.total
    return result
  }

  async function getDoneApprovals(params: {
    page: number
    pageSize: number
  }): Promise<PageResult<Approval>> {
    const result: PageResult<Approval> = await request.get('/approvals/done', { params })
    doneApprovals.value = result.list
    return result
  }

  async function getApprovalDetail(id: number): Promise<Approval> {
    const result: Approval = await request.get(`/approvals/${id}`)
    currentApproval.value = result
    return result
  }

  async function createApproval(params: CreateApprovalParams): Promise<Approval> {
    const result: Approval = await request.post('/approvals', params)
    return result
  }

  async function approvalAction(params: ApprovalActionParams): Promise<void> {
    await request.post(`/approvals/${params.approvalId}/action`, {
      action: params.action,
      comment: params.comment,
      targetUserId: params.transferTo
    })
  }

  async function withdrawApproval(id: number): Promise<void> {
    await request.post(`/approvals/${id}/withdraw`)
  }

  return {
    myApprovals,
    pendingApprovals,
    doneApprovals,
    currentApproval,
    unreadCount,
    getMyApprovals,
    getPendingApprovals,
    getDoneApprovals,
    getApprovalDetail,
    createApproval,
    approvalAction,
    withdrawApproval
  }
})
