<template>
  <div class="tenant-info-page">
    <div class="page-title-bar">
      <h2 class="page-title">企业信息</h2>
    </div>

    <n-spin :show="loading">
      <n-card>
        <n-form
          ref="formRef"
          :model="formData"
          :rules="formRules"
          label-placement="left"
          label-width="80px"
        >
          <n-grid :cols="2" :x-gap="24">
            <n-form-item-gi label="企业名称" path="name">
              <n-input v-model:value="formData.name" placeholder="请输入企业名称" />
            </n-form-item-gi>
            <n-form-item-gi label="企业标识">
              <n-input :value="tenantStore.tenantInfo?.slug" disabled />
            </n-form-item-gi>
            <n-form-item-gi label="联系人" path="contactName">
              <n-input v-model:value="formData.contactName" placeholder="请输入联系人" />
            </n-form-item-gi>
            <n-form-item-gi label="联系电话" path="contactPhone">
              <n-input v-model:value="formData.contactPhone" placeholder="请输入联系电话" />
            </n-form-item-gi>
            <n-form-item-gi label="联系邮箱" path="contactEmail">
              <n-input v-model:value="formData.contactEmail" placeholder="请输入联系邮箱" />
            </n-form-item-gi>
            <n-form-item-gi label="当前套餐">
              <n-tag :type="planStatusType" size="small">
                {{ tenantStore.currentPlan?.name || '-' }}
              </n-tag>
            </n-form-item-gi>
            <n-form-item-gi label="用户数">
              <span class="user-count">
                {{ tenantStore.currentUsers }} / {{ tenantStore.maxUsers }}
              </span>
              <n-progress
                type="line"
                :percentage="userPercentage"
                :show-indicator="false"
                style="width: 120px; margin-left: 12px"
                :color="userPercentage > 90 ? '#FF4D4F' : '#1677FF'"
              />
            </n-form-item-gi>
            <n-form-item-gi label="到期时间">
              <span class="expire-date">{{ tenantStore.tenantInfo?.planExpireAt?.substring(0, 10) || '-' }}</span>
            </n-form-item-gi>
          </n-grid>

          <div class="form-actions">
            <n-button type="primary" :loading="saving" @click="handleSave">保存修改</n-button>
          </div>
        </n-form>
      </n-card>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  NCard,
  NForm,
  NFormItemGi,
  NGrid,
  NInput,
  NButton,
  NTag,
  NProgress,
  NSpin,
  useMessage
} from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import { useTenantStore } from '@/stores/tenant'
import type { TenantUpdateForm } from '@/types/tenant'

const message = useMessage()
const tenantStore = useTenantStore()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const saving = ref(false)

const formData = ref<TenantUpdateForm>({
  name: '',
  contactName: '',
  contactPhone: '',
  contactEmail: ''
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入企业名称', trigger: 'blur' }],
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

const planStatusType = computed(() => {
  if (tenantStore.isExpired) return 'error' as const
  if (tenantStore.isTrial) return 'warning' as const
  return 'success' as const
})

const userPercentage = computed(() => {
  if (!tenantStore.maxUsers) return 0
  return Math.round((tenantStore.currentUsers / tenantStore.maxUsers) * 100)
})

async function handleSave() {
  try {
    await formRef.value?.validate()
    saving.value = true
    await tenantStore.updateTenantInfo(formData.value)
    message.success('企业信息更新成功')
  } catch (err: unknown) {
    if (err instanceof Error) {
      message.error(err.message)
    }
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  loading.value = true
  try {
    const info = await tenantStore.fetchTenantInfo()
    formData.value = {
      name: info.name,
      contactName: info.contactName,
      contactPhone: info.contactPhone,
      contactEmail: info.contactEmail
    }
  } catch {
    message.error('获取企业信息失败')
  } finally {
    loading.value = false
  }
})
</script>

<style lang="scss" scoped>
.tenant-info-page {
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

.user-count {
  font-size: 14px;
  color: $text-color-2;
}

.expire-date {
  font-size: 14px;
  color: $text-color-2;
}

.form-actions {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid $border-color-dark;
}
</style>
