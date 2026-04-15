<template>
  <n-layout has-sider class="oa-layout">
    <!-- 侧边栏 -->
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
        <span v-show="!collapsed" class="logo-text">企业OA系统</span>
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
      <!-- 头部 -->
      <n-layout-header bordered class="oa-header">
        <div class="header-left">
          <n-breadcrumb>
            <n-breadcrumb-item v-for="item in breadcrumbs" :key="item.path || item.title">
              {{ item.title }}
            </n-breadcrumb-item>
          </n-breadcrumb>
        </div>

        <div class="header-right">
          <!-- 待办提醒 -->
          <n-badge :value="pendingCount" :max="99" class="header-badge">
            <n-button quaternary circle @click="router.push('/approval/pending')">
              <template #icon>
                <n-icon><MailOutline /></n-icon>
              </template>
            </n-button>
          </n-badge>

          <!-- 公告提醒 -->
          <n-badge :value="noticeCount" :max="99" class="header-badge">
            <n-button quaternary circle @click="router.push('/notice/list')">
              <template #icon>
                <n-icon><MegaphoneOutline /></n-icon>
              </template>
            </n-button>
          </n-badge>

          <!-- 用户下拉 -->
          <n-dropdown :options="userOptions" @select="handleUserSelect">
            <div class="user-info">
              <n-avatar round :src="userAvatar" :size="32" />
              <span class="user-name">{{ userName }}</span>
            </div>
          </n-dropdown>
        </div>
      </n-layout-header>

      <!-- 内容区 -->
      <n-layout-content class="oa-content">
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
  import { ref, computed, h, onMounted } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import {
    NLayout,
    NLayoutSider,
    NLayoutHeader,
    NLayoutContent,
    NMenu,
    NBreadcrumb,
    NBreadcrumbItem,
    NButton,
    NIcon,
    NBadge,
    NAvatar,
    NDropdown
  } from 'naive-ui'
  import {
    HomeOutline,
    DocumentTextOutline,
    MegaphoneOutline,
    CalendarOutline,
    SettingsOutline,
    PersonOutline,
    MailOutline,
    LogOutOutline,
    AddCircleOutline,
    PaperPlaneOutline,
    ListOutline,
    CheckmarkDoneOutline,
    PeopleOutline,
    GitBranchOutline,
    ShieldCheckmarkOutline,
    GitCompareOutline,
    LockClosedOutline,
    BusinessOutline,
    PricetagOutline,
    ReceiptOutline
  } from '@vicons/ionicons5'
  import type { MenuOption } from 'naive-ui'
  import { useUserStore } from '@/stores/user'

  const router = useRouter()
  const route = useRoute()
  const userStore = useUserStore()
  const collapsed = ref(false)
  const pendingCount = ref(5)
  const noticeCount = ref(3)

  const userName = computed(() => userStore.userName)
  const userAvatar = computed(() => userStore.userAvatar)
  const breadcrumbs = computed(() => {
    const matched = route.matched.filter((item) => item.meta && item.meta.title)
    return matched.map((item) => ({
      title: item.meta.title as string,
      path: item.path
    }))
  })

  const activeKey = computed(() => route.path)

  // 渲染图标
  const renderIcon = (icon: typeof HomeOutline) => {
    return () => h(NIcon, null, { default: () => h(icon) })
  }

  // 菜单配置
  const menuOptions: MenuOption[] = [
    {
      label: '工作台',
      key: '/dashboard',
      icon: renderIcon(HomeOutline)
    },
    {
      label: '审批中心',
      key: 'approval',
      icon: renderIcon(DocumentTextOutline),
      children: [
        { label: '发起申请', key: '/approval/apply', icon: renderIcon(AddCircleOutline) },
        { label: '我的申请', key: '/approval/my-apply', icon: renderIcon(PaperPlaneOutline) },
        { label: '待我审批', key: '/approval/pending', icon: renderIcon(MailOutline) },
        { label: '已办审批', key: '/approval/done', icon: renderIcon(CheckmarkDoneOutline) }
      ]
    },
    {
      label: '公告通知',
      key: 'notice',
      icon: renderIcon(MegaphoneOutline),
      children: [{ label: '公告列表', key: '/notice/list', icon: renderIcon(ListOutline) }]
    },
    {
      label: '日程管理',
      key: 'schedule',
      icon: renderIcon(CalendarOutline),
      children: [
        { label: '日程日历', key: '/schedule/calendar', icon: renderIcon(CalendarOutline) },
        { label: '日程列表', key: '/schedule/list', icon: renderIcon(ListOutline) }
      ]
    },
    {
      label: '系统管理',
      key: 'system',
      icon: renderIcon(SettingsOutline),
      children: [
        { label: '用户管理', key: '/system/user', icon: renderIcon(PeopleOutline) },
        { label: '部门管理', key: '/system/dept', icon: renderIcon(GitBranchOutline) },
        { label: '角色管理', key: '/system/role', icon: renderIcon(ShieldCheckmarkOutline) },
        { label: '流程配置', key: '/system/flow', icon: renderIcon(GitCompareOutline) }
      ]
    },
    {
      label: '企业管理',
      key: 'tenant',
      icon: renderIcon(BusinessOutline),
      children: [
        { label: '企业信息', key: '/tenant/info', icon: renderIcon(BusinessOutline) },
        { label: '套餐管理', key: '/tenant/plan', icon: renderIcon(PricetagOutline) },
        { label: '账单管理', key: '/tenant/invoices', icon: renderIcon(ReceiptOutline) }
      ]
    }
  ]

  // 用户下拉选项
  const userOptions = [
    { label: '个人信息', key: 'profile', icon: renderIcon(PersonOutline) },
    { label: '修改密码', key: 'password', icon: renderIcon(LockClosedOutline) },
    { type: 'divider', key: 'd1' },
    { label: '退出登录', key: 'logout', icon: renderIcon(LogOutOutline) }
  ]

  // 菜单选择
  const handleMenuSelect = (key: string) => {
    if (key.startsWith('/')) {
      router.push(key)
    }
  }

  // 用户菜单选择
  const handleUserSelect = (key: string) => {
    switch (key) {
      case 'profile':
        router.push('/profile/info')
        break
      case 'password':
        router.push('/profile/password')
        break
      case 'logout':
        userStore.logout()
        router.push('/login')
        break
    }
  }

  onMounted(() => {
    // 可以在这里获取待办数量
  })
</script>

<style lang="scss" scoped>
  .oa-layout {
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

  .oa-header {
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

  .header-badge {
    margin-right: 8px;
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

  .oa-content {
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
