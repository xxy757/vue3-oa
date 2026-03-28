import { defineStore } from 'pinia'
import { ref } from 'vue'
import type {
  Schedule,
  SchedulePriority,
  ScheduleRemind,
  CreateScheduleParams,
  ScheduleListParams
} from '@/types/schedule'
import { request } from '@/utils/request'

export const useScheduleStore = defineStore('schedule', () => {
  const schedules = ref<Schedule[]>([])
  const currentSchedule = ref<Schedule | null>(null)
  const weekSchedules = ref<Schedule[]>([])
  const loading = ref(false)

  // 获取日程列表
  async function getSchedules(params: ScheduleListParams): Promise<Schedule[]> {
    loading.value = true
    try {
      const result: Schedule[] = await request.get('/schedule/list', { params })
      schedules.value = result
      return result
    } finally {
      loading.value = false
    }
  }

  // 获取日程详情
  async function getScheduleDetail(id: number): Promise<Schedule> {
    loading.value = true
    try {
      const result: Schedule = await request.get(`/schedule/${id}`)
      currentSchedule.value = result
      return result
    } finally {
      loading.value = false
    }
  }

  // 创建日程
  async function createSchedule(params: CreateScheduleParams): Promise<Schedule> {
    loading.value = true
    try {
      const result: Schedule = await request.post('/schedule/create', params)
      schedules.value.push(result)
      return result
    } finally {
      loading.value = false
    }
  }

  // 更新日程
  async function updateSchedule(id: number, params: Partial<CreateScheduleParams>): Promise<Schedule> {
    loading.value = true
    try {
      const result: Schedule = await request.put(`/schedule/${id}`, params)
      const index = schedules.value.findIndex(s => s.id === id)
      if (index > -1) {
        schedules.value[index] = result
      }
      if (currentSchedule.value?.id === id) {
        currentSchedule.value = result
      }
      return result
    } finally {
      loading.value = false
    }
  }

  // 删除日程
  async function deleteSchedule(id: number): Promise<void> {
    loading.value = true
    try {
      await request.delete(`/schedule/${id}`)
      schedules.value = schedules.value.filter(s => s.id !== id)
      if (currentSchedule.value?.id === id) {
        currentSchedule.value = null
      }
    } finally {
      loading.value = false
    }
  }

  // 获取本周日程
  async function getWeekSchedules(): Promise<Schedule[]> {
    loading.value = true
    try {
      const result: Schedule[] = await request.get('/schedule/week')
      weekSchedules.value = result
      return result
    } finally {
      loading.value = false
    }
  }

  // 根据日期获取日程
  function getSchedulesByDate(date: string): Schedule[] {
    return schedules.value.filter(s => s.startDate === date)
  }

  // 获取日期的日程数量
  function getScheduleCountByDate(date: string): number {
    return schedules.value.filter(s => s.startDate === date).length
  }

  // 获取优先级颜色
  function getPriorityColor(priority: SchedulePriority): string {
    const colors: Record<SchedulePriority, string> = {
      low: '#18a058',
      medium: '#f0a020',
      high: '#d03050'
    }
    return colors[priority]
  }

  // 获取优先级标签
  function getPriorityLabel(priority: SchedulePriority): string {
    const labels: Record<SchedulePriority, string> = {
      low: '低',
      medium: '中',
      high: '高'
    }
    return labels[priority]
  }

  // 获取提醒标签
  function getRemindLabel(remind: ScheduleRemind): string {
    const labels: Record<ScheduleRemind, string> = {
      none: '不提醒',
      '5min': '5分钟前',
      '15min': '15分钟前',
      '30min': '30分钟前',
      '1hour': '1小时前',
      '1day': '1天前'
    }
    return labels[remind]
  }

  return {
    schedules,
    currentSchedule,
    weekSchedules,
    loading,
    getSchedules,
    getScheduleDetail,
    createSchedule,
    updateSchedule,
    deleteSchedule,
    getWeekSchedules,
    getSchedulesByDate,
    getScheduleCountByDate,
    getPriorityColor,
    getPriorityLabel,
    getRemindLabel
  }
})
