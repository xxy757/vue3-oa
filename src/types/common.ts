// 通用类型

// 分页请求参数
export interface PaginationParams {
  page: number
  pageSize: number
}

// 分页响应
export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

// API 响应
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

// 通用选项类型
export interface SelectOption {
  label: string
  value: string | number
  disabled?: boolean
}

// 树形节点
export interface TreeNode {
  id: number
  name: string
  parentId: number | null
  children?: TreeNode[]
  [key: string]: unknown
}
