export interface User {
  id: number
  tenantId: number
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  deptId: number
  deptName: string
  roleId: number
  roleName: string
  permissions: string[]
  status: 0 | 1
  createTime: string
}

export interface LoginForm {
  username: string
  password: string
  tenantSlug?: string
  remember?: boolean
}

export interface LoginResult {
  token: string
  user: User
  tenant: TenantInfo
}

export interface TenantInfo {
  id: number
  name: string
  slug: string
  status: string
  plan: PlanInfo
}

export interface PlanInfo {
  id: number
  name: string
  code: string
  price: number
  features: Record<string, unknown>
  maxUsers: number
}

export interface Dept {
  id: number
  name: string
  parentId: number | null
  children?: Dept[]
  sort: number
  leader: string
  phone: string
  email: string
  status: 0 | 1
}

export interface Role {
  id: number
  name: string
  code: string
  description: string
  permissions: string[]
  status: 0 | 1
  createTime: string
}
