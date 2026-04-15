<template>
  <div class="tenant-invoices-page">
    <div class="page-title-bar">
      <h2 class="page-title">账单管理</h2>
    </div>

    <n-card>
      <n-spin :show="loading">
        <n-table :bordered="false" :single-line="false">
          <thead>
            <tr>
              <th>账单编号</th>
              <th>周期</th>
              <th>用户数</th>
              <th>金额</th>
              <th>状态</th>
              <th>支付时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="invoice in tenantStore.invoices" :key="invoice.id">
              <td>{{ invoice.invoiceNo }}</td>
              <td>{{ invoice.periodStart?.substring(0, 10) }} ~ {{ invoice.periodEnd?.substring(0, 10) }}</td>
              <td>{{ invoice.userCount }}</td>
              <td>¥{{ invoice.amount.toFixed(2) }}</td>
              <td>
                <n-tag :type="invoiceStatusType(invoice.status)" size="small">
                  {{ invoiceStatusLabel(invoice.status) }}
                </n-tag>
              </td>
              <td>{{ invoice.paidAt?.substring(0, 10) || '-' }}</td>
              <td>
                <n-button text type="primary" size="small" @click="handleView(invoice)">详情</n-button>
              </td>
            </tr>
            <tr v-if="tenantStore.invoices.length === 0">
              <td colspan="7" class="empty-row">暂无账单记录</td>
            </tr>
          </tbody>
        </n-table>
      </n-spin>
    </n-card>

    <n-modal
      v-model:show="showDetail"
      preset="dialog"
      title="账单详情"
      :show-icon="false"
      style="width: 520px"
    >
      <template v-if="currentInvoice">
        <n-descriptions label-placement="left" bordered :column="1">
          <n-descriptions-item label="账单编号">{{ currentInvoice.invoiceNo }}</n-descriptions-item>
          <n-descriptions-item label="计费周期">
            {{ currentInvoice.periodStart?.substring(0, 10) }} ~ {{ currentInvoice.periodEnd?.substring(0, 10) }}
          </n-descriptions-item>
          <n-descriptions-item label="用户数">{{ currentInvoice.userCount }}</n-descriptions-item>
          <n-descriptions-item label="金额">¥{{ currentInvoice.amount.toFixed(2) }}</n-descriptions-item>
          <n-descriptions-item label="状态">
            <n-tag :type="invoiceStatusType(currentInvoice.status)" size="small">
              {{ invoiceStatusLabel(currentInvoice.status) }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="支付方式">{{ currentInvoice.paymentMethod || '-' }}</n-descriptions-item>
          <n-descriptions-item label="支付时间">{{ currentInvoice.paidAt?.substring(0, 10) || '-' }}</n-descriptions-item>
        </n-descriptions>
      </template>
      <template #action>
        <n-button @click="showDetail = false">关闭</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  NCard,
  NTable,
  NTag,
  NButton,
  NModal,
  NDescriptions,
  NDescriptionsItem,
  NSpin,
  useMessage
} from 'naive-ui'
import { useTenantStore } from '@/stores/tenant'
import type { Invoice } from '@/types/tenant'
import { INVOICE_STATUS_MAP } from '@/utils/plan'

const message = useMessage()
const tenantStore = useTenantStore()

const loading = ref(false)
const showDetail = ref(false)
const currentInvoice = ref<Invoice | null>(null)

function invoiceStatusType(status: string) {
  return INVOICE_STATUS_MAP[status]?.type || 'default'
}

function invoiceStatusLabel(status: string) {
  return INVOICE_STATUS_MAP[status]?.label || status
}

function handleView(invoice: Invoice) {
  currentInvoice.value = invoice
  showDetail.value = true
}

onMounted(async () => {
  loading.value = true
  try {
    await tenantStore.fetchInvoices()
  } catch {
    message.error('获取账单列表失败')
  } finally {
    loading.value = false
  }
})
</script>

<style lang="scss" scoped>
.tenant-invoices-page {
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

.empty-row {
  text-align: center;
  color: $text-color-3;
  padding: 40px 0 !important;
}
</style>
