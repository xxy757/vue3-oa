import { MockMethod } from 'vite-plugin-mock'
import Mock from 'mockjs'

const Random = Mock.Random

// 模拟用户数据
const users = [
  {
    id: 1,
    username: 'admin',
    password: '123456',
    nickname: '管理员',
    email: 'admin@company.com',
    phone: '13800138000',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin',
    deptId: 1,
    deptName: '技术部',
    roleId: 1,
    roleName: '管理员',
    status: 1,
    createTime: '2024-01-01 00:00:00'
  },
  {
    id: 2,
    username: 'user',
    password: '123456',
    nickname: '张三',
    email: 'zhangsan@company.com',
    phone: '13800138001',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhangsan',
    deptId: 1,
    deptName: '技术部',
    roleId: 2,
    roleName: '员工',
    status: 1,
    createTime: '2024-01-15 00:00:00'
  },
  {
    id: 3,
    username: 'manager',
    password: '123456',
    nickname: '李经理',
    email: 'limanager@company.com',
    phone: '13800138002',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=manager',
    deptId: 1,
    deptName: '技术部',
    roleId: 3,
    roleName: '部门经理',
    status: 1,
    createTime: '2024-01-10 00:00:00'
  }
]

// 部门数据
const departments = [
  { id: 1, name: '技术部', parentId: null, sort: 1, leader: '李经理', phone: '010-12345678', email: 'tech@company.com', status: 1 },
  { id: 2, name: '产品部', parentId: null, sort: 2, leader: '王经理', phone: '010-12345679', email: 'product@company.com', status: 1 },
  { id: 3, name: '市场部', parentId: null, sort: 3, leader: '赵经理', phone: '010-12345680', email: 'market@company.com', status: 1 },
  { id: 4, name: '人事部', parentId: null, sort: 4, leader: '钱经理', phone: '010-12345681', email: 'hr@company.com', status: 1 },
  { id: 5, name: '财务部', parentId: null, sort: 5, leader: '孙经理', phone: '010-12345682', email: 'finance@company.com', status: 1 },
  { id: 6, name: '前端组', parentId: 1, sort: 1, leader: '周组长', phone: '010-12345683', email: 'frontend@company.com', status: 1 },
  { id: 7, name: '后端组', parentId: 1, sort: 2, leader: '吴组长', phone: '010-12345684', email: 'backend@company.com', status: 1 }
]

// 角色数据
const roles = [
  { id: 1, name: '管理员', code: 'admin', description: '系统管理员，拥有所有权限', permissions: ['*'], status: 1, createTime: '2024-01-01 00:00:00' },
  { id: 2, name: '员工', code: 'employee', description: '普通员工', permissions: ['approval:apply', 'notice:view', 'schedule:view'], status: 1, createTime: '2024-01-01 00:00:00' },
  { id: 3, name: '部门经理', code: 'manager', description: '部门经理，可审批申请', permissions: ['approval:apply', 'approval:approve', 'notice:view', 'notice:create', 'schedule:view'], status: 1, createTime: '2024-01-01 00:00:00' }
]

export default [
  // 登录
  {
    url: '/api/auth/login',
    method: 'post',
    response: ({ body }: { body: { username: string; password: string } }) => {
      const { username, password } = body
      const user = users.find(u => u.username === username && u.password === password)

      if (user) {
        const { password: _, ...userInfo } = user
        return {
          code: 200,
          message: '登录成功',
          data: {
            token: Random.guid(),
            user: userInfo
          }
        }
      }

      return {
        code: 401,
        message: '用户名或密码错误',
        data: null
      }
    }
  },

  // 获取用户信息
  {
    url: '/api/auth/info',
    method: 'get',
    response: () => {
      const { password: _, ...userInfo } = users[0]
      return {
        code: 200,
        message: '成功',
        data: userInfo
      }
    }
  },

  // 修改密码
  {
    url: '/api/auth/password',
    method: 'put',
    response: () => {
      return {
        code: 200,
        message: '密码修改成功',
        data: null
      }
    }
  },

  // 获取用户列表
  {
    url: '/api/user/list',
    method: 'get',
    response: ({ query }: { query: { page: number; pageSize: number; keyword?: string } }) => {
      const { page = 1, pageSize = 10, keyword = '' } = query

      let filteredUsers = users.map(u => {
        const { password: _, ...rest } = u
        return rest
      })

      if (keyword) {
        filteredUsers = filteredUsers.filter(u =>
          u.username.includes(keyword) || u.nickname.includes(keyword)
        )
      }

      const start = (page - 1) * pageSize
      const list = filteredUsers.slice(start, start + pageSize)

      return {
        code: 200,
        message: '成功',
        data: {
          list,
          total: filteredUsers.length,
          page,
          pageSize,
          totalPages: Math.ceil(filteredUsers.length / pageSize)
        }
      }
    }
  },

  // 获取部门列表
  {
    url: '/api/dept/list',
    method: 'get',
    response: () => {
      return {
        code: 200,
        message: '成功',
        data: departments
      }
    }
  },

  // 获取角色列表
  {
    url: '/api/role/list',
    method: 'get',
    response: () => {
      return {
        code: 200,
        message: '成功',
        data: roles
      }
    }
  }
] as MockMethod[]
