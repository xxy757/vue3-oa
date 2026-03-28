import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.locale('zh-cn')
dayjs.extend(relativeTime)

export function formatDate(date: string | Date, format = 'YYYY-MM-DD'): string {
  return dayjs(date).format(format)
}

export function formatDateTime(date: string | Date, format = 'YYYY-MM-DD HH:mm:ss'): string {
  return dayjs(date).format(format)
}

export function formatTime(date: string | Date, format = 'HH:mm:ss'): string {
  return dayjs(date).format(format)
}

export function fromNow(date: string | Date): string {
  return dayjs(date).fromNow()
}

export function getDayOfWeek(date: string | Date): string {
  const days = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  return days[dayjs(date).day()]
}

export function getWeekDates(date: string | Date): string[] {
  const current = dayjs(date)
  const dayOfWeek = current.day()
  const dates: string[] = []

  for (let i = 0; i < 7; i++) {
    dates.push(current.add(i - dayOfWeek, 'day').format('YYYY-MM-DD'))
  }
  return dates
}

export function getMonthDates(year: number, month: number): (string | null)[][] {
  const firstDay = dayjs(`${year}-${month}-01`)
  const daysInMonth = firstDay.daysInMonth()
  const startDayOfWeek = firstDay.day()

  const weeks: (string | null)[][] = []
  let week: (string | null)[] = []

  // 填充月初空白
  for (let i = 0; i < startDayOfWeek; i++) {
    week.push(null)
  }

  // 填充日期
  for (let day = 1; day <= daysInMonth; day++) {
    week.push(firstDay.date(day).format('YYYY-MM-DD'))

    if (week.length === 7) {
      weeks.push(week)
      week = []
    }
  }

  // 填充月末空白
  if (week.length > 0) {
    while (week.length < 7) {
      week.push(null)
    }
    weeks.push(week)
  }

  return weeks
}

export function isToday(date: string | Date): boolean {
  return dayjs(date).isSame(dayjs(), 'day')
}

export function isSameDay(date1: string | Date, date2: string | Date): boolean {
  return dayjs(date1).isSame(dayjs(date2), 'day')
}

export function calculateDays(startDate: string | Date, endDate: string | Date): number {
  return dayjs(endDate).diff(dayjs(startDate), 'day') + 1
}

export function calculateHours(startTime: string, endTime: string): number {
  const start = dayjs(`2000-01-01 ${startTime}`)
  const end = dayjs(`2000-01-01 ${endTime}`)
  return end.diff(start, 'hour', true)
}
