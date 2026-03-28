<template>
  <div class="profile-info">
    <n-grid :cols="24" :x-gap="16">
      <!-- 左侧头像卡片 -->
      <n-gi :span="8">
        <n-card title="个人信息" size="small">
          <div class="avatar-section">
            <n-avatar
              :src="userInfo?.avatar"
              :size="120"
              round
              class="user-avatar"
            />
            <n-h3 class="user-name">{{ userInfo?.nickname }}</n-h3>
            <n-space justify="center">
              <n-tag type="primary">{{ userInfo?.roleName }}</n-tag>
              <n-tag type="info">{{ userInfo?.deptName }}</n-tag>
            </n-space>
          </div>
          <n-divider />
          <n-descriptions :column="1" label-placement="left">
            <n-descriptions-item label="用户名">{{ userInfo?.username }}</n-descriptions-item>
            <n-descriptions-item label="邮箱">{{ userInfo?.email }}</n-descriptions-item>
            <n-descriptions-item label="电话">{{ userInfo?.phone }}</n-descriptions-item>
            <n-descriptions-item label="创建时间">{{ userInfo?.createTime }}</n-descriptions-item>
          </n-descriptions>
        </n-card>
      </n-gi>

      <!-- 右侧编辑卡片 -->
      <n-gi :span="16">
        <n-card title="编辑信息" size="small">
          <n-form
            ref="formRef"
            :model="formData"
            :rules="formRules"
            label-placement="left"
            label-width="80"
            style="max-width: 500px"
          >
            <n-form-item label="昵称" path="nickname">
              <n-input v-model:value="formData.nickname" placeholder="请输入昵称" />
            </n-form-item>
            <n-form-item label="邮箱" path="email">
              <n-input v-model:value="formData.email" placeholder="请输入邮箱" />
            </n-form-item>
            <n-form-item label="电话" path="phone">
              <n-input v-model:value="formData.phone" placeholder="请输入电话" />
            </n-form-item>
            <n-form-item label="头像" path="avatar">
              <div class="avatar-upload">
                <n-input v-model:value="formData.avatar" placeholder="请输入头像URL" style="flex: 1" />
                <n-avatar
                  v-if="formData.avatar"
                  :src="formData.avatar"
                  :size="60"
                  round
                  style="margin-left: 16px"
                />
              </div>
            </n-form-item>
            <n-form-item>
              <n-space>
                <n-button type="primary" :loading="loading" @click="handleSubmit">
                  保存修改
                </n-button>
                <n-button @click="handleReset">重置</n-button>
              </n-space>
            </n-form-item>
          </n-form>
        </n-card>

        <!-- 账号安全 -->
        <n-card title="账号安全" size="small" style="margin-top: 16px">
          <n-list bordered>
            <n-list-item>
              <template #prefix>
                <n-icon size="24"><LockClosedOutline /></n-icon>
              </template>
              <n-thing title="登录密码" description="定期更换密码可以提高账号安全性">
                <template #action>
                  <n-button text type="primary" @click="router.push('/profile/password')">
                    修改密码
                  </n-button>
                </template>
              </n-thing>
            </n-list-item>
            <n-list-item>
              <template #prefix>
                <n-icon size="24"><ShieldCheckmarkOutline /></n-icon>
              </template>
              <n-thing title="账号状态" description="当前账号状态正常">
                <template #action>
                  <n-tag type="success">正常</n-tag>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  NGrid,
  NGi,
  NCard,
  NAvatar,
  NH3,
  NSpace,
  NTag,
  NDivider,
  NDescriptions,
  NDescriptionsItem,
  NForm,
  NFormItem,
  NInput,
  NButton,
  NList,
  NListItem,
  NThing,
  NIcon,
  useMessage
} from 'naive-ui'
import { LockClosedOutline, ShieldCheckmarkOutline } from '@vicons/ionicons5'
import type { FormInst, FormRules } from 'naive-ui'
import { useUserStore } from '@/stores/user'
import { request } from '@/utils/request'
import type { User } from '@/types/user'

const router = useRouter()
const message = useMessage()
const userStore = useUserStore()

// 状态
const loading = ref(false)
const formRef = ref<FormInst | null>(null)

// 用户信息
const userInfo = computed(() => userStore.userInfo)

// 表单数据
const formData = reactive({
  nickname: '',
  email: '',
  phone: '',
  avatar: ''
})

// 表单验证规则
const formRules: FormRules = {
  nickname: { required: true, message: '请输入昵称', trigger: 'blur' },
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码', trigger: 'blur' }
  ]
}

// 初始化表单数据
function initFormData() {
  if (userInfo.value) {
    formData.nickname = userInfo.value.nickname
    formData.email = userInfo.value.email
    formData.phone = userInfo.value.phone
    formData.avatar = userInfo.value.avatar
  }
}

// 重置表单
function handleReset() {
  initFormData()
}

// 提交表单
async function handleSubmit() {
  try {
    await formRef.value?.validate()
    loading.value = true

    // 调用接口更新用户信息
    await request.put('/auth/info', formData)

    // 更新本地状态
    userStore.updateUserInfo(formData)

    message.success('保存成功')
  } catch (error) {
    if (error) {
      message.error('保存失败')
    }
  } finally {
    loading.value = false
  }
}

// 初始化
onMounted(() => {
  initFormData()
})
</script>

<style lang="scss" scoped>
.profile-info {
  .avatar-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 16px 0;

    .user-avatar {
      margin-bottom: 16px;
    }

    .user-name {
      margin: 0 0 12px 0;
    }
  }

  .avatar-upload {
    display: flex;
    align-items: center;
  }
}
</style>
