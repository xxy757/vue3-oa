<template>
  <div class="apply-page">
    <n-card title="发起申请">
      <n-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-placement="left"
        label-width="100"
        require-mark-placement="right-hanging"
      >
        <!-- 申请类型 -->
        <n-form-item label="申请类型" path="type">
          <n-select
            v-model:value="formData.type"
            :options="typeOptions"
            placeholder="请选择申请类型"
            @update:value="handleTypeChange"
          />
        </n-form-item>

        <!-- 申请标题 -->
        <n-form-item label="申请标题" path="title">
          <n-input
            v-model:value="formData.title"
            placeholder="请输入申请标题"
            maxlength="50"
            show-count
          />
        </n-form-item>

        <!-- 请假申请表单 -->
        <template v-if="formData.type === 'leave'">
          <n-form-item label="请假类型" path="leaveType">
            <n-select
              v-model:value="formData.leaveType"
              :options="leaveTypeOptions"
              placeholder="请选择请假类型"
            />
          </n-form-item>
          <n-form-item label="开始日期" path="startDate">
            <n-date-picker
              v-model:value="dateRange.start"
              type="date"
              placeholder="请选择开始日期"
              style="width: 100%"
            />
          </n-form-item>
          <n-form-item label="结束日期" path="endDate">
            <n-date-picker
              v-model:value="dateRange.end"
              type="date"
              placeholder="请选择结束日期"
              style="width: 100%"
            />
          </n-form-item>
          <n-form-item label="请假天数">
            <n-input-number
              v-model:value="formData.days"
              :min="0.5"
              :max="30"
              :step="0.5"
              disabled
            />
          </n-form-item>
          <n-form-item label="请假原因" path="reason">
            <n-input
              v-model:value="formData.reason"
              type="textarea"
              placeholder="请输入请假原因"
              :rows="4"
              maxlength="200"
              show-count
            />
          </n-form-item>
        </template>

        <!-- 报销申请表单 -->
        <template v-if="formData.type === 'expense'">
          <n-form-item label="报销类型" path="expenseType">
            <n-select
              v-model:value="formData.expenseType"
              :options="expenseTypeOptions"
              placeholder="请选择报销类型"
            />
          </n-form-item>
          <n-form-item label="报销金额" path="amount">
            <n-input-number
              v-model:value="formData.amount"
              :min="0"
              :precision="2"
              placeholder="请输入报销金额"
              style="width: 100%"
            >
              <template #prefix>¥</template>
            </n-input-number>
          </n-form-item>
          <n-form-item label="报销说明" path="description">
            <n-input
              v-model:value="formData.description"
              type="textarea"
              placeholder="请输入报销说明"
              :rows="4"
              maxlength="200"
              show-count
            />
          </n-form-item>
          <n-form-item label="附件">
            <n-upload :file-list="fileList" :max="5" @update:file-list="handleFileChange">
              <n-button>上传附件</n-button>
            </n-upload>
          </n-form-item>
        </template>

        <!-- 加班申请表单 -->
        <template v-if="formData.type === 'overtime'">
          <n-form-item label="加班日期" path="overtimeDate">
            <n-date-picker
              v-model:value="dateRange.overtime"
              type="date"
              placeholder="请选择加班日期"
              style="width: 100%"
            />
          </n-form-item>
          <n-form-item label="开始时间" path="startTime">
            <n-time-picker
              v-model:value="timeRange.start"
              format="HH:mm"
              placeholder="请选择开始时间"
              style="width: 100%"
            />
          </n-form-item>
          <n-form-item label="结束时间" path="endTime">
            <n-time-picker
              v-model:value="timeRange.end"
              format="HH:mm"
              placeholder="请选择结束时间"
              style="width: 100%"
            />
          </n-form-item>
          <n-form-item label="加班时长">
            <n-input-number
              v-model:value="formData.hours"
              :min="0.5"
              :max="24"
              :step="0.5"
              disabled
            />
            <span style="margin-left: 8px">小时</span>
          </n-form-item>
          <n-form-item label="加班原因" path="reason">
            <n-input
              v-model:value="formData.reason"
              type="textarea"
              placeholder="请输入加班原因"
              :rows="4"
              maxlength="200"
              show-count
            />
          </n-form-item>
        </template>

        <!-- 出差申请表单 -->
        <template v-if="formData.type === 'travel'">
          <n-form-item label="目的地" path="destination">
            <n-input
              v-model:value="formData.destination"
              placeholder="请输入目的地"
              maxlength="50"
            />
          </n-form-item>
          <n-form-item label="开始日期" path="startDate">
            <n-date-picker
              v-model:value="dateRange.start"
              type="date"
              placeholder="请选择开始日期"
              style="width: 100%"
            />
          </n-form-item>
          <n-form-item label="结束日期" path="endDate">
            <n-date-picker
              v-model:value="dateRange.end"
              type="date"
              placeholder="请选择结束日期"
              style="width: 100%"
            />
          </n-form-item>
          <n-form-item label="出差天数">
            <n-input-number v-model:value="formData.days" :min="1" :max="30" disabled />
          </n-form-item>
          <n-form-item label="出差原因" path="reason">
            <n-input
              v-model:value="formData.reason"
              type="textarea"
              placeholder="请输入出差原因"
              :rows="4"
              maxlength="200"
              show-count
            />
          </n-form-item>
          <n-form-item label="预算金额" path="budget">
            <n-input-number
              v-model:value="formData.budget"
              :min="0"
              :precision="2"
              placeholder="请输入预算金额"
              style="width: 100%"
            >
              <template #prefix>¥</template>
            </n-input-number>
          </n-form-item>
        </template>

        <!-- 通用审批表单 -->
        <template v-if="formData.type === 'general'">
          <n-form-item label="申请内容" path="content">
            <n-input
              v-model:value="formData.content"
              type="textarea"
              placeholder="请输入申请内容"
              :rows="6"
              maxlength="500"
              show-count
            />
          </n-form-item>
          <n-form-item label="附件">
            <n-upload :file-list="fileList" :max="5" @update:file-list="handleFileChange">
              <n-button>上传附件</n-button>
            </n-upload>
          </n-form-item>
        </template>

        <!-- 操作按钮 -->
        <n-form-item>
          <n-space>
            <n-button type="primary" :loading="loading" @click="handleSubmit">提交申请</n-button>
            <n-button @click="handleReset">重置</n-button>
            <n-button @click="router.back()">返回</n-button>
          </n-space>
        </n-form-item>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import {
    NCard,
    NForm,
    NFormItem,
    NSelect,
    NInput,
    NInputNumber,
    NDatePicker,
    NTimePicker,
    NUpload,
    NButton,
    NSpace,
    useMessage
  } from 'naive-ui'
  import type { FormInst, FormRules, UploadFileInfo } from 'naive-ui'
  import { useApprovalStore } from '@/stores/approval'
  import type { ApprovalType, CreateApprovalParams } from '@/types/approval'

  const router = useRouter()
  const message = useMessage()
  const approvalStore = useApprovalStore()

  const formRef = ref<FormInst | null>(null)
  const loading = ref(false)
  const fileList = ref<UploadFileInfo[]>([])

  // 申请类型选项
  const typeOptions = [
    { label: '请假申请', value: 'leave' },
    { label: '报销申请', value: 'expense' },
    { label: '加班申请', value: 'overtime' },
    { label: '出差申请', value: 'travel' },
    { label: '通用审批', value: 'general' }
  ]

  // 请假类型选项
  const leaveTypeOptions = [
    { label: '年假', value: 'annual' },
    { label: '病假', value: 'sick' },
    { label: '事假', value: 'personal' },
    { label: '产假', value: 'maternity' },
    { label: '婚假', value: 'marriage' },
    { label: '丧假', value: 'bereavement' }
  ]

  // 报销类型选项
  const expenseTypeOptions = [
    { label: '差旅费', value: 'travel' },
    { label: '办公费', value: 'office' },
    { label: '招待费', value: 'entertainment' },
    { label: '其他', value: 'other' }
  ]

  // 日期范围
  const dateRange = reactive({
    start: null as number | null,
    end: null as number | null,
    overtime: null as number | null
  })

  // 时间范围
  const timeRange = reactive({
    start: null as number | null,
    end: null as number | null
  })

  // 表单数据
  const formData = reactive<{
    type: ApprovalType | null
    title: string
    leaveType: string | null
    startDate: string
    endDate: string
    days: number
    reason: string
    expenseType: string | null
    amount: number | null
    description: string
    attachments: string[]
    overtimeDate: string
    startTime: string
    endTime: string
    hours: number
    destination: string
    budget: number | null
    content: string
  }>({
    type: null,
    title: '',
    leaveType: null,
    startDate: '',
    endDate: '',
    days: 0,
    reason: '',
    expenseType: null,
    amount: null,
    description: '',
    attachments: [],
    overtimeDate: '',
    startTime: '',
    endTime: '',
    hours: 0,
    destination: '',
    budget: null,
    content: ''
  })

  // 表单验证规则
  const rules: FormRules = {
    type: { required: true, message: '请选择申请类型', trigger: 'change' },
    title: { required: true, message: '请输入申请标题', trigger: 'blur' }
  }

  // 计算请假/出差天数
  watch(
    () => [dateRange.start, dateRange.end] as const,
    ([start, end]) => {
      if (start && end) {
        const diff = Math.ceil((end - start) / (1000 * 60 * 60 * 24)) + 1
        formData.days = Math.max(1, diff)
        formData.startDate = formatDate(start)
        formData.endDate = formatDate(end)
      }
    }
  )

  // 计算加班时长
  watch(
    () => [timeRange.start, timeRange.end] as const,
    ([start, end]) => {
      if (start && end) {
        const diff = (end - start) / (1000 * 60 * 60)
        formData.hours = Math.max(0.5, Math.round(diff * 2) / 2)
        formData.startTime = formatTime(start)
        formData.endTime = formatTime(end)
      }
    }
  )

  // 监听加班日期
  watch(
    () => dateRange.overtime,
    (val: number | null) => {
      if (val) {
        formData.overtimeDate = formatDate(val)
      }
    }
  )

  // 格式化日期
  function formatDate(timestamp: number): string {
    const date = new Date(timestamp)
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
  }

  // 格式化时间
  function formatTime(timestamp: number): string {
    const date = new Date(timestamp)
    return `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
  }

  // 类型切换
  function handleTypeChange() {
    // 重置表单数据
    formData.title = ''
    formData.leaveType = null
    formData.startDate = ''
    formData.endDate = ''
    formData.days = 0
    formData.reason = ''
    formData.expenseType = null
    formData.amount = null
    formData.description = ''
    formData.attachments = []
    formData.overtimeDate = ''
    formData.startTime = ''
    formData.endTime = ''
    formData.hours = 0
    formData.destination = ''
    formData.budget = null
    formData.content = ''
    dateRange.start = null
    dateRange.end = null
    dateRange.overtime = null
    timeRange.start = null
    timeRange.end = null
    fileList.value = []
  }

  // 文件变化
  function handleFileChange(files: UploadFileInfo[]) {
    fileList.value = files
    formData.attachments = files.map((f) => f.url || f.name).filter(Boolean)
  }

  // 提交申请
  async function handleSubmit() {
    if (!formRef.value) return

    try {
      await formRef.value.validate()
    } catch {
      message.error('请完善表单信息')
      return
    }

    if (!formData.type) {
      message.error('请选择申请类型')
      return
    }

    loading.value = true

    try {
      const params: CreateApprovalParams = {
        type: formData.type,
        title: formData.title,
        ...getExtraParams()
      }

      await approvalStore.createApproval(params)
      message.success('申请提交成功')
      router.push('/approval/my-apply')
    } catch (error) {
      message.error((error as Error).message || '提交失败')
    } finally {
      loading.value = false
    }
  }

  // 获取额外参数
  function getExtraParams(): Record<string, unknown> {
    switch (formData.type) {
      case 'leave':
        return {
          leaveType: formData.leaveType,
          startDate: formData.startDate,
          endDate: formData.endDate,
          days: formData.days,
          reason: formData.reason
        }
      case 'expense':
        return {
          expenseType: formData.expenseType,
          amount: formData.amount,
          description: formData.description,
          attachments: formData.attachments
        }
      case 'overtime':
        return {
          overtimeDate: formData.overtimeDate,
          startTime: formData.startTime,
          endTime: formData.endTime,
          hours: formData.hours,
          reason: formData.reason
        }
      case 'travel':
        return {
          destination: formData.destination,
          startDate: formData.startDate,
          endDate: formData.endDate,
          days: formData.days,
          reason: formData.reason,
          budget: formData.budget
        }
      case 'general':
        return {
          content: formData.content,
          attachments: formData.attachments
        }
      default:
        return {}
    }
  }

  // 重置表单
  function handleReset() {
    formRef.value?.restoreValidation()
    handleTypeChange()
  }
</script>

<style lang="scss" scoped>
  .apply-page {
    max-width: 800px;
    margin: 0 auto;
  }
</style>
