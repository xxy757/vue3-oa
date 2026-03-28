// 用户相关类型
export interface User {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  deptId: number
  deptName: string
  roleId: number
  roleName: string
  status: 0 | 1 // 0: 禁用, 1: 启用
  createTime: string
}

export interface LoginForm {
  username: string
  password: string
  remember?: boolean
}

export interface LoginResult {
  token: string
  user: User
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
