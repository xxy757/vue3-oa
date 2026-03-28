<template>
  <div class="my-apply-page">
    <n-card title="我的申请">
      <!-- 筛选区域 -->
      <div class="filter-section">
        <n-space>
          <n-select
            v-model:value="filterStatus"
            :options="statusOptions"
            placeholder="审批状态"
            clearable
            style="width: 120px"
            @update:value="handleSearch"
          />
          <n-select
            v-model:value="filterType"
            :options="typeOptions"
            placeholder="申请类型"
            clearable
            style="width: 120px"
            @update:value="handleSearch"
          />
          <n-button type="primary" @click="router.push('/approval/apply')">发起申请</n-button>
        </n-space>
      </div>

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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  NCard,
  NSpace,
  NSelect,
  NButton,
  NDataTable,
  NPagination,
  NTag,
  NPopconfirm,
  useMessage
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { useApprovalStore } from '@/stores/approval'
import type { Approval, ApprovalType, ApprovalStatus } from '@/types/approval'
import { ApprovalTypeLabels, ApprovalStatusLabels } from '@/types/approval'

const router = useRouter()
const message = useMessage()
const approvalStore = useApprovalStore()

const loading = ref(false)
const tableData = ref<Approval[]>([])
const filterStatus = ref<ApprovalStatus | null>(null)
const filterType = ref<ApprovalType | null>(null)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 状态选项
const statusOptions = [
  { label: '待审批', value: 'pending' },
  { label: '已通过', value: 'approved' },
  { label: '已驳回', value: 'rejected' },
  { label: '已撤回', value: 'withdrawn' },
  { label: '已转交', value: 'transferred' }
]

// 类型选项
const typeOptions = [
  { label: '请假申请', value: 'leave' },
  { label: '报销申请', value: 'expense' },
  { label: '加班申请', value: 'overtime' },
  { label: '出差申请', value: 'travel' },
  { label: '通用审批', value: 'general' }
]

// 状态颜色映射
const statusColorMap: Record<ApprovalStatus, 'default' | 'success' | 'warning' | 'error' | 'info'> = {
  pending: 'warning',
  approved: 'success',
  rejected: 'error',
  withdrawn: 'default',
  transferred: 'info'
}

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
    title: '审批状态',
    key: 'status',
    width: 100,
    render: (row) => h(NTag, { type: statusColorMap[row.status], size: 'small' }, () => ApprovalStatusLabels[row.status])
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
    title: '更新时间',
    key: 'updateTime',
    width: 160
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
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
        row.status === 'pending'
          ? h(
              NPopconfirm,
              {
                onPositiveClick: () => handleWithdraw(row)
              },
              {
                trigger: () =>
                  h(
                    NButton,
                    {
                      size: 'small',
                      type: 'warning'
                    },
                    () => '撤回'
                  ),
                default: () => '确定撤回该申请吗？'
              }
            )
          : null
      ])
    }
  }
]

// 获取列表数据
async function fetchData() {
  loading.value = true
  try {
    const result = await approvalStore.getMyApprovals({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: filterStatus.value || undefined,
      type: filterType.value || undefined
    })
    tableData.value = result.list
    pagination.total = result.total
  } catch (error) {
    message.error((error as Error).message || '获取数据失败')
  } finally {
    loading.value = false
  }
}

// 搜索
function handleSearch() {
  pagination.page = 1
  fetchData()
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

// 撤回申请
async function handleWithdraw(row: Approval) {
  try {
    await approvalStore.withdrawApproval(row.id)
    message.success('撤回成功')
    fetchData()
  } catch (error) {
    message.error((error as Error).message || '撤回失败')
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.my-apply-page {
  .filter-section {
    margin-bottom: 16px;
  }

  .pagination-section {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
