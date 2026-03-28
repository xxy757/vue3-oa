<template>
  <div class="flow-management">
    <!-- 操作栏 -->
    <n-card class="action-card" size="small">
      <n-space justify="end">
        <n-button type="primary" @click="handleAdd">
          <template #icon><n-icon><AddOutline /></n-icon></template>
          新增流程
        </n-button>
      </n-space>
    </n-card>

    <!-- 流程列表 -->
    <n-card class="table-card">
      <n-data-table
        :columns="columns"
        :data="flowList"
        :loading="loading"
        :row-key="(row: FlowConfig) => row.id"
      />
    </n-card>

    <!-- 新增/编辑流程弹窗 -->
    <n-modal
      v-model:show="modalVisible"
      preset="dialog"
      :title="isEdit ? '编辑流程' : '新增流程'"
      style="width: 800px"
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
        <n-form-item label="流程名称" path="name">
          <n-input v-model:value="formData.name" placeholder="请输入流程名称" />
        </n-form-item>
        <n-form-item label="流程编码" path="code">
          <n-input v-model:value="formData.code" placeholder="请输入流程编码" :disabled="isEdit" />
        </n-form-item>
        <n-form-item label="描述" path="description">
          <n-input
            v-model:value="formData.description"
            type="textarea"
            placeholder="请输入流程描述"
            :rows="2"
          />
        </n-form-item>

        <!-- 流程节点配置 -->
        <n-form-item label="流程节点">
          <div class="flow-nodes">
            <div v-for="(node, index) in formData.nodes" :key="index" class="flow-node">
              <n-card size="small" :title="`节点 ${index + 1}`">
                <template #header-extra>
                  <n-button
                    v-if="formData.nodes.length > 1"
                    text
                    type="error"
                    @click="removeNode(index)"
                  >
                    <template #icon><n-icon><CloseOutline /></n-icon></template>
                  </n-button>
                </template>
                <n-space vertical>
                  <n-input v-model:value="node.name" placeholder="节点名称" />
                  <n-select
                    v-model:value="node.type"
                    :options="nodeTypeOptions"
                    placeholder="节点类型"
                  />
                  <n-select
                    v-if="node.type === 'approval'"
                    v-model:value="node.approver"
                    :options="approverOptions"
                    placeholder="审批人"
                    multiple
                  />
                </n-space>
              </n-card>
              <div v-if="index < formData.nodes.length - 1" class="node-arrow">
                <n-icon size="24"><ArrowDownOutline /></n-icon>
              </div>
            </div>
            <n-button dashed block @click="addNode">
              <template #icon><n-icon><AddOutline /></n-icon></template>
              添加节点
            </n-button>
          </div>
        </n-form-item>

        <n-form-item label="状态" path="status">
          <n-switch v-model:value="formData.status" :checked-value="1" :unchecked-value="0">
            <template #checked>启用</template>
            <template #unchecked>禁用</template>
          </n-switch>
        </n-form-item>
      </n-form>
    </n-modal>

    <!-- 查看流程弹窗 -->
    <n-modal
      v-model:show="previewVisible"
      preset="card"
      title="流程预览"
      style="width: 600px"
    >
      <div v-if="currentFlow" class="flow-preview">
        <n-h4>{{ currentFlow.name }}</n-h4>
        <n-p depth="3">{{ currentFlow.description }}</n-p>
        <n-divider />
        <div class="flow-diagram">
          <div class="flow-start">
            <n-tag type="success">开始</n-tag>
          </div>
          <div v-for="(node, index) in currentFlow.nodes" :key="index" class="flow-node-preview">
            <div class="node-arrow-preview">
              <n-icon size="20"><ArrowDownOutline /></n-icon>
            </div>
            <n-card size="small" :class="['node-card', `node-${node.type}`]">
              <template #header>
                <n-space align="center">
                  <n-icon v-if="node.type === 'approval'"><CheckmarkCircleOutline /></n-icon>
                  <n-icon v-else-if="node.type === 'notify'"><NotificationsOutline /></n-icon>
                  <n-icon v-else><GitBranchOutline /></n-icon>
                  <span>{{ node.name }}</span>
                </n-space>
              </template>
              <n-text v-if="node.type === 'approval'" depth="3">
                审批人: {{ getApproverNames(node.approver) }}
              </n-text>
              <n-text v-else-if="node.type === 'notify'" depth="3">
                通知人: {{ getApproverNames(node.approver) }}
              </n-text>
            </n-card>
          </div>
          <div class="flow-end">
            <div class="node-arrow-preview">
              <n-icon size="20"><ArrowDownOutline /></n-icon>
            </div>
            <n-tag type="error">结束</n-tag>
          </div>
        </div>
      </div>
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
  NSelect,
  NSwitch,
  NTag,
  NPopconfirm,
  NDivider,
  NH4,
  NP,
  NText,
  useMessage
} from 'naive-ui'
import {
  AddOutline,
  CreateOutline,
  TrashOutline,
  CloseOutline,
  ArrowDownOutline,
  EyeOutline,
  CheckmarkCircleOutline,
  NotificationsOutline,
  GitBranchOutline
} from '@vicons/ionicons5'
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'
import { request } from '@/utils/request'
import type { Role, User } from '@/types/user'

interface FlowNode {
  name: string
  type: 'submit' | 'approval' | 'notify' | 'condition'
  approver: number[]
}

interface FlowConfig {
  id: number
  name: string
  code: string
  description: string
  nodes: FlowNode[]
  status: 0 | 1
  createTime: string
}

const message = useMessage()

// 节点类型选项
const nodeTypeOptions = [
  { label: '提交节点', value: 'submit' },
  { label: '审批节点', value: 'approval' },
  { label: '通知节点', value: 'notify' },
  { label: '条件分支', value: 'condition' }
]

// 状态
const loading = ref(false)
const flowList = ref<FlowConfig[]>([])
const modalVisible = ref(false)
const modalLoading = ref(false)
const previewVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInst | null>(null)
const currentFlow = ref<FlowConfig | null>(null)

// 审批人选项
const approverOptions = ref<{ label: string; value: number }[]>([])

// 模拟流程数据
const mockFlows: FlowConfig[] = [
  {
    id: 1,
    name: '请假审批流程',
    code: 'leave',
    description: '员工请假申请审批流程',
    nodes: [
      { name: '提交申请', type: 'submit', approver: [] },
      { name: '部门经理审批', type: 'approval', approver: [3] },
      { name: '人事审批', type: 'approval', approver: [1] },
      { name: '通知申请人', type: 'notify', approver: [] }
    ],
    status: 1,
    createTime: '2024-01-01 00:00:00'
  },
  {
    id: 2,
    name: '报销审批流程',
    code: 'expense',
    description: '费用报销申请审批流程',
    nodes: [
      { name: '提交申请', type: 'submit', approver: [] },
      { name: '部门经理审批', type: 'approval', approver: [3] },
      { name: '财务审批', type: 'approval', approver: [1] },
      { name: '通知申请人', type: 'notify', approver: [] }
    ],
    status: 1,
    createTime: '2024-01-01 00:00:00'
  },
  {
    id: 3,
    name: '出差审批流程',
    code: 'travel',
    description: '出差申请审批流程',
    nodes: [
      { name: '提交申请', type: 'submit', approver: [] },
      { name: '部门经理审批', type: 'approval', approver: [3] },
      { name: '通知人事', type: 'notify', approver: [1] }
    ],
    status: 1,
    createTime: '2024-01-01 00:00:00'
  }
]

// 表单数据
const formData = reactive({
  id: 0,
  name: '',
  code: '',
  description: '',
  nodes: [] as FlowNode[],
  status: 1 as 0 | 1
})

// 表单验证规则
const formRules: FormRules = {
  name: { required: true, message: '请输入流程名称', trigger: 'blur' },
  code: [
    { required: true, message: '请输入流程编码', trigger: 'blur' },
    { pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/, message: '编码只能包含字母、数字和下划线，且以字母开头', trigger: 'blur' }
  ]
}

// 表格列配置
const columns: DataTableColumns<FlowConfig> = [
  { title: '流程名称', key: 'name', width: 150 },
  { title: '流程编码', key: 'code', width: 120 },
  { title: '描述', key: 'description', ellipsis: { tooltip: true } },
  {
    title: '节点数',
    key: 'nodes',
    width: 100,
    render: (row) => row.nodes.length
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
    width: 200,
    fixed: 'right',
    render: (row) => {
      return h(NSpace, null, {
        default: () => [
          h(
            NButton,
            {
              size: 'small',
              quaternary: true,
              type: 'info',
              onClick: () => handlePreview(row)
            },
            { icon: () => h(NIcon, null, { default: () => h(EyeOutline) }), default: () => '预览' }
          ),
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
                  { size: 'small', quaternary: true, type: 'error' },
                  { icon: () => h(NIcon, null, { default: () => h(TrashOutline) }), default: () => '删除' }
                ),
              default: () => '确定删除该流程吗？'
            }
          )
        ]
      })
    }
  }
]

// 获取审批人名称
function getApproverNames(ids: number[]): string {
  if (!ids || ids.length === 0) return '自动'
  const names = ids.map(id => {
    const option = approverOptions.value.find(opt => opt.value === id)
    return option?.label || '未知'
  })
  return names.join(', ')
}

// 获取流程列表
async function fetchFlowList() {
  loading.value = true
  try {
    // 模拟接口调用
    await new Promise(resolve => setTimeout(resolve, 300))
    flowList.value = mockFlows
  } catch (error) {
    message.error('获取流程列表失败')
  } finally {
    loading.value = false
  }
}

// 获取用户列表（作为审批人选项）
async function fetchUserList() {
  try {
    const result = await request.get<{ list: User[] }>('/user/list', { params: { pageSize: 100 } })
    approverOptions.value = result.list.map(user => ({
      label: `${user.nickname} (${user.deptName})`,
      value: user.id
    }))
  } catch (error) {
    // 使用模拟数据
    approverOptions.value = [
      { label: '管理员 (技术部)', value: 1 },
      { label: '张三 (技术部)', value: 2 },
      { label: '李经理 (技术部)', value: 3 }
    ]
  }
}

// 重置表单
function resetForm() {
  formData.id = 0
  formData.name = ''
  formData.code = ''
  formData.description = ''
  formData.nodes = [{ name: '提交申请', type: 'submit', approver: [] }]
  formData.status = 1
}

// 添加节点
function addNode() {
  formData.nodes.push({ name: '', type: 'approval', approver: [] })
}

// 移除节点
function removeNode(index: number) {
  formData.nodes.splice(index, 1)
}

// 新增流程
function handleAdd() {
  isEdit.value = false
  resetForm()
  modalVisible.value = true
}

// 编辑流程
function handleEdit(row: FlowConfig) {
  isEdit.value = true
  formData.id = row.id
  formData.name = row.name
  formData.code = row.code
  formData.description = row.description
  formData.nodes = JSON.parse(JSON.stringify(row.nodes))
  formData.status = row.status
  modalVisible.value = true
}

// 预览流程
function handlePreview(row: FlowConfig) {
  currentFlow.value = row
  previewVisible.value = true
}

// 提交表单
async function handleSubmit() {
  try {
    await formRef.value?.validate()

    // 验证节点
    if (formData.nodes.length === 0) {
      message.error('请至少添加一个节点')
      return false
    }

    for (const node of formData.nodes) {
      if (!node.name) {
        message.error('请填写所有节点名称')
        return false
      }
    }

    modalLoading.value = true

    if (isEdit.value) {
      // 模拟更新
      const index = mockFlows.findIndex(f => f.id === formData.id)
      if (index !== -1) {
        mockFlows[index] = { ...mockFlows[index], ...formData }
      }
      message.success('编辑成功')
    } else {
      // 模拟新增
      const newFlow: FlowConfig = {
        ...formData,
        id: Date.now(),
        createTime: new Date().toISOString()
      }
      mockFlows.push(newFlow)
      message.success('新增成功')
    }

    modalVisible.value = false
    fetchFlowList()
  } catch (error) {
    if (error) {
      message.error(isEdit.value ? '编辑失败' : '新增失败')
    }
  } finally {
    modalLoading.value = false
  }
  return true
}

// 取消
function handleCancel() {
  modalVisible.value = false
  resetForm()
}

// 删除流程
async function handleDelete(row: FlowConfig) {
  try {
    const index = mockFlows.findIndex(f => f.id === row.id)
    if (index !== -1) {
      mockFlows.splice(index, 1)
    }
    message.success('删除成功')
    fetchFlowList()
  } catch (error) {
    message.error('删除失败')
  }
}

// 初始化
onMounted(() => {
  fetchFlowList()
  fetchUserList()
})
</script>

<style lang="scss" scoped>
.flow-management {
  .action-card {
    margin-bottom: 16px;
  }

  .flow-nodes {
    width: 100%;

    .flow-node {
      margin-bottom: 8px;
    }

    .node-arrow {
      display: flex;
      justify-content: center;
      padding: 8px 0;
      color: #999;
    }
  }

  .flow-preview {
    .flow-diagram {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 16px 0;

      .flow-start,
      .flow-end {
        padding: 8px 0;
      }

      .flow-node-preview {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 100%;
        max-width: 300px;

        .node-arrow-preview {
          padding: 8px 0;
          color: #999;
        }

        .node-card {
          width: 100%;

          &.node-approval {
            border-left: 3px solid #18a058;
          }

          &.node-notify {
            border-left: 3px solid #2080f0;
          }

          &.node-condition {
            border-left: 3px solid #f0a020;
          }
        }
      }
    }
  }
}
</style>
