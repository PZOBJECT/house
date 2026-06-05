<template>
  <div class="bill-history">
    <!-- 头部 -->
    <div class="page-header">
      <h2>历史账单</h2>
      <div class="header-actions">
        <el-date-picker
          v-model="selectedMonth"
          type="month"
          placeholder="选择月份（可清空）"
          format="YYYY年MM月"
          value-format="YYYY-MM"
          clearable
          style="width: 180px"
          @change="onFilterChange"
        />
        <el-select v-model="filterPaid" placeholder="收费状态" clearable style="width: 130px" @change="onFilterChange">
          <el-option label="全部" :value="undefined" />
          <el-option label="已收费" :value="1" />
          <el-option label="未收费" :value="0" />
        </el-select>
      </div>
    </div>

    <!-- 楼层筛选 — 大卡片按钮 -->
    <div class="floor-filter">
      <span class="floor-label">楼层</span>
      <div class="floor-buttons">
        <div
          v-for="f in floorOptions"
          :key="f.value"
          class="floor-card"
          :class="{ active: filterFloor === f.value }"
          @click="selectFloor(f.value)"
        >
          <span class="floor-num">{{ f.label }}</span>
          <span v-if="f.desc" class="floor-desc">{{ f.desc }}</span>
          <span v-if="filterFloor === f.value" class="floor-check">✓</span>
        </div>
      </div>
    </div>

    <!-- 表格区域 -->
    <div class="table-wrapper">
      <el-table
        :data="pagedBills"
        border
        stripe
        v-loading="loading"
        empty-text="暂无历史账单"
        :row-class-name="tableRowClassName"
        show-summary
        :summary-method="getSummaries"
        height="100%"
      >
        <el-table-column prop="room_no" label="房间号" width="90" align="center" />
        <el-table-column label="归属月份" width="100" align="center">
          <template #default="{ row }">
            {{ row.year }}年{{ String(row.month).padStart(2, '0') }}月
          </template>
        </el-table-column>
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
            {{ formatDate(row.paid_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" align="center">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="openReceipt(row)">收据</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 分页 -->
    <div class="page-footer">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="9"
        :total="bills.length"
        layout="total, prev, pager, next"
        background
        small
      />
    </div>

    <ReceiptDialog ref="receiptDialogRef" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getBillList } from '../api'
import ReceiptDialog from './ReceiptDialog.vue'

const bills = ref([])
const loading = ref(false)
const receiptDialogRef = ref(null)
const selectedMonth = ref('')
const filterPaid = ref(undefined)
const filterFloor = ref(undefined)
const currentPage = ref(1)

const floorOptions = [
  { value: undefined, label: '全部', desc: '2-5楼' },
  { value: 2, label: '2楼', desc: '201-209' },
  { value: 3, label: '3楼', desc: '301-309' },
  { value: 4, label: '4楼', desc: '401-409' },
  { value: 5, label: '5楼', desc: '501-509' }
]

const pagedBills = computed(() => {
  const start = (currentPage.value - 1) * 9
  return bills.value.slice(start, start + 9)
})

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

function onFilterChange() {
  currentPage.value = 1
  fetchBills()
}

function selectFloor(val) {
  filterFloor.value = val
  currentPage.value = 1
  fetchBills()
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

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

function openReceipt(row) {
  receiptDialogRef.value.open(row.id)
}

onMounted(fetchBills)
</script>

<style scoped>
.bill-history {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 24px);
  padding: 0 2px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 8px;
  flex-shrink: 0;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

/* ── 楼层卡片筛选 ── */
.floor-filter {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 10px;
  flex-shrink: 0;
}

.floor-label {
  font-size: 15px;
  font-weight: 600;
  color: #606266;
  white-space: nowrap;
}

.floor-buttons {
  display: flex;
  gap: 10px;
}

.floor-card {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 88px;
  height: 56px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  user-select: none;
  background: #f5f7fa;
  border: 2px solid transparent;
  color: #606266;
}

.floor-card:hover {
  background: #ecf5ff;
  border-color: #b3d8ff;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.12);
}

.floor-card.active {
  background: #409eff;
  border-color: #409eff;
  color: #fff;
  box-shadow: 0 3px 10px rgba(64, 158, 255, 0.3);
  transform: translateY(-1px);
}

.floor-num {
  font-size: 17px;
  font-weight: 700;
  line-height: 1.2;
}

.floor-desc {
  font-size: 11px;
  opacity: 0.75;
  line-height: 1.2;
}

.floor-card.active .floor-desc {
  opacity: 0.9;
}

.floor-check {
  position: absolute;
  top: -6px;
  right: -6px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #67c23a;
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 1px 4px rgba(103, 194, 58, 0.4);
}

.table-wrapper {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.table-wrapper :deep(.el-table__row) {
  height: 48px;
}

.table-wrapper :deep(.el-table__row td) {
  padding: 6px 0;
}

.page-footer {
  display: flex;
  justify-content: center;
  padding: 12px 0 4px;
  flex-shrink: 0;
}

:deep(.paid-row) {
  background-color: #f4f4f5 !important;
  color: #909399;
}
</style>
