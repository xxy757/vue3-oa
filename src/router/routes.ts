import type { RouteRecordRaw } from 'vue-router'

export const constantRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: {
      title: '登录',
      requiresAuth: false
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue'),
    meta: {
      title: '企业注册',
      requiresAuth: false
    }
  },
  {
    path: '/choose-plan',
    name: 'ChoosePlan',
    component: () => import('@/views/auth/ChoosePlan.vue'),
    meta: {
      title: '选择套餐',
      requiresAuth: true
    }
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: {
      title: '页面不存在',
      requiresAuth: false
    }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404'
  }
]

export const asyncRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/components/Layout/OALayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: '工作台',
          icon: 'home-outline',
          requiresAuth: true
        }
      }
    ]
  },
  {
    path: '/approval',
    name: 'ApprovalLayout',
    component: () => import('@/components/Layout/OALayout.vue'),
    redirect: '/approval/my-apply',
    meta: {
      title: '审批中心',
      icon: 'document-text-outline',
      requiresAuth: true
    },
    children: [
      {
        path: 'apply',
        name: 'ApprovalApply',
        component: () => import('@/views/approval/Apply.vue'),
        meta: {
          title: '发起申请',
          icon: 'add-circle-outline',
          requiresAuth: true
        }
      },
      {
        path: 'my-apply',
        name: 'MyApply',
        component: () => import('@/views/approval/MyApply.vue'),
        meta: {
          title: '我的申请',
          icon: 'paper-plane-outline',
          requiresAuth: true
        }
      },
      {
        path: 'pending',
        name: 'PendingApproval',
        component: () => import('@/views/approval/Pending.vue'),
        meta: {
          title: '待我审批',
          icon: 'mail-unread-outline',
          requiresAuth: true
        }
      },
      {
        path: 'done',
        name: 'DoneApproval',
        component: () => import('@/views/approval/Done.vue'),
        meta: {
          title: '已办审批',
          icon: 'checkmark-done-outline',
          requiresAuth: true
        }
      },
      {
        path: 'detail/:id',
        name: 'ApprovalDetail',
        component: () => import('@/views/approval/Detail.vue'),
        meta: {
          title: '审批详情',
          requiresAuth: true,
          hidden: true
        }
      }
    ]
  },
  {
    path: '/notice',
    name: 'NoticeLayout',
    component: () => import('@/components/Layout/OALayout.vue'),
    redirect: '/notice/list',
    meta: {
      title: '公告通知',
      icon: 'megaphone-outline',
      requiresAuth: true
    },
    children: [
      {
        path: 'list',
        name: 'NoticeList',
        component: () => import('@/views/notice/List.vue'),
        meta: {
          title: '公告列表',
          icon: 'list-outline',
          requiresAuth: true
        }
      },
      {
        path: 'detail/:id',
        name: 'NoticeDetail',
        component: () => import('@/views/notice/Detail.vue'),
        meta: {
          title: '公告详情',
          requiresAuth: true,
          hidden: true
        }
      }
    ]
  },
  {
    path: '/schedule',
    name: 'ScheduleLayout',
    component: () => import('@/components/Layout/OALayout.vue'),
    redirect: '/schedule/calendar',
    meta: {
      title: '日程管理',
      icon: 'calendar-outline',
      requiresAuth: true
    },
    children: [
      {
        path: 'calendar',
        name: 'ScheduleCalendar',
        component: () => import('@/views/schedule/Calendar.vue'),
        meta: {
          title: '日程日历',
          icon: 'calendar-outline',
          requiresAuth: true
        }
      },
      {
        path: 'list',
        name: 'ScheduleList',
        component: () => import('@/views/schedule/List.vue'),
        meta: {
          title: '日程列表',
          icon: 'list-outline',
          requiresAuth: true
        }
      }
    ]
  },
  {
    path: '/system',
    name: 'SystemLayout',
    component: () => import('@/components/Layout/OALayout.vue'),
    redirect: '/system/user',
    meta: {
      title: '系统管理',
      icon: 'settings-outline',
      requiresAuth: true,
      roles: ['admin']
    },
    children: [
      {
        path: 'user',
        name: 'SystemUser',
        component: () => import('@/views/system/User.vue'),
        meta: {
          title: '用户管理',
          icon: 'people-outline',
          requiresAuth: true,
          roles: ['admin']
        }
      },
      {
        path: 'dept',
        name: 'SystemDept',
        component: () => import('@/views/system/Dept.vue'),
        meta: {
          title: '部门管理',
          icon: 'git-branch-outline',
          requiresAuth: true,
          roles: ['admin']
        }
      },
      {
        path: 'role',
        name: 'SystemRole',
        component: () => import('@/views/system/Role.vue'),
        meta: {
          title: '角色管理',
          icon: 'shield-checkmark-outline',
          requiresAuth: true,
          roles: ['admin']
        }
      },
      {
        path: 'flow',
        name: 'SystemFlow',
        component: () => import('@/views/system/Flow.vue'),
        meta: {
          title: '流程配置',
          icon: 'git-compare-outline',
          requiresAuth: true,
          roles: ['admin']
        }
      }
    ]
  },
  {
    path: '/profile',
    name: 'ProfileLayout',
    component: () => import('@/components/Layout/OALayout.vue'),
    redirect: '/profile/info',
    meta: {
      title: '个人中心',
      icon: 'person-outline',
      requiresAuth: true
    },
    children: [
      {
        path: 'info',
        name: 'ProfileInfo',
        component: () => import('@/views/profile/Info.vue'),
        meta: {
          title: '个人信息',
          icon: 'person-outline',
          requiresAuth: true
        }
      },
      {
        path: 'password',
        name: 'ProfilePassword',
        component: () => import('@/views/profile/Password.vue'),
        meta: {
          title: '修改密码',
          icon: 'lock-closed-outline',
          requiresAuth: true
        }
      }
    ]
  },
  {
    path: '/tenant',
    name: 'TenantLayout',
    component: () => import('@/components/Layout/OALayout.vue'),
    redirect: '/tenant/info',
    meta: {
      title: '企业管理',
      icon: 'business-outline',
      requiresAuth: true,
      roles: ['admin']
    },
    children: [
      {
        path: 'info',
        name: 'TenantInfo',
        component: () => import('@/views/tenant/Info.vue'),
        meta: {
          title: '企业信息',
          icon: 'business-outline',
          requiresAuth: true,
          roles: ['admin']
        }
      },
      {
        path: 'plan',
        name: 'TenantPlan',
        component: () => import('@/views/tenant/Plan.vue'),
        meta: {
          title: '套餐管理',
          icon: 'pricetag-outline',
          requiresAuth: true,
          roles: ['admin']
        }
      },
      {
        path: 'invoices',
        name: 'TenantInvoices',
        component: () => import('@/views/tenant/Invoices.vue'),
        meta: {
          title: '账单管理',
          icon: 'receipt-outline',
          requiresAuth: true,
          roles: ['admin']
        }
      }
    ]
  },
  {
    path: '/admin',
    name: 'AdminLayout',
    component: () => import('@/components/Layout/AdminLayout.vue'),
    redirect: '/admin/dashboard',
    meta: {
      title: '平台管理',
      requiresAuth: true,
      roles: ['super_admin']
    },
    children: [
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/Dashboard.vue'),
        meta: {
          title: '管理概览',
          icon: 'home-outline',
          requiresAuth: true,
          roles: ['super_admin']
        }
      },
      {
        path: 'tenants',
        name: 'AdminTenants',
        component: () => import('@/views/admin/Tenants.vue'),
        meta: {
          title: '租户管理',
          icon: 'business-outline',
          requiresAuth: true,
          roles: ['super_admin']
        }
      },
      {
        path: 'plans',
        name: 'AdminPlans',
        component: () => import('@/views/admin/Plans.vue'),
        meta: {
          title: '套餐管理',
          icon: 'pricetag-outline',
          requiresAuth: true,
          roles: ['super_admin']
        }
      }
    ]
  }
]
