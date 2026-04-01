import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginForm, LoginResult, TenantInfo } from '@/types/user'
import { request } from '@/utils/request'
import {
  getToken,
  setToken,
  getStoredUser,
  setStoredUser,
  clearAuth,
  setTenantSlug,
  getStoredTenant,
  setStoredTenant
} from '@/utils/storage'

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(getToken())
  const userInfo = ref<User | null>(getStoredUser<User>())
  const tenantInfo = ref<TenantInfo | null>(getStoredTenant<TenantInfo>())
  const isLoggedIn = computed(() => !!token.value)
  const userName = computed(() => userInfo.value?.nickname || userInfo.value?.username || '')
  const userAvatar = computed(() => userInfo.value?.avatar || '')
  const userRole = computed(() => userInfo.value?.roleName || '')
  const userPermissions = computed(() => userInfo.value?.permissions || [])
  const tenantName = computed(() => tenantInfo.value?.name || '')
  const tenantSlug = computed(() => tenantInfo.value?.slug || '')
  const tenantStatus = computed(() => tenantInfo.value?.status || '')
  const tenantPlan = computed(() => tenantInfo.value?.plan || null)

  const hasPermission = computed(() => (perm: string) => {
    if (!userInfo.value?.permissions) return false
    const permissions = userInfo.value.permissions
    return (
      permissions.includes('*') ||
      permissions.includes(perm) ||
      permissions.some((p) => p.startsWith(perm))
    )
  })

  async function login(loginForm: LoginForm): Promise<void> {
    const result: LoginResult = await request.post('/auth/login', loginForm)
    token.value = result.token
    userInfo.value = result.user
    tenantInfo.value = result.tenant
    setToken(result.token)
    setStoredUser(result.user)
    setStoredTenant(result.tenant)
    if (result.tenant?.slug) {
      setTenantSlug(result.tenant.slug)
    }
  }

  function logout(): void {
    token.value = null
    userInfo.value = null
    tenantInfo.value = null
    clearAuth()
  }

  async function getUserInfo(): Promise<User> {
    const data: any = await request.get('/auth/info')
    userInfo.value = data
    setStoredUser(data)
    if (data.tenant) {
      tenantInfo.value = data.tenant
      setStoredTenant(data.tenant)
    }
    return data
  }

  async function changePassword(oldPassword: string, newPassword: string): Promise<void> {
    await request.put('/auth/password', { oldPassword, newPassword })
  }

  function updateUserInfo(info: Partial<User>): void {
    if (userInfo.value) {
      userInfo.value = { ...userInfo.value, ...info }
      setStoredUser(userInfo.value)
    }
  }

  return {
    token,
    userInfo,
    tenantInfo,
    isLoggedIn,
    userName,
    userAvatar,
    userRole,
    userPermissions,
    tenantName,
    tenantSlug,
    tenantStatus,
    tenantPlan,
    hasPermission,
    login,
    logout,
    getUserInfo,
    changePassword,
    updateUserInfo
  }
})
