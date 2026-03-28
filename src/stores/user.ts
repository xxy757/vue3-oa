import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginForm, LoginResult } from '@/types/user'
import { request } from '@/utils/request'
import { getToken, setToken, removeToken, getStoredUser, setStoredUser, clearAuth } from '@/utils/storage'

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(getToken())
  const userInfo = ref<User | null>(getStoredUser<User>())

  const isLoggedIn = computed(() => !!token.value)
  const userName = computed(() => userInfo.value?.nickname || userInfo.value?.username || '')
  const userAvatar = computed(() => userInfo.value?.avatar || '')
  const userRole = computed(() => userInfo.value?.roleName || '')

  // 登录
  async function login(loginForm: LoginForm): Promise<void> {
    const result: LoginResult = await request.post('/auth/login', loginForm)
    token.value = result.token
    userInfo.value = result.user
    setToken(result.token)
    setStoredUser(result.user)
  }

  // 退出登录
  function logout(): void {
    token.value = null
    userInfo.value = null
    clearAuth()
  }

  // 获取用户信息
  async function getUserInfo(): Promise<User> {
    const user: User = await request.get('/auth/info')
    userInfo.value = user
    setStoredUser(user)
    return user
  }

  // 修改密码
  async function changePassword(oldPassword: string, newPassword: string): Promise<void> {
    await request.put('/auth/password', { oldPassword, newPassword })
  }

  // 更新用户信息
  function updateUserInfo(info: Partial<User>): void {
    if (userInfo.value) {
      userInfo.value = { ...userInfo.value, ...info }
      setStoredUser(userInfo.value)
    }
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    userName,
    userAvatar,
    userRole,
    login,
    logout,
    getUserInfo,
    changePassword,
    updateUserInfo
  }
})
