<template>
  <div class="choose-plan-page">
    <div class="choose-plan-container">
      <div class="page-header">
        <h1 class="page-title">选择适合您的套餐</h1>
        <p class="page-desc">所有套餐均享 14 天免费试用，无需信用卡</p>
      </div>

      <div class="plan-grid">
        <div
          v-for="plan in plans"
          :key="plan.id"
          class="plan-card"
          :class="{ active: currentPlanId === plan.id }"
        >
          <div class="plan-badge" v-if="plan.code === 'standard'">
            <n-tag size="small" type="info">推荐</n-tag>
          </div>
          <div class="plan-header">
            <div class="plan-name">{{ plan.name }}</div>
            <div class="plan-price">
              <span v-if="plan.price === 0" class="price-free">免费</span>
              <template v-else>
                <span class="price-symbol">¥</span>
                <span class="price-amount">{{ plan.price }}</span>
                <span class="price-unit">/人/月</span>
              </template>
            </div>
            <div class="plan-users">适合 {{ plan.minUsers }}-{{ plan.maxUsers }} 人团队</div>
          </div>

          <div class="plan-features">
            <div class="feature-title">功能特性</div>
            <div v-for="(val, key) in plan.features" :key="key" class="feature-item">
              <span class="feature-check" :class="{ enabled: val }">✓</span>
              <span>{{ PLAN_LABELS[key as string] || key }}</span>
              <span v-if="typeof val === 'number'" class="feature-value">{{ formatStorage(val as number) }}</span>
            </div>
          </div>

          <div class="plan-action">
            <n-button
              v-if="currentPlanId === plan.id"
              type="primary"
              disabled
              block
            >
              当前套餐
            </n-button>
            <n-button
              v-else
              :type="plan.code === 'standard' ? 'primary' : 'default'"
              block
              @click="handleUpgrade(plan)"
            >
              {{ currentPlanId ? '升级' : '开始使用' }}
            </n-button>
          </div>
        </div>
      </div>

      <div class="plan-footer">
        <n-button text @click="router.back()">返回</n-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NTag, useMessage, useDialog } from 'naive-ui'
import { useTenantStore } from '@/stores/tenant'
import type { Plan } from '@/types/tenant'
import { PLAN_LABELS, formatStorage } from '@/utils/plan'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const tenantStore = useTenantStore()

const plans = ref<Plan[]>([])
const currentPlanId = ref<number>(0)

function handleUpgrade(plan: Plan) {
  dialog.warning({
    title: '确认升级',
    content: `确定要升级到「${plan.name}」套餐吗？费用为 ¥${plan.price}/人/月`,
    positiveText: '确认升级',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await tenantStore.upgradePlan(plan.id)
        currentPlanId.value = plan.id
        message.success('套餐升级成功')
      } catch (err: unknown) {
        const errorMsg = err instanceof Error ? err.message : '升级失败'
        message.error(errorMsg)
      }
    }
  })
}

onMounted(async () => {
  try {
    await tenantStore.fetchPlans()
    plans.value = tenantStore.plans
    if (tenantStore.currentPlan) {
      currentPlanId.value = tenantStore.currentPlan.id
    }
  } catch {
    // silently fail
  }
})
</script>

<style lang="scss" scoped>
.choose-plan-page {
  min-height: 100vh;
  background-color: $bg-color-3;
  padding: 40px 24px;
}

.choose-plan-container {
  max-width: 960px;
  margin: 0 auto;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;

  .page-title {
    font-size: 24px;
    font-weight: 600;
    color: $text-color-1;
    margin: 0 0 8px 0;
  }

  .page-desc {
    font-size: 14px;
    color: $text-color-3;
    margin: 0;
  }
}

.plan-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.plan-card {
  background: $bg-color-1;
  border: 1px solid $border-color;
  border-radius: $border-radius;
  padding: 24px 20px;
  display: flex;
  flex-direction: column;
  position: relative;
  transition: all $transition-duration ease;

  &:hover {
    box-shadow: $box-shadow-light;
  }

  &.active {
    border-color: $primary-color;
    background: $primary-color-suppl;
  }

  .plan-badge {
    position: absolute;
    top: -8px;
    right: 12px;
  }

  .plan-header {
    margin-bottom: 20px;

    .plan-name {
      font-size: 18px;
      font-weight: 500;
      color: $text-color-1;
      margin-bottom: 12px;
    }

    .plan-price {
      margin-bottom: 4px;

      .price-free {
        font-size: 24px;
        font-weight: 600;
        color: $success-color;
      }

      .price-symbol {
        font-size: 18px;
        color: $text-color-1;
      }

      .price-amount {
        font-size: 32px;
        font-weight: 600;
        color: $text-color-1;
      }

      .price-unit {
        font-size: 12px;
        color: $text-color-3;
      }
    }

    .plan-users {
      font-size: 12px;
      color: $text-color-3;
    }
  }

  .plan-features {
    flex: 1;
    margin-bottom: 20px;

    .feature-title {
      font-size: 13px;
      font-weight: 500;
      color: $text-color-2;
      margin-bottom: 8px;
    }

    .feature-item {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 13px;
      color: $text-color-2;
      padding: 4px 0;

      .feature-check {
        color: $border-color;
        font-size: 14px;

        &.enabled {
          color: $success-color;
        }
      }

      .feature-value {
        margin-left: auto;
        color: $text-color-3;
        font-size: 12px;
      }
    }
  }

  .plan-action {
    margin-top: auto;
  }
}

.plan-footer {
  text-align: center;
  margin-top: 24px;
}

@media (max-width: 1024px) {
  .plan-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .plan-grid {
    grid-template-columns: 1fr;
  }
}
</style>
