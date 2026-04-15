import { useTenantStore } from '@/stores/tenant'

export const PLAN_LABELS: Record<string, string> = {
  approval: '审批中心',
  schedule: '日程管理',
  notice: '公告通知',
  storage: '存储空间',
  api: 'API 接口',
  sso: '单点登录'
}

export const PLAN_STATUS_MAP: Record<string, { label: string; type: 'success' | 'warning' | 'error' | 'info' | 'default' }> = {
  trial: { label: '试用中', type: 'warning' },
  active: { label: '已激活', type: 'success' },
  suspended: { label: '已暂停', type: 'error' },
  cancelled: { label: '已注销', type: 'default' }
}

export const INVOICE_STATUS_MAP: Record<string, { label: string; type: 'success' | 'warning' | 'error' | 'info' | 'default' }> = {
  pending: { label: '待支付', type: 'warning' },
  paid: { label: '已支付', type: 'success' },
  overdue: { label: '已逾期', type: 'error' },
  cancelled: { label: '已取消', type: 'default' }
}

export function usePlanPermission() {
  const tenantStore = useTenantStore()

  function hasFeature(feature: string): boolean {
    return tenantStore.hasFeature(feature)
  }

  function checkUserLimit(): { exceeded: boolean; current: number; max: number } {
    return tenantStore.checkUserLimit()
  }

  return { hasFeature, checkUserLimit }
}

export function formatPrice(price: number): string {
  return price === 0 ? '免费' : `¥${price.toFixed(2)}/人/月`
}

export function formatStorage(mb: number): string {
  if (mb >= 1024) {
    return `${(mb / 1024).toFixed(mb % 1024 === 0 ? 0 : 1)}GB`
  }
  return `${mb}MB`
}

export function getPlanColor(code: string): string {
  const colorMap: Record<string, string> = {
    free: '#8c8c8c',
    standard: '#1677FF',
    professional: '#722ED1',
    enterprise: '#FA8C16'
  }
  return colorMap[code] || '#1677FF'
}
