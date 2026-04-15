<template>
  <div class="register-page">
    <div class="register-container">
      <div class="register-header">
        <h1 class="register-title">企业注册</h1>
        <p class="register-desc">注册即享 14 天免费试用</p>
      </div>

      <n-steps :current="currentStep" class="register-steps">
        <n-step title="企业信息" />
        <n-step title="选择套餐" />
        <n-step title="注册完成" />
      </n-steps>

      <div class="register-content">
        <n-form
          v-show="currentStep === 1"
          ref="formRef"
          :model="formData"
          :rules="formRules"
          label-placement="left"
          label-width="80px"
        >
          <n-form-item label="企业名称" path="name">
            <n-input v-model:value="formData.name" placeholder="请输入企业名称" />
          </n-form-item>
          <n-form-item label="企业标识" path="slug">
            <n-input v-model:value="formData.slug" placeholder="用于子域名，如 mycompany">
              <template #suffix>.oa-saas.com</template>
            </n-input>
          </n-form-item>
          <n-form-item label="联系人" path="contactName">
            <n-input v-model:value="formData.contactName" placeholder="请输入联系人姓名" />
          </n-form-item>
          <n-form-item label="联系电话" path="contactPhone">
            <n-input v-model:value="formData.contactPhone" placeholder="请输入联系电话" />
          </n-form-item>
          <n-form-item label="联系邮箱" path="contactEmail">
            <n-input v-model:value="formData.contactEmail" placeholder="请输入联系邮箱" />
          </n-form-item>
        </n-form>

        <div v-show="currentStep === 2" class="plan-select">
          <div class="plan-cards">
            <div
              v-for="plan in plans"
              :key="plan.id"
              class="plan-card"
              :class="{ active: selectedPlanId === plan.id }"
              @click="selectedPlanId = plan.id"
            >
              <div class="plan-name">{{ plan.name }}</div>
              <div class="plan-price">
                <span v-if="plan.price === 0" class="price-free">免费</span>
                <template v-else>
                  <span class="price-symbol">¥</span>
                  <span class="price-amount">{{ plan.price }}</span>
                  <span class="price-unit">/人/月</span>
                </template>
              </div>
              <div class="plan-users">{{ plan.minUsers }}-{{ plan.maxUsers }} 人</div>
              <div class="plan-features">
                <div v-for="(val, key) in plan.features" :key="key" class="plan-feature">
                  <n-icon :color="val ? '#52C41A' : '#D9D9D9'" size="16">
                    <CheckmarkCircleOutline v-if="val" />
                    <CloseCircleOutline v-else />
                  </n-icon>
                  <span>{{ PLAN_LABELS[key as string] || key }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-show="currentStep === 3" class="register-result">
          <n-result status="success" title="注册成功" :description="`企业 ${registerData?.name} 已创建成功`">
            <template #footer>
              <div class="result-info">
                <n-card size="small">
                  <n-descriptions label-placement="left" bordered :column="1">
                    <n-descriptions-item label="企业标识">{{ registerData?.slug }}</n-descriptions-item>
                    <n-descriptions-item label="管理员账号">{{ registerData?.adminUser?.username }}</n-descriptions-item>
                    <n-descriptions-item label="初始密码">
                      <n-tag type="warning">{{ registerData?.adminUser?.tempPassword }}</n-tag>
                    </n-descriptions-item>
                    <n-descriptions-item label="试用到期">{{ registerData?.trialEndsAt?.substring(0, 10) }}</n-descriptions-item>
                  </n-descriptions>
                </n-card>
                <n-button type="primary" @click="goLogin">立即登录</n-button>
              </div>
            </template>
          </n-result>
        </div>
      </div>

      <div v-if="currentStep !== 3" class="register-footer">
        <n-button v-if="currentStep > 1" @click="currentStep--">上一步</n-button>
        <span v-else />
        <div class="footer-right">
          <n-button v-if="currentStep < 2" type="primary" @click="handleNext">下一步</n-button>
          <n-button v-else type="primary" :loading="loading" @click="handleRegister">完成注册</n-button>
          <n-button text type="primary" @click="goLogin">已有账号？去登录</n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  NForm,
  NFormItem,
  NInput,
  NButton,
  NSteps,
  NStep,
  NCard,
  NResult,
  NDescriptions,
  NDescriptionsItem,
  NIcon,
  NTag,
  useMessage
} from 'naive-ui'
import { CheckmarkCircleOutline, CloseCircleOutline } from '@vicons/ionicons5'
import type { FormInst, FormRules } from 'naive-ui'
import { request } from '@/utils/request'
import type { Plan, RegisterForm, RegisterResult } from '@/types/tenant'
import { PLAN_LABELS } from '@/utils/plan'

const router = useRouter()
const message = useMessage()

const formRef = ref<FormInst | null>(null)
const currentStep = ref(1)
const loading = ref(false)
const plans = ref<Plan[]>([])
const selectedPlanId = ref<number>(0)
const registerData = ref<RegisterResult | null>(null)

const formData = ref<RegisterForm>({
  name: '',
  slug: '',
  contactName: '',
  contactPhone: '',
  contactEmail: ''
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入企业名称', trigger: 'blur' }],
  slug: [
    { required: true, message: '请输入企业标识', trigger: 'blur' },
    { pattern: /^[a-z0-9][a-z0-9-]{1,48}[a-z0-9]$/, message: '只能包含小写字母、数字和连字符', trigger: 'blur' }
  ],
  contactName: [{ required: true, message: '请输入联系人', trigger: 'blur' }],
  contactPhone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  contactEmail: [
    { required: true, message: '请输入联系邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

async function handleNext() {
  if (currentStep.value === 1) {
    try {
      await formRef.value?.validate()
      currentStep.value = 2
    } catch {
      // validation failed
    }
  }
}

async function handleRegister() {
  if (!selectedPlanId.value) {
    message.warning('请选择一个套餐')
    return
  }

  loading.value = true
  try {
    const data = await request.post<RegisterResult>('/tenant/register', {
      ...formData.value,
      planId: selectedPlanId.value
    })
    registerData.value = data
    currentStep.value = 3
    message.success('注册成功')
  } catch (err: unknown) {
    const errorMsg = err instanceof Error ? err.message : '注册失败'
    message.error(errorMsg)
  } finally {
    loading.value = false
  }
}

function goLogin() {
  router.push('/login')
}

onMounted(async () => {
  try {
    const data = await request.get<Plan[]>('/plans')
    plans.value = data
    if (data.length > 0) {
      selectedPlanId.value = data[0].id
    }
  } catch {
    message.error('获取套餐列表失败，请刷新重试')
  }
})
</script>

<style lang="scss" scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: $bg-color-3;
  padding: 24px;
}

.register-container {
  width: 720px;
  background: $bg-color-1;
  border-radius: $border-radius;
  box-shadow: $box-shadow-light;
  padding: 40px;
}

.register-header {
  text-align: center;
  margin-bottom: 32px;

  .register-title {
    font-size: 24px;
    font-weight: 600;
    color: $text-color-1;
    margin: 0 0 8px 0;
  }

  .register-desc {
    font-size: 14px;
    color: $text-color-3;
    margin: 0;
  }
}

.register-steps {
  margin-bottom: 32px;
}

.register-content {
  min-height: 300px;
}

.plan-select {
  .plan-cards {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }

  .plan-card {
    border: 1px solid $border-color;
    border-radius: $border-radius;
    padding: 20px;
    cursor: pointer;
    transition: all $transition-duration ease;

    &:hover {
      border-color: $primary-color;
      box-shadow: $box-shadow-light;
    }

    &.active {
      border-color: $primary-color;
      background: $primary-color-suppl;
    }

    .plan-name {
      font-size: 16px;
      font-weight: 500;
      color: $text-color-1;
      margin-bottom: 12px;
    }

    .plan-price {
      margin-bottom: 8px;

      .price-free {
        font-size: 20px;
        font-weight: 600;
        color: $success-color;
      }

      .price-symbol {
        font-size: 16px;
        color: $text-color-1;
      }

      .price-amount {
        font-size: 28px;
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
      margin-bottom: 12px;
    }

    .plan-features {
      .plan-feature {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 13px;
        color: $text-color-2;
        padding: 4px 0;
      }
    }
  }
}

.register-result {
  .result-info {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
  }

  .result-btn {
    margin-top: 8px;
  }
}

.register-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid $border-color-dark;

  .footer-right {
    display: flex;
    align-items: center;
    gap: 16px;
  }
}
</style>
