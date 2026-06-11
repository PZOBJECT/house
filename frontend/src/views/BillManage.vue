<template>
  <div class="bill-manage">
    <!-- 头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2>月度账单</h2>
        <span class="paid-badge">
          <svg class="badge-icon" viewBox="0 0 24 24" width="14" height="14"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41L9 16.17z" fill="currentColor"/></svg>
          <span class="badge-num">{{ paidCount }}</span>
          <span class="badge-text">已缴费</span>
        </span>
        <span class="unpaid-badge">
          <svg class="badge-icon" viewBox="0 0 24 24" width="14" height="14"><path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12 19 6.41z" fill="currentColor"/></svg>
          <span class="badge-num">{{ unpaidCount }}</span>
          <span class="badge-text">未缴费</span>
        </span>
      </div>
      <div class="filter-bar">
        <el-date-picker
          v-model="selectedMonth"
          type="month"
          placeholder="选择月份"
          format="YYYY年MM月"
          value-format="YYYY-MM"
          :clearable="false"
          class="filter-month"
          @change="onMonthChange"
        />
        <el-button type="primary" class="filter-btn" @click="handleGenerate">+ 生成当月账单</el-button>
      </div>
    </div>

    <!-- 统计栏 -->
    <div class="stats-bar">
      <div class="stat-item">
        <span class="stat-label">房间</span>
        <span class="stat-value">{{ bills.length }}<span class="stat-unit">间</span></span>
      </div>
      <div class="stat-divider" />
      <div class="stat-item">
        <span class="stat-label">⚡ 用电量</span>
        <span class="stat-value stat-purple">{{ calcTotal('elec_usage') }}<span class="stat-unit">度</span></span>
      </div>
      <div class="stat-divider" />
      <div class="stat-item">
        <span class="stat-label">电费合计</span>
        <span class="stat-value stat-blue">{{ calcTotal('elec_cost') }}<span class="stat-unit">元</span></span>
      </div>
      <div class="stat-divider" />
      <div class="stat-item">
        <span class="stat-label">💧 用水量</span>
        <span class="stat-value stat-aqua">{{ calcTotal('water_usage') }}<span class="stat-unit">吨</span></span>
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
        ref="tableRef"
        :data="bills"
        border
        stripe
        v-loading="loading"
        empty-text="暂无账单数据，请先生成账单"
        :row-class-name="tableRowClassName"
        height="100%"
      >
        <el-table-column prop="room_no" label="房间号" width="90" align="center" fixed />
        <el-table-column label="电表读数" align="center">
          <el-table-column prop="elec_last" label="上月" width="90" align="center" />
          <el-table-column prop="elec_current" label="本月" width="90" align="center">
            <template #default="{ row }">
              <template v-if="row.is_paid">
                {{ row.elec_current }}
              </template>
              <template v-else>
                <el-button type="primary" link size="small" @click="openMeterDialog(row, 'elec')">
                  {{ row.elec_current ?? '点击录入' }}
                </el-button>
              </template>
            </template>
          </el-table-column>
          <el-table-column prop="elec_usage" label="用量" width="80" align="center" />
          <el-table-column prop="elec_cost" label="电费" width="90" align="center" />
        </el-table-column>
        <el-table-column label="水表读数" align="center">
          <el-table-column prop="water_last" label="上月" width="90" align="center" />
          <el-table-column prop="water_current" label="本月" width="90" align="center">
            <template #default="{ row }">
              <template v-if="row.is_paid">
                {{ row.water_current }}
              </template>
              <template v-else>
                <el-button type="primary" link size="small" @click="openMeterDialog(row, 'water')">
                  {{ row.water_current ?? '点击录入' }}
                </el-button>
              </template>
            </template>
          </el-table-column>
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
        <el-table-column label="操作" width="220" align="center" fixed="right">
          <template #default="{ row }">
            <template v-if="row.is_paid">
              <el-button type="warning" size="small" @click="handleUnpay(row)">取消收费</el-button>
            </template>
            <template v-else>
              <el-button type="success" size="small" @click="handlePay(row)">标记收费</el-button>
            </template>
            <el-button type="primary" size="small" @click="openReceipt(row)">收据</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 修改读数对话框 -->
    <el-dialog v-model="meterDialog.visible" :title="meterDialog.title" width="400px">
      <el-form label-width="120px">
        <el-form-item label="上月读数">
          <el-input v-model="meterDialog.last" placeholder="请输入上月读数" />
        </el-form-item>
        <el-form-item label="本月读数">
          <el-input v-model="meterDialog.current" placeholder="请输入本月读数" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="meterDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="saveMeter">保存</el-button>
      </template>
    </el-dialog>

    <!-- 收据对话框 -->
    <ReceiptDialog ref="receiptDialogRef" />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { getBillList, generateBills, updateBill, payBill, unpayBill } from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'
import ReceiptDialog from './ReceiptDialog.vue'

const bills = ref([])
const loading = ref(false)
const receiptDialogRef = ref(null)
const filterFloor = ref(2)

const floorOptions = [
  { value: undefined, label: '全部', desc: '2-5楼' },
  { value: 2, label: '2楼', desc: '201-209' },
  { value: 3, label: '3楼', desc: '301-309' },
  { value: 4, label: '4楼', desc: '401-409' },
  { value: 5, label: '5楼', desc: '501-509' }
]

const now = new Date()
const selectedMonth = ref(`${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`)

const queryParams = computed(() => {
  const [y, m] = selectedMonth.value.split('-')
  return { year: parseInt(y), month: parseInt(m) }
})

const meterDialog = reactive({
  visible: false,
  title: '',
  billId: null,
  type: '',
  last: '',
  current: '',
  lastEditable: false
})

const paidCount = computed(() => bills.value.filter(b => b.is_paid).length)
const unpaidCount = computed(() => bills.value.filter(b => !b.is_paid).length)

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

function onMonthChange() {
  fetchBills()
}

function selectFloor(val) {
  filterFloor.value = val
  fetchBills()
}

async function fetchBills() {
  loading.value = true
  try {
    const params = { ...queryParams.value }
    if (filterFloor.value !== undefined) {
      params.floor = filterFloor.value
    }
    const res = await getBillList(params)
    bills.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function handleGenerate() {
  try {
    await ElMessageBox.confirm(
      `确认生成 ${queryParams.value.year}年${queryParams.value.month}月 的账单？`,
      '确认生成',
      { confirmButtonText: '确认', cancelButtonText: '取消', type: 'warning' }
    )
  } catch {
    return
  }
  await generateBills(queryParams.value.year, queryParams.value.month)
  ElMessage.success('账单生成成功')
  await fetchBills()
}

function openMeterDialog(row, type) {
  meterDialog.billId = row.id
  meterDialog.type = type
  if (type === 'elec') {
    meterDialog.title = `修改电表读数 - ${row.room_no}`
    meterDialog.lastEditable = true
    meterDialog.last = row.elec_last ?? ''
    meterDialog.current = row.elec_current ?? ''
  } else {
    meterDialog.title = `修改水表读数 - ${row.room_no}`
    meterDialog.lastEditable = true
    meterDialog.last = row.water_last ?? ''
    meterDialog.current = row.water_current ?? ''
  }
  meterDialog.visible = true
}

async function saveMeter() {
  const currentVal = parseFloat(meterDialog.current)
  if (isNaN(currentVal) || currentVal < 0) {
    ElMessage.warning('请输入有效的本月读数')
    return
  }

  const lastVal = meterDialog.last === '' ? null : parseFloat(meterDialog.last)
  if (lastVal !== null && (isNaN(lastVal) || lastVal < 0)) {
    ElMessage.warning('请输入有效的上月读数')
    return
  }

  if (lastVal !== null && currentVal < lastVal) {
    ElMessage.warning('本月读数不能低于上月读数')
    return
  }

  const payload = {}
  if (meterDialog.type === 'elec') {
    payload.elec_current = currentVal
    if (lastVal !== null) payload.elec_last = lastVal
  } else {
    payload.water_current = currentVal
    if (lastVal !== null) payload.water_last = lastVal
  }
  await updateBill(meterDialog.billId, payload)
  ElMessage.success('读数已更新')
  meterDialog.visible = false
  await fetchBills()
}

async function handlePay(row) {
  try {
    await ElMessageBox.confirm(
      `确认将 ${row.room_no} 标记为已收费？`,
      '确认收费',
      { confirmButtonText: '确认', cancelButtonText: '取消', type: 'warning' }
    )
  } catch {
    return
  }
  await payBill(row.id)
  ElMessage.success('已标记收费')
  await fetchBills()
}

async function handleUnpay(row) {
  try {
    await ElMessageBox.confirm(
      `确认取消 ${row.room_no} 的收费状态？`,
      '确认取消收费',
      { confirmButtonText: '确认', cancelButtonText: '取消', type: 'warning' }
    )
  } catch {
    return
  }
  await unpayBill(row.id)
  ElMessage.success('已取消收费')
  await fetchBills()
}

function openReceipt(row) {
  receiptDialogRef.value.open(row.id)
}

onMounted(fetchBills)
</script>

<style scoped>
.bill-manage {
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

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

/* 已缴费/未缴费徽章 */
.paid-badge,
.unpaid-badge {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  padding: 3px 12px 3px 10px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
  white-space: nowrap;
  line-height: 1.5;
}

.paid-badge {
  background: linear-gradient(135deg, #f0f9eb, #e8f5e0);
  color: #67c23a;
  border: 1px solid #c8e6b5;
  box-shadow: 0 1px 3px rgba(103, 194, 58, 0.15);
}

.unpaid-badge {
  background: linear-gradient(135deg, #fef0f0, #fde8e8);
  color: #f56c6c;
  border: 1px solid #fbc4c4;
  box-shadow: 0 1px 3px rgba(245, 108, 108, 0.15);
}

.badge-icon {
  flex-shrink: 0;
  opacity: 0.9;
}

.badge-num {
  font-size: 16px;
  font-weight: 800;
  min-width: 16px;
  text-align: center;
  margin-right: 1px;
}

.badge-text {
  font-size: 12px;
  font-weight: 500;
  opacity: 0.85;
}

/* 日期+按钮筛选栏 */
.filter-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 5px 10px;
  background: #f5f2ed;
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

.filter-btn {
  font-weight: 600;
  padding: 8px 16px;
  border-radius: 8px;
}

/* ── 统计栏 ── */
.stats-bar {
  display: flex;
  align-items: center;
  gap: 0;
  margin-bottom: 10px;
  flex-shrink: 0;
  background: #f5f2ed;
  border-radius: 10px;
  padding: 8px 12px;
  border: 1px solid #e8ecf1;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1px;
  flex: 1;
  min-width: 0;
}

.stat-label {
  font-size: 11px;
  font-weight: 500;
  color: #909399;
  letter-spacing: 0.2px;
  white-space: nowrap;
}

.stat-value {
  font-size: 16px;
  font-weight: 700;
  color: #303133;
  line-height: 1.3;
}

.stat-unit {
  font-size: 11px;
  font-weight: 500;
  margin-left: 2px;
  opacity: 0.7;
}

.stat-blue { color: #409eff; }
.stat-purple { color: #7c3aed; }
.stat-aqua { color: #0ea5e9; }
.stat-orange { color: #e6a23c; }
.stat-red { color: #f56c6c; }

.stat-divider {
  width: 1px;
  height: 30px;
  background: #e0e4ea;
  margin: 0 3px;
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
  background: #f3f1ec;
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
