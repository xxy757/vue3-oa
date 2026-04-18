<template>
  <div class="user-layout">
    <div class="user-layout-content">
      <div class="top">
        <div class="header">
          <div class="logo">
            <n-icon :size="32" color="#1677FF"><BusinessOutline /></n-icon>
          </div>
          <span class="title">Ant OA</span>
        </div>
        <div class="desc">
          企业级智能协同办公平台
        </div>
      </div>

      <div class="main">
        <n-tabs
          v-model:value="activeTab"
          type="line"
          size="large"
          :tab-style="{ padding: '12px 0', fontSize: '16px' }"
          class="login-tabs"
        >
          <n-tab-pane name="account" tab="账号登录">
            <n-form
              ref="formRef"
              :model="formValue"
              :rules="rules"
              class="login-form"
            >
              <n-form-item path="tenantSlug">
                <n-input
                  v-model:value="formValue.tenantSlug"
                  placeholder="企业标识: demo"
                  size="large"
                >
                  <template #prefix>
                    <n-icon :component="BusinessOutline" :size="16" />
                  </template>
                </n-input>
              </n-form-item>

              <n-form-item path="username">
                <n-input
                  v-model:value="formValue.username"
                  placeholder="用户名: admin"
                  size="large"
                  :input-props="{ autocomplete: 'username' }"
                >
                  <template #prefix>
                    <n-icon :component="PersonOutline" :size="16" />
                  </template>
                </n-input>
              </n-form-item>

              <n-form-item path="password">
                <n-input
                  v-model:value="formValue.password"
                  type="password"
                  placeholder="密码: 123456"
                  size="large"
                  show-password-on="click"
                  :input-props="{ autocomplete: 'current-password' }"
                  @keyup.enter="handleLogin"
                >
                  <template #prefix>
                    <n-icon :component="LockClosedOutline" :size="16" />
                  </template>
                </n-input>
              </n-form-item>

              <div class="form-actions">
                <n-checkbox v-model:checked="formValue.remember">
                  自动登录
                </n-checkbox>
                <n-button text size="small" type="primary">忘记密码</n-button>
              </div>

              <n-button
                type="primary"
                block
                size="large"
                :loading="loading"
                @click="handleLogin"
                class="submit-btn"
              >
                登 录
              </n-button>
            </n-form>
          </n-tab-pane>

          <n-tab-pane name="phone" tab="手机号登录">
            <div class="phone-placeholder">
              <n-form class="login-form">
                <n-form-item>
                  <n-input placeholder="手机号" size="large">
                    <template #prefix>
                      <n-icon :component="MailOutline" :size="16" />
                    </template>
                  </n-input>
                </n-form-item>
                <n-form-item>
                  <n-input placeholder="验证码" size="large">
                    <template #prefix>
                      <n-icon :component="ShieldCheckmarkOutline" :size="16" />
                    </template>
                    <template #suffix>
                      <n-button text type="primary" size="small">获取验证码</n-button>
                    </template>
                  </n-input>
                </n-form-item>
                <n-button
                  block
                  size="large"
                  class="submit-btn"
                >
                  登 录
                </n-button>
              </n-form>
            </div>
          </n-tab-pane>
        </n-tabs>

        <div class="other-login">
          <n-divider class="other-divider">其他登录方式</n-divider>
          <div class="other-icons">
            <n-icon :size="24" class="icon-item"><MegaphoneOutline /></n-icon>
            <n-icon :size="24" class="icon-item"><SettingsOutline /></n-icon>
            <n-icon :size="24" class="icon-item"><MailOutline /></n-icon>
          </div>
        </div>
      </div>
    </div>

    <div class="footer">
      <div class="links">
        <span>帮助</span>
        <span>隐私</span>
        <span>条款</span>
      </div>
      <div class="copyright">Copyright © 2026 Ant OA Team</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  NForm,
  NFormItem,
  NInput,
  NButton,
  NCheckbox,
  NIcon,
  NTabs,
  NTabPane,
  NDivider,
  useMessage,
  type FormInst,
  type FormRules
} from 'naive-ui'
import {
  PersonOutline,
  LockClosedOutline,
  BusinessOutline,
  ShieldCheckmarkOutline,
  MailOutline,
  MegaphoneOutline,
  SettingsOutline
} from '@vicons/ionicons5'
import { useUserStore } from '@/stores/user'
import { setTenantSlug } from '@/utils/storage'
import type { LoginForm } from '@/types/user'

const router = useRouter()
const route = useRoute()
const message = useMessage()
const userStore = useUserStore()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const activeTab = ref('account')

const formValue = reactive<LoginForm>({
  username: '',
  password: '',
  tenantSlug: '',
  remember: false
})

const rules: FormRules = {
  tenantSlug: {
    required: true,
    message: '请输入企业标识',
    trigger: ['blur', 'input']
  },
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
    if (formValue.tenantSlug) {
      setTenantSlug(formValue.tenantSlug)
    }
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
.user-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: $bg-color-1;
  padding: 32px 0 24px;
}

.user-layout-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 24px;
}

.top {
  text-align: center;
  margin-bottom: 40px;
}

.header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 12px;
}

.logo {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.title {
  font-size: 30px;
  font-weight: 600;
  color: $text-color-1;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'PingFang SC', 'Microsoft YaHei', sans-serif;
  letter-spacing: 1px;
}

.desc {
  font-size: 14px;
  color: $text-color-3;
  margin-top: 8px;
  letter-spacing: 1px;
}

.main {
  width: 368px;
  max-width: 100%;
}

.login-tabs {
  :deep(.n-tabs-nav) {
    justify-content: center;
    margin-bottom: 24px;
  }
}

.login-form {
  margin-top: 8px;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.submit-btn {
  height: 40px;
  font-size: 16px;
  font-weight: 500;
}

.other-login {
  margin-top: 24px;
}

.other-divider {
  margin: 16px 0;

  :deep(.n-divider__title) {
    font-size: 12px;
    color: $text-color-4;
  }
}

.other-icons {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;
}

.icon-item {
  cursor: pointer;
  color: $text-color-4;
  transition: color $transition-duration ease;

  &:hover {
    color: $primary-color;
  }
}

.footer {
  text-align: center;
  padding: 0 16px;

  .links {
    display: flex;
    justify-content: center;
    gap: 24px;
    margin-bottom: 8px;

    span {
      font-size: 14px;
      color: $text-color-3;
      cursor: pointer;
      transition: color $transition-duration ease;

      &:hover {
        color: $primary-color;
      }
    }
  }

  .copyright {
    font-size: 12px;
    color: $text-color-4;
  }
}

@media (max-width: 480px) {
  .main {
    width: 100%;
    padding: 0 16px;
  }

  .title {
    font-size: 24px;
  }
}
</style>
