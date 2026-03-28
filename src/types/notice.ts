// 公告相关类型

export type NoticeType = 'notice' | 'announcement' | 'policy' | 'urgent'
export type NoticeStatus = 'draft' | 'published' | 'archived'

export interface Notice {
  id: number
  title: string
  type: NoticeType
  content: string
  summary: string
  coverImage: string
  publisherId: number
  publisherName: string
  publisherAvatar: string
  status: NoticeStatus
  isTop: boolean
  readCount: number
  createTime: string
  updateTime: string
  isRead?: boolean
}

export interface NoticeListParams {
  page: number
  pageSize: number
  type?: NoticeType
  keyword?: string
}

export interface CreateNoticeParams {
  title: string
  type: NoticeType
  content: string
  summary: string
  coverImage?: string
  isTop?: boolean
}

export const NoticeTypeLabels: Record<NoticeType, string> = {
  notice: '通知',
  announcement: '公告',
  policy: '制度',
  urgent: '紧急'
}

export const NoticeStatusLabels: Record<NoticeStatus, string> = {
  draft: '草稿',
  published: '已发布',
  archived: '已归档'
}
