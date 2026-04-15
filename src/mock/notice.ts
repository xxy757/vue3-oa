import { MockMethod } from 'vite-plugin-mock'
import Mock from 'mockjs'

const Random = Mock.Random

// 生成公告数据
const notices = Array.from({ length: 15 }, (_, i) => ({
  id: i + 1,
  title: Random.ctitle(10, 25),
  type: Random.pick(['notice', 'announcement', 'policy', 'urgent']),
  content: `<h2>${Random.ctitle(5, 10)}</h2><p>${Random.cparagraph(3, 6)}</p><p>${Random.cparagraph(3, 6)}</p>`,
  summary: Random.cparagraph(1, 2),
  coverImage: Random.image('800x400', Random.color(), '#FFF', 'Notice'),
  publisherId: 1,
  publisherName: '管理员',
  publisherAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin',
  status: 'published',
  isTop: i < 2,
  readCount: Random.integer(10, 500),
  createTime: Random.datetime('yyyy-MM-dd HH:mm:ss'),
  updateTime: Random.datetime('yyyy-MM-dd HH:mm:ss'),
  isRead: Random.boolean()
}))

export default [
  // 获取公告列表
  {
    url: '/api/v1/notices',
    method: 'get',
    response: ({
      query
    }: {
      query: { page: number; pageSize: number; type?: string; keyword?: string }
    }) => {
      const { page = 1, pageSize = 10, type, keyword } = query

      let filteredNotices = [...notices]

      if (type && type !== 'all') {
        filteredNotices = filteredNotices.filter((n) => n.type === type)
      }

      if (keyword) {
        filteredNotices = filteredNotices.filter((n) => n.title.includes(keyword))
      }

      // 置顶排序
      filteredNotices.sort((a, b) => {
        if (a.isTop && !b.isTop) return -1
        if (!a.isTop && b.isTop) return 1
        return new Date(b.createTime).getTime() - new Date(a.createTime).getTime()
      })

      const start = (page - 1) * pageSize
      const list = filteredNotices.slice(start, start + pageSize)

      return {
        code: 200,
        message: '成功',
        data: {
          list,
          total: filteredNotices.length,
          page,
          pageSize,
          totalPages: Math.ceil(filteredNotices.length / pageSize)
        }
      }
    }
  },

  // 获取公告详情
  {
    url: '/api/notices/:id',
    method: 'get',
    response: ({ query }: { query: { id: string } }) => {
      const id = parseInt(query.id)
      const notice = notices.find((n) => n.id === id)

      if (notice) {
        // 增加阅读次数
        notice.readCount += 1
        notice.isRead = true

        return {
          code: 200,
          message: '成功',
          data: notice
        }
      }

      return {
        code: 404,
        message: '公告不存在',
        data: null
      }
    }
  },

  // 发布公告
  {
    url: '/api/notices',
    method: 'post',
    response: ({ body }: { body: Record<string, unknown> }) => {
      const newNotice = {
        id: notices.length + 1,
        title: body.title as string,
        type: body.type as string,
        content: body.content as string,
        summary: (body.summary as string) || (body.content as string).substring(0, 100),
        coverImage: body.coverImage || '',
        publisherId: 1,
        publisherName: '管理员',
        publisherAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin',
        status: 'published',
        isTop: body.isTop || false,
        readCount: 0,
        createTime: Random.datetime('yyyy-MM-dd HH:mm:ss'),
        updateTime: Random.datetime('yyyy-MM-dd HH:mm:ss'),
        isRead: false
      }

      notices.unshift(newNotice as (typeof notices)[0])

      return {
        code: 200,
        message: '发布成功',
        data: newNotice
      }
    }
  },

  // 标记已读
  {
    url: '/api/notices/:id/read',
    method: 'post',
    response: ({ query }: { query: { id: string } }) => {
      const id = parseInt(query.id)
      const notice = notices.find((n) => n.id === id)

      if (notice) {
        notice.isRead = true
      }

      return {
        code: 200,
        message: '成功',
        data: null
      }
    }
  },

  // 获取未读数量
  {
    url: '/api/v1/notices/unread-count',
    method: 'get',
    response: () => {
      const count = notices.filter((n) => !n.isRead).length

      return {
        code: 200,
        message: '成功',
        data: { count }
      }
    }
  }
] as MockMethod[]
