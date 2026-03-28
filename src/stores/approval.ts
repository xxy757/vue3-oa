import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Approval, ApprovalType, ApprovalStatus, CreateApprovalParams, ApprovalActionParams } from '@/types/approval'
import type { PageResult } from '@/types/common'
import { request } from '@/utils/request'

export const useApprovalStore = defineStore('approval', () => {
  const myApprovals = ref<Approval[]>([])
  const pendingApprovals = ref<Approval[]>([])
  const doneApprovals = ref<Approval[]>([])
  const currentApproval = ref<Approval | null>(null)
  const unreadCount = ref(0)

  // 获取我的申请列表
  async function getMyApprovals(params: { page: number; pageSize: number; status?: ApprovalStatus; type?: ApprovalType }): Promise<PageResult<Approval>> {
    const result: PageResult<Approval> = await request.get('/approval/my', { params })
    myApprovals.value = result.list
    return result
  }

  // 获取待审批列表
  async function getPendingApprovals(params: { page: number; pageSize: number }): Promise<PageResult<Approval>> {
    const result: PageResult<Approval> = await request.get('/approval/pending', { params })
    pendingApprovals.value = result.list
    unreadCount.value = result.total
    return result
  }

  // 获取已办审批列表
  async function getDoneApprovals(params: { page: number; pageSize: number }): Promise<PageResult<Approval>> {
    const result: PageResult<Approval> = await request.get('/approval/done', { params })
    doneApprovals.value = result.list
    return result
  }

  // 获取审批详情
  async function getApprovalDetail(id: number): Promise<Approval> {
    const result: Approval = await request.get(`/approval/${id}`)
    currentApproval.value = result
    return result
  }

  // 发起申请
  async function createApproval(params: CreateApprovalParams): Promise<Approval> {
    const result: Approval = await request.post('/approval/create', params)
    return result
  }

  // 审批操作
  async function approvalAction(params: ApprovalActionParams): Promise<void> {
    await request.post('/approval/action', params)
  }

  // 撤回申请
  async function withdrawApproval(id: number): Promise<void> {
    await request.post('/approval/withdraw', { id })
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
