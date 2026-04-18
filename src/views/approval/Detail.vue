<template>
  <div class="detail-page">
    <n-spin :show="loading">
      <n-space vertical size="large">
        <!-- 返回按钮 -->
        <n-button text @click="router.back()">
          <template #icon>
            <n-icon><ArrowBackOutline /></n-icon>
          </template>
          返回
        </n-button>

        <!-- 申请信息卡片 -->
        <n-card title="申请信息">
          <n-descriptions label-placement="left" :column="2">
            <n-descriptions-item label="申请标题">{{ approval?.title }}</n-descriptions-item>
            <n-descriptions-item label="申请类型">{{ approval ? ApprovalTypeLabels[approval.type] : '' }}</n-descriptions-item>
            <n-descriptions-item label="审批状态">
              <n-tag :type="statusColorMap[approval?.status || 'pending']" size="small">
                {{ approval ? ApprovalStatusLabels[approval.status] : '' }}
              </n-tag>
            </n-descriptions-item>
            <n-descriptions-item label="当前进度">{{ approval ? `${approval.currentStep}/${approval.totalStep}` : '' }}</n-descriptions-item>
            <n-descriptions-item label="申请人">{{ approval?.applicantName }}</n-descriptions-item>
            <n-descriptions-item label="所属部门">{{ approval?.applicantDept }}</n-descriptions-item>
            <n-descriptions-item label="申请时间">{{ approval?.createTime }}</n-descriptions-item>
            <n-descriptions-item label="更新时间">{{ approval?.updateTime }}</n-descriptions-item>
          </n-descriptions>

          <!-- 类型特定字段 -->
          <n-divider />

          <!-- 请假信息 -->
          <template v-if="approval?.type === 'leave'">
            <n-descriptions label-placement="left" :column="2">
              <n-descriptions-item label="请假类型">{{ leaveTypeLabels[approval.leaveType] }}</n-descriptions-item>
              <n-descriptions-item label="请假天数">{{ approval.days }}天</n-descriptions-item>
              <n-descriptions-item label="开始日期">{{ approval.startDate }}</n-descriptions-item>
              <n-descriptions-item label="结束日期">{{ approval.endDate }}</n-descriptions-item>
              <n-descriptions-item label="请假原因" :span="2">{{ approval.reason }}</n-descriptions-item>
            </n-descriptions>
          </template>

          <!-- 报销信息 -->
          <template v-if="approval?.type === 'expense'">
            <n-descriptions label-placement="left" :column="2">
              <n-descriptions-item label="报销类型">{{ expenseTypeLabels[approval.expenseType] }}</n-descriptions-item>
              <n-descriptions-item label="报销金额">¥{{ approval.amount.toFixed(2) }}</n-descriptions-item>
              <n-descriptions-item label="报销说明" :span="2">{{ approval.description }}</n-descriptions-item>
            </n-descriptions>
          </template>

          <!-- 加班信息 -->
          <template v-if="approval?.type === 'overtime'">
            <n-descriptions label-placement="left" :column="2">
              <n-descriptions-item label="加班日期">{{ approval.overtimeDate }}</n-descriptions-item>
              <n-descriptions-item label="加班时长">{{ approval.hours }}小时</n-descriptions-item>
              <n-descriptions-item label="开始时间">{{ approval.startTime }}</n-descriptions-item>
              <n-descriptions-item label="结束时间">{{ approval.endTime }}</n-descriptions-item>
              <n-descriptions-item label="加班原因" :span="2">{{ approval.reason }}</n-descriptions-item>
            </n-descriptions>
          </template>

          <!-- 出差信息 -->
          <template v-if="approval?.type === 'travel'">
            <n-descriptions label-placement="left" :column="2">
              <n-descriptions-item label="目的地">{{ approval.destination }}</n-descriptions-item>
              <n-descriptions-item label="出差天数">{{ approval.days }}天</n-descriptions-item>
              <n-descriptions-item label="开始日期">{{ approval.startDate }}</n-descriptions-item>
              <n-descriptions-item label="结束日期">{{ approval.endDate }}</n-descriptions-item>
              <n-descriptions-item label="预算金额">¥{{ approval.budget.toFixed(2) }}</n-descriptions-item>
              <n-descriptions-item label="出差原因" :span="2">{{ approval.reason }}</n-descriptions-item>
            </n-descriptions>
          </template>

          <!-- 通用审批信息 -->
          <template v-if="approval?.type === 'general'">
            <n-descriptions label-placement="left" :column="1">
              <n-descriptions-item label="申请内容">{{ approval.content }}</n-descriptions-item>
            </n-descriptions>
          </template>
        </n-card>

        <!-- 审批流程 -->
        <n-card title="审批流程">
          <n-timeline>
            <n-timeline-item
              v-for="node in approval?.flowNodes"
              :key="node.id"
              :type="getTimelineType(node.status)"
              :title="node.nodeName"
            >
              <div class="timeline-content">
                <div class="approver-info">
                  <n-avatar round :src="node.approverAvatar" :size="24" />
                  <span class="approver-name">{{ node.approverName }}</span>
                  <n-tag :type="statusColorMap[node.status]" size="small">{{ ApprovalStatusLabels[node.status] }}</n-tag>
                </div>
                <div v-if="node.comment" class="comment">{{ node.comment }}</div>
                <div class="time">{{ node.createTime }}</div>
              </div>
            </n-timeline-item>
          </n-timeline>
        </n-card>

        <!-- 操作区域 -->
        <n-card v-if="showActions" title="审批操作">
          <n-space>
            <template v-if="isPendingApproval">
              <n-button type="success" :loading="actionLoading" @click="handleApprove">批准</n-button>
              <n-button type="error" :loading="actionLoading" @click="handleReject">驳回</n-button>
              <n-button type="info" :loading="actionLoading" @click="showTransferModal = true">转交</n-button>
            </template>
            <template v-if="canWithdraw">
              <n-popconfirm @positive-click="handleWithdraw">
                <template #trigger>
                  <n-button type="warning" :loading="actionLoading">撤回申请</n-button>
                </template>
                确定要撤回该申请吗？
              </n-popconfirm>
            </template>
          </n-space>
        </n-card>
      </n-space>
    </n-spin>

    <!-- 审批意见弹窗 -->
    <n-modal v-model:show="showActionModal" preset="dialog" :title="actionTitle" positive-text="确定" negative-text="取消" @positive-click="submitAction">
      <n-form ref="actionFormRef" :model="actionForm" label-placement="left" label-width="80">
        <n-form-item label="审批意见">
          <n-input v-model:value="actionForm.comment" type="textarea" placeholder="请输入审批意见（选填）" :rows="4" maxlength="200" show-count />
        </n-form-item>
      </n-form>
    </n-modal>

    <!-- 转交弹窗 -->
    <n-modal v-model:show="showTransferModal" preset="dialog" title="转交审批" positive-text="确定" negative-text="取消" @positive-click="submitTransfer">
      <n-form ref="transferFormRef" :model="transferForm" label-placement="left" label-width="80">
        <n-form-item label="转交给" path="transferTo">
          <n-select v-model:value="transferForm.transferTo" :options="userOptions" placeholder="请选择转交人员" />
        </n-form-item>
        <n-form-item label="转交原因">
          <n-input v-model:value="transferForm.comment" type="textarea" placeholder="请输入转交原因（选填）" :rows="4" maxlength="200" show-count />
        </n-form-item>
      </n-form>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  NSpin,
  NSpace,
  NButton,
  NIcon,
  NCard,
  NDescriptions,
  NDescriptionsItem,
  NDivider,
  NTag,
  NAvatar,
  NTimeline,
  NTimelineItem,
  NModal,
  NForm,
  NFormItem,
  NInput,
  NSelect,
  NPopconfirm,
  useMessage
} from 'naive-ui'
import { ArrowBackOutline } from '@vicons/ionicons5'
import { useApprovalStore } from '@/stores/approval'
import { useUserStore } from '@/stores/user'
import type { Approval, ApprovalStatus } from '@/types/approval'
import { ApprovalTypeLabels, ApprovalStatusLabels } from '@/types/approval'
import { request } from '@/utils/request'
import type { User } from '@/types/user'

const router = useRouter()
const route = useRoute()
const message = useMessage()
const approvalStore = useApprovalStore()
const userStore = useUserStore()

const loading = ref(false)
const actionLoading = ref(false)
const approval = ref<Approval | null>(null)

// 审批操作弹窗
const showActionModal = ref(false)
const actionTitle = ref('')
const actionForm = reactive({
  action: '' as 'approve' | 'reject',
  comment: ''
})

// 转交弹窗
const showTransferModal = ref(false)
const transferForm = reactive({
  transferTo: null as number | null,
  comment: ''
})

// 用户选项
const userOptions = ref<{ label: string; value: number }[]>([])

// 请假类型标签
const leaveTypeLabels: Record<string, string> = {
  annual: '年假',
  sick: '病假',
  personal: '事假',
  maternity: '产假',
  marriage: '婚假',
  bereavement: '丧假'
}

// 报销类型标签
const expenseTypeLabels: Record<string, string> = {
  travel: '差旅费',
  office: '办公费',
  entertainment: '招待费',
  other: '其他'
}

// 状态颜色映射
const statusColorMap: Record<ApprovalStatus, 'default' | 'success' | 'warning' | 'error' | 'info'> = {
  pending: 'warning',
  approved: 'success',
  rejected: 'error',
  withdrawn: 'default',
  transferred: 'info'
}

// 时间线颜色
function getTimelineType(status: ApprovalStatus): 'default' | 'success' | 'error' | 'warning' | 'info' {
  const map: Record<ApprovalStatus, 'default' | 'success' | 'error' | 'warning' | 'info'> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'error',
    withdrawn: 'default',
    transferred: 'info'
  }
  return map[status]
}

// 判断是否为待审批人
const isPendingApproval = computed(() => {
  if (!approval.value) return false
  const currentNode = approval.value.flowNodes.find((n) => n.status === 'pending')
  if (!currentNode) return false
  // 这里简化处理，实际应判断当前用户是否为审批人
  return approval.value.status === 'pending'
})

// 判断是否可以撤回
const canWithdraw = computed(() => {
  if (!approval.value) return false
  return approval.value.status === 'pending' && approval.value.applicantId === userStore.userInfo?.id
})

// 是否显示操作区域
const showActions = computed(() => {
  return isPendingApproval.value || canWithdraw.value
})

// 获取详情
async function fetchDetail() {
  const id = parseInt(route.params.id as string)
  if (!id) {
    message.error('无效的申请ID')
    router.back()
    return
  }

  loading.value = true
  try {
    approval.value = await approvalStore.getApprovalDetail(id)
  } catch (error) {
    message.error((error as Error).message || '获取详情失败')
    router.back()
  } finally {
    loading.value = false
  }
}

// 批准
function handleApprove() {
  actionTitle.value = '批准申请'
  actionForm.action = 'approve'
  actionForm.comment = ''
  showActionModal.value = true
}

// 驳回
function handleReject() {
  actionTitle.value = '驳回申请'
  actionForm.action = 'reject'
  actionForm.comment = ''
  showActionModal.value = true
}

// 提交审批操作
async function submitAction() {
  if (!approval.value) return

  actionLoading.value = true
  try {
    await approvalStore.approvalAction({
      approvalId: approval.value.id,
      action: actionForm.action,
      comment: actionForm.comment
    })
    message.success(actionForm.action === 'approve' ? '批准成功' : '驳回成功')
    showActionModal.value = false
    fetchDetail()
  } catch (error) {
    message.error((error as Error).message || '操作失败')
    return false
  } finally {
    actionLoading.value = false
  }
  return true
}

// 提交转交
async function submitTransfer() {
  if (!approval.value || !transferForm.transferTo) {
    message.error('请选择转交人员')
    return false
  }

  actionLoading.value = true
  try {
    await approvalStore.approvalAction({
      approvalId: approval.value.id,
      action: 'transfer',
      comment: transferForm.comment,
      transferTo: transferForm.transferTo
    })
    message.success('转交成功')
    showTransferModal.value = false
    fetchDetail()
  } catch (error) {
    message.error((error as Error).message || '转交失败')
    return false
  } finally {
    actionLoading.value = false
  }
  return true
}

// 撤回申请
async function handleWithdraw() {
  if (!approval.value) return

  actionLoading.value = true
  try {
    await approvalStore.withdrawApproval(approval.value.id)
    message.success('撤回成功')
    fetchDetail()
  } catch (error) {
    message.error((error as Error).message || '撤回失败')
  } finally {
    actionLoading.value = false
  }
}

onMounted(() => {
  fetchDetail()
  fetchUserOptions()
})

async function fetchUserOptions() {
  try {
    const result = await request.get<{ list: User[] }>('/user/list', {
      params: { pageSize: 100 }
    })
    userOptions.value = result.list.map((user) => ({
      label: `${user.nickname} (${user.deptName})`,
      value: user.id
    }))
  } catch {
    userOptions.value = []
  }
}
</script>

<style lang="scss" scoped>
.detail-page {
  max-width: 900px;
  margin: 0 auto;

  .timeline-content {
    .approver-info {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .approver-name {
      font-weight: 500;
    }

    .comment {
      margin-top: 8px;
      color: $text-color-3;
    }

    .time {
      margin-top: 4px;
      font-size: 12px;
      color: $text-color-4;
    }
  }
}
</style>
