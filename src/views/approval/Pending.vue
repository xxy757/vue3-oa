<template>
  <div class="pending-page">
    <n-card title="待我审批">
      <!-- 表格 -->
      <n-data-table
        :columns="columns"
        :data="tableData"
        :loading="loading"
        :pagination="false"
        :row-key="(row: Approval) => row.id"
      />

      <!-- 分页 -->
      <div class="pagination-section">
        <n-pagination
          v-model:page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :item-count="pagination.total"
          :page-sizes="[10, 20, 50]"
          show-size-picker
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </n-card>

    <!-- 审批操作弹窗 -->
    <n-modal v-model:show="actionModalVisible" preset="dialog" :title="actionTitle" positive-text="确定" negative-text="取消" @positive-click="handleActionSubmit">
      <n-form ref="actionFormRef" :model="actionForm" label-placement="left" label-width="80">
        <n-form-item label="审批意见">
          <n-input v-model:value="actionForm.comment" type="textarea" placeholder="请输入审批意见（选填）" :rows="4" maxlength="200" show-count />
        </n-form-item>
      </n-form>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  NCard,
  NDataTable,
  NPagination,
  NTag,
  NButton,
  NSpace,
  NModal,
  NForm,
  NFormItem,
  NInput,
  useMessage
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { useApprovalStore } from '@/stores/approval'
import type { Approval } from '@/types/approval'
import { ApprovalTypeLabels, ApprovalStatusLabels } from '@/types/approval'

const router = useRouter()
const message = useMessage()
const approvalStore = useApprovalStore()

const loading = ref(false)
const tableData = ref<Approval[]>([])
const actionModalVisible = ref(false)
const actionTitle = ref('')
const actionFormRef = ref()
const actionForm = reactive({
  approvalId: 0,
  action: '' as 'approve' | 'reject',
  comment: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 表格列配置
const columns: DataTableColumns<Approval> = [
  {
    title: '申请标题',
    key: 'title',
    ellipsis: { tooltip: true }
  },
  {
    title: '申请类型',
    key: 'type',
    width: 100,
    render: (row) => ApprovalTypeLabels[row.type]
  },
  {
    title: '申请人',
    key: 'applicantName',
    width: 100
  },
  {
    title: '申请人部门',
    key: 'applicantDept',
    width: 100
  },
  {
    title: '当前进度',
    key: 'progress',
    width: 100,
    render: (row) => `${row.currentStep}/${row.totalStep}`
  },
  {
    title: '申请时间',
    key: 'createTime',
    width: 160
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render: (row) => {
      return h(NSpace, null, () => [
        h(
          NButton,
          {
            size: 'small',
            onClick: () => handleDetail(row)
          },
          () => '详情'
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'success',
            onClick: () => handleAction(row, 'approve')
          },
          () => '批准'
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            onClick: () => handleAction(row, 'reject')
          },
          () => '驳回'
        )
      ])
    }
  }
]

// 获取列表数据
async function fetchData() {
  loading.value = true
  try {
    const result = await approvalStore.getPendingApprovals({
      page: pagination.page,
      pageSize: pagination.pageSize
    })
    tableData.value = result.list
    pagination.total = result.total
  } catch (error) {
    message.error((error as Error).message || '获取数据失败')
  } finally {
    loading.value = false
  }
}

// 分页变化
function handlePageChange() {
  fetchData()
}

function handlePageSizeChange() {
  pagination.page = 1
  fetchData()
}

// 查看详情
function handleDetail(row: Approval) {
  router.push(`/approval/detail/${row.id}`)
}

// 打开审批操作弹窗
function handleAction(row: Approval, action: 'approve' | 'reject') {
  actionForm.approvalId = row.id
  actionForm.action = action
  actionForm.comment = ''
  actionTitle.value = action === 'approve' ? '批准申请' : '驳回申请'
  actionModalVisible.value = true
}

// 提交审批操作
async function handleActionSubmit() {
  try {
    await approvalStore.approvalAction({
      approvalId: actionForm.approvalId,
      action: actionForm.action,
      comment: actionForm.comment
    })
    message.success(actionForm.action === 'approve' ? '批准成功' : '驳回成功')
    actionModalVisible.value = false
    fetchData()
  } catch (error) {
    message.error((error as Error).message || '操作失败')
    return false
  }
  return true
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.pending-page {
  .pagination-section {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
