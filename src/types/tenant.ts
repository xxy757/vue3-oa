export interface Plan {
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

export interface TenantDetail {
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
  plan: Plan
  createTime: string
  updateTime: string
}

export interface TenantUpdateForm {
  name?: string
  logo?: string
  contactName?: string
  contactPhone?: string
  contactEmail?: string
}

export interface RegisterForm {
  name: string
  slug: string
  contactName: string
  contactPhone: string
  contactEmail: string
  planId?: number
}

export interface RegisterResult {
  tenantId: number
  name: string
  slug: string
  trialEndsAt: string
  adminUser: {
    id: number
    username: string
    tempPassword: string
  }
}

export interface UpgradeResult {
  planId: number
  planName: string
  expireAt: string
  invoiceNo: string
}

export interface Invoice {
  id: number
  tenantId: number
  planId: number
  invoiceNo: string
  periodStart: string
  periodEnd: string
  userCount: number
  amount: number
  status: string
  paidAt: string | null
  paymentMethod: string
  paymentTransactionId: string
  createTime: string
  updateTime: string
}
