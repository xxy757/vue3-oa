<template>
  <div class="tenant-plan-page">
    <div class="page-title-bar">
      <h2 class="page-title">套餐管理</h2>
    </div>

    <n-spin :show="loading">
      <n-card class="current-plan-card">
        <div class="section-title">当前套餐</div>
        <div class="current-plan-info">
          <div class="plan-detail">
            <div class="plan-name-row">
              <span class="plan-name">{{ tenantStore.currentPlan?.name || '-' }}</span>
              <n-tag :type="planStatusType" size="small">{{ planStatusLabel }}</n-tag>
            </div>
            <div class="plan-meta">
              <span class="meta-item">
                <span class="meta-label">用户数</span>
                <span class="meta-value">{{ tenantStore.currentUsers }} / {{ tenantStore.maxUsers }}</span>
              </span>
              <span class="meta-item">
                <span class="meta-label">费用</span>
                <span class="meta-value">¥{{ tenantStore.currentPlan?.price || 0 }}/人/月</span>
              </span>
              <span class="meta-item">
                <span class="meta-label">到期时间</span>
                <span class="meta-value">{{ tenantStore.tenantInfo?.planExpireAt?.substring(0, 10) || '-' }}</span>
              </span>
            </div>
          </div>
          <n-button type="primary" @click="showUpgrade = true">升级套餐</n-button>
        </div>
      </n-card>

      <n-card style="margin-top: 16px">
        <div class="section-title">功能对比</div>
        <n-table :bordered="false" :single-line="false" size="small">
          <thead>
            <tr>
              <th>功能</th>
              <th v-for="plan in tenantStore.plans" :key="plan.id">{{ plan.name }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="featureKey in featureKeys" :key="featureKey">
              <td>{{ PLAN_LABELS[featureKey] || featureKey }}</td>
              <td v-for="plan in tenantStore.plans" :key="plan.id">
                <template v-if="typeof plan.features[featureKey] === 'number'">
                  {{ formatStorage(plan.features[featureKey] as number) }}
                </template>
                <template v-else>
                  <span :style="{ color: plan.features[featureKey] ? '#52C41A' : '#D9D9D9' }">
                    {{ plan.features[featureKey] ? '✓' : '✗' }}
                  </span>
                </template>
              </td>
            </tr>
            <tr>
              <td>价格</td>
              <td v-for="plan in tenantStore.plans" :key="plan.id">
                {{ formatPrice(plan.price) }}
              </td>
            </tr>
          </tbody>
        </n-table>
      </n-card>
    </n-spin>

    <n-modal
      v-model:show="showUpgrade"
      preset="dialog"
      title="升级套餐"
      :show-icon="false"
      style="width: 520px"
    >
      <div class="upgrade-plans">
        <div
          v-for="plan in upgradePlans"
          :key="plan.id"
          class="upgrade-item"
          :class="{ selected: selectedPlanId === plan.id }"
          @click="selectedPlanId = plan.id"
        >
          <div class="upgrade-name">{{ plan.name }}</div>
          <div class="upgrade-price">¥{{ plan.price }}/人/月</div>
          <div class="upgrade-users">{{ plan.minUsers }}-{{ plan.maxUsers }} 人</div>
        </div>
      </div>
      <template #action>
        <n-button @click="showUpgrade = false">取消</n-button>
        <n-button type="primary" :loading="upgrading" :disabled="!selectedPlanId" @click="handleUpgrade">
          确认升级
        </n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  NCard,
  NButton,
  NTag,
  NTable,
  NModal,
  NSpin,
  useMessage
} from 'naive-ui'
import { useTenantStore } from '@/stores/tenant'
import { PLAN_LABELS, formatPrice, formatStorage } from '@/utils/plan'

const message = useMessage()
const tenantStore = useTenantStore()

const loading = ref(false)
const showUpgrade = ref(false)
const upgrading = ref(false)
const selectedPlanId = ref<number>(0)

const planStatusType = computed(() => {
  if (tenantStore.isExpired) return 'error' as const
  if (tenantStore.isTrial) return 'warning' as const
  return 'success' as const
})

const planStatusLabel = computed(() => {
  const status = tenantStore.tenantStatus
  if (status === 'trial') return '试用中'
  if (status === 'active') return '已激活'
  if (status === 'suspended') return '已暂停'
  return status
})

const featureKeys = computed(() => {
  const keys = new Set<string>()
  tenantStore.plans.forEach(plan => {
    Object.keys(plan.features).forEach(k => keys.add(k))
  })
  return Array.from(keys)
})

const upgradePlans = computed(() => {
  const currentId = tenantStore.currentPlan?.id
  return tenantStore.plans.filter(p => p.id !== currentId)
})

async function handleUpgrade() {
  if (!selectedPlanId.value) return
  upgrading.value = true
  try {
    await tenantStore.upgradePlan(selectedPlanId.value)
    message.success('套餐升级成功')
    showUpgrade.value = false
  } catch (err: unknown) {
    const errorMsg = err instanceof Error ? err.message : '升级失败'
    message.error(errorMsg)
  } finally {
    upgrading.value = false
  }
}

onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([
      tenantStore.fetchTenantInfo(),
      tenantStore.fetchPlans()
    ])
  } catch {
    message.error('获取套餐信息失败')
  } finally {
    loading.value = false
  }
})
</script>

<style lang="scss" scoped>
.tenant-plan-page {
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

.section-title {
  font-size: 16px;
  font-weight: 500;
  color: $text-color-1;
  margin-bottom: 16px;
}

.current-plan-card {
  .current-plan-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .plan-detail {
    .plan-name-row {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 12px;

      .plan-name {
        font-size: 16px;
        font-weight: 500;
        color: $text-color-1;
      }
    }

    .plan-meta {
      display: flex;
      gap: 32px;

      .meta-item {
        .meta-label {
          font-size: 12px;
          color: $text-color-3;
          display: block;
          margin-bottom: 4px;
        }

        .meta-value {
          font-size: 14px;
          color: $text-color-2;
        }
      }
    }
  }
}

.upgrade-plans {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;

  .upgrade-item {
    border: 1px solid $border-color;
    border-radius: $border-radius;
    padding: 16px;
    cursor: pointer;
    transition: all $transition-duration ease;

    &:hover {
      border-color: $primary-color;
    }

    &.selected {
      border-color: $primary-color;
      background: $primary-color-suppl;
    }

    .upgrade-name {
      font-size: 14px;
      font-weight: 500;
      color: $text-color-1;
      margin-bottom: 8px;
    }

    .upgrade-price {
      font-size: 14px;
      color: $primary-color;
      margin-bottom: 4px;
    }

    .upgrade-users {
      font-size: 12px;
      color: $text-color-3;
    }
  }
}
</style>
