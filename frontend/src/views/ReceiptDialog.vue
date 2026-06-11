<template>
  <el-dialog v-model="visible" title="房租收据" width="560px" :close-on-click-modal="false" @opened="onOpened">
    <div ref="receiptRef" class="receipt-container">
      <div class="receipt-title">房租收据</div>
      <div class="receipt-info">
        <div class="info-row">
          <span>房间号：</span>
          <span class="info-value">{{ receipt.room_no }}</span>
        </div>
        <div class="info-row">
          <span>账期：</span>
          <span class="info-value">{{ receipt.year }}年{{ receipt.month }}月</span>
        </div>
        <div class="info-row" v-if="receipt.paid_at">
          <span>收费日期：</span>
          <span class="info-value">{{ formatDate(receipt.paid_at) }}</span>
        </div>
      </div>
      <table class="receipt-table">
        <thead>
          <tr>
            <th>项目</th>
            <th>上月读数</th>
            <th>本月读数</th>
            <th>用量</th>
            <th>单价</th>
            <th>金额（元）</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>电费</td>
            <td>{{ receipt.elec_last }}</td>
            <td>{{ receipt.elec_current }}</td>
            <td>{{ receipt.elec_usage }}</td>
            <td>{{ receipt.elec_price }}元/度</td>
            <td>{{ receipt.elec_cost }}</td>
          </tr>
          <tr>
            <td>水费</td>
            <td>{{ receipt.water_last }}</td>
            <td>{{ receipt.water_current }}</td>
            <td>{{ receipt.water_usage }}</td>
            <td>{{ receipt.water_price }}元/吨</td>
            <td>{{ receipt.water_cost }}</td>
          </tr>
          <tr>
            <td>房租</td>
            <td>-</td>
            <td>-</td>
            <td>-</td>
            <td>{{ receipt.rent_price }}元/月</td>
            <td>{{ receipt.rent_cost }}</td>
          </tr>
        </tbody>
        <tfoot>
          <tr>
            <td colspan="5" class="total-label">合计</td>
            <td class="total-amount">{{ receipt.total_cost }}元</td>
          </tr>
        </tfoot>
      </table>
      <div class="receipt-footer">
        <div>收款人：___________</div>
        <div>日期：___________</div>
      </div>
    </div>
    <template #footer>
      <el-button @click="visible = false">关闭</el-button>
      <el-button type="primary" @click="downloadReceipt" :loading="downloading">下载收据图片</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { getReceipt } from '../api'
import html2canvas from 'html2canvas'
import { saveAs } from 'file-saver'

const visible = ref(false)
const downloading = ref(false)
const receiptRef = ref(null)

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const receipt = reactive({
  room_no: '',
  year: '',
  month: '',
  elec_last: '',
  elec_current: '',
  elec_usage: '',
  elec_cost: '',
  elec_price: '',
  water_last: '',
  water_current: '',
  water_usage: '',
  water_cost: '',
  water_price: '',
  rent_price: '',
  rent_cost: '',
  total_cost: '',
  paid_at: ''
})

function onOpened() {
  // 对话框打开后无需额外操作
}

async function open(billId) {
  visible.value = true
  const res = await getReceipt(billId)
  Object.assign(receipt, res.data)
}

async function downloadReceipt() {
  if (!receiptRef.value) return
  downloading.value = true
  try {
    const canvas = await html2canvas(receiptRef.value, {
      backgroundColor: '#ffffff',
      scale: 2
    })
    const blob = await new Promise((resolve) => canvas.toBlob(resolve, 'image/png'))
    const filename = `收据_${receipt.room_no}_${receipt.year}年${receipt.month}月.png`
    saveAs(blob, filename)
  } catch (e) {
    console.error('生成收据图片失败', e)
  } finally {
    downloading.value = false
  }
}

defineExpose({ open })
</script>

<style scoped>
.receipt-container {
  padding: 28px 24px;
  background: #fff;
  font-size: 16px;
  color: #333;
}

.receipt-title {
  text-align: center;
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 24px;
  letter-spacing: 4px;
}

.receipt-info {
  margin-bottom: 20px;
}

.info-row {
  margin-bottom: 8px;
  font-size: 16px;
}

.info-value {
  font-weight: bold;
}

.receipt-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 20px;
}

.receipt-table th,
.receipt-table td {
  border: 1px solid #dcdfe6;
  padding: 10px 8px;
  text-align: center;
  font-size: 15px;
}

.receipt-table thead th {
  background: #f5f2ed;
  font-weight: bold;
}

.receipt-table tfoot td {
  font-weight: bold;
}

.total-label {
  text-align: right;
  padding-right: 14px;
}

.total-amount {
  color: #e6a23c;
  font-size: 18px;
}

.receipt-footer {
  display: flex;
  justify-content: space-between;
  margin-top: 28px;
  padding: 0 10px;
  font-size: 15px;
}
</style>
