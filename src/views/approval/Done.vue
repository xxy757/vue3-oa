<template>
  <div class="done-page">
    <n-card title="已办审批">
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
import { NCard, NDataTable, NPagination, NTag, NButton, NSpace, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { useApprovalStore } from '@/stores/approval'
import type { Approval, ApprovalStatus } from '@/types/approval'
import { ApprovalTypeLabels, ApprovalStatusLabels } from '@/types/approval'

const router = useRouter()
const message = useMessage()
const approvalStore = useApprovalStore()

const loading = ref(false)
const tableData = ref<Approval[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

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
    title: '审批状态',
    key: 'status',
    width: 100,
    render: (row) => h(NTag, { type: statusColorMap[row.status], size: 'small' }, () => ApprovalStatusLabels[row.status])
  },
  {
    title: '处理时间',
    key: 'updateTime',
    width: 160
  },
  {
    title: '操作',
    key: 'actions',
    width: 80,
    render: (row) => {
      return h(
        NButton,
        {
          size: 'small',
          onClick: () => handleDetail(row)
        },
        () => '详情'
      )
    }
  }
]

// 获取列表数据
async function fetchData() {
  loading.value = true
  try {
    const result = await approvalStore.getDoneApprovals({
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

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.done-page {
  .pagination-section {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
