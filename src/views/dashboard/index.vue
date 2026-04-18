<template>
  <div class="dashboard">
    <!-- 欢迎区域 -->
    <n-card class="welcome-card" :bordered="false">
      <div class="welcome-content">
        <div class="welcome-text">
          <h2 class="welcome-title">{{ getGreeting() }}，{{ userStore.userName }}</h2>
          <p class="welcome-date">
            {{ currentDate }}
          </p>
        </div>
        <div class="welcome-weather">
          <n-icon size="48" :component="SunnyOutline" />
        </div>
      </div>
    </n-card>

    <!-- 统计卡片 -->
    <n-grid :x-gap="16" :y-gap="16" :cols="4" item-responsive responsive="screen">
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card stat-card--pending" hoverable>
          <div class="stat-content">
            <div class="stat-icon">
              <n-icon size="32" :component="DocumentTextOutline" />
            </div>
            <div class="stat-info">
              <n-statistic label="待审批">
                <template #default>
                  <n-number-animation :from="0" :to="stats.todoApproval" :duration="1000" />
                </template>
              </n-statistic>
            </div>
          </div>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card stat-card--notice" hoverable @click="router.push('/notice/list')">
          <div class="stat-content">
            <div class="stat-icon">
              <n-icon size="32" :component="MegaphoneOutline" />
            </div>
            <div class="stat-info">
              <n-statistic label="未读公告">
                <template #default>
                  <n-number-animation :from="0" :to="stats.unreadNotice" :duration="1000" />
                </template>
              </n-statistic>
            </div>
          </div>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card
          class="stat-card stat-card--schedule"
          hoverable
          @click="router.push('/schedule/calendar')"
        >
          <div class="stat-content">
            <div class="stat-icon">
              <n-icon size="32" :component="CalendarOutline" />
            </div>
            <div class="stat-info">
              <n-statistic label="今日日程">
                <template #default>
                  <n-number-animation :from="0" :to="stats.todaySchedule" :duration="1000" />
                </template>
              </n-statistic>
            </div>
          </div>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card
          class="stat-card stat-card--apply"
          hoverable
          @click="router.push('/approval/my-apply')"
        >
          <div class="stat-content">
            <div class="stat-icon">
              <n-icon size="32" :component="PaperPlaneOutline" />
            </div>
            <div class="stat-info">
              <n-statistic label="我的申请">
                <template #default>
                  <n-number-animation :from="0" :to="stats.myPending" :duration="1000" />
                </template>
              </n-statistic>
            </div>
          </div>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- 主要内容区域 -->
    <n-grid
      :x-gap="16"
      :y-gap="16"
      :cols="3"
      item-responsive
      responsive="screen"
      class="main-content"
    >
      <!-- 待办事项 -->
      <n-gi span="3 l:2">
        <n-card title="待办审批" :bordered="false">
          <template #header-extra>
            <n-button text type="primary" @click="router.push('/approval/pending')">
              查看更多
            </n-button>
          </template>
          <n-spin :show="pendingLoading">
            <n-empty v-if="pendingList.length === 0" description="暂无待办事项" />
            <n-list v-else hoverable clickable>
              <n-list-item
                v-for="item in pendingList"
                :key="item.id"
                @click="goToApprovalDetail(item.id)"
              >
                <template #prefix>
                  <n-tag :type="getStatusType(item.status)" size="small">
                    {{ ApprovalTypeLabels[item.type] }}
                  </n-tag>
                </template>
                <n-thing :title="item.title" :description="item.createTime">
                  <template #header-extra>
                    <n-text depth="3">
                      {{ item.applicantName }}
                    </n-text>
                  </template>
                </n-thing>
              </n-list-item>
            </n-list>
          </n-spin>
        </n-card>
      </n-gi>

      <!-- 快捷入口 -->
      <n-gi span="3 l:1">
        <n-card title="快捷入口" :bordered="false">
          <n-grid :x-gap="12" :y-gap="12" :cols="2">
            <n-gi v-for="shortcut in shortcuts" :key="shortcut.path">
              <div class="shortcut-item" @click="router.push(shortcut.path)">
                <n-icon size="24" :component="shortcut.icon" />
                <span class="shortcut-label">{{ shortcut.label }}</span>
              </div>
            </n-gi>
          </n-grid>
        </n-card>

        <!-- 本周日程 -->
        <n-card title="本周日程" :bordered="false" class="week-schedule-card">
          <template #header-extra>
            <n-button text type="primary" @click="router.push('/schedule/calendar')">
              查看更多
            </n-button>
          </template>
          <n-spin :show="scheduleLoading">
            <n-empty v-if="weekSchedules.length === 0" description="暂无日程安排" />
            <n-timeline v-else>
              <n-timeline-item
                v-for="schedule in weekSchedules.slice(0, 5)"
                :key="schedule.id"
                :type="getPriorityType(schedule.priority)"
                :title="schedule.title"
                :time="formatScheduleTime(schedule)"
                :line-type="schedule.priority === 'high' ? 'dashed' : 'default'"
              >
                <template #default>
                  <n-text depth="3">
                    {{ schedule.startTime }} - {{ schedule.endTime }}
                    <span v-if="schedule.location"> | {{ schedule.location }}</span>
                  </n-text>
                </template>
              </n-timeline-item>
            </n-timeline>
          </n-spin>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- 图表区域 -->
    <n-grid
      :x-gap="16"
      :y-gap="16"
      :cols="2"
      item-responsive
      responsive="screen"
      class="chart-section"
    >
      <n-gi span="2 l:1">
        <n-card title="审批概览" :bordered="false">
          <div ref="approvalChartRef" class="chart-container"></div>
        </n-card>
      </n-gi>
      <n-gi span="2 l:1">
        <n-card title="审批状态分布" :bordered="false">
          <div ref="typeChartRef" class="chart-container"></div>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, computed } from 'vue'
  import { useRouter } from 'vue-router'
  import {
    NCard,
    NGrid,
    NGi,
    NStatistic,
    NNumberAnimation,
    NIcon,
    NList,
    NListItem,
    NThing,
    NTag,
    NText,
    NButton,
    NEmpty,
    NSpin,
    NTimeline,
    NTimelineItem
  } from 'naive-ui'
  import {
    DocumentTextOutline,
    MegaphoneOutline,
    CalendarOutline,
    PaperPlaneOutline,
    SunnyOutline,
    AddCircleOutline
  } from '@vicons/ionicons5'
  import * as echarts from 'echarts'
  import { useUserStore } from '@/stores/user'
  import { useApprovalStore } from '@/stores/approval'
  import { useNoticeStore } from '@/stores/notice'
  import { useScheduleStore } from '@/stores/schedule'
  import { ApprovalTypeLabels, type ApprovalType } from '@/types/approval'
  import { formatDate } from '@/utils/date'
  import { request } from '@/utils/request'

  const router = useRouter()
  const userStore = useUserStore()
  const approvalStore = useApprovalStore()
  const noticeStore = useNoticeStore()
  const scheduleStore = useScheduleStore()

  // 日期显示
  const currentDate = computed(() => {
    const now = new Date()
    const weekDays = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
    return `${formatDate(now, 'YYYY年MM月DD日')} ${weekDays[now.getDay()]}`
  })

  const stats = ref({
    todoApproval: 0,
    unreadNotice: 0,
    todaySchedule: 0,
    myPending: 0
  })

  const approvalStatsData = ref<{
    myPending: number
    myApproved: number
    myRejected: number
    pendingApproval: number
  }>({ myPending: 0, myApproved: 0, myRejected: 0, pendingApproval: 0 })

  // 待审批列表
  const pendingList = ref<ReturnType<typeof extractApprovalInfo>[]>([])
  const pendingLoading = ref(false)

  // 本周日程
  const weekSchedules = computed(() => scheduleStore.weekSchedules)
  const scheduleLoading = ref(false)

  // 图表引用
  const approvalChartRef = ref<HTMLElement | null>(null)
  const typeChartRef = ref<HTMLElement | null>(null)
  let approvalChart: echarts.ECharts | null = null
  let typeChart: echarts.ECharts | null = null

  // 问候语
  function getGreeting(): string {
    const hour = new Date().getHours()
    if (hour < 6) return '凌晨好'
    if (hour < 9) return '早上好'
    if (hour < 12) return '上午好'
    if (hour < 14) return '中午好'
    if (hour < 17) return '下午好'
    if (hour < 19) return '傍晚好'
    if (hour < 22) return '晚上好'
    return '夜深了'
  }

  // 获取状态类型
  function getStatusType(
    status: string
  ): 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error' {
    const map: Record<string, 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error'> = {
      pending: 'warning',
      approved: 'success',
      rejected: 'error',
      withdrawn: 'default',
      transferred: 'info'
    }
    return map[status] || 'default'
  }

  // 获取优先级类型
  function getPriorityType(priority: string): 'default' | 'info' | 'success' | 'warning' | 'error' {
    const map: Record<string, 'default' | 'info' | 'success' | 'warning' | 'error'> = {
      low: 'info',
      medium: 'warning',
      high: 'error'
    }
    return map[priority] || 'default'
  }

  // 格式化日程时间
  function formatScheduleTime(schedule: {
    startDate: string
    startTime: string
    endTime: string
  }): string {
    return `${formatDate(schedule.startDate, 'MM/DD')} ${schedule.startTime} - ${schedule.endTime}`
  }

  // 跳转到审批详情
  function goToApprovalDetail(id: number): void {
    router.push(`/approval/detail/${id}`)
  }

  // 提取审批信息
  function extractApprovalInfo(approval: {
    id: number
    type: ApprovalType
    title: string
    applicantName: string
    status: string
    createTime: string
  }) {
    return {
      id: approval.id,
      type: approval.type,
      title: approval.title,
      applicantName: approval.applicantName,
      status: approval.status,
      createTime: approval.createTime
    }
  }

  // 快捷入口
  const shortcuts = [
    { label: '发起申请', path: '/approval/apply', icon: AddCircleOutline },
    { label: '公告列表', path: '/notice/list', icon: MegaphoneOutline },
    { label: '日程管理', path: '/schedule/calendar', icon: CalendarOutline },
    { label: '我的申请', path: '/approval/my-apply', icon: PaperPlaneOutline }
  ]

  async function fetchStats(): Promise<void> {
    try {
      const approvalStats = await approvalStore.getPendingApprovals({ page: 1, pageSize: 5 })
      stats.value.todoApproval = approvalStats.total

      const unreadCount = await noticeStore.getUnreadCount()
      stats.value.unreadNotice = unreadCount

      const schedules = await scheduleStore.getWeekSchedules()

      const today = formatDate(new Date(), 'YYYY-MM-DD')
      stats.value.todaySchedule = schedules.filter((s) => s.startDate === today).length

      const myApprovals = await approvalStore.getMyApprovals({ page: 1, pageSize: 100 })
      stats.value.myPending = myApprovals.list.filter((a) => a.status === 'pending').length

      pendingList.value = approvalStore.pendingApprovals.map(extractApprovalInfo)

      approvalStatsData.value = await request.get<{
        myPending: number
        myApproved: number
        myRejected: number
        pendingApproval: number
      }>('/approvals/stats')
    } catch {
      // 静默处理
    }
  }

  async function fetchPendingList(): Promise<void> {
    pendingLoading.value = true
    try {
      const result = await approvalStore.getPendingApprovals({ page: 1, pageSize: 5 })
      pendingList.value = result.list.map(extractApprovalInfo)
    } catch {
      // 静默处理
    } finally {
      pendingLoading.value = false
    }
  }

  function initApprovalChart(): void {
    if (!approvalChartRef.value) return

    approvalChart = echarts.init(approvalChartRef.value)

    const s = approvalStatsData.value
    const categories = ['已通过', '已驳回', '待审批']
    const values = [s.myApproved, s.myRejected, s.myPending]

    const option: echarts.EChartsOption = {
      tooltip: {
        trigger: 'axis'
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: categories
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '审批数量',
          type: 'bar',
          barWidth: '40%',
          data: [
            {
              value: values[0],
              itemStyle: { color: '#52C41A' }
            },
            {
              value: values[1],
              itemStyle: { color: '#FF4D4F' }
            },
            {
              value: values[2],
              itemStyle: { color: '#FAAD14' }
            }
          ]
        }
      ]
    }

    approvalChart.setOption(option)
  }

  function initTypeChart(): void {
    if (!typeChartRef.value) return

    typeChart = echarts.init(typeChartRef.value)

    const s = approvalStatsData.value

    const option: echarts.EChartsOption = {
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        left: 'left'
      },
      series: [
        {
          name: '审批状态',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: false,
            position: 'center'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 14,
              fontWeight: 'bold'
            }
          },
          labelLine: {
            show: false
          },
          data: [
            { value: s.myApproved, name: '已通过', itemStyle: { color: '#52C41A' } },
            { value: s.myRejected, name: '已驳回', itemStyle: { color: '#FF4D4F' } },
            { value: s.myPending, name: '待审批', itemStyle: { color: '#FAAD14' } }
          ]
        }
      ]
    }

    typeChart.setOption(option)
  }

  // 处理窗口大小变化
  function handleResize(): void {
    approvalChart?.resize()
    typeChart?.resize()
  }

  onMounted(async () => {
    await Promise.all([fetchStats(), fetchPendingList()])
    initApprovalChart()
    initTypeChart()
    window.addEventListener('resize', handleResize)
  })
</script>

<style lang="scss" scoped>
  .dashboard {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .welcome-card {
    background: $primary-color !important;
    border: none !important;
    box-shadow: $box-shadow !important;

    :deep(.n-card__content) {
      padding: 24px;
    }

    .welcome-content {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .welcome-text {
      color: #fff;
    }

    .welcome-title {
      font-size: 20px;
      font-weight: 600;
      margin: 0 0 8px 0;
      color: #fff;
    }

    .welcome-date {
      font-size: 14px;
      margin: 0;
      opacity: 0.9;
      color: #fff;
    }

    .welcome-weather {
      color: #fff;
      opacity: 0.8;
    }
  }

  .stat-card {
    cursor: pointer;
    border: 1px solid $border-color !important;
    box-shadow: $box-shadow-light !important;
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-4px);
      box-shadow: $box-shadow !important;
    }

    :deep(.n-card__content) {
      padding: 20px;
    }

    .stat-content {
      display: flex;
      align-items: center;
      gap: 16px;
    }

    .stat-icon {
      width: 56px;
      height: 56px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    &--pending .stat-icon {
      background: rgba($warning-color, 0.15);
      color: $warning-color;
    }

    &--notice .stat-icon {
      background: rgba($success-color, 0.15);
      color: $success-color;
    }

    &--schedule .stat-icon {
      background: rgba($primary-color, 0.15);
      color: $primary-color;
    }

    &--apply .stat-icon {
      background: rgba($error-color, 0.15);
      color: $error-color;
    }
  }

  .shortcut-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 20px 16px;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    background: $bg-color-1;
    border: 1px solid $border-color;

    &:hover {
      background: $bg-color-2;
      transform: translateY(-2px);
      box-shadow: $box-shadow-light;
    }

    .shortcut-label {
      margin-top: 8px;
      font-size: 12px;
      color: $text-color-3;
    }
  }

  .week-schedule-card {
    margin-top: 16px;
    border: 1px solid $border-color !important;
  }

  .chart-section {
    margin-top: 16px;

    :deep(.n-card) {
      border: 1px solid $border-color !important;
    }
  }

  .chart-container {
    width: 100%;
    height: 300px;
  }

  .main-content {
    margin-top: 16px;

    :deep(.n-card) {
      border: 1px solid $border-color !important;
    }
  }

  @media (max-width: 768px) {
    .welcome-title {
      font-size: 16px !important;
    }

    .stat-card .stat-icon {
      width: 48px;
      height: 48px;
    }

    .chart-container {
      height: 250px;
    }
  }
</style>
