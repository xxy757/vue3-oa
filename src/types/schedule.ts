// 日程相关类型

export type SchedulePriority = 'low' | 'medium' | 'high'
export type ScheduleRemind = 'none' | '5min' | '15min' | '30min' | '1hour' | '1day'

export interface Schedule {
  id: number
  title: string
  description: string
  startDate: string
  endDate: string
  startTime: string
  endTime: string
  isAllDay: boolean
  priority: SchedulePriority
  remind: ScheduleRemind
  location: string
  creatorId: number
  creatorName: string
  participantIds: number[]
  participants: { id: number; name: string; avatar: string }[]
  color: string
  createTime: string
  updateTime: string
}

export interface CreateScheduleParams {
  title: string
  description?: string
  startDate: string
  endDate?: string
  startTime?: string
  endTime?: string
  isAllDay?: boolean
  priority?: SchedulePriority
  remind?: ScheduleRemind
  location?: string
  participantIds?: number[]
  color?: string
}

export interface ScheduleListParams {
  startDate: string
  endDate: string
  creatorId?: number
}

export const SchedulePriorityLabels: Record<SchedulePriority, string> = {
  low: '低',
  medium: '中',
  high: '高'
}

export const ScheduleRemindLabels: Record<ScheduleRemind, string> = {
  none: '不提醒',
  '5min': '5分钟前',
  '15min': '15分钟前',
  '30min': '30分钟前',
  '1hour': '1小时前',
  '1day': '1天前'
}

export const ScheduleColors = [
  '#2080f0',
  '#18a058',
  '#f0a020',
  '#d03050',
  '#8a2be2',
  '#00ced1',
  '#ff69b4',
  '#708090'
]
