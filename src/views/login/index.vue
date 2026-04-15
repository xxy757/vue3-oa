<template>
  <div class="login-container">
    <div class="login-background">
      <div class="login-box">
        <n-card class="login-card" :bordered="false" :content-style="{ padding: '20px' }">
          <div class="login-header">
            <h1 class="login-title">企业OA办公系统</h1>
            <p class="login-subtitle">Enterprise Office Automation System</p>
          </div>

          <n-form
            ref="formRef"
            :model="formValue"
            :rules="rules"
            class="login-form"
          >
            <n-form-item path="username">
              <n-input
                v-model:value="formValue.username"
                placeholder="请输入用户名"
                :input-props="{ autocomplete: 'username' }"
              >
                <template #prefix>
                  <n-icon :component="PersonOutline" />
                </template>
              </n-input>
            </n-form-item>

            <n-form-item path="password">
              <n-input
                v-model:value="formValue.password"
                type="password"
                placeholder="请输入密码"
                show-password-on="click"
                :input-props="{ autocomplete: 'current-password' }"
                @keyup.enter="handleLogin"
              >
                <template #prefix>
                  <n-icon :component="LockClosedOutline" />
                </template>
              </n-input>
            </n-form-item>

            <n-form-item path="remember">
              <n-checkbox v-model:checked="formValue.remember">
                记住我
              </n-checkbox>
            </n-form-item>

            <n-form-item>
              <n-button
                type="primary"
                block
                :loading="loading"
                @click="handleLogin"
              >
                登录
              </n-button>
            </n-form-item>
          </n-form>

          <div class="login-footer">
            <p class="test-account">测试账号: admin/123456, user/123456, manager/123456</p>
          </div>
        </n-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  NCard,
  NForm,
  NFormItem,
  NInput,
  NButton,
  NCheckbox,
  NIcon,
  useMessage,
  type FormInst,
  type FormRules
} from 'naive-ui'
import { PersonOutline, LockClosedOutline } from '@vicons/ionicons5'
import { useUserStore } from '@/stores/user'
import type { LoginForm } from '@/types/user'

const router = useRouter()
const route = useRoute()
const message = useMessage()
const userStore = useUserStore()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)

const formValue = reactive<LoginForm>({
  username: '',
  password: '',
  remember: false
})

const rules: FormRules = {
  username: {
    required: true,
    message: '请输入用户名',
    trigger: ['blur', 'input']
  },
  password: {
    required: true,
    message: '请输入密码',
    trigger: ['blur', 'input']
  }
}

async function handleLogin() {
  try {
    await formRef.value?.validate()

    loading.value = true

    await userStore.login(formValue)

    message.success('登录成功')

    const redirect = (route.query.redirect as string) || '/dashboard'
    router.push(redirect)
  } catch (error) {
    if (error instanceof Error) {
      message.error(error.message || '登录失败，请检查用户名和密码')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.login-container {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}

.login-background {
  width: 100%;
  height: 100%;
  background: $bg-color-3;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-box {
  width: 400px;
  max-width: 90vw;
}

.login-card {
  border-radius: $border-radius;
  box-shadow: $box-shadow;
}

.login-header {
  text-align: center;
  margin-bottom: 24px;
}

.login-title {
  font-size: 20px;
  font-weight: 600;
  color: $text-color-1;
  margin: 0 0 8px 0;
  line-height: 28px;
}

.login-subtitle {
  font-size: 12px;
  color: $text-color-3;
  margin: 0;
  line-height: 20px;
}

.login-form {
  margin-top: 20px;
}

.login-footer {
  text-align: center;
  margin-top: 16px;
}

.test-account {
  font-size: 12px;
  color: $text-color-3;
  margin: 0;
  line-height: 20px;
}

@media (max-width: 480px) {
  .login-box {
    width: 90%;
    max-width: 360px;
  }
}
</style>
