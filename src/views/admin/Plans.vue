<template>
  <div class="admin-plans">
    <div class="page-title-bar">
      <h2 class="page-title">套餐管理</h2>
    </div>

    <n-card>
      <div class="filter-bar">
        <span class="filter-info">共 {{ plans.length }} 个套餐</span>
        <n-button type="primary" @click="handleAdd">
          <template #icon><n-icon><AddOutline /></n-icon></template>
          新增套餐
        </n-button>
      </div>

      <n-spin :show="loading">
        <n-table :bordered="false" :single-line="false">
          <thead>
            <tr>
              <th>套餐名称</th>
              <th>标识</th>
              <th>价格</th>
              <th>用户范围</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="plan in plans" :key="plan.id">
              <td>
                <span class="plan-name">{{ plan.name }}</span>
              </td>
              <td>
                <n-tag size="small">{{ plan.code }}</n-tag>
              </td>
              <td>{{ formatPrice(plan.price) }}</td>
              <td>{{ plan.minUsers }} - {{ plan.maxUsers }} 人</td>
              <td>
                <n-tag :type="plan.isActive ? 'success' : 'default'" size="small">
                  {{ plan.isActive ? '启用' : '停用' }}
                </n-tag>
              </td>
              <td class="action-cell">
                <n-button text type="primary" size="small" @click="handleEdit(plan)">编辑</n-button>
                <n-button
                  v-if="plan.isActive"
                  text type="warning" size="small"
                  @click="handleToggle(plan, false)"
                >停用</n-button>
                <n-button
                  v-else
                  text type="success" size="small"
                  @click="handleToggle(plan, true)"
                >启用</n-button>
              </td>
            </tr>
            <tr v-if="plans.length === 0">
              <td colspan="6" class="empty-row">暂无套餐数据</td>
            </tr>
          </tbody>
        </n-table>
      </n-spin>
    </n-card>

    <n-modal
      v-model:show="showModal"
      preset="dialog"
      :title="editingPlan ? '编辑套餐' : '新增套餐'"
      style="width: 640px"
    >
      <n-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-placement="left"
        label-width="100px"
      >
        <n-form-item label="套餐名称" path="name">
          <n-input v-model:value="formData.name" placeholder="如：标准版" />
        </n-form-item>
        <n-form-item label="套餐标识" path="code">
          <n-input v-model:value="formData.code" placeholder="如：standard" :disabled="!!editingPlan" />
        </n-form-item>
        <n-form-item label="价格(元)" path="price">
          <n-input-number v-model:value="formData.price" :min="0" :precision="2" style="width: 100%" />
        </n-form-item>
        <n-form-item label="最小用户数" path="minUsers">
          <n-input-number v-model:value="formData.minUsers" :min="1" style="width: 100%" />
        </n-form-item>
        <n-form-item label="最大用户数" path="maxUsers">
          <n-input-number v-model:value="formData.maxUsers" :min="1" style="width: 100%" />
        </n-form-item>
        <n-form-item label="功能特性">
          <div class="features-form">
            <div v-for="key in featureKeys" :key="key" class="feature-row">
              <n-checkbox v-model:checked="formData.features[key]">
                {{ PLAN_LABELS[key] || key }}
              </n-checkbox>
            </div>
          </div>
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
import { ref, reactive, onMounted } from 'vue'
import {
  NCard,
  NTable,
  NInput,
  NInputNumber,
  NButton,
  NTag,
  NModal,
  NForm,
  NFormItem,
  NIcon,
  NCheckbox,
  NSpin,
  useMessage
} from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import { AddOutline } from '@vicons/ionicons5'
import { request } from '@/utils/request'
import type { Plan } from '@/types/tenant'
import { PLAN_LABELS, formatPrice } from '@/utils/plan'

const message = useMessage()

const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const plans = ref<Plan[]>([])
const editingPlan = ref<Plan | null>(null)
const formRef = ref<FormInst | null>(null)

const featureKeys = ['approval', 'schedule', 'notice', 'api', 'sso']

const formData = reactive({
  name: '',
  code: '',
  price: 0,
  minUsers: 1,
  maxUsers: 50,
  features: {} as Record<string, boolean | number>
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入套餐名称', trigger: 'blur' }],
  code: [
    { required: true, message: '请输入套餐标识', trigger: 'blur' },
    { pattern: /^[a-z][a-z0-9-]*$/, message: '只能包含小写字母、数字和连字符', trigger: 'blur' }
  ],
  price: [{ required: true, type: 'number', message: '请输入价格', trigger: 'blur' }],
  minUsers: [{ required: true, type: 'number', message: '请输入最小用户数', trigger: 'blur' }],
  maxUsers: [{ required: true, type: 'number', message: '请输入最大用户数', trigger: 'blur' }]
}

function handleAdd() {
  editingPlan.value = null
  Object.assign(formData, {
    name: '',
    code: '',
    price: 0,
    minUsers: 1,
    maxUsers: 50,
    features: { approval: true, schedule: true, notice: true, api: false, sso: false }
  })
  showModal.value = true
}

function handleEdit(plan: Plan) {
  editingPlan.value = plan
  Object.assign(formData, {
    name: plan.name,
    code: plan.code,
    price: plan.price,
    minUsers: plan.minUsers,
    maxUsers: plan.maxUsers,
    features: { ...plan.features }
  })
  showModal.value = true
}

async function handleSave() {
  try {
    await formRef.value?.validate()
    saving.value = true
    if (editingPlan.value) {
      await request.put(`/admin/plans/${editingPlan.value.id}`, formData)
      message.success('套餐更新成功')
    } else {
      await request.post('/admin/plans', formData)
      message.success('套餐创建成功')
    }
    showModal.value = false
    await fetchPlans()
  } catch (err: unknown) {
    if (err instanceof Error) {
      message.error(err.message)
    }
  } finally {
    saving.value = false
  }
}

async function handleToggle(plan: Plan, active: boolean) {
  try {
    await request.put(`/admin/plans/${plan.id}`, { isActive: active ? 1 : 0 })
    message.success(active ? '套餐已启用' : '套餐已停用')
    await fetchPlans()
  } catch {
    message.error('操作失败')
  }
}

async function fetchPlans() {
  loading.value = true
  try {
    const data = await request.get<Plan[]>('/admin/plans')
    plans.value = Array.isArray(data) ? data : []
  } catch {
    message.error('获取套餐列表失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchPlans()
})
</script>

<style lang="scss" scoped>
.admin-plans {
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

  .filter-info {
    font-size: 14px;
    color: $text-color-3;
  }
}

.plan-name {
  font-weight: 500;
  color: $text-color-1;
}

.action-cell {
  .n-button {
    margin-right: 8px;
  }
}

.features-form {
  .feature-row {
    padding: 4px 0;
  }
}

.empty-row {
  text-align: center;
  color: $text-color-3;
  padding: 40px 0 !important;
}
</style>
