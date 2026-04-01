import { MockMethod } from 'vite-plugin-mock'
import Mock from 'mockjs'

const Random = Mock.Random

const colors = ['#2080f0', '#18a058', '#f0a020', '#d03050', '#8a2be2', '#00ced1']

// 生成日程数据
const schedules = Array.from({ length: 30 }, (_, i) => {
  const baseDate = new Date()
  baseDate.setDate(baseDate.getDate() + Random.integer(-7, 14))

  const isAllDay = Random.boolean()
  const startHour = Random.integer(8, 18)
  const endHour = Math.min(startHour + Random.integer(1, 3), 22)

  return {
    id: i + 1,
    title: Random.ctitle(4, 12),
    description: Random.cparagraph(1, 2),
    startDate: Random.date('yyyy-MM-dd'),
    endDate: Random.date('yyyy-MM-dd'),
    startTime: isAllDay ? '' : `${startHour.toString().padStart(2, '0')}:00`,
    endTime: isAllDay ? '' : `${endHour.toString().padStart(2, '0')}:00`,
    isAllDay,
    priority: Random.pick(['low', 'medium', 'high']),
    remind: Random.pick(['none', '5min', '15min', '30min', '1hour', '1day']),
    location: Random.city(),
    creatorId: 1,
    creatorName: '管理员',
    participantIds: Random.shuffle([1, 2, 3, 4, 5]).slice(0, Random.integer(0, 3)),
    participants: [],
    color: Random.pick(colors),
    createTime: Random.datetime('yyyy-MM-dd HH:mm:ss'),
    updateTime: Random.datetime('yyyy-MM-dd HH:mm:ss')
  }
})

// 生成当月的日程
const today = new Date()
const currentMonth = today.getMonth()
const currentYear = today.getFullYear()

for (let day = 1; day <= 28; day += Random.integer(1, 4)) {
  const date = `${currentYear}-${(currentMonth + 1).toString().padStart(2, '0')}-${day.toString().padStart(2, '0')}`

  const schedule = {
    id: schedules.length + 1,
    title: Random.ctitle(4, 12),
    description: Random.cparagraph(1, 2),
    startDate: date,
    endDate: date,
    startTime: `${Random.integer(8, 18).toString().padStart(2, '0')}:00`,
    endTime: `${Random.integer(19, 22).toString().padStart(2, '0')}:00`,
    isAllDay: false,
    priority: Random.pick(['low', 'medium', 'high']),
    remind: Random.pick(['none', '15min', '30min', '1hour']),
    location: Random.pick(['会议室A', '会议室B', '办公室', '线上会议', '']),
    creatorId: 1,
    creatorName: '管理员',
    participantIds: [],
    participants: [],
    color: Random.pick(colors),
    createTime: Random.datetime('yyyy-MM-dd HH:mm:ss'),
    updateTime: Random.datetime('yyyy-MM-dd HH:mm:ss')
  }

  schedules.push(schedule)
}

export default [
  // 获取日程列表
  {
    url: '/api/schedules',
    method: 'get',
    response: ({
      query
    }: {
      query: { startDate: string; endDate: string; creatorId?: number }
    }) => {
      const { startDate, endDate, creatorId } = query

      let filteredSchedules = schedules.filter((s) => {
        const scheduleDate = new Date(s.startDate)
        const start = new Date(startDate)
        const end = new Date(endDate)
        return scheduleDate >= start && scheduleDate <= end
      })

      if (creatorId) {
        filteredSchedules = filteredSchedules.filter((s) => s.creatorId === creatorId)
      }

      return {
        code: 200,
        message: '成功',
        data: filteredSchedules
      }
    }
  },

  // 获取日程详情
  {
    url: '/api/schedules/:id',
    method: 'get',
    response: ({ query }: { query: { id: string } }) => {
      const id = parseInt(query.id)
      const schedule = schedules.find((s) => s.id === id)

      if (schedule) {
        return {
          code: 200,
          message: '成功',
          data: schedule
        }
      }

      return {
        code: 404,
        message: '日程不存在',
        data: null
      }
    }
  },

  // 创建日程
  {
    url: '/api/schedules',
    method: 'post',
    response: ({ body }: { body: Record<string, unknown> }) => {
      const newSchedule = {
        id: schedules.length + 1,
        title: body.title as string,
        description: body.description || '',
        startDate: body.startDate as string,
        endDate: body.endDate || body.startDate,
        startTime: body.startTime || '',
        endTime: body.endTime || '',
        isAllDay: body.isAllDay || false,
        priority: body.priority || 'medium',
        remind: body.remind || 'none',
        location: body.location || '',
        creatorId: 1,
        creatorName: '管理员',
        participantIds: body.participantIds || [],
        participants: [],
        color: body.color || Random.pick(colors),
        createTime: Random.datetime('yyyy-MM-dd HH:mm:ss'),
        updateTime: Random.datetime('yyyy-MM-dd HH:mm:ss')
      }

      schedules.push(newSchedule as (typeof schedules)[0])

      return {
        code: 200,
        message: '创建成功',
        data: newSchedule
      }
    }
  },

  // 更新日程
  {
    url: '/api/schedules/:id',
    method: 'put',
    response: ({ query, body }: { query: { id: string }; body: Record<string, unknown> }) => {
      const id = parseInt(query.id)
      const index = schedules.findIndex((s) => s.id === id)

      if (index > -1) {
        schedules[index] = {
          ...schedules[index],
          ...body,
          updateTime: Random.datetime('yyyy-MM-dd HH:mm:ss')
        } as (typeof schedules)[0]

        return {
          code: 200,
          message: '更新成功',
          data: schedules[index]
        }
      }

      return {
        code: 404,
        message: '日程不存在',
        data: null
      }
    }
  },

  // 删除日程
  {
    url: '/api/schedules/:id',
    method: 'delete',
    response: ({ query }: { query: { id: string } }) => {
      const id = parseInt(query.id)
      const index = schedules.findIndex((s) => s.id === id)

      if (index > -1) {
        schedules.splice(index, 1)

        return {
          code: 200,
          message: '删除成功',
          data: null
        }
      }

      return {
        code: 404,
        message: '日程不存在',
        data: null
      }
    }
  },

  // 获取本周日程
  {
    url: '/api/schedules/week',
    method: 'get',
    response: () => {
      const today = new Date()
      const dayOfWeek = today.getDay()
      const startOfWeek = new Date(today)
      startOfWeek.setDate(today.getDate() - dayOfWeek)
      const endOfWeek = new Date(today)
      endOfWeek.setDate(today.getDate() + (6 - dayOfWeek))

      const weekSchedules = schedules.filter((s) => {
        const scheduleDate = new Date(s.startDate)
        return scheduleDate >= startOfWeek && scheduleDate <= endOfWeek
      })

      return {
        code: 200,
        message: '成功',
        data: weekSchedules
      }
    }
  }
] as MockMethod[]
