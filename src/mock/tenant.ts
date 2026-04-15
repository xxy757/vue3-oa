import type { MockMethod } from 'vite-plugin-mock'

interface MockPlan {
  id: number
  name: string
  code: string
  price: number
  minUsers: number
  maxUsers: number
  features: Record<string, boolean | number>
  isActive: number
  createTime: string
  updateTime: string
}

interface MockTenant {
  id: number
  name: string
  slug: string
  logo: string
  contactName: string
  contactPhone: string
  contactEmail: string
  currentUsers: number
  maxUsers: number
  status: string
  trialEndsAt: string | null
  planExpireAt: string | null
  plan: MockPlan
  createTime: string
  updateTime: string
}

const plans: MockPlan[] = [
  {
    id: 1,
    name: '免费版',
    code: 'free',
    price: 0,
    minUsers: 1,
    maxUsers: 5,
    features: { approval: true, schedule: false, notice: false, storage: 100, api: false, sso: false },
    isActive: 1,
    createTime: '2025-01-01T00:00:00Z',
    updateTime: '2025-01-01T00:00:00Z'
  },
  {
    id: 2,
    name: '标准版',
    code: 'standard',
    price: 29,
    minUsers: 1,
    maxUsers: 50,
    features: { approval: true, schedule: true, notice: true, storage: 1024, api: false, sso: false },
    isActive: 1,
    createTime: '2025-01-01T00:00:00Z',
    updateTime: '2025-01-01T00:00:00Z'
  },
  {
    id: 3,
    name: '专业版',
    code: 'professional',
    price: 59,
    minUsers: 5,
    maxUsers: 200,
    features: { approval: true, schedule: true, notice: true, storage: 5120, api: true, sso: false },
    isActive: 1,
    createTime: '2025-01-01T00:00:00Z',
    updateTime: '2025-01-01T00:00:00Z'
  },
  {
    id: 4,
    name: '企业版',
    code: 'enterprise',
    price: 99,
    minUsers: 10,
    maxUsers: 999,
    features: { approval: true, schedule: true, notice: true, storage: 10240, api: true, sso: true },
    isActive: 1,
    createTime: '2025-01-01T00:00:00Z',
    updateTime: '2025-01-01T00:00:00Z'
  }
]

const tenants: MockTenant[] = [
  {
    id: 1,
    name: '示例科技有限公司',
    slug: 'demo-corp',
    logo: '',
    contactName: '张三',
    contactPhone: '13800138001',
    contactEmail: 'admin@demo.com',
    currentUsers: 12,
    maxUsers: 50,
    status: 'active',
    trialEndsAt: null,
    planExpireAt: '2026-06-30T23:59:59Z',
    plan: plans[1],
    createTime: '2025-01-15T10:00:00Z',
    updateTime: '2025-03-01T08:30:00Z'
  }
]

const invoices = [
  {
    id: 1,
    tenantId: 1,
    planId: 2,
    invoiceNo: 'INV-2025-001',
    periodStart: '2025-03-01T00:00:00Z',
    periodEnd: '2025-04-01T00:00:00Z',
    userCount: 12,
    amount: 348.00,
    status: 'paid',
    paidAt: '2025-03-01T09:15:00Z',
    paymentMethod: 'bank_transfer',
    paymentTransactionId: 'TXN-20250301-001',
    createTime: '2025-03-01T00:00:00Z',
    updateTime: '2025-03-01T09:15:00Z'
  },
  {
    id: 2,
    tenantId: 1,
    planId: 2,
    invoiceNo: 'INV-2025-002',
    periodStart: '2025-04-01T00:00:00Z',
    periodEnd: '2025-05-01T00:00:00Z',
    userCount: 12,
    amount: 348.00,
    status: 'pending',
    paidAt: null,
    paymentMethod: '',
    paymentTransactionId: '',
    createTime: '2025-04-01T00:00:00Z',
    updateTime: '2025-04-01T00:00:00Z'
  }
]

let nextTenantId = 2
let nextInvoiceId = 3

export default [
  {
    url: '/api/v1/plans',
    method: 'get',
    response: () => ({
      code: 200,
      data: plans.filter(p => p.isActive === 1),
      message: 'ok'
    })
  },
  {
    url: '/api/v1/tenant/register',
    method: 'post',
    response: ({ body }: { body: Record<string, unknown> }) => {
      const newTenant = {
        id: nextTenantId++,
        name: String(body.name || ''),
        slug: String(body.slug || ''),
        adminUser: {
          id: nextTenantId * 100,
          username: 'admin',
          tempPassword: 'Temp' + Math.random().toString(36).substring(2, 10)
        },
        trialEndsAt: new Date(Date.now() + 14 * 24 * 3600 * 1000).toISOString()
      }
      return {
        code: 200,
        data: newTenant,
        message: 'ok'
      }
    }
  },
  {
    url: '/api/v1/tenant/info',
    method: 'get',
    response: () => ({
      code: 200,
      data: tenants[0],
      message: 'ok'
    })
  },
  {
    url: '/api/v1/tenant/info',
    method: 'put',
    response: ({ body }: { body: Record<string, unknown> }) => {
      Object.assign(tenants[0], body)
      return { code: 200, data: null, message: 'ok' }
    }
  },
  {
    url: '/api/v1/tenant/plan/upgrade',
    method: 'post',
    response: ({ body }: { body: Record<string, unknown> }) => {
      const planId = body.planId as number
      const plan = plans.find(p => p.id === planId)
      if (plan) {
        tenants[0].plan = plan
        tenants[0].maxUsers = plan.maxUsers
        tenants[0].planExpireAt = new Date(Date.now() + 365 * 24 * 3600 * 1000).toISOString()
      }
      const invoiceNo = `INV-${new Date().getFullYear()}-${String(nextInvoiceId++).padStart(3, '0')}`
      return {
        code: 200,
        data: {
          planId,
          planName: plan?.name || '',
          expireAt: tenants[0].planExpireAt,
          invoiceNo
        },
        message: 'ok'
      }
    }
  },
  {
    url: '/api/v1/tenant/invoices',
    method: 'get',
    response: () => ({
      code: 200,
      data: invoices,
      message: 'ok'
    })
  },
  {
    url: '/api/v1/admin/dashboard',
    method: 'get',
    response: () => ({
      code: 200,
      data: {
        totalTenants: 128,
        activeTenants: 96,
        newTenantsThisMonth: 12,
        totalUsers: 1580,
        monthlyRevenue: 89760.00,
        revenueGrowth: 15.3,
        planDistribution: [
          { name: '免费版', count: 32, percentage: 25, color: '#8c8c8c' },
          { name: '标准版', count: 56, percentage: 44, color: '#1677FF' },
          { name: '专业版', count: 28, percentage: 22, color: '#722ED1' },
          { name: '企业版', count: 12, percentage: 9, color: '#FA8C16' }
        ],
        recentTenants: [
          { id: 1, name: '星辰科技有限公司', status: 'active', createTime: '2025-03-28T10:00:00Z' },
          { id: 2, name: '未来数字科技', status: 'trial', createTime: '2025-03-27T14:30:00Z' },
          { id: 3, name: '云端智能科技', status: 'trial', createTime: '2025-03-26T09:15:00Z' },
          { id: 4, name: '创新软件公司', status: 'active', createTime: '2025-03-25T16:45:00Z' },
          { id: 5, name: '数据驱动科技', status: 'active', createTime: '2025-03-24T11:20:00Z' }
        ]
      },
      message: 'ok'
    })
  },
  {
    url: '/api/v1/admin/tenants',
    method: 'get',
    response: () => ({
      code: 200,
      data: [
        ...tenants,
        { id: 2, name: '未来数字科技', slug: 'future-digital', logo: '', contactName: '李四', contactPhone: '13900139002', contactEmail: 'admin@future.com', currentUsers: 5, maxUsers: 200, status: 'trial', trialEndsAt: '2025-04-10T23:59:59Z', planExpireAt: null, plan: plans[2], createTime: '2025-03-27T14:30:00Z', updateTime: '2025-03-27T14:30:00Z' },
        { id: 3, name: '云端智能科技', slug: 'cloud-ai', logo: '', contactName: '王五', contactPhone: '13700137003', contactEmail: 'admin@cloudai.com', currentUsers: 30, maxUsers: 50, status: 'active', trialEndsAt: null, planExpireAt: '2025-12-31T23:59:59Z', plan: plans[1], createTime: '2025-02-10T08:00:00Z', updateTime: '2025-03-15T10:20:00Z' },
        { id: 4, name: '创新软件公司', slug: 'innovate-soft', logo: '', contactName: '赵六', contactPhone: '13600136004', contactEmail: 'admin@innovate.com', currentUsers: 8, maxUsers: 5, status: 'suspended', trialEndsAt: null, planExpireAt: '2025-02-28T23:59:59Z', plan: plans[0], createTime: '2025-01-20T12:00:00Z', updateTime: '2025-03-01T00:00:00Z' }
      ],
      message: 'ok'
    })
  },
  {
    url: '/api/v1/admin/tenants',
    method: 'post',
    response: ({ body }: { body: Record<string, unknown> }) => {
      const newTenant = {
        id: nextTenantId++,
        name: String(body.name || ''),
        slug: String(body.slug || ''),
        logo: '',
        contactName: String(body.contactName || ''),
        contactPhone: String(body.contactPhone || ''),
        contactEmail: String(body.contactEmail || ''),
        currentUsers: 0,
        maxUsers: plans.find(p => p.id === body.planId)?.maxUsers || 5,
        status: 'trial',
        trialEndsAt: new Date(Date.now() + 14 * 24 * 3600 * 1000).toISOString(),
        planExpireAt: null,
        plan: plans.find(p => p.id === body.planId) || plans[0],
        createTime: new Date().toISOString(),
        updateTime: new Date().toISOString()
      }
      tenants.push(newTenant)
      return { code: 200, data: newTenant, message: 'ok' }
    }
  },
  {
    url: '/api/v1/admin/tenants/\\d+',
    method: 'put',
    response: ({ body, url }: { body: Record<string, unknown>; url: string }) => {
      const id = Number((url.match(/\/tenants\/(\d+)/) || [])[1])
      const tenant = tenants.find(t => t.id === id)
      if (tenant) Object.assign(tenant, body)
      return { code: 200, data: null, message: 'ok' }
    }
  },
  {
    url: '/api/v1/admin/tenants/\\d+/activate',
    method: 'put',
    response: ({ url }: { url: string }) => {
      const id = Number((url.match(/\/tenants\/(\d+)/) || [])[1])
      const tenant = tenants.find(t => t.id === id)
      if (tenant) tenant.status = 'active'
      return { code: 200, data: null, message: 'ok' }
    }
  },
  {
    url: '/api/v1/admin/tenants/\\d+/suspend',
    method: 'put',
    response: ({ url }: { url: string }) => {
      const id = Number((url.match(/\/tenants\/(\d+)/) || [])[1])
      const tenant = tenants.find(t => t.id === id)
      if (tenant) tenant.status = 'suspended'
      return { code: 200, data: null, message: 'ok' }
    }
  },
  {
    url: '/api/v1/admin/plans',
    method: 'get',
    response: () => ({
      code: 200,
      data: plans,
      message: 'ok'
    })
  },
  {
    url: '/api/v1/admin/plans',
    method: 'post',
    response: ({ body }: { body: Record<string, unknown> }) => {
      const newPlan = {
        id: plans.length + 1,
        name: String(body.name || ''),
        code: String(body.code || ''),
        price: Number(body.price || 0),
        minUsers: Number(body.minUsers || 1),
        maxUsers: Number(body.maxUsers || 50),
        features: (body.features || {}) as Record<string, boolean | number>,
        isActive: 1,
        createTime: new Date().toISOString(),
        updateTime: new Date().toISOString()
      }
      plans.push(newPlan)
      return { code: 200, data: newPlan, message: 'ok' }
    }
  },
  {
    url: '/api/v1/admin/plans/\\d+',
    method: 'put',
    response: ({ body, url }: { body: Record<string, unknown>; url: string }) => {
      const id = Number((url.match(/\/plans\/(\d+)/) || [])[1])
      const plan = plans.find(p => p.id === id)
      if (plan) {
        if (body.isActive !== undefined) plan.isActive = Number(body.isActive)
        if (body.name !== undefined) plan.name = String(body.name)
        if (body.price !== undefined) plan.price = Number(body.price)
      }
      return { code: 200, data: null, message: 'ok' }
    }
  }
] as MockMethod[]
