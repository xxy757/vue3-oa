<template>
  <n-layout has-sider class="admin-layout">
    <n-layout-sider
      bordered
      :collapsed="collapsed"
      collapse-mode="width"
      :collapsed-width="64"
      :width="220"
      show-trigger
      @collapse="collapsed = true"
      @expand="collapsed = false"
    >
      <div class="logo">
        <img src="@/assets/logo.svg" alt="Logo" class="logo-icon" />
        <span v-show="!collapsed" class="logo-text">SaaS 管理平台</span>
      </div>
      <n-menu
        :collapsed="collapsed"
        :collapsed-width="64"
        :collapsed-icon-size="22"
        :options="menuOptions"
        :value="activeKey"
        @update:value="handleMenuSelect"
      />
    </n-layout-sider>

    <n-layout>
      <n-layout-header bordered class="admin-header">
        <div class="header-left">
          <n-breadcrumb>
            <n-breadcrumb-item v-for="item in breadcrumbs" :key="item.path || item.title">
              {{ item.title }}
            </n-breadcrumb-item>
          </n-breadcrumb>
        </div>
        <div class="header-right">
          <n-dropdown :options="userOptions" @select="handleUserSelect">
            <div class="user-info">
              <n-avatar round :src="userAvatar" :size="32" />
              <span class="user-name">{{ userName }}</span>
            </div>
          </n-dropdown>
        </div>
      </n-layout-header>

      <n-layout-content class="admin-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  NLayout,
  NLayoutSider,
  NLayoutHeader,
  NLayoutContent,
  NMenu,
  NBreadcrumb,
  NBreadcrumbItem,
  NAvatar,
  NDropdown,
  NIcon
} from 'naive-ui'
import {
  HomeOutline,
  BusinessOutline,
  PersonOutline,
  LogOutOutline,
  PricetagOutline
} from '@vicons/ionicons5'
import type { MenuOption } from 'naive-ui'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const collapsed = ref(false)

const userName = computed(() => userStore.userName)
const userAvatar = computed(() => userStore.userAvatar)
const activeKey = computed(() => route.path)

const breadcrumbs = computed(() => {
  const matched = route.matched.filter(item => item.meta && item.meta.title)
  return matched.map(item => ({
    title: item.meta.title as string,
    path: item.path
  }))
})

const renderIcon = (icon: typeof HomeOutline) => {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuOptions: MenuOption[] = [
  {
    label: '管理概览',
    key: '/admin/dashboard',
    icon: renderIcon(HomeOutline)
  },
  {
    label: '租户管理',
    key: '/admin/tenants',
    icon: renderIcon(BusinessOutline)
  },
  {
    label: '套餐管理',
    key: '/admin/plans',
    icon: renderIcon(PricetagOutline)
  }
]

const userOptions = [
  { label: '个人信息', key: 'profile', icon: renderIcon(PersonOutline) },
  { type: 'divider', key: 'd1' },
  { label: '退出登录', key: 'logout', icon: renderIcon(LogOutOutline) }
]

function handleMenuSelect(key: string) {
  if (key.startsWith('/')) {
    router.push(key)
  }
}

function handleUserSelect(key: string) {
  switch (key) {
    case 'profile':
      router.push('/profile/info')
      break
    case 'logout':
      userStore.logout()
      router.push('/login')
      break
  }
}
</script>

<style lang="scss" scoped>
.admin-layout {
  height: 100vh;
  background-color: $bg-color-3;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 16px;
  border-bottom: 1px solid $border-color-dark;
  background: $bg-color-1;

  .logo-icon {
    width: 32px;
    height: 32px;
  }

  .logo-text {
    margin-left: 12px;
    font-size: 16px;
    font-weight: 600;
    color: $primary-color;
    white-space: nowrap;
  }
}

.admin-header {
  height: 60px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: $bg-color-1;
  border-bottom: 1px solid $border-color-dark;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 6px;
  transition: background $transition-duration ease;
  border: 1px solid transparent;

  &:hover {
    background: $bg-color-2;
    border-color: $border-color;
  }

  .user-name {
    margin-left: 8px;
    font-size: 14px;
    color: $text-color-2;
  }
}

.admin-content {
  padding: 24px;
  background-color: $bg-color-3;
  min-height: calc(100vh - 60px);
  overflow-y: auto;
}

:deep(.n-layout-sider) {
  background: $bg-color-1 !important;
  border-right: 1px solid $border-color-dark !important;
}

:deep(.n-menu) {
  .n-menu-item {
    margin: 4px 8px;
    border-radius: 6px;

    &.n-menu-item--selected {
      background: $primary-color-suppl !important;
      color: $primary-color !important;
    }
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
