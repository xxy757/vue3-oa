<template>
  <div class="role-management">
    <!-- 操作栏 -->
    <n-card class="action-card" size="small">
      <n-space justify="end">
        <n-button type="primary" @click="handleAdd">
          <template #icon><n-icon><AddOutline /></n-icon></template>
          新增角色
        </n-button>
      </n-space>
    </n-card>

    <!-- 角色列表 -->
    <n-card class="table-card">
      <n-data-table
        :columns="columns"
        :data="roleList"
        :loading="loading"
        :row-key="(row: Role) => row.id"
      />
    </n-card>

    <!-- 新增/编辑角色弹窗 -->
    <n-modal
      v-model:show="modalVisible"
      preset="dialog"
      :title="isEdit ? '编辑角色' : '新增角色'"
      style="width: 600px"
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
        <n-form-item label="角色名称" path="name">
          <n-input v-model:value="formData.name" placeholder="请输入角色名称" />
        </n-form-item>
        <n-form-item label="角色编码" path="code">
          <n-input v-model:value="formData.code" placeholder="请输入角色编码" :disabled="isEdit" />
        </n-form-item>
        <n-form-item label="描述" path="description">
          <n-input
            v-model:value="formData.description"
            type="textarea"
            placeholder="请输入角色描述"
            :rows="3"
          />
        </n-form-item>
        <n-form-item label="权限配置" path="permissions">
          <n-card size="small" style="width: 100%">
            <n-checkbox-group v-model:value="formData.permissions">
              <n-grid :cols="2" :x-gap="12" :y-gap="8">
                <n-gi v-for="perm in permissionOptions" :key="perm.value">
                  <n-checkbox :value="perm.value" :label="perm.label" />
                </n-gi>
              </n-grid>
            </n-checkbox-group>
          </n-card>
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
import { ref, reactive, h, onMounted } from 'vue'
import {
  NCard,
  NSpace,
  NButton,
  NIcon,
  NDataTable,
  NModal,
  NForm,
  NFormItem,
  NInput,
  NSwitch,
  NTag,
  NPopconfirm,
  NCheckbox,
  NCheckboxGroup,
  NGrid,
  NGi,
  useMessage
} from 'naive-ui'
import { AddOutline, CreateOutline, TrashOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'
import { request } from '@/utils/request'
import type { Role } from '@/types/user'

const message = useMessage()

// 权限选项
const permissionOptions = [
  { label: '所有权限', value: '*' },
  { label: '发起申请', value: 'approval:apply' },
  { label: '审批申请', value: 'approval:approve' },
  { label: '查看公告', value: 'notice:view' },
  { label: '发布公告', value: 'notice:create' },
  { label: '管理公告', value: 'notice:manage' },
  { label: '查看日程', value: 'schedule:view' },
  { label: '管理日程', value: 'schedule:manage' },
  { label: '用户管理', value: 'system:user' },
  { label: '部门管理', value: 'system:dept' },
  { label: '角色管理', value: 'system:role' },
  { label: '流程配置', value: 'system:flow' }
]

// 状态
const loading = ref(false)
const roleList = ref<Role[]>([])
const modalVisible = ref(false)
const modalLoading = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInst | null>(null)

// 表单数据
const formData = reactive({
  id: 0,
  name: '',
  code: '',
  description: '',
  permissions: [] as string[],
  status: 1 as 0 | 1
})

// 表单验证规则
const formRules: FormRules = {
  name: { required: true, message: '请输入角色名称', trigger: 'blur' },
  code: [
    { required: true, message: '请输入角色编码', trigger: 'blur' },
    { pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/, message: '编码只能包含字母、数字和下划线，且以字母开头', trigger: 'blur' }
  ],
  permissions: { required: true, type: 'array', message: '请选择权限', trigger: 'change' }
}

// 表格列配置
const columns: DataTableColumns<Role> = [
  { title: '角色名称', key: 'name', width: 120 },
  { title: '角色编码', key: 'code', width: 120 },
  { title: '描述', key: 'description', ellipsis: { tooltip: true } },
  {
    title: '权限',
    key: 'permissions',
    width: 300,
    render: (row) => {
      const tags = row.permissions.map(p => {
        const perm = permissionOptions.find(opt => opt.value === p)
        return h(
          NTag,
          { type: p === '*' ? 'error' : 'default', size: 'small', style: { margin: '2px' } },
          { default: () => perm?.label || p }
        )
      })
      return h('div', { style: { display: 'flex', flexWrap: 'wrap' } }, tags)
    }
  },
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
              onPositiveClick: () => handleDelete(row)
            },
            {
              trigger: () =>
                h(
                  NButton,
                  { size: 'small', quaternary: true, type: 'error', disabled: row.code === 'admin' },
                  { icon: () => h(NIcon, null, { default: () => h(TrashOutline) }), default: () => '删除' }
                ),
              default: () => '确定删除该角色吗？'
            }
          )
        ]
      })
    }
  }
]

// 获取角色列表
async function fetchRoleList() {
  loading.value = true
  try {
    const result: Role[] = await request.get('/role/list')
    roleList.value = result
  } catch (error) {
    message.error('获取角色列表失败')
  } finally {
    loading.value = false
  }
}

// 重置表单
function resetForm() {
  formData.id = 0
  formData.name = ''
  formData.code = ''
  formData.description = ''
  formData.permissions = []
  formData.status = 1
}

// 新增角色
function handleAdd() {
  isEdit.value = false
  resetForm()
  modalVisible.value = true
}

// 编辑角色
function handleEdit(row: Role) {
  isEdit.value = true
  formData.id = row.id
  formData.name = row.name
  formData.code = row.code
  formData.description = row.description
  formData.permissions = [...row.permissions]
  formData.status = row.status
  modalVisible.value = true
}

// 提交表单
async function handleSubmit() {
  try {
    await formRef.value?.validate()
    modalLoading.value = true

    if (isEdit.value) {
      await request.put(`/role/${formData.id}`, formData)
      message.success('编辑成功')
    } else {
      await request.post('/role', formData)
      message.success('新增成功')
    }

    modalVisible.value = false
    fetchRoleList()
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

// 删除角色
async function handleDelete(row: Role) {
  try {
    await request.delete(`/role/${row.id}`)
    message.success('删除成功')
    fetchRoleList()
  } catch (error) {
    message.error('删除失败')
  }
}

// 初始化
onMounted(() => {
  fetchRoleList()
})
</script>

<style lang="scss" scoped>
.role-management {
  .action-card {
    margin-bottom: 16px;
  }
}
</style>
