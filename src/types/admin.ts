export interface DashboardStats {
  totalTenants: number
  activeTenants: number
  newTenantsThisMonth: number
  totalUsers: number
  monthlyRevenue: number
  revenueGrowth: number
  planDistribution: Array<{ name: string; count: number; percentage: number; color: string }>
  recentTenants: Array<{ id: number; name: string; status: string; createTime: string }>
}
