const TOKEN_KEY = 'oa_token'
const USER_KEY = 'oa_user'
const TENANT_KEY = 'oa_tenant'
const TENANT_SLUG_KEY = 'oa_tenant_slug'

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export function removeToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

export function getStoredUser<T>(): T | null {
  const user = localStorage.getItem(USER_KEY)
  return user ? JSON.parse(user) : null
}

export function setStoredUser<T>(user: T): void {
  localStorage.setItem(USER_KEY, JSON.stringify(user))
}

export function removeStoredUser(): void {
  localStorage.removeItem(USER_KEY)
}
export function getTenantSlug(): string | null {
  return localStorage.getItem(TENANT_SLUG_KEY)
}
export function setTenantSlug(slug: string): void {
  localStorage.setItem(TENANT_SLUG_KEY, slug)
}
export function removeTenantSlug(): void {
  localStorage.removeItem(TENANT_SLUG_KEY)
}
export function getStoredTenant<T>(): T | null {
  const tenant = localStorage.getItem(TENANT_KEY)
  return tenant ? JSON.parse(tenant) : null
}
export function setStoredTenant<T>(tenant: T): void {
  localStorage.setItem(TENANT_KEY, JSON.stringify(tenant))
}
export function removeStoredTenant(): void {
  localStorage.removeItem(TENANT_KEY)
}
export function clearAuth(): void {
  removeToken()
  removeStoredUser()
  removeStoredTenant()
  removeTenantSlug()
}
