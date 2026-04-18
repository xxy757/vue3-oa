<template>
  <div class="schedule-calendar">
    <n-card title="日程日历">
      <div class="calendar-container">
        <!-- 左侧日历 -->
        <div class="calendar-left">
          <n-calendar
            v-model:value="selectedDate"
            :is-date-disabled="isDateDisabled"
            @update:value="handleDateChange"
          >
            <template #default="{ year, month, date }">
              <div class="calendar-cell">
                <span class="date-number">{{ date }}</span>
                <div v-if="getDateScheduleCount(year, month, date) > 0" class="schedule-dots">
                  <n-badge
                    :value="getDateScheduleCount(year, month, date)"
                    :max="9"
                    color="#1677FF"
                  />
                </div>
              </div>
            </template>
          </n-calendar>
        </div>

        <!-- 右侧日程列表 -->
        <div class="calendar-right">
          <div class="schedule-header">
            <n-h3>{{ formatSelectedDate }}</n-h3>
            <n-button type="primary" @click="handleAddSchedule">
              <template #icon>
                <n-icon><AddOutline /></n-icon>
              </template>
              新增日程
            </n-button>
          </div>

          <n-divider />

          <div class="schedule-list">
            <n-spin :show="loading">
              <template v-if="daySchedules.length > 0">
                <div
                  v-for="schedule in daySchedules"
                  :key="schedule.id"
                  class="schedule-item"
                  :style="{ borderLeftColor: getPriorityColor(schedule.priority) }"
                  @click="handleViewSchedule(schedule)"
                >
                  <div class="schedule-title">
                    <n-tag
                      :color="{ color: getPriorityColor(schedule.priority), textColor: '#fff' }"
                      size="small"
                    >
                      {{ getPriorityLabel(schedule.priority) }}
                    </n-tag>
                    <span class="title-text">{{ schedule.title }}</span>
                  </div>
                  <div class="schedule-time">
                    <n-icon><TimeOutline /></n-icon>
                    <span v-if="schedule.isAllDay">全天</span>
                    <span v-else>{{ schedule.startTime }} - {{ schedule.endTime }}</span>
                  </div>
                  <div v-if="schedule.location" class="schedule-location">
                    <n-icon><LocationOutline /></n-icon>
                    <span>{{ schedule.location }}</span>
                  </div>
                </div>
              </template>
              <n-empty v-else description="当天暂无日程" />
            </n-spin>
          </div>
        </div>
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

        <n-form-item label="优先级" path="priority">
          <n-select
            v-model:value="formData.priority"
            :options="priorityOptions"
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
  import { ref, computed, onMounted } from 'vue'
  import {
    NCard,
    NCalendar,
    NButton,
    NIcon,
    NH3,
    NDivider,
    NSpin,
    NEmpty,
    NTag,
    NBadge,
    NModal,
    NForm,
    NFormItem,
    NInput,
    NDatePicker,
    NTimePicker,
    NSelect,
    NSwitch,
    NDescriptions,
    NDescriptionsItem,
    NSpace,
    useMessage,
    useDialog
  } from 'naive-ui'
  import { AddOutline, TimeOutline, LocationOutline } from '@vicons/ionicons5'
  import { useScheduleStore } from '@/stores/schedule'
  import type { Schedule, SchedulePriority, CreateScheduleParams } from '@/types/schedule'
  import { SchedulePriorityLabels, ScheduleRemindLabels, ScheduleColors } from '@/types/schedule'
  import { formatDate } from '@/utils/date'

  const message = useMessage()
  const dialog = useDialog()
  const scheduleStore = useScheduleStore()

  const loading = ref(false)
  const submitting = ref(false)
  const selectedDate = ref(Date.now())
  const showModal = ref(false)
  const showDetail = ref(false)
  const isEdit = ref(false)
  const currentSchedule = ref<Schedule | null>(null)
  const formRef = ref()

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
    color: '#1677FF'
  })

  const formRules = {
    title: { required: true, message: '请输入日程标题', trigger: 'blur' },
    startDate: { required: true, message: '请选择开始日期', trigger: 'change' }
  }

  // 优先级选项
  const priorityOptions = Object.entries(SchedulePriorityLabels).map(([value, label]) => ({
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

  // 格式化选中的日期
  const formatSelectedDate = computed(() => {
    return formatDate(selectedDate.value, 'YYYY年MM月DD日')
  })

  // 获取选中日期的日程
  const daySchedules = computed(() => {
    const dateStr = formatDate(selectedDate.value, 'YYYY-MM-DD')
    return scheduleStore.getSchedulesByDate(dateStr)
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

  // 禁用日期
  const isDateDisabled = (_timestamp: number) => {
    return false // 不禁用任何日期
  }

  // 获取日期的日程数量
  const getDateScheduleCount = (year: number, month: number, date: number) => {
    const dateStr = `${year}-${(month + 1).toString().padStart(2, '0')}-${date.toString().padStart(2, '0')}`
    return scheduleStore.getScheduleCountByDate(dateStr)
  }

  // 日期变化
  const handleDateChange = (timestamp: number) => {
    selectedDate.value = timestamp
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

  // 新增日程
  const handleAddSchedule = () => {
    isEdit.value = false
    currentSchedule.value = null
    formData.value = {
      title: '',
      description: '',
      startDate: formatDate(selectedDate.value, 'YYYY-MM-DD'),
      endDate: formatDate(selectedDate.value, 'YYYY-MM-DD'),
      startTime: '09:00',
      endTime: '10:00',
      isAllDay: false,
      priority: 'medium',
      remind: 'none',
      location: '',
      color: '#1677FF'
    }
    showModal.value = true
  }

  // 查看日程
  const handleViewSchedule = (schedule: Schedule) => {
    currentSchedule.value = schedule
    showDetail.value = true
  }

  // 编辑日程
  const handleEditSchedule = () => {
    if (!currentSchedule.value) return

    isEdit.value = true
    formData.value = {
      title: currentSchedule.value.title,
      description: currentSchedule.value.description,
      startDate: currentSchedule.value.startDate,
      endDate: currentSchedule.value.endDate,
      startTime: currentSchedule.value.startTime,
      endTime: currentSchedule.value.endTime,
      isAllDay: currentSchedule.value.isAllDay,
      priority: currentSchedule.value.priority,
      remind: currentSchedule.value.remind,
      location: currentSchedule.value.location,
      color: currentSchedule.value.color
    }
    showDetail.value = false
    showModal.value = true
  }

  // 删除日程
  const handleDeleteSchedule = () => {
    if (!currentSchedule.value) return

    dialog.warning({
      title: '确认删除',
      content: '确定要删除这个日程吗？删除后将无法恢复。',
      positiveText: '确定删除',
      negativeText: '取消',
      onPositiveClick: async () => {
        try {
          await scheduleStore.deleteSchedule(currentSchedule.value!.id)
          message.success('删除成功')
          showDetail.value = false
          currentSchedule.value = null
        } catch (error) {
          message.error('删除失败')
        }
      }
    })
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
      // 刷新数据
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
      const today = new Date()
      const startDate = new Date(today.getFullYear(), today.getMonth(), 1)
      const endDate = new Date(today.getFullYear(), today.getMonth() + 1, 0)

      await scheduleStore.getSchedules({
        startDate: formatDate(startDate, 'YYYY-MM-DD'),
        endDate: formatDate(endDate, 'YYYY-MM-DD')
      })
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    fetchSchedules()
  })
</script>

<style lang="scss" scoped>
  .schedule-calendar {
    .calendar-container {
      display: flex;
      gap: 24px;

      .calendar-left {
        flex: 1;
        min-width: 400px;

        .calendar-cell {
          display: flex;
          flex-direction: column;
          align-items: center;
          height: 100%;
          padding: 4px;

          .date-number {
            font-size: 14px;
          }

          .schedule-dots {
            margin-top: 2px;
          }
        }
      }

      .calendar-right {
        width: 400px;
        flex-shrink: 0;

        .schedule-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
        }

        .schedule-list {
          max-height: 500px;
          overflow-y: auto;

          .schedule-item {
            padding: 12px 16px;
            margin-bottom: 12px;
            background: $bg-color-2;
            border-radius: 8px;
            border-left: 4px solid $primary-color;
            cursor: pointer;
            transition: all 0.2s;

            &:hover {
              background: $bg-color-4;
            }

            .schedule-title {
              display: flex;
              align-items: center;
              gap: 8px;
              margin-bottom: 8px;

              .title-text {
                font-weight: 500;
              }
            }

            .schedule-time,
            .schedule-location {
              display: flex;
              align-items: center;
              gap: 4px;
              color: $text-color-3;
              font-size: 12px;
              margin-top: 4px;
            }
          }
        }
      }
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
          border-color: $text-color-1;
        }
      }
    }
  }
</style>
