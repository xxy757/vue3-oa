<template>
  <div class="notice-list">
    <n-card title="公告通知">
      <!-- 搜索和筛选区域 -->
      <div class="filter-bar">
        <n-space>
          <n-select
            v-model:value="filterType"
            :options="typeOptions"
            placeholder="选择类型"
            clearable
            style="width: 120px"
            @update:value="handleFilterChange"
          />
          <n-input
            v-model:value="keyword"
            placeholder="搜索公告标题"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
            @clear="handleSearch"
          >
            <template #prefix>
              <n-icon><SearchOutline /></n-icon>
            </template>
          </n-input>
          <n-button type="primary" @click="handleSearch">
            搜索
          </n-button>
        </n-space>
      </div>

      <!-- 公告列表 -->
      <n-data-table
        :columns="columns"
        :data="noticeList"
        :loading="loading"
        :row-key="(row: Notice) => row.id"
        :row-class-name="getRowClass"
        @update:page="handlePageChange"
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  NCard,
  NSpace,
  NSelect,
  NInput,
  NButton,
  NDataTable,
  NTag,
  NPagination,
  NIcon,
  NBadge
} from 'naive-ui'
import { SearchOutline, PinOutline } from '@vicons/ionicons5'
import type { DataTableColumns } from 'naive-ui'
import { useNoticeStore } from '@/stores/notice'
import type { Notice, NoticeType } from '@/types/notice'
import { NoticeTypeLabels } from '@/types/notice'

const router = useRouter()
const noticeStore = useNoticeStore()

const loading = ref(false)
const filterType = ref<string | null>(null)
const keyword = ref('')
const noticeList = ref<Notice[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
  totalPages: 0
})

// 类型选项
const typeOptions = [
  { label: '全部', value: 'all' },
  { label: '通知', value: 'notice' },
  { label: '公告', value: 'announcement' },
  { label: '制度', value: 'policy' },
  { label: '紧急', value: 'urgent' }
]

// 类型标签颜色
const getTypeTagType = (type: NoticeType) => {
  const typeMap: Record<NoticeType, 'default' | 'info' | 'success' | 'warning' | 'error'> = {
    notice: 'info',
    announcement: 'default',
    policy: 'success',
    urgent: 'error'
  }
  return typeMap[type]
}

// 表格列配置
const columns: DataTableColumns<Notice> = [
  {
    title: '',
    key: 'isTop',
    width: 40,
    render(row) {
      return row.isTop
        ? h(NIcon, { color: '#f0a020', size: 18 }, { default: () => h(PinOutline) })
        : null
    }
  },
  {
    title: '标题',
    key: 'title',
    ellipsis: { tooltip: true },
    render(row) {
      return h('div', { class: 'notice-title-cell' }, [
        h('span', {
          class: row.isRead ? '' : 'unread',
          style: row.isRead ? {} : { fontWeight: 'bold' }
        }, row.title)
      ])
    }
  },
  {
    title: '类型',
    key: 'type',
    width: 100,
    render(row) {
      return h(NTag, { type: getTypeTagType(row.type), size: 'small' }, {
        default: () => NoticeTypeLabels[row.type]
      })
    }
  },
  {
    title: '发布者',
    key: 'publisherName',
    width: 100
  },
  {
    title: '发布时间',
    key: 'createTime',
    width: 160
  },
  {
    title: '阅读',
    key: 'readCount',
    width: 80,
    render(row) {
      return h('span', {}, `${row.readCount}次`)
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 100,
    render(row) {
      return h(NButton, {
        text: true,
        type: 'primary',
        onClick: () => handleViewDetail(row)
      }, { default: () => '查看详情' })
    }
  }
]

// 获取行样式类名
const getRowClass = (row: Notice) => {
  return row.isRead ? '' : 'unread-row'
}

// 获取公告列表
const fetchNoticeList = async () => {
  loading.value = true
  try {
    const result = await noticeStore.getNoticeList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      type: filterType.value && filterType.value !== 'all' ? filterType.value as NoticeType : undefined,
      keyword: keyword.value || undefined
    })
    noticeList.value = result.list
    pagination.total = result.total
    pagination.totalPages = result.totalPages
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchNoticeList()
}

// 筛选变化
const handleFilterChange = () => {
  pagination.page = 1
  fetchNoticeList()
}

// 分页变化
const handlePageChange = (page: number) => {
  pagination.page = page
  fetchNoticeList()
}

// 每页条数变化
const handlePageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  pagination.page = 1
  fetchNoticeList()
}

// 查看详情
const handleViewDetail = (row: Notice) => {
  router.push(`/notice/detail/${row.id}`)
}

onMounted(() => {
  fetchNoticeList()
})
</script>

<style lang="scss" scoped>
.notice-list {
  .filter-bar {
    margin-bottom: 16px;
  }

  .pagination-wrapper {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
  }

  .notice-title-cell {
    .unread {
      font-weight: bold;
    }
  }

  :deep(.unread-row) {
    background-color: rgba(24, 144, 255, 0.05);
  }
}
</style>
