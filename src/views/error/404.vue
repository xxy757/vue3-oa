<template>
  <div class="error-container">
    <div class="error-content">
      <div class="error-code">404</div>
      <div class="error-title">页面不存在</div>
      <div class="error-description">
        抱歉，您访问的页面不存在或已被移除
      </div>
      <n-button type="primary" @click="goHome">
        <template #icon>
          <n-icon :component="HomeOutline" />
        </template>
        返回首页
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { NButton, NIcon } from 'naive-ui'
import { HomeOutline } from '@vicons/ionicons5'

const router = useRouter()

function goHome() {
  router.push('/dashboard')
}
</script>

<style lang="scss" scoped>
.error-container {
  width: 100%;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: $bg-color-3;
}

.error-content {
  text-align: center;
}

.error-code {
  font-size: 72px;
  font-weight: 600;
  color: $primary-color;
  line-height: 1;
  margin-bottom: 16px;
}

.error-title {
  font-size: 20px;
  font-weight: 600;
  color: $text-color-1;
  margin-bottom: 8px;
}

.error-description {
  font-size: 14px;
  color: $text-color-3;
  margin-bottom: 32px;
}
</style>
