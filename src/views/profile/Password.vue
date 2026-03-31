<template>
  <div class="change-password">
    <n-card title="修改密码" size="small" style="max-width: 500px">
      <n-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-placement="left"
        label-width="100"
      >
        <n-form-item label="旧密码" path="oldPassword">
          <n-input
            v-model:value="formData.oldPassword"
            type="password"
            placeholder="请输入旧密码"
            show-password-on="click"
          />
        </n-form-item>
        <n-form-item label="新密码" path="newPassword">
          <n-input
            v-model:value="formData.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password-on="click"
          />
        </n-form-item>
        <n-form-item label="确认新密码" path="confirmPassword">
          <n-input
            v-model:value="formData.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password-on="click"
          />
        </n-form-item>

        <!-- 密码强度提示 -->
        <n-form-item label="密码强度">
          <div class="password-strength">
            <n-progress
              :percentage="strengthPercentage"
              :status="strengthStatus"
              :show-indicator="false"
              :height="8"
            />
            <n-text :type="strengthTextType" style="margin-top: 8px">
              {{ strengthText }}
            </n-text>
          </div>
        </n-form-item>

        <n-form-item label="密码要求">
          <n-text depth="3">
            <ul class="password-tips">
              <li>密码长度至少6位</li>
              <li>建议包含大小写字母、数字和特殊字符</li>
              <li>不要使用与用户名相同的密码</li>
              <li>不要使用过于简单的密码（如123456）</li>
            </ul>
          </n-text>
        </n-form-item>

        <n-form-item>
          <n-space>
            <n-button type="primary" :loading="loading" @click="handleSubmit"> 确认修改 </n-button>
            <n-button @click="handleReset">重置</n-button>
          </n-space>
        </n-form-item>
      </n-form>
    </n-card>

    <!-- 安全提示 -->
    <n-card title="安全提示" size="small" style="max-width: 500px; margin-top: 16px">
      <n-alert type="info">
        <template #header>密码安全须知</template>
        <ul class="security-tips">
          <li>请定期更换密码，建议每3个月更换一次</li>
          <li>不要在多个网站使用相同的密码</li>
          <li>不要将密码告知他人</li>
          <li>如果发现账号异常，请立即修改密码</li>
        </ul>
      </n-alert>
    </n-card>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, computed } from 'vue'
  import { useRouter } from 'vue-router'
  import {
    NCard,
    NForm,
    NFormItem,
    NInput,
    NButton,
    NSpace,
    NProgress,
    NText,
    NAlert,
    useMessage
  } from 'naive-ui'
  import type { FormInst, FormRules, FormItemRule } from 'naive-ui'
  import { useUserStore } from '@/stores/user'

  const router = useRouter()
  const message = useMessage()
  const userStore = useUserStore()

  // 状态
  const loading = ref(false)
  const formRef = ref<FormInst | null>(null)

  // 表单数据
  const formData = reactive({
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  })

  // 密码强度计算
  const strengthPercentage = computed(() => {
    const password = formData.newPassword
    if (!password) return 0

    let score = 0

    // 长度
    if (password.length >= 6) score += 20
    if (password.length >= 10) score += 10

    // 包含小写字母
    if (/[a-z]/.test(password)) score += 15

    // 包含大写字母
    if (/[A-Z]/.test(password)) score += 15

    // 包含数字
    if (/[0-9]/.test(password)) score += 20

    // 包含特殊字符
    if (/[!@#$%^&*(),.?":{}|<>]/.test(password)) score += 20

    return Math.min(score, 100)
  })

  const strengthStatus = computed(() => {
    const percentage = strengthPercentage.value
    if (percentage < 40) return 'error'
    if (percentage < 70) return 'warning'
    return 'success'
  })

  const strengthText = computed(() => {
    const percentage = strengthPercentage.value
    if (percentage === 0) return '请输入密码'
    if (percentage < 40) return '弱 - 建议增加复杂度'
    if (percentage < 70) return '中 - 可以更复杂'
    return '强 - 密码安全'
  })

  const strengthTextType = computed(() => {
    const percentage = strengthPercentage.value
    if (percentage < 40) return 'error'
    if (percentage < 70) return 'warning'
    return 'success'
  })

  // 确认密码验证
  const validateConfirmPassword = (_rule: FormItemRule, value: string): boolean => {
    return value === formData.newPassword
  }

  // 表单验证规则
  const formRules: FormRules = {
    oldPassword: { required: true, message: '请输入旧密码', trigger: 'blur' },
    newPassword: [
      { required: true, message: '请输入新密码', trigger: 'blur' },
      { min: 6, message: '密码长度至少6位', trigger: 'blur' }
    ],
    confirmPassword: [
      { required: true, message: '请再次输入新密码', trigger: 'blur' },
      { validator: validateConfirmPassword, message: '两次输入的密码不一致', trigger: 'blur' }
    ]
  }

  // 重置表单
  function handleReset() {
    formData.oldPassword = ''
    formData.newPassword = ''
    formData.confirmPassword = ''
    formRef.value?.restoreValidation()
  }

  // 提交表单
  async function handleSubmit() {
    try {
      await formRef.value?.validate()
      loading.value = true

      // 调用修改密码接口
      await userStore.changePassword(formData.oldPassword, formData.newPassword)

      message.success('密码修改成功，请重新登录')

      // 退出登录并跳转到登录页
      setTimeout(() => {
        userStore.logout()
        router.push('/login')
      }, 1500)
    } catch (error) {
      if (error) {
        message.error('密码修改失败')
      }
    } finally {
      loading.value = false
    }
  }
</script>

<style lang="scss" scoped>
  .change-password {
    .password-strength {
      width: 100%;
    }

    .password-tips,
    .security-tips {
      margin: 0;
      padding-left: 20px;
      line-height: 1.8;
    }

    .security-tips {
      list-style-type: disc;
    }
  }
</style>
