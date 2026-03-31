import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Notice, CreateNoticeParams, NoticeListParams } from '@/types/notice'
import type { PageResult } from '@/types/common'
import { request } from '@/utils/request'

export const useNoticeStore = defineStore('notice', () => {
  const noticeList = ref<Notice[]>([])
  const currentNotice = ref<Notice | null>(null)
  const unreadCount = ref(0)

  // 获取公告列表
  async function getNoticeList(params: NoticeListParams): Promise<PageResult<Notice>> {
    const result: PageResult<Notice> = await request.get('/notices', { params })
    noticeList.value = result.list
    return result
  }

  async function getNoticeDetail(id: number): Promise<Notice> {
    const result: Notice = await request.get(`/notices/${id}`)
    currentNotice.value = result
    return result
  }

  async function createNotice(params: CreateNoticeParams): Promise<Notice> {
    const result: Notice = await request.post('/notices', params)
    return result
  }

  async function markAsRead(id: number): Promise<void> {
    await request.post(`/notices/${id}/read`)
  }

  async function getUnreadCount(): Promise<number> {
    const result: { count: number } = await request.get('/notices/unread-count')
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
