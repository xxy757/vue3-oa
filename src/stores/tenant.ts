import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { TenantDetail, Plan, Invoice, TenantUpdateForm, UpgradeResult } from '@/types/tenant'
import { request } from '@/utils/request'

export const useTenantStore = defineStore('tenant', () => {
  const tenantInfo = ref<TenantDetail | null>(null)
  const plans = ref<Plan[]>([])
  const invoices = ref<Invoice[]>([])
  const loading = ref(false)

  const tenantName = computed(() => tenantInfo.value?.name || '')
  const tenantLogo = computed(() => tenantInfo.value?.logo || '')
  const currentPlan = computed(() => tenantInfo.value?.plan || null)
  const currentUsers = computed(() => tenantInfo.value?.currentUsers || 0)
  const maxUsers = computed(() => tenantInfo.value?.maxUsers || 0)
  const tenantStatus = computed(() => tenantInfo.value?.status || '')
  const planFeatures = computed(() => tenantInfo.value?.plan?.features || {})
  const isTrial = computed(() => tenantStatus.value === 'trial')
  const isExpired = computed(() => {
    if (!tenantInfo.value?.planExpireAt) return false
    return new Date(tenantInfo.value.planExpireAt) < new Date()
  })
  const userLimitReached = computed(() => currentUsers.value >= maxUsers.value)

  async function fetchTenantInfo(): Promise<TenantDetail> {
    loading.value = true
    try {
      const data = await request.get<TenantDetail>('/tenant/info')
      tenantInfo.value = data
      return data
    } finally {
      loading.value = false
    }
  }

  async function updateTenantInfo(form: TenantUpdateForm): Promise<void> {
    await request.put('/tenant/info', form)
    if (tenantInfo.value) {
      tenantInfo.value = { ...tenantInfo.value, ...form }
    }
  }

  async function fetchPlans(): Promise<Plan[]> {
    const data = await request.get<Plan[]>('/plans')
    plans.value = data
    return data
  }

  async function upgradePlan(planId: number): Promise<UpgradeResult> {
    const data = await request.post<UpgradeResult>('/tenant/plan/upgrade', { planId })
    await fetchTenantInfo()
    return data
  }

  async function fetchInvoices(): Promise<Invoice[]> {
    const data = await request.get<Invoice[]>('/tenant/invoices')
    invoices.value = data
    return data
  }

  function hasFeature(feature: string): boolean {
    const features = planFeatures.value
    if (!features) return false
    return features[feature] === true
  }

  function checkUserLimit(): { exceeded: boolean; current: number; max: number } {
    return {
      exceeded: userLimitReached.value,
      current: currentUsers.value,
      max: maxUsers.value
    }
  }

  return {
    tenantInfo,
    plans,
    invoices,
    loading,
    tenantName,
    tenantLogo,
    currentPlan,
    currentUsers,
    maxUsers,
    tenantStatus,
    planFeatures,
    isTrial,
    isExpired,
    userLimitReached,
    fetchTenantInfo,
    updateTenantInfo,
    fetchPlans,
    upgradePlan,
    fetchInvoices,
    hasFeature,
    checkUserLimit
  }
})
