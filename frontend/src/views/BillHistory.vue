<template>
  <div class="bill-history">
    <!-- 头部 -->
    <div class="page-header">
      <h2>历史账单</h2>
      <div class="filter-bar">
        <el-date-picker
          v-model="selectedMonth"
          type="month"
          placeholder="选择月份（可清空）"
          format="YYYY年MM月"
          value-format="YYYY-MM"
          clearable
          class="filter-month"
          @change="onFilterChange"
        />
        <el-select v-model="filterPaid" placeholder="收费状态" clearable class="filter-select" @change="onFilterChange">
          <el-option label="全部" :value="undefined" />
          <el-option label="已收费" :value="1" />
          <el-option label="未收费" :value="0" />
        </el-select>
      </div>
    </div>

    <!-- 统计栏 -->
    <div class="stats-bar">
      <div class="stat-item">
        <span class="stat-label">记录</span>
        <span class="stat-value">{{ bills.length }}<span class="stat-unit">条</span></span>
      </div>
      <div class="stat-divider" />
      <div class="stat-item">
        <span class="stat-label">电费合计</span>
        <span class="stat-value stat-blue">{{ calcTotal('elec_cost') }}<span class="stat-unit">元</span></span>
      </div>
      <div class="stat-divider" />
      <div class="stat-item">
        <span class="stat-label">水费合计</span>
        <span class="stat-value stat-blue">{{ calcTotal('water_cost') }}<span class="stat-unit">元</span></span>
      </div>
      <div class="stat-divider" />
      <div class="stat-item">
        <span class="stat-label">房租合计</span>
        <span class="stat-value stat-orange">{{ calcTotal('rent_cost') }}<span class="stat-unit">元</span></span>
      </div>
      <div class="stat-divider" />
      <div class="stat-item">
        <span class="stat-label">总合计</span>
        <span class="stat-value stat-red">{{ calcTotal('total_cost') }}<span class="stat-unit">元</span></span>
      </div>
    </div>

    <!-- 楼层筛选 — 大卡片按钮 -->
    <div class="floor-filter">
      <span class="floor-label">
        <svg class="floor-icon" viewBox="0 0 24 24" width="16" height="16">
          <path d="M19 9.3V4h-3v2.6L12 3 2 12h3v8h5v-6h4v6h5v-8h3l-3-2.7z" fill="currentColor"/>
        </svg>
        楼层
      </span>
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
        :data="bills"
        border
        stripe
        v-loading="loading"
        empty-text="暂无历史账单"
        :row-class-name="tableRowClassName"
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
const now = new Date()
const selectedMonth = ref(`${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`)
const filterPaid = ref(undefined)
const filterFloor = ref(undefined)

const floorOptions = [
  { value: undefined, label: '全部', desc: '2-5楼' },
  { value: 2, label: '2楼', desc: '201-209' },
  { value: 3, label: '3楼', desc: '301-309' },
  { value: 4, label: '4楼', desc: '401-409' },
  { value: 5, label: '5楼', desc: '501-509' }
]

function calcTotal(prop) {
  let total = 0
  bills.value.forEach(row => {
    const val = Number(row[prop])
    if (!isNaN(val)) total += val
  })
  return total.toFixed(2)
}

function tableRowClassName({ row }) {
  return row.is_paid ? 'paid-row' : ''
}

function onFilterChange() {
  fetchBills()
}

function selectFloor(val) {
  filterFloor.value = val
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

/* ── 头部 ── */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 8px;
  flex-shrink: 0;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

/* 筛选栏容器 */
.filter-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 5px 10px;
  background: #e8edf3;
  border-radius: 10px;
  border: 1px solid #e8ecf1;
}

.filter-month :deep(.el-input__wrapper) {
  background: #fff;
  border-radius: 8px;
  box-shadow: none;
  border: 1px solid #e0e4ea;
}

.filter-month :deep(.el-input__wrapper:hover) {
  border-color: #409eff;
}

.filter-select {
  width: 130px;
}

.filter-select :deep(.el-input__wrapper) {
  background: #fff;
  border-radius: 8px;
  box-shadow: none;
  border: 1px solid #e0e4ea;
}

.filter-select :deep(.el-input__wrapper:hover) {
  border-color: #409eff;
}

/* ── 统计栏 ── */
.stats-bar {
  display: flex;
  align-items: center;
  gap: 0;
  margin-bottom: 10px;
  flex-shrink: 0;
  background: #e8edf3;
  border-radius: 10px;
  padding: 10px 16px;
  border: 1px solid #e8ecf1;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  flex: 1;
  min-width: 0;
}

.stat-label {
  font-size: 12px;
  font-weight: 500;
  color: #909399;
  letter-spacing: 0.3px;
}

.stat-value {
  font-size: 18px;
  font-weight: 700;
  color: #303133;
  line-height: 1.3;
}

.stat-unit {
  font-size: 12px;
  font-weight: 500;
  margin-left: 2px;
  opacity: 0.7;
}

.stat-blue { color: #409eff; }
.stat-orange { color: #e6a23c; }
.stat-red { color: #f56c6c; }

.stat-divider {
  width: 1px;
  height: 34px;
  background: #e0e4ea;
  margin: 0 4px;
  flex-shrink: 0;
}

/* ── 楼层卡片筛选 ── */
.floor-filter {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
  flex-shrink: 0;
}

.floor-label {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 6px 13px;
  background: #edf2f7;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #4a5568;
  white-space: nowrap;
  letter-spacing: 0.5px;
}

.floor-icon {
  flex-shrink: 0;
  opacity: 0.8;
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
  background: #e2e8f0;
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

/* ── 表格 ── */
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

:deep(.paid-row) {
  background-color: #f4f4f5 !important;
  color: #909399;
}
</style>
