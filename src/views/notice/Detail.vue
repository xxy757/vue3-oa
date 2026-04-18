<template>
  <div class="notice-detail">
    <n-spin :show="loading">
      <n-card v-if="notice">
        <!-- 头部信息 -->
        <template #header>
          <div class="detail-header">
            <div class="header-top">
              <n-tag :type="getTypeTagType(notice.type)" size="small">
                {{ NoticeTypeLabels[notice.type] }}
              </n-tag>
              <n-tag v-if="notice.isTop" type="warning" size="small">
                <template #icon>
                  <n-icon><PinOutline /></n-icon>
                </template>
                置顶
              </n-tag>
            </div>
            <h1 class="detail-title">{{ notice.title }}</h1>
            <div class="detail-meta">
              <n-space size="large">
                <span class="meta-item">
                  <n-icon><PersonOutline /></n-icon>
                  {{ notice.publisherName }}
                </span>
                <span class="meta-item">
                  <n-icon><TimeOutline /></n-icon>
                  {{ notice.createTime }}
                </span>
                <span class="meta-item">
                  <n-icon><EyeOutline /></n-icon>
                  阅读 {{ notice.readCount }} 次
                </span>
              </n-space>
            </div>
          </div>
        </template>

        <!-- 公告内容 -->
        <div class="detail-content" v-html="notice.content"></div>

        <!-- 底部操作 -->
        <template #action>
          <n-space>
            <n-button @click="handleGoBack">
              <template #icon>
                <n-icon><ArrowBackOutline /></n-icon>
              </template>
              返回列表
            </n-button>
          </n-space>
        </template>
      </n-card>

      <!-- 公告不存在 -->
      <n-empty v-else-if="!loading" description="公告不存在">
        <template #extra>
          <n-button @click="handleGoBack">返回列表</n-button>
        </template>
      </n-empty>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NCard,
  NSpin,
  NTag,
  NIcon,
  NButton,
  NSpace,
  NEmpty
} from 'naive-ui'
import {
  PinOutline,
  PersonOutline,
  TimeOutline,
  EyeOutline,
  ArrowBackOutline
} from '@vicons/ionicons5'
import { useNoticeStore } from '@/stores/notice'
import type { Notice, NoticeType } from '@/types/notice'
import { NoticeTypeLabels } from '@/types/notice'

const route = useRoute()
const router = useRouter()
const noticeStore = useNoticeStore()

const loading = ref(false)
const notice = ref<Notice | null>(null)

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

// 获取公告详情
const fetchNoticeDetail = async () => {
  const id = Number(route.params.id)
  if (!id) {
    return
  }

  loading.value = true
  try {
    notice.value = await noticeStore.getNoticeDetail(id)
  } finally {
    loading.value = false
  }
}

// 返回列表
const handleGoBack = () => {
  router.push('/notice/list')
}

onMounted(() => {
  fetchNoticeDetail()
})
</script>

<style lang="scss" scoped>
.notice-detail {
  .detail-header {
    .header-top {
      margin-bottom: 8px;
      display: flex;
      gap: 8px;
    }

    .detail-title {
      margin: 0 0 12px 0;
      font-size: 20px;
      font-weight: 600;
      color: $text-color-1;
    }

    .detail-meta {
      color: $text-color-3;
      font-size: 14px;

      .meta-item {
        display: inline-flex;
        align-items: center;
        gap: 4px;
      }
    }
  }

  .detail-content {
    padding: 24px 0;
    line-height: 1.8;
    font-size: 14px;
    color: $text-color-2;

    :deep(h1),
    :deep(h2),
    :deep(h3),
    :deep(h4),
    :deep(h5),
    :deep(h6) {
      margin: 16px 0 8px;
      color: $text-color-1;
    }

    :deep(p) {
      margin: 12px 0;
    }

    :deep(img) {
      max-width: 100%;
      border-radius: 4px;
    }

    :deep(ul),
    :deep(ol) {
      padding-left: 24px;
      margin: 12px 0;
    }

    :deep(blockquote) {
      margin: 12px 0;
      padding: 8px 16px;
      border-left: 4px solid $primary-color;
      background-color: $bg-color-2;
      border-radius: 0 4px 4px 0;
    }
  }
}
</style>
