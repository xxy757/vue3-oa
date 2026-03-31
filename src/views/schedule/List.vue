<template>
  <div class="schedule-list">
    <n-card title="日程列表">
      <!-- 搜索和筛选区域 -->
      <div class="filter-bar">
        <n-space>
          <n-date-picker
            v-model:value="dateRange"
            type="daterange"
            clearable
            @update:value="handleDateRangeChange"
          />
          <n-select
            v-model:value="filterPriority"
            :options="priorityOptions"
            placeholder="优先级"
            clearable
            style="width: 120px"
            @update:value="handleFilterChange"
          />
          <n-input
            v-model:value="keyword"
            placeholder="搜索日程标题"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
            @clear="handleSearch"
          />
          <n-button type="primary" @click="handleSearch"> 搜索 </n-button>
          <n-button type="primary" @click="handleAddSchedule">
            <template #icon>
              <n-icon><AddOutline /></n-icon>
            </template>
            新增日程
          </n-button>
        </n-space>
      </div>

      <!-- 日程列表 -->
      <n-data-table
        :columns="columns"
        :data="filteredSchedules"
        :loading="loading"
        :row-key="(row: Schedule) => row.id"
      />

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <n-pagination
          v-model:page="pagination.page"
          :page-count="pagination.totalPages"
          :page-size="pagination.pageSize"
          show-size-picker
          :page-sizes="[10, 20, 50]"
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </n-card>

    <!-- 日程表单弹窗 -->
    <n-modal
      v-model:show="showModal"
      preset="dialog"
      :title="isEdit ? '编辑日程' : '新增日程'"
      :positive-text="isEdit ? '保存' : '创建'"
      :negative-text="'取消'"
      :loading="submitting"
      style="width: 600px"
      @positive-click="handleSubmit"
      @negative-click="handleCancel"
    >
      <n-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-placement="left"
        label-width="80"
      >
        <n-form-item label="标题" path="title">
          <n-input v-model:value="formData.title" placeholder="请输入日程标题" />
        </n-form-item>

        <n-form-item label="描述" path="description">
          <n-input
            v-model:value="formData.description"
            type="textarea"
            placeholder="请输入日程描述"
            :rows="3"
          />
        </n-form-item>

        <n-form-item label="全天" path="isAllDay">
          <n-switch v-model:value="formData.isAllDay" @update:value="handleAllDayChange" />
        </n-form-item>

        <n-space>
          <n-form-item label="开始日期" path="startDate">
            <n-date-picker
              v-model:value="startDateValue"
              type="date"
              clearable
              @update:value="handleStartDateChange"
            />
          </n-form-item>

          <n-form-item v-if="!formData.isAllDay" label="开始时间" path="startTime">
            <n-time-picker
              v-model:value="startTimeValue"
              format="HH:mm"
              clearable
              @update:value="handleStartTimeChange"
            />
          </n-form-item>
        </n-space>

        <n-space>
          <n-form-item label="结束日期" path="endDate">
            <n-date-picker
              v-model:value="endDateValue"
              type="date"
              clearable
              @update:value="handleEndDateChange"
            />
          </n-form-item>

          <n-form-item v-if="!formData.isAllDay" label="结束时间" path="endTime">
            <n-time-picker
              v-model:value="endTimeValue"
              format="HH:mm"
              clearable
              @update:value="handleEndTimeChange"
            />
          </n-form-item>
        </n-space>

        <n-form-item label="优先级" path="priority">
          <n-select
            v-model:value="formData.priority"
            :options="prioritySelectOptions"
            placeholder="选择优先级"
          />
        </n-form-item>

        <n-form-item label="提醒" path="remind">
          <n-select
            v-model:value="formData.remind"
            :options="remindOptions"
            placeholder="选择提醒时间"
          />
        </n-form-item>

        <n-form-item label="地点" path="location">
          <n-input v-model:value="formData.location" placeholder="请输入地点" />
        </n-form-item>

        <n-form-item label="颜色" path="color">
          <div class="color-picker">
            <div
              v-for="color in colorOptions"
              :key="color"
              class="color-item"
              :class="{ active: formData.color === color }"
              :style="{ backgroundColor: color }"
              @click="formData.color = color"
            />
          </div>
        </n-form-item>
      </n-form>
    </n-modal>

    <!-- 日程详情弹窗 -->
    <n-modal v-model:show="showDetail" preset="card" title="日程详情" style="width: 500px">
      <template v-if="currentSchedule">
        <n-descriptions :column="1" label-placement="left">
          <n-descriptions-item label="标题">
            {{ currentSchedule.title }}
          </n-descriptions-item>
          <n-descriptions-item label="优先级">
            <n-tag
              :color="{ color: getPriorityColor(currentSchedule.priority), textColor: '#fff' }"
            >
              {{ getPriorityLabel(currentSchedule.priority) }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="时间">
            <template v-if="currentSchedule.isAllDay">
              {{ currentSchedule.startDate }} 全天
            </template>
            <template v-else>
              {{ currentSchedule.startDate }} {{ currentSchedule.startTime }} -
              {{ currentSchedule.endTime }}
            </template>
          </n-descriptions-item>
          <n-descriptions-item v-if="currentSchedule.location" label="地点">
            {{ currentSchedule.location }}
          </n-descriptions-item>
          <n-descriptions-item label="提醒">
            {{ getRemindLabel(currentSchedule.remind) }}
          </n-descriptions-item>
          <n-descriptions-item v-if="currentSchedule.description" label="描述">
            {{ currentSchedule.description }}
          </n-descriptions-item>
          <n-descriptions-item label="创建人">
            {{ currentSchedule.creatorName }}
          </n-descriptions-item>
          <n-descriptions-item label="创建时间">
            {{ currentSchedule.createTime }}
          </n-descriptions-item>
        </n-descriptions>
      </template>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showDetail = false">关闭</n-button>
          <n-button type="primary" @click="handleEditSchedule">编辑</n-button>
          <n-button type="error" @click="handleDeleteSchedule">删除</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted, h } from 'vue'
  import {
    NCard,
    NSpace,
    NDatePicker,
    NSelect,
    NInput,
    NButton,
    NDataTable,
    NTag,
    NPagination,
    NIcon,
    NModal,
    NForm,
    NFormItem,
    NSwitch,
    NTimePicker,
    NDescriptions,
    NDescriptionsItem,
    useMessage,
    useDialog
  } from 'naive-ui'
  import { AddOutline, CreateOutline, TrashOutline, EyeOutline } from '@vicons/ionicons5'
  import type { DataTableColumns } from 'naive-ui'
  import { useScheduleStore } from '@/stores/schedule'
  import type { Schedule, SchedulePriority, CreateScheduleParams } from '@/types/schedule'
  import { SchedulePriorityLabels, ScheduleRemindLabels, ScheduleColors } from '@/types/schedule'
  import { formatDate } from '@/utils/date'

  const message = useMessage()
  const dialog = useDialog()
  const scheduleStore = useScheduleStore()

  const loading = ref(false)
  const submitting = ref(false)
  const showModal = ref(false)
  const showDetail = ref(false)
  const isEdit = ref(false)
  const currentSchedule = ref<Schedule | null>(null)
  const formRef = ref()

  const dateRange = ref<[number, number] | null>(null)
  const filterPriority = ref<string | null>(null)
  const keyword = ref('')

  const pagination = reactive({
    page: 1,
    pageSize: 10,
    total: 0,
    totalPages: 0
  })

  import { reactive } from 'vue'

  const formData = ref<CreateScheduleParams>({
    title: '',
    description: '',
    startDate: '',
    endDate: '',
    startTime: '',
    endTime: '',
    isAllDay: false,
    priority: 'medium',
    remind: 'none',
    location: '',
    color: '#2080f0'
  })

  const formRules = {
    title: { required: true, message: '请输入日程标题', trigger: 'blur' },
    startDate: { required: true, message: '请选择开始日期', trigger: 'change' }
  }

  // 优先级筛选选项
  const priorityOptions = [
    { label: '全部', value: '' },
    { label: '高', value: 'high' },
    { label: '中', value: 'medium' },
    { label: '低', value: 'low' }
  ]

  // 优先级表单选项
  const prioritySelectOptions = Object.entries(SchedulePriorityLabels).map(([value, label]) => ({
    label,
    value
  }))

  // 提醒选项
  const remindOptions = Object.entries(ScheduleRemindLabels).map(([value, label]) => ({
    label,
    value
  }))

  // 颜色选项
  const colorOptions = ScheduleColors

  // 筛选后的日程
  const filteredSchedules = computed(() => {
    let result = scheduleStore.schedules

    if (filterPriority.value) {
      result = result.filter((s) => s.priority === filterPriority.value)
    }

    if (keyword.value) {
      result = result.filter((s) => s.title.toLowerCase().includes(keyword.value.toLowerCase()))
    }

    // 计算分页
    pagination.total = result.length
    pagination.totalPages = Math.ceil(result.length / pagination.pageSize)

    // 分页截取
    const start = (pagination.page - 1) * pagination.pageSize
    const end = start + pagination.pageSize

    return result.slice(start, end)
  })

  // 开始日期值
  const startDateValue = computed({
    get: () => (formData.value.startDate ? new Date(formData.value.startDate).getTime() : null),
    set: () => {}
  })

  // 开始时间值
  const startTimeValue = computed({
    get: () => {
      if (formData.value.startTime) {
        const [hours, minutes] = formData.value.startTime.split(':').map(Number)
        return hours * 60 * 60 * 1000 + minutes * 60 * 1000
      }
      return null
    },
    set: () => {}
  })

  // 结束日期值
  const endDateValue = computed({
    get: () => (formData.value.endDate ? new Date(formData.value.endDate).getTime() : null),
    set: () => {}
  })

  // 结束时间值
  const endTimeValue = computed({
    get: () => {
      if (formData.value.endTime) {
        const [hours, minutes] = formData.value.endTime.split(':').map(Number)
        return hours * 60 * 60 * 1000 + minutes * 60 * 1000
      }
      return null
    },
    set: () => {}
  })

  // 表格列配置
  const columns: DataTableColumns<Schedule> = [
    {
      title: '标题',
      key: 'title',
      ellipsis: { tooltip: true },
      render(row) {
        return h('div', { style: { display: 'flex', alignItems: 'center', gap: '8px' } }, [
          h('div', {
            style: {
              width: '4px',
              height: '16px',
              borderRadius: '2px',
              backgroundColor: row.color || '#2080f0'
            }
          }),
          h('span', {}, row.title)
        ])
      }
    },
    {
      title: '优先级',
      key: 'priority',
      width: 80,
      render(row) {
        return h(
          NTag,
          {
            color: { color: getPriorityColor(row.priority), textColor: '#fff' },
            size: 'small'
          },
          { default: () => getPriorityLabel(row.priority) }
        )
      }
    },
    {
      title: '开始时间',
      key: 'startTime',
      width: 160,
      render(row) {
        if (row.isAllDay) {
          return `${row.startDate} 全天`
        }
        return `${row.startDate} ${row.startTime}`
      }
    },
    {
      title: '结束时间',
      key: 'endTime',
      width: 160,
      render(row) {
        if (row.isAllDay) {
          return row.endDate
        }
        return `${row.endDate} ${row.endTime}`
      }
    },
    {
      title: '地点',
      key: 'location',
      width: 120,
      ellipsis: { tooltip: true },
      render(row) {
        return row.location || '-'
      }
    },
    {
      title: '提醒',
      key: 'remind',
      width: 100,
      render(row) {
        return getRemindLabel(row.remind)
      }
    },
    {
      title: '操作',
      key: 'actions',
      width: 200,
      render(row) {
        return h(
          NSpace,
          {},
          {
            default: () => [
              h(
                NButton,
                {
                  text: true,
                  type: 'primary',
                  onClick: () => handleViewSchedule(row)
                },
                {
                  default: () => '查看',
                  icon: () => h(NIcon, {}, { default: () => h(EyeOutline) })
                }
              ),
              h(
                NButton,
                {
                  text: true,
                  type: 'primary',
                  onClick: () => handleEditRow(row)
                },
                {
                  default: () => '编辑',
                  icon: () => h(NIcon, {}, { default: () => h(CreateOutline) })
                }
              ),
              h(
                NButton,
                {
                  text: true,
                  type: 'error',
                  onClick: () => handleDeleteRow(row)
                },
                {
                  default: () => '删除',
                  icon: () => h(NIcon, {}, { default: () => h(TrashOutline) })
                }
              )
            ]
          }
        )
      }
    }
  ]

  // 获取优先级颜色
  const getPriorityColor = (priority: SchedulePriority) => {
    return scheduleStore.getPriorityColor(priority)
  }

  // 获取优先级标签
  const getPriorityLabel = (priority: SchedulePriority) => {
    return scheduleStore.getPriorityLabel(priority)
  }

  // 获取提醒标签
  const getRemindLabel = (remind: string) => {
    return scheduleStore.getRemindLabel(remind as never) || '不提醒'
  }

  // 日期范围变化
  const handleDateRangeChange = (value: [number, number] | null) => {
    dateRange.value = value
    pagination.page = 1
    fetchSchedules()
  }

  // 筛选变化
  const handleFilterChange = () => {
    pagination.page = 1
  }

  // 搜索
  const handleSearch = () => {
    pagination.page = 1
  }

  // 分页变化
  const handlePageChange = (page: number) => {
    pagination.page = page
  }

  // 每页条数变化
  const handlePageSizeChange = (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
  }

  // 开始日期变化
  const handleStartDateChange = (timestamp: number | null) => {
    if (timestamp) {
      formData.value.startDate = formatDate(timestamp, 'YYYY-MM-DD')
      if (!formData.value.endDate) {
        formData.value.endDate = formData.value.startDate
      }
    }
  }

  // 开始时间变化
  const handleStartTimeChange = (timestamp: number | null) => {
    if (timestamp !== null) {
      const hours = Math.floor(timestamp / (60 * 60 * 1000))
      const minutes = Math.floor((timestamp % (60 * 60 * 1000)) / (60 * 1000))
      formData.value.startTime = `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`
    }
  }

  // 结束日期变化
  const handleEndDateChange = (timestamp: number | null) => {
    if (timestamp) {
      formData.value.endDate = formatDate(timestamp, 'YYYY-MM-DD')
    }
  }

  // 结束时间变化
  const handleEndTimeChange = (timestamp: number | null) => {
    if (timestamp !== null) {
      const hours = Math.floor(timestamp / (60 * 60 * 1000))
      const minutes = Math.floor((timestamp % (60 * 60 * 1000)) / (60 * 1000))
      formData.value.endTime = `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`
    }
  }

  // 全天变化
  const handleAllDayChange = (value: boolean) => {
    if (value) {
      formData.value.startTime = ''
      formData.value.endTime = ''
    }
  }

  // 新增日程
  const handleAddSchedule = () => {
    isEdit.value = false
    currentSchedule.value = null
    const today = formatDate(new Date(), 'YYYY-MM-DD')
    formData.value = {
      title: '',
      description: '',
      startDate: today,
      endDate: today,
      startTime: '09:00',
      endTime: '10:00',
      isAllDay: false,
      priority: 'medium',
      remind: 'none',
      location: '',
      color: '#2080f0'
    }
    showModal.value = true
  }

  // 查看日程
  const handleViewSchedule = (schedule: Schedule) => {
    currentSchedule.value = schedule
    showDetail.value = true
  }

  // 编辑行
  const handleEditRow = (schedule: Schedule) => {
    currentSchedule.value = schedule
    isEdit.value = true
    formData.value = {
      title: schedule.title,
      description: schedule.description,
      startDate: schedule.startDate,
      endDate: schedule.endDate,
      startTime: schedule.startTime,
      endTime: schedule.endTime,
      isAllDay: schedule.isAllDay,
      priority: schedule.priority,
      remind: schedule.remind,
      location: schedule.location,
      color: schedule.color
    }
    showModal.value = true
  }

  // 编辑日程（从详情弹窗）
  const handleEditSchedule = () => {
    if (!currentSchedule.value) return
    handleEditRow(currentSchedule.value)
    showDetail.value = false
  }

  // 删除行
  const handleDeleteRow = (schedule: Schedule) => {
    dialog.warning({
      title: '确认删除',
      content: `确定要删除日程「${schedule.title}」吗？删除后将无法恢复。`,
      positiveText: '确定删除',
      negativeText: '取消',
      onPositiveClick: async () => {
        try {
          await scheduleStore.deleteSchedule(schedule.id)
          message.success('删除成功')
        } catch (error) {
          message.error('删除失败')
        }
      }
    })
  }

  // 删除日程（从详情弹窗）
  const handleDeleteSchedule = () => {
    if (!currentSchedule.value) return
    handleDeleteRow(currentSchedule.value)
    showDetail.value = false
  }

  // 提交表单
  const handleSubmit = async () => {
    try {
      await formRef.value?.validate()

      submitting.value = true

      if (isEdit.value && currentSchedule.value) {
        await scheduleStore.updateSchedule(currentSchedule.value.id, formData.value)
        message.success('更新成功')
      } else {
        await scheduleStore.createSchedule(formData.value)
        message.success('创建成功')
      }

      showModal.value = false
      await fetchSchedules()
    } catch (error) {
      if (error) {
        message.error(isEdit.value ? '更新失败' : '创建失败')
      }
    } finally {
      submitting.value = false
    }
  }

  // 取消
  const handleCancel = () => {
    showModal.value = false
  }

  // 获取日程列表
  const fetchSchedules = async () => {
    loading.value = true
    try {
      let startDate: string
      let endDate: string

      if (dateRange.value) {
        startDate = formatDate(dateRange.value[0], 'YYYY-MM-DD')
        endDate = formatDate(dateRange.value[1], 'YYYY-MM-DD')
      } else {
        const today = new Date()
        startDate = formatDate(new Date(today.getFullYear(), today.getMonth() - 1, 1), 'YYYY-MM-DD')
        endDate = formatDate(new Date(today.getFullYear(), today.getMonth() + 2, 0), 'YYYY-MM-DD')
      }

      await scheduleStore.getSchedules({ startDate, endDate })
      pagination.total = scheduleStore.schedules.length
      pagination.totalPages = Math.ceil(pagination.total / pagination.pageSize)
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    fetchSchedules()
  })
</script>

<style lang="scss" scoped>
  .schedule-list {
    .filter-bar {
      margin-bottom: 16px;
    }

    .pagination-wrapper {
      margin-top: 16px;
      display: flex;
      justify-content: flex-end;
    }

    .color-picker {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;

      .color-item {
        width: 24px;
        height: 24px;
        border-radius: 4px;
        cursor: pointer;
        border: 2px solid transparent;
        transition: all 0.2s;

        &:hover {
          transform: scale(1.1);
        }

        &.active {
          border-color: #333;
        }
      }
    }
  }
</style>
