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

const whiteList = ['/login', '/register', '/404']

  router.beforeEach((to, _from, next) => {
    document.title = to.meta.title ? `${to.meta.title} - 企业OA办公系统` : '企业OA办公系统'
    const token = getToken()
    if (token) {
      if (to.path === '/login') {
        next('/dashboard')
      } else {
        next()
      }
    } else {
      if (!to.meta.requiresAuth || whiteList.includes(to.path)) {
        next()
      } else {
        next(`/login?redirect=${encodeURIComponent(to.fullPath)}`)
      }
    }
  })

export default router
