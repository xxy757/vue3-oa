<template>
  <div class="admin-dashboard">
    <div class="page-title-bar">
      <h2 class="page-title">管理概览</h2>
    </div>

    <n-grid :cols="4" :x-gap="16" :y-gap="16">
      <n-gi>
        <n-card class="stat-card">
          <div class="stat-content">
            <div class="stat-label">总租户数</div>
            <div class="stat-value">{{ stats.totalTenants }}</div>
            <div class="stat-extra">
              <span class="stat-trend up">+{{ stats.newTenantsThisMonth }}</span>
              <span class="stat-period">本月新增</span>
            </div>
          </div>
          <div class="stat-icon" style="background: #E6F4FF;">
            <n-icon :size="28" color="#1677FF"><BusinessOutline /></n-icon>
          </div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="stat-card">
          <div class="stat-content">
            <div class="stat-label">活跃租户</div>
            <div class="stat-value">{{ stats.activeTenants }}</div>
            <div class="stat-extra">
              <span class="stat-period">占比 {{ activeRate }}%</span>
            </div>
          </div>
          <div class="stat-icon" style="background: #F6FFED;">
            <n-icon :size="28" color="#52C41A"><CheckmarkCircleOutline /></n-icon>
          </div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="stat-card">
          <div class="stat-content">
            <div class="stat-label">总用户数</div>
            <div class="stat-value">{{ stats.totalUsers }}</div>
            <div class="stat-extra">
              <span class="stat-period">跨所有租户</span>
            </div>
          </div>
          <div class="stat-icon" style="background: #FFF7E6;">
            <n-icon :size="28" color="#FA8C16"><PeopleOutline /></n-icon>
          </div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="stat-card">
          <div class="stat-content">
            <div class="stat-label">本月收入</div>
            <div class="stat-value">¥{{ stats.monthlyRevenue.toFixed(2) }}</div>
            <div class="stat-extra">
              <span class="stat-trend up">+{{ stats.revenueGrowth }}%</span>
              <span class="stat-period">环比增长</span>
            </div>
          </div>
          <div class="stat-icon" style="background: #F9F0FF;">
            <n-icon :size="28" color="#722ED1"><WalletOutline /></n-icon>
          </div>
        </n-card>
      </n-gi>
    </n-grid>

    <n-grid :cols="2" :x-gap="16" :y-gap="16" style="margin-top: 16px">
      <n-gi>
        <n-card title="套餐分布">
          <div class="plan-distribution">
            <div v-for="item in stats.planDistribution" :key="item.name" class="distribution-item">
              <div class="dist-label">{{ item.name }}</div>
              <n-progress
                type="line"
                :percentage="item.percentage"
                :show-indicator="false"
                :height="16"
                :color="item.color"
                style="flex: 1; margin: 0 16px"
              />
              <div class="dist-count">{{ item.count }} 家</div>
            </div>
          </div>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card title="最近注册租户">
          <div class="recent-tenants">
            <div v-for="tenant in stats.recentTenants" :key="tenant.id" class="recent-tenant-item">
              <div class="tenant-main">
                <span class="tenant-name">{{ tenant.name }}</span>
                <n-tag size="small" :type="tenant.status === 'active' ? 'success' : 'warning'">
                  {{ tenant.status === 'active' ? '已激活' : '试用中' }}
                </n-tag>
              </div>
              <div class="tenant-time">{{ tenant.createTime?.substring(0, 10) }}</div>
            </div>
            <div v-if="stats.recentTenants.length === 0" class="empty-tip">暂无数据</div>
          </div>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  NCard,
  NGrid,
  NGi,
  NIcon,
  NProgress,
  NTag,
  useMessage
} from 'naive-ui'
import {
  BusinessOutline,
  CheckmarkCircleOutline,
  PeopleOutline,
  WalletOutline
} from '@vicons/ionicons5'
import { request } from '@/utils/request'
import type { DashboardStats } from '@/types/admin'

const message = useMessage()

const stats = ref<DashboardStats>({
  totalTenants: 0,
  activeTenants: 0,
  newTenantsThisMonth: 0,
  totalUsers: 0,
  monthlyRevenue: 0,
  revenueGrowth: 0,
  planDistribution: [],
  recentTenants: []
})

const activeRate = computed(() => {
  if (!stats.value.totalTenants) return 0
  return Math.round((stats.value.activeTenants / stats.value.totalTenants) * 100)
})

onMounted(async () => {
  try {
    const data = await request.get<DashboardStats>('/admin/dashboard')
    stats.value = data
  } catch {
    message.error('获取管理概览数据失败')
  }
})
</script>

<style lang="scss" scoped>
.admin-dashboard {
  padding: 0;
}

.page-title-bar {
  margin-bottom: 16px;

  .page-title {
    font-size: 20px;
    font-weight: 600;
    color: $text-color-1;
    margin: 0;
  }
}

.stat-card {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .stat-content {
    flex: 1;

    .stat-label {
      font-size: 12px;
      color: $text-color-3;
      margin-bottom: 8px;
    }

    .stat-value {
      font-size: 30px;
      font-weight: 600;
      color: $text-color-1;
      line-height: 38px;
      margin-bottom: 4px;
    }

    .stat-extra {
      font-size: 12px;
      color: $text-color-3;

      .stat-trend {
        margin-right: 4px;

        &.up {
          color: $success-color;
        }

        &.down {
          color: $error-color;
        }
      }
    }
  }

  .stat-icon {
    width: 56px;
    height: 56px;
    border-radius: $border-radius;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }
}

.plan-distribution {
  .distribution-item {
    display: flex;
    align-items: center;
    padding: 8px 0;

    .dist-label {
      width: 80px;
      font-size: 14px;
      color: $text-color-2;
    }

    .dist-count {
      width: 60px;
      text-align: right;
      font-size: 14px;
      color: $text-color-2;
    }
  }
}

.recent-tenants {
  .recent-tenant-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 0;
    border-bottom: 1px solid $border-color-dark;

    &:last-child {
      border-bottom: none;
    }

    .tenant-main {
      display: flex;
      align-items: center;
      gap: 8px;

      .tenant-name {
        font-size: 14px;
        color: $text-color-1;
      }
    }

    .tenant-time {
      font-size: 12px;
      color: $text-color-3;
    }
  }
}

.empty-tip {
  text-align: center;
  color: $text-color-3;
  padding: 32px 0;
  font-size: 14px;
}
</style>
