<template>
  <div class="dept-management">
    <n-grid :cols="24" :x-gap="16">
      <!-- 左侧部门树 -->
      <n-gi :span="8">
        <n-card title="部门结构" size="small">
          <template #header-extra>
            <n-button type="primary" size="small" @click="handleAdd(null)">
              <template #icon
                ><n-icon><AddOutline /></n-icon
              ></template>
              新增
            </n-button>
          </template>
          <n-spin :show="treeLoading">
            <n-tree
              :data="treeData"
              :node-props="nodeProps"
              block-line
              expand-on-click
              :default-expanded-keys="expandedKeys"
              :selected-keys="selectedKeys"
              @update:selected-keys="handleSelect"
            />
          </n-spin>
        </n-card>
      </n-gi>

      <!-- 右侧部门信息 -->
      <n-gi :span="16">
        <n-card title="部门信息" size="small">
          <template v-if="currentDept">
            <n-descriptions :column="2" label-placement="left" bordered>
              <n-descriptions-item label="部门名称">{{ currentDept.name }}</n-descriptions-item>
              <n-descriptions-item label="上级部门">
                {{ currentDept.parentId ? getDeptName(currentDept.parentId) : '无' }}
              </n-descriptions-item>
              <n-descriptions-item label="负责人">{{ currentDept.leader }}</n-descriptions-item>
              <n-descriptions-item label="联系电话">{{ currentDept.phone }}</n-descriptions-item>
              <n-descriptions-item label="邮箱">{{ currentDept.email }}</n-descriptions-item>
              <n-descriptions-item label="排序">{{ currentDept.sort }}</n-descriptions-item>
              <n-descriptions-item label="状态">
                <n-tag :type="currentDept.status === 1 ? 'success' : 'error'" size="small">
                  {{ currentDept.status === 1 ? '启用' : '禁用' }}
                </n-tag>
              </n-descriptions-item>
            </n-descriptions>
            <n-space style="margin-top: 16px">
              <n-button type="primary" @click="handleEdit(currentDept)">
                <template #icon
                  ><n-icon><CreateOutline /></n-icon
                ></template>
                编辑
              </n-button>
              <n-button @click="handleAdd(currentDept)">
                <template #icon
                  ><n-icon><AddOutline /></n-icon
                ></template>
                添加子部门
              </n-button>
              <n-popconfirm @positive-click="handleDelete(currentDept)">
                <template #trigger>
                  <n-button type="error" :disabled="!canDelete(currentDept)">
                    <template #icon
                      ><n-icon><TrashOutline /></n-icon
                    ></template>
                    删除
                  </n-button>
                </template>
                确定删除该部门吗？
              </n-popconfirm>
            </n-space>
          </template>
          <n-empty v-else description="请选择部门" />
        </n-card>
      </n-gi>
    </n-grid>

    <!-- 新增/编辑部门弹窗 -->
    <n-modal
      v-model:show="modalVisible"
      preset="dialog"
      :title="isEdit ? '编辑部门' : '新增部门'"
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
        <n-form-item label="部门名称" path="name">
          <n-input v-model:value="formData.name" placeholder="请输入部门名称" />
        </n-form-item>
        <n-form-item label="上级部门" path="parentId">
          <n-tree-select
            v-model:value="formData.parentId"
            :options="deptTreeOptions"
            placeholder="请选择上级部门"
            label-field="name"
            key-field="id"
            children-field="children"
            clearable
          />
        </n-form-item>
        <n-form-item label="负责人" path="leader">
          <n-input v-model:value="formData.leader" placeholder="请输入负责人" />
        </n-form-item>
        <n-form-item label="联系电话" path="phone">
          <n-input v-model:value="formData.phone" placeholder="请输入联系电话" />
        </n-form-item>
        <n-form-item label="邮箱" path="email">
          <n-input v-model:value="formData.email" placeholder="请输入邮箱" />
        </n-form-item>
        <n-form-item label="排序" path="sort">
          <n-input-number v-model:value="formData.sort" :min="1" style="width: 100%" />
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
  import { ref, reactive, computed, onMounted } from 'vue'
  import {
    NGrid,
    NGi,
    NCard,
    NButton,
    NIcon,
    NTree,
    NModal,
    NForm,
    NFormItem,
    NInput,
    NInputNumber,
    NTreeSelect,
    NSwitch,
    NTag,
    NSpace,
    NDescriptions,
    NDescriptionsItem,
    NEmpty,
    NSpin,
    NPopconfirm,
    useMessage
  } from 'naive-ui'
  import { AddOutline, CreateOutline, TrashOutline } from '@vicons/ionicons5'
  import type { TreeSelectOption, TreeOption, FormInst, FormRules } from 'naive-ui'
  import { request } from '@/utils/request'
  import type { Dept } from '@/types/user'

  const message = useMessage()

  // 状态
  const treeLoading = ref(false)
  const modalVisible = ref(false)
  const modalLoading = ref(false)
  const isEdit = ref(false)
  const formRef = ref<FormInst | null>(null)

  // 部门数据
  const deptList = ref<Dept[]>([])
  const currentDept = ref<Dept | null>(null)
  const selectedKeys = ref<number[]>([])
  const expandedKeys = ref<number[]>([])

  // 表单数据
  const formData = reactive({
    id: 0,
    name: '',
    parentId: null as number | null,
    leader: '',
    phone: '',
    email: '',
    sort: 1,
    status: 1 as 0 | 1
  })

  // 表单验证规则
  const formRules: FormRules = {
    name: { required: true, message: '请输入部门名称', trigger: 'blur' },
    leader: { required: true, message: '请输入负责人', trigger: 'blur' },
    phone: { required: true, message: '请输入联系电话', trigger: 'blur' },
    email: { required: true, message: '请输入邮箱', trigger: 'blur' }
  }

  // 构建树形数据
  function buildTree(items: Dept[], parentId: number | null = null): TreeOption[] {
    return items
      .filter((item) => item.parentId === parentId)
      .sort((a, b) => a.sort - b.sort)
      .map((item) => ({
        key: item.id,
        label: item.name,
        children: buildTree(items, item.id) as any,
        ...item
      })) as any as TreeOption[]
  }

  // 树形数据
  const treeData = computed(() => buildTree(deptList.value) as TreeOption[])

  // 部门选择器选项（排除当前编辑的部门）
  const deptTreeOptions = computed(() => {
    const buildOptions = (
      items: Dept[],
      parentId: number | null = null,
      excludeId?: number
    ): Dept[] => {
      return items
        .filter((item) => item.parentId === parentId && item.id !== excludeId)
        .map((item) => ({
          ...item,
          children: buildOptions(items, item.id, excludeId)
        }))
    }
    return buildOptions(
      deptList.value,
      null,
      isEdit.value ? formData.id : undefined
    ) as any as TreeSelectOption[]
  })

  // 获取部门名称
  function getDeptName(id: number): string {
    const dept = deptList.value.find((d) => d.id === id)
    return dept?.name || ''
  }

  // 是否可删除
  function canDelete(dept: Dept): boolean {
    // 有子部门不能删除
    return !deptList.value.some((d) => d.parentId === dept.id)
  }

  // 节点属性
  function nodeProps({ option }: { option: TreeOption }) {
    return {
      onClick() {
        const dept = deptList.value.find((d) => d.id === option.key)
        if (dept) {
          currentDept.value = dept
          selectedKeys.value = [dept.id]
        }
      }
    }
  }

  // 选择节点
  function handleSelect(keys: number[]) {
    if (keys.length > 0) {
      const dept = deptList.value.find((d) => d.id === keys[0])
      if (dept) {
        currentDept.value = dept
      }
    }
  }

  // 获取部门列表
  async function fetchDeptList() {
    treeLoading.value = true
    try {
      const result: Dept[] = await request.get('/dept/list')
      deptList.value = result
      // 展开第一层
      expandedKeys.value = result.filter((d) => d.parentId === null).map((d) => d.id)
      // 默认选中第一个
      if (result.length > 0 && !currentDept.value) {
        currentDept.value = result[0]
        selectedKeys.value = [result[0].id]
      }
    } catch (error) {
      message.error('获取部门列表失败')
    } finally {
      treeLoading.value = false
    }
  }

  // 重置表单
  function resetForm() {
    formData.id = 0
    formData.name = ''
    formData.parentId = null
    formData.leader = ''
    formData.phone = ''
    formData.email = ''
    formData.sort = 1
    formData.status = 1
  }

  // 新增部门
  function handleAdd(parent: Dept | null) {
    isEdit.value = false
    resetForm()
    if (parent) {
      formData.parentId = parent.id
    }
    modalVisible.value = true
  }

  // 编辑部门
  function handleEdit(dept: Dept) {
    isEdit.value = true
    formData.id = dept.id
    formData.name = dept.name
    formData.parentId = dept.parentId
    formData.leader = dept.leader
    formData.phone = dept.phone
    formData.email = dept.email
    formData.sort = dept.sort
    formData.status = dept.status
    modalVisible.value = true
  }

  // 提交表单
  async function handleSubmit() {
    try {
      await formRef.value?.validate()
      modalLoading.value = true

      if (isEdit.value) {
        await request.put(`/dept/${formData.id}`, formData)
        message.success('编辑成功')
      } else {
        await request.post('/dept', formData)
        message.success('新增成功')
      }

      modalVisible.value = false
      fetchDeptList()
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

  // 删除部门
  async function handleDelete(dept: Dept) {
    try {
      await request.delete(`/dept/${dept.id}`)
      message.success('删除成功')
      currentDept.value = null
      selectedKeys.value = []
      fetchDeptList()
    } catch (error) {
      message.error('删除失败')
    }
  }

  // 初始化
  onMounted(() => {
    fetchDeptList()
  })
</script>

<style lang="scss" scoped>
  .dept-management {
    // 样式
  }
</style>
