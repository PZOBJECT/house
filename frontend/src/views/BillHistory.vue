<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px">
      <h2 style="margin: 0">历史账单</h2>
      <div style="display: flex; gap: 12px; align-items: center; flex-wrap: wrap">
        <el-date-picker
          v-model="selectedMonth"
          type="month"
          placeholder="选择月份（可清空）"
          format="YYYY年MM月"
          value-format="YYYY-MM"
          clearable
          style="width: 180px"
          @change="fetchBills"
        />
        <el-select v-model="filterPaid" placeholder="收费状态" clearable style="width: 130px" @change="fetchBills">
          <el-option label="全部" :value="undefined" />
          <el-option label="已收费" :value="1" />
          <el-option label="未收费" :value="0" />
        </el-select>
        <el-select v-model="filterFloor" placeholder="楼层筛选" clearable style="width: 120px" @change="fetchBills">
          <el-option label="全部" :value="undefined" />
          <el-option label="2楼" :value="2" />
          <el-option label="3楼" :value="3" />
          <el-option label="4楼" :value="4" />
          <el-option label="5楼" :value="5" />
        </el-select>
      </div>
    </div>

    <el-table :data="bills" border stripe v-loading="loading" empty-text="暂无历史账单" :row-class-name="tableRowClassName" show-summary :summary-method="getSummaries">
      <el-table-column prop="room_no" label="房间号" width="90" align="center" />
      <el-table-column label="电表读数" align="center">
        <el-table-column prop="elec_last" label="上月" width="90" align="center" />
        <el-table-column prop="elec_current" label="本月" width="90" align="center" />
        <el-table-column prop="elec_usage" label="用量" width="80" align="center" />
        <el-table-column prop="elec_cost" label="电费" width="90" align="center" />
      </el-table-column>
      <el-table-column label="水表读数" align="center">
        <el-table-column prop="water_last" label="上月" width="90" align="center" />
        <el-table-column prop="water_current" label="本月" width="90" align="center" />
        <el-table-column prop="water_usage" label="用量" width="80" align="center" />
        <el-table-column prop="water_cost" label="水费" width="90" align="center" />
      </el-table-column>
      <el-table-column prop="rent_cost" label="房租（元）" width="110" align="center" />
      <el-table-column prop="total_cost" label="总费用（元）" width="120" align="center">
        <template #default="{ row }">
          <span style="font-weight: bold; color: #e6a23c">{{ row.total_cost }}</span>
        </template>
      </el-table-column>
      <el-table-column label="收费状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_paid ? 'success' : 'danger'">
            {{ row.is_paid ? '已收费' : '未收费' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="收费日期" width="160" align="center">
        <template #default="{ row }">
          {{ row.paid_at || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100" align="center">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="openReceipt(row)">收据</el-button>
        </template>
      </el-table-column>
    </el-table>

    <ReceiptDialog ref="receiptDialogRef" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getBillList } from '../api'
import ReceiptDialog from './ReceiptDialog.vue'

const bills = ref([])
const loading = ref(false)
const receiptDialogRef = ref(null)
const selectedMonth = ref('')
const filterPaid = ref(undefined)
const filterFloor = ref(undefined)

function tableRowClassName({ row }) {
  return row.is_paid ? 'paid-row' : ''
}

function getSummaries({ columns, data }) {
  const sums = []
  columns.forEach((column, index) => {
    if (index === 0) {
      sums[index] = '合计'
      return
    }
    const prop = column.property
    if (['elec_cost', 'water_cost', 'rent_cost', 'total_cost'].includes(prop)) {
      let total = 0
      data.forEach(row => {
        const val = Number(row[prop])
        if (!isNaN(val)) {
          total += val
        }
      })
      sums[index] = total.toFixed(2)
    } else {
      sums[index] = ''
    }
  })
  sums[0] = `共 ${data.length} 间`
  return sums
}

async function fetchBills() {
  loading.value = true
  try {
    const params = {}
    if (selectedMonth.value) {
      const [y, m] = selectedMonth.value.split('-')
      params.year = parseInt(y)
      params.month = parseInt(m)
    }
    if (filterPaid.value !== undefined) {
      params.is_paid = filterPaid.value
    }
    if (filterFloor.value !== undefined) {
      params.floor = filterFloor.value
    }
    const res = await getBillList(params)
    bills.value = res.data || []
  } finally {
    loading.value = false
  }
}

function openReceipt(row) {
  receiptDialogRef.value.open(row.id)
}

onMounted(fetchBills)
</script>

<style scoped>
:deep(.paid-row) {
  background-color: #f0f0f0;
  color: #666;
}
</style>
