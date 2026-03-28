import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Notice, NoticeType, CreateNoticeParams, NoticeListParams } from '@/types/notice'
import type { PageResult } from '@/types/common'
import { request } from '@/utils/request'

export const useNoticeStore = defineStore('notice', () => {
  const noticeList = ref<Notice[]>([])
  const currentNotice = ref<Notice | null>(null)
  const unreadCount = ref(0)

  // 获取公告列表
  async function getNoticeList(params: NoticeListParams): Promise<PageResult<Notice>> {
    const result: PageResult<Notice> = await request.get('/notice/list', { params })
    noticeList.value = result.list
    return result
  }

  // 获取公告详情
  async function getNoticeDetail(id: number): Promise<Notice> {
    const result: Notice = await request.get(`/notice/${id}`)
    currentNotice.value = result
    return result
  }

  // 发布公告
  async function createNotice(params: CreateNoticeParams): Promise<Notice> {
    const result: Notice = await request.post('/notice/create', params)
    return result
  }

  // 标记已读
  async function markAsRead(id: number): Promise<void> {
    await request.post(`/notice/${id}/read`)
  }

  // 获取未读数量
  async function getUnreadCount(): Promise<number> {
    const result: { count: number } = await request.get('/notice/unread-count')
    unreadCount.value = result.count
    return result.count
  }

  return {
    noticeList,
    currentNotice,
    unreadCount,
    getNoticeList,
    getNoticeDetail,
    createNotice,
    markAsRead,
    getUnreadCount
  }
})
