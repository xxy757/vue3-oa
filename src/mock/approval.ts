import { MockMethod } from 'vite-plugin-mock'
import Mock from 'mockjs'

const Random = Mock.Random

// 生成审批数据
const generateApproval = (id: number, type: string, status: string) => {
  const baseApproval = {
    id,
    type,
    title: Random.ctitle(5, 15),
    applicantId: Random.integer(2, 3),
    applicantName: Random.cname(),
    applicantDept: '技术部',
    status,
    currentStep: status === 'pending' ? Random.integer(1, 2) : 2,
    totalStep: 2,
    createTime: Random.datetime('yyyy-MM-dd HH:mm:ss'),
    updateTime: Random.datetime('yyyy-MM-dd HH:mm:ss'),
    flowNodes: [
      {
        id: 1,
        nodeId: 1,
        nodeName: '部门经理审批',
        approverId: 3,
        approverName: '李经理',
        approverAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=manager',
        status: status === 'pending' ? 'pending' : 'approved',
        comment: status !== 'pending' ? '同意' : '',
        sort: 1,
        createTime: Random.datetime('yyyy-MM-dd HH:mm:ss')
      },
      {
        id: 2,
        nodeId: 2,
        nodeName: '人事审批',
        approverId: 4,
        approverName: '人事专员',
        approverAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=hr',
        status:
          status === 'pending'
            ? 'pending'
            : Random.pick(['approved', 'approved', 'approved', 'rejected']),
        comment: '',
        sort: 2,
        createTime: Random.datetime('yyyy-MM-dd HH:mm:ss')
      }
    ]
  }

  switch (type) {
    case 'leave':
      return {
        ...baseApproval,
        leaveType: Random.pick(['annual', 'sick', 'personal']),
        startDate: Random.date('yyyy-MM-dd'),
        endDate: Random.date('yyyy-MM-dd'),
        days: Random.integer(1, 5),
        reason: Random.csentence(10, 30)
      }
    case 'expense':
      return {
        ...baseApproval,
        expenseType: Random.pick(['travel', 'office', 'entertainment', 'other']),
        amount: Random.float(100, 10000, 2, 2),
        description: Random.csentence(10, 30),
        attachments: []
      }
    case 'overtime':
      return {
        ...baseApproval,
        overtimeDate: Random.date('yyyy-MM-dd'),
        startTime: '18:00',
        endTime: '21:00',
        hours: Random.float(1, 4, 1, 1),
        reason: Random.csentence(10, 30)
      }
    case 'travel':
      return {
        ...baseApproval,
        destination: Random.city(),
        startDate: Random.date('yyyy-MM-dd'),
        endDate: Random.date('yyyy-MM-dd'),
        days: Random.integer(1, 7),
        reason: Random.csentence(10, 30),
        budget: Random.float(1000, 20000, 2, 2)
      }
    default:
      return {
        ...baseApproval,
        content: Random.csentence(20, 50),
        attachments: []
      }
  }
}

// 预生成一些审批数据
const myApprovals = [
  generateApproval(1, 'leave', 'approved'),
  generateApproval(2, 'expense', 'pending'),
  generateApproval(3, 'overtime', 'approved'),
  generateApproval(4, 'travel', 'rejected'),
  generateApproval(5, 'leave', 'pending')
]

const pendingApprovals = [
  generateApproval(101, 'leave', 'pending'),
  generateApproval(102, 'expense', 'pending'),
  generateApproval(103, 'overtime', 'pending'),
  generateApproval(104, 'travel', 'pending'),
  generateApproval(105, 'leave', 'pending')
]

const doneApprovals = [
  generateApproval(201, 'leave', 'approved'),
  generateApproval(202, 'expense', 'approved'),
  generateApproval(203, 'overtime', 'rejected'),
  generateApproval(204, 'travel', 'approved'),
  generateApproval(205, 'leave', 'approved')
]

export default [
  // 获取我的申请列表
  {
    url: '/api/approvals/my',
    method: 'get',
    response: ({
      query
    }: {
      query: { page: number; pageSize: number; status?: string; type?: string }
    }) => {
      const { page = 1, pageSize = 10 } = query

      const start = (page - 1) * pageSize
      const list = myApprovals.slice(start, start + pageSize)

      return {
        code: 200,
        message: '成功',
        data: {
          list,
          total: myApprovals.length,
          page,
          pageSize,
          totalPages: Math.ceil(myApprovals.length / pageSize)
        }
      }
    }
  },

  // 获取待审批列表
  {
    url: '/api/approvals/pending',
    method: 'get',
    response: ({ query }: { query: { page: number; pageSize: number } }) => {
      const { page = 1, pageSize = 10 } = query

      const start = (page - 1) * pageSize
      const list = pendingApprovals.slice(start, start + pageSize)

      return {
        code: 200,
        message: '成功',
        data: {
          list,
          total: pendingApprovals.length,
          page,
          pageSize,
          totalPages: Math.ceil(pendingApprovals.length / pageSize)
        }
      }
    }
  },

  // 获取已办审批列表
  {
    url: '/api/approvals/done',
    method: 'get',
    response: ({ query }: { query: { page: number; pageSize: number } }) => {
      const { page = 1, pageSize = 10 } = query

      const start = (page - 1) * pageSize
      const list = doneApprovals.slice(start, start + pageSize)

      return {
        code: 200,
        message: '成功',
        data: {
          list,
          total: doneApprovals.length,
          page,
          pageSize,
          totalPages: Math.ceil(doneApprovals.length / pageSize)
        }
      }
    }
  },

  // 获取审批详情
  {
    url: '/api/approvals/:id',
    method: 'get',
    response: ({ query }: { query: { id: string } }) => {
      const id = parseInt(query.id)
      const approval = [...myApprovals, ...pendingApprovals, ...doneApprovals].find(
        (a) => a.id === id
      )

      if (approval) {
        return {
          code: 200,
          message: '成功',
          data: approval
        }
      }

      return {
        code: 404,
        message: '审批不存在',
        data: null
      }
    }
  },

  // 发起申请
  {
    url: '/api/approvals',
    method: 'post',
    response: ({ body }: { body: Record<string, unknown> }) => {
      const newApproval = generateApproval(
        Mock.Random.integer(1000, 9999),
        body.type as string,
        'pending'
      )
      newApproval.title = body.title as string
      myApprovals.unshift(newApproval as (typeof myApprovals)[0])

      return {
        code: 200,
        message: '申请提交成功',
        data: newApproval
      }
    }
  },

  // 审批操作
  {
    url: '/api/approvals/:id/action',
    method: 'post',
    response: () => {
      return {
        code: 200,
        message: '操作成功',
        data: null
      }
    }
  },

  // 撤回申请
  {
    url: '/api/approvals/:id/withdraw',
    method: 'post',
    response: () => {
      return {
        code: 200,
        message: '撤回成功',
        data: null
      }
    }
  },

  // 获取待办统计
  {
    url: '/api/approvals/stats',
    method: 'get',
    response: () => {
      return {
        code: 200,
        message: '成功',
        data: {
          myPending: myApprovals.filter((a) => a.status === 'pending').length,
          myApproved: myApprovals.filter((a) => a.status === 'approved').length,
          myRejected: myApprovals.filter((a) => a.status === 'rejected').length,
          todoApproval: pendingApprovals.length
        }
      }
    }
  }
] as MockMethod[]
