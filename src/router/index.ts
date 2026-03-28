import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { getToken } from '@/utils/storage'
import { constantRoutes, asyncRoutes } from './routes'

const routes: RouteRecordRaw[] = [...constantRoutes, ...asyncRoutes]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

// 白名单路由
const whiteList = ['/login', '/404']

router.beforeEach((to, _from, next) => {
  // 设置页面标题
  document.title = `${to.meta.title || ''} - 企业OA办公系统`

  const token = getToken()

  if (token) {
    if (to.path === '/login') {
      // 已登录，跳转到首页
      next('/dashboard')
    } else {
      next()
    }
  } else {
    if (whiteList.includes(to.path)) {
      // 在白名单中，直接进入
      next()
    } else {
      // 未登录，跳转到登录页
      next('/login')
    }
  }
})

export default router
