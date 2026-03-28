<template>
  <div class="user-management">
    <!-- 搜索和操作栏 -->
    <n-card class="search-card" size="small">
      <n-space justify="space-between">
        <n-space>
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索用户名/昵称"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <n-icon><SearchOutline /></n-icon>
            </template>
          </n-input>
          <n-button type="primary" @click="handleSearch">
            <template #icon><n-icon><SearchOutline /></n-icon></template>
            搜索
          </n-button>
        </n-space>
        <n-button type="primary" @click="handleAdd">
          <template #icon><n-icon><AddOutline /></n-icon></template>
          新增用户
        </n-button>
      </n-space>
    </n-card>

    <!-- 用户列表 -->
    <n-card class="table-card">
      <n-data-table
        :columns="columns"
        :data="userList"
        :loading="loading"
        :pagination="pagination"
        :row-key="(row: User) => row.id"
        @update:page="handlePageChange"
        @update:page-size="handlePageSizeChange"
      />
    </n-card>

    <!-- 新增/编辑用户弹窗 -->
    <n-modal
      v-model:show="modalVisible"
      preset="dialog"
      :title="isEdit ? '编辑用户' : '新增用户'"
      positive-text="确定"
      negative-text="取消"
      :loading="modalLoading"
      @positive-click="handleSubmit"
      @negative-click="handleCancel"
    >
      <n-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-placement="left"
        label-width="80"
        style="margin-top: 20px"
      >
        <n-form-item label="用户名" path="username">
          <n-input v-model:value="formData.username" placeholder="请输入用户名" :disabled="isEdit" />
        </n-form-item>
        <n-form-item label="昵称" path="nickname">
          <n-input v-model:value="formData.nickname" placeholder="请输入昵称" />
        </n-form-item>
        <n-form-item v-if="!isEdit" label="密码" path="password">
          <n-input
            v-model:value="formData.password"
            type="password"
            placeholder="请输入密码"
            show-password-on="click"
          />
        </n-form-item>
        <n-form-item label="部门" path="deptId">
          <n-tree-select
            v-model:value="formData.deptId"
            :options="deptOptions"
            placeholder="请选择部门"
            label-field="name"
            key-field="id"
            children-field="children"
          />
        </n-form-item>
        <n-form-item label="角色" path="roleId">
          <n-select
            v-model:value="formData.roleId"
            :options="roleOptions"
            placeholder="请选择角色"
            label-field="name"
            value-field="id"
          />
        </n-form-item>
        <n-form-item label="邮箱" path="email">
          <n-input v-model:value="formData.email" placeholder="请输入邮箱" />
        </n-form-item>
        <n-form-item label="电话" path="phone">
          <n-input v-model:value="formData.phone" placeholder="请输入电话" />
        </n-form-item>
        <n-form-item label="状态" path="status">
          <n-switch v-model:value="formData.status" :checked-value="1" :unchecked-value="0">
            <template #checked>启用</template>
            <template #unchecked>禁用</template>
          </n-switch>
        </n-form-item>
      </n-form>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted, computed } from 'vue'
import {
  NCard,
  NSpace,
  NInput,
  NButton,
  NIcon,
  NDataTable,
  NModal,
  NForm,
  NFormItem,
  NSelect,
  NTreeSelect,
  NSwitch,
  NTag,
  NPopconfirm,
  useMessage
} from 'naive-ui'
import { SearchOutline, AddOutline, CreateOutline, TrashOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'
import { request } from '@/utils/request'
import type { User, Dept, Role } from '@/types/user'
import type { PageResult } from '@/types/common'

const message = useMessage()

// 状态
const loading = ref(false)
const searchKeyword = ref('')
const userList = ref<User[]>([])
const modalVisible = ref(false)
const modalLoading = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInst | null>(null)

// 部门和角色数据
const deptList = ref<Dept[]>([])
const roleList = ref<Role[]>([])

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50]
})

// 表单数据
const formData = reactive({
  id: 0,
  username: '',
  nickname: '',
  password: '',
  deptId: null as number | null,
  roleId: null as number | null,
  email: '',
  phone: '',
  status: 1 as 0 | 1
})

// 表单验证规则
const formRules: FormRules = {
  username: { required: true, message: '请输入用户名', trigger: 'blur' },
  nickname: { required: true, message: '请输入昵称', trigger: 'blur' },
  password: { required: true, message: '请输入密码', trigger: 'blur' },
  deptId: { required: true, type: 'number', message: '请选择部门', trigger: 'change' },
  roleId: { required: true, type: 'number', message: '请选择角色', trigger: 'change' },
  email: { required: true, message: '请输入邮箱', trigger: 'blur' }
}

// 部门选项 - 转换为树形结构
const deptOptions = computed(() => {
  const buildTree = (items: Dept[], parentId: number | null = null): Dept[] => {
    return items
      .filter(item => item.parentId === parentId)
      .map(item => ({
        ...item,
        children: buildTree(items, item.id)
      }))
  }
  return buildTree(deptList.value)
})

// 角色选项
const roleOptions = computed(() => {
  return roleList.value.map(role => ({
    label: role.name,
    value: role.id,
    name: role.name,
    id: role.id
  }))
})

// 表格列配置
const columns: DataTableColumns<User> = [
  { title: '用户名', key: 'username', width: 120 },
  { title: '昵称', key: 'nickname', width: 120 },
  { title: '部门', key: 'deptName', width: 120 },
  { title: '角色', key: 'roleName', width: 120 },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row) => {
      return h(
        NTag,
        { type: row.status === 1 ? 'success' : 'error', size: 'small' },
        { default: () => (row.status === 1 ? '启用' : '禁用') }
      )
    }
  },
  { title: '邮箱', key: 'email', width: 180 },
  { title: '电话', key: 'phone', width: 130 },
  { title: '创建时间', key: 'createTime', width: 180 },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    fixed: 'right',
    render: (row) => {
      return h(NSpace, null, {
        default: () => [
          h(
            NButton,
            {
              size: 'small',
              quaternary: true,
              type: 'primary',
              onClick: () => handleEdit(row)
            },
            { icon: () => h(NIcon, null, { default: () => h(CreateOutline) }), default: () => '编辑' }
          ),
          h(
            NPopconfirm,
            {
              onPositiveClick: () => handleToggleStatus(row)
            },
            {
              trigger: () =>
                h(
                  NButton,
                  {
                    size: 'small',
                    quaternary: true,
                    type: row.status === 1 ? 'warning' : 'success'
                  },
                  { default: () => (row.status === 1 ? '禁用' : '启用') }
                ),
              default: () => `确定${row.status === 1 ? '禁用' : '启用'}该用户吗？`
            }
          ),
          h(
            NPopconfirm,
            {
              onPositiveClick: () => handleDelete(row)
            },
            {
              trigger: () =>
                h(
                  NButton,
                  { size: 'small', quaternary: true, type: 'error' },
                  { icon: () => h(NIcon, null, { default: () => h(TrashOutline) }), default: () => '删除' }
                ),
              default: () => '确定删除该用户吗？'
            }
          )
        ]
      })
    }
  }
]

// 获取用户列表
async function fetchUserList() {
  loading.value = true
  try {
    const result: PageResult<User> = await request.get('/user/list', {
      params: {
        page: pagination.page,
        pageSize: pagination.pageSize,
        keyword: searchKeyword.value
      }
    })
    userList.value = result.list
    pagination.itemCount = result.total
  } catch (error) {
    message.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 获取部门列表
async function fetchDeptList() {
  try {
    const result: Dept[] = await request.get('/dept/list')
    deptList.value = result
  } catch (error) {
    message.error('获取部门列表失败')
  }
}

// 获取角色列表
async function fetchRoleList() {
  try {
    const result: Role[] = await request.get('/role/list')
    roleList.value = result
  } catch (error) {
    message.error('获取角色列表失败')
  }
}

// 搜索
function handleSearch() {
  pagination.page = 1
  fetchUserList()
}

// 分页变化
function handlePageChange(page: number) {
  pagination.page = page
  fetchUserList()
}

function handlePageSizeChange(pageSize: number) {
  pagination.pageSize = pageSize
  pagination.page = 1
  fetchUserList()
}

// 重置表单
function resetForm() {
  formData.id = 0
  formData.username = ''
  formData.nickname = ''
  formData.password = ''
  formData.deptId = null
  formData.roleId = null
  formData.email = ''
  formData.phone = ''
  formData.status = 1
}

// 新增用户
function handleAdd() {
  isEdit.value = false
  resetForm()
  modalVisible.value = true
}

// 编辑用户
function handleEdit(row: User) {
  isEdit.value = true
  formData.id = row.id
  formData.username = row.username
  formData.nickname = row.nickname
  formData.deptId = row.deptId
  formData.roleId = row.roleId
  formData.email = row.email
  formData.phone = row.phone
  formData.status = row.status
  modalVisible.value = true
}

// 提交表单
async function handleSubmit() {
  try {
    await formRef.value?.validate()
    modalLoading.value = true

    if (isEdit.value) {
      // 编辑用户
      await request.put(`/user/${formData.id}`, formData)
      message.success('编辑成功')
    } else {
      // 新增用户
      await request.post('/user', formData)
      message.success('新增成功')
    }

    modalVisible.value = false
    fetchUserList()
  } catch (error) {
    if (error) {
      message.error(isEdit.value ? '编辑失败' : '新增失败')
    }
  } finally {
    modalLoading.value = false
  }
}

// 取消
function handleCancel() {
  modalVisible.value = false
  resetForm()
}

// 切换用户状态
async function handleToggleStatus(row: User) {
  try {
    await request.put(`/user/${row.id}/status`, { status: row.status === 1 ? 0 : 1 })
    message.success('操作成功')
    fetchUserList()
  } catch (error) {
    message.error('操作失败')
  }
}

// 删除用户
async function handleDelete(row: User) {
  try {
    await request.delete(`/user/${row.id}`)
    message.success('删除成功')
    fetchUserList()
  } catch (error) {
    message.error('删除失败')
  }
}

// 初始化
onMounted(() => {
  fetchUserList()
  fetchDeptList()
  fetchRoleList()
})
</script>

<style lang="scss" scoped>
.user-management {
  .search-card {
    margin-bottom: 16px;
  }

  .table-card {
    // 表格样式
  }
}
</style>
