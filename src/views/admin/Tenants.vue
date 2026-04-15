<template>
  <div class="admin-tenants">
    <div class="page-title-bar">
      <h2 class="page-title">租户管理</h2>
    </div>

    <n-card>
      <div class="filter-bar">
        <div class="filter-left">
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索企业名称"
            clearable
            style="width: 240px"
          >
            <template #prefix>
              <n-icon><SearchOutline /></n-icon>
            </template>
          </n-input>
          <n-select
            v-model:value="filterStatus"
            placeholder="状态筛选"
            clearable
            :options="statusOptions"
            style="width: 140px"
          />
        </div>
        <n-button type="primary" @click="handleAdd">
          <template #icon><n-icon><AddOutline /></n-icon></template>
          新增租户
        </n-button>
      </div>

      <n-spin :show="loading">
        <n-table :bordered="false" :single-line="false">
          <thead>
            <tr>
              <th>企业名称</th>
              <th>企业标识</th>
              <th>套餐</th>
              <th>用户数</th>
              <th>状态</th>
              <th>到期时间</th>
              <th>注册时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="tenant in pagedTenants" :key="tenant.id">
              <td>{{ tenant.name }}</td>
              <td>{{ tenant.slug }}</td>
              <td>
                <n-tag size="small" :color="{ color: getPlanColor(tenant.plan?.code || ''), textColor: '#fff' }">
                  {{ tenant.plan?.name || '-' }}
                </n-tag>
              </td>
              <td>{{ tenant.currentUsers }} / {{ tenant.maxUsers }}</td>
              <td>
                <n-tag :type="tenantStatusType(tenant.status)" size="small">
                  {{ tenantStatusLabel(tenant.status) }}
                </n-tag>
              </td>
              <td>{{ tenant.planExpireAt?.substring(0, 10) || '-' }}</td>
              <td>{{ tenant.createTime?.substring(0, 10) }}</td>
              <td class="action-cell">
                <n-button text type="primary" size="small" @click="handleEdit(tenant)">编辑</n-button>
                <n-button
                  v-if="tenant.status === 'suspended'"
                  text type="success" size="small"
                  @click="handleActivate(tenant)"
                >启用</n-button>
                <n-button
                  v-if="tenant.status === 'active'"
                  text type="warning" size="small"
                  @click="handleSuspend(tenant)"
                >暂停</n-button>
              </td>
            </tr>
            <tr v-if="pagedTenants.length === 0">
              <td colspan="8" class="empty-row">暂无租户数据</td>
            </tr>
          </tbody>
        </n-table>
      </n-spin>

      <div class="pagination-bar">
        <n-pagination
          v-model:page="currentPage"
          v-model:page-size="pageSize"
          :item-count="totalCount"
          show-size-picker
          :page-sizes="[10, 20, 50]"
        />
      </div>
    </n-card>

    <n-modal
      v-model:show="showModal"
      preset="dialog"
      :title="editingTenant ? '编辑租户' : '新增租户'"
      style="width: 520px"
    >
      <n-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-placement="left"
        label-width="80px"
      >
        <n-form-item label="企业名称" path="name">
          <n-input v-model:value="formData.name" placeholder="请输入企业名称" />
        </n-form-item>
        <n-form-item label="企业标识" path="slug">
          <n-input v-model:value="formData.slug" placeholder="请输入企业标识" :disabled="!!editingTenant" />
        </n-form-item>
        <n-form-item label="联系人" path="contactName">
          <n-input v-model:value="formData.contactName" placeholder="请输入联系人" />
        </n-form-item>
        <n-form-item label="联系电话" path="contactPhone">
          <n-input v-model:value="formData.contactPhone" placeholder="请输入联系电话" />
        </n-form-item>
        <n-form-item label="联系邮箱" path="contactEmail">
          <n-input v-model:value="formData.contactEmail" placeholder="请输入联系邮箱" />
        </n-form-item>
        <n-form-item label="套餐" path="planId">
          <n-select v-model:value="formData.planId" :options="planOptions" placeholder="请选择套餐" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" :loading="saving" @click="handleSave">确认</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  NCard,
  NTable,
  NInput,
  NSelect,
  NButton,
  NTag,
  NModal,
  NForm,
  NFormItem,
  NIcon,
  NPagination,
  NSpin,
  useMessage,
  useDialog
} from 'naive-ui'
import type { FormInst, FormRules, SelectOption } from 'naive-ui'
import { SearchOutline, AddOutline } from '@vicons/ionicons5'
import { request } from '@/utils/request'
import type { TenantDetail, Plan } from '@/types/tenant'
import { PLAN_STATUS_MAP, getPlanColor } from '@/utils/plan'

const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const searchKeyword = ref('')
const filterStatus = ref<string | null>(null)
const currentPage = ref(1)
const pageSize = ref(10)
const tenants = ref<TenantDetail[]>([])
const editingTenant = ref<TenantDetail | null>(null)
const allPlans = ref<Plan[]>([])
const formRef = ref<FormInst | null>(null)

const formData = ref({
  name: '',
  slug: '',
  contactName: '',
  contactPhone: '',
  contactEmail: '',
  planId: null as number | null
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入企业名称', trigger: 'blur' }],
  slug: [
    { required: true, message: '请输入企业标识', trigger: 'blur' },
    { pattern: /^[a-z0-9][a-z0-9-]{1,48}[a-z0-9]$/, message: '只能包含小写字母、数字和连字符', trigger: 'blur' }
  ],
  contactName: [{ required: true, message: '请输入联系人', trigger: 'blur' }],
  contactPhone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  contactEmail: [
    { required: true, message: '请输入联系邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  planId: [{ required: true, message: '请选择套餐', trigger: 'change' }]
}

const statusOptions = [
  { label: '试用中', value: 'trial' },
  { label: '已激活', value: 'active' },
  { label: '已暂停', value: 'suspended' },
  { label: '已注销', value: 'cancelled' }
]

const planOptions = computed<SelectOption[]>(() =>
  allPlans.value.map(p => ({ label: p.name, value: p.id }))
)

const filteredTenants = computed(() => {
  let result = tenants.value
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(t => t.name.toLowerCase().includes(kw) || t.slug.toLowerCase().includes(kw))
  }
  if (filterStatus.value) {
    result = result.filter(t => t.status === filterStatus.value)
  }
  return result
})

const totalCount = computed(() => filteredTenants.value.length)

const pagedTenants = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredTenants.value.slice(start, start + pageSize.value)
})

function tenantStatusType(status: string) {
  return PLAN_STATUS_MAP[status]?.type || 'default'
}

function tenantStatusLabel(status: string) {
  return PLAN_STATUS_MAP[status]?.label || status
}

function handleAdd() {
  editingTenant.value = null
  formData.value = { name: '', slug: '', contactName: '', contactPhone: '', contactEmail: '', planId: null }
  showModal.value = true
}

function handleEdit(tenant: TenantDetail) {
  editingTenant.value = tenant
  formData.value = {
    name: tenant.name,
    slug: tenant.slug,
    contactName: tenant.contactName,
    contactPhone: tenant.contactPhone,
    contactEmail: tenant.contactEmail,
    planId: tenant.plan?.id || null
  }
  showModal.value = true
}

async function handleSave() {
  try {
    await formRef.value?.validate()
    saving.value = true
    if (editingTenant.value) {
      await request.put(`/admin/tenants/${editingTenant.value.id}`, formData.value)
      message.success('租户更新成功')
    } else {
      await request.post('/admin/tenants', formData.value)
      message.success('租户创建成功')
    }
    showModal.value = false
    await fetchTenants()
  } catch (err: unknown) {
    if (err instanceof Error) {
      message.error(err.message)
    }
  } finally {
    saving.value = false
  }
}

function handleActivate(tenant: TenantDetail) {
  dialog.info({
    title: '确认启用',
    content: `确定要启用租户「${tenant.name}」吗？`,
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await request.put(`/admin/tenants/${tenant.id}/activate`)
        message.success('租户已启用')
        await fetchTenants()
      } catch {
        message.error('操作失败')
      }
    }
  })
}

function handleSuspend(tenant: TenantDetail) {
  dialog.warning({
    title: '确认暂停',
    content: `确定要暂停租户「${tenant.name}」吗？暂停后该租户下所有用户将无法使用系统`,
    positiveText: '确认暂停',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await request.put(`/admin/tenants/${tenant.id}/suspend`)
        message.success('租户已暂停')
        await fetchTenants()
      } catch {
        message.error('操作失败')
      }
    }
  })
}

async function fetchTenants() {
  loading.value = true
  try {
    const data = await request.get<TenantDetail[]>(`/admin/tenants?page=${currentPage.value}&pageSize=${pageSize.value}`)
    tenants.value = Array.isArray(data) ? data : []
  } catch {
    message.error('获取租户列表失败')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await Promise.all([
    fetchTenants(),
    request.get<Plan[]>('/plans').then(data => { allPlans.value = data }).catch(() => {})
  ])
})
</script>

<style lang="scss" scoped>
.admin-tenants {
  padding: 0;
}

.page-title-bar {
  margin-bottom: 16px;

  .page-title {
    font-size: 20px;
    font-weight: 600;
    color: $text-color-1;
    margin: 0;
  }
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  .filter-left {
    display: flex;
    gap: 12px;
  }
}

.action-cell {
  .n-button {
    margin-right: 8px;
  }
}

.pagination-bar {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.empty-row {
  text-align: center;
  color: $text-color-3;
  padding: 40px 0 !important;
}
</style>
