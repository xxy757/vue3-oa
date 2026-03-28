<template>
  <div class="login-container">
    <div class="login-background">
      <div class="login-box">
        <n-card class="login-card" :bordered="false">
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
                size="large"
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
                size="large"
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
                size="large"
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
import { useRouter } from 'vue-router'
import { useMessage, type FormInst, type FormRules } from 'naive-ui'
import { PersonOutline, LockClosedOutline } from '@vicons/ionicons5'
import { useUserStore } from '@/stores/user'
import type { LoginForm } from '@/types/user'

const router = useRouter()
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

    router.push('/dashboard')
  } catch (error) {
    if (error instanceof Error) {
      message.error(error.message || '登录失败，请检查用户名和密码')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}

.login-background {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.login-background::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image:
    radial-gradient(circle at 20% 80%, rgba(255, 255, 255, 0.1) 0%, transparent 50%),
    radial-gradient(circle at 80% 20%, rgba(255, 255, 255, 0.1) 0%, transparent 50%);
}

.login-box {
  position: relative;
  z-index: 1;
}

.login-card {
  width: 400px;
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-title {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.login-subtitle {
  font-size: 14px;
  color: #999;
  margin: 0;
  line-height: 1.5;
}

.login-form {
  margin-top: 20px;
}

.login-footer {
  text-align: center;
  margin-top: 20px;
}

.test-account {
  font-size: 12px;
  color: #999;
  margin: 0;
  line-height: 1.5;
}

@media (max-width: 480px) {
  .login-card {
    width: 90%;
    max-width: 360px;
  }

  .login-title {
    font-size: 24px;
  }
}
</style>
