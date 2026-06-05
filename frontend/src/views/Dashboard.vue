<template>
  <div class="dashboard">
    <!-- 顶部标题 -->
    <div class="dash-header">
      <div class="dash-title">
        <svg viewBox="0 0 24 24" width="24" height="24">
          <path d="M19 9.3V4h-3v2.6L12 3 2 12h3v8h5v-6h4v6h5v-8h3l-3-2.7z" fill="currentColor"/>
        </svg>
        <span>{{ currentMonth }} 数据概览</span>
      </div>
      <div class="dash-time">{{ nowStr }}</div>
    </div>

    <!-- KPI 卡片行 -->
    <div class="kpi-row">
      <div class="kpi-card card-income">
        <div class="kpi-icon-wrap">
          <svg viewBox="0 0 24 24" width="28" height="28"><path d="M11.8 10.9c-2.27-.59-3-1.2-3-2.15 0-1.09 1.01-1.85 2.7-1.85 1.78 0 2.44.85 2.5 2.1h2.21c-.07-1.72-1.12-3.3-3.21-3.81V3h-3v2.16c-1.94.42-3.5 1.68-3.5 3.61 0 2.31 1.91 3.46 4.7 4.13 2.5.6 3 1.48 3 2.41 0 .69-.49 1.79-2.7 1.79-2.06 0-2.87-.92-2.98-2.1h-2.2c.12 2.19 1.76 3.42 3.68 3.83V21h3v-2.15c1.95-.37 3.5-1.5 3.5-3.55 0-2.84-2.43-3.81-4.7-4.4z" fill="currentColor"/></svg>
        </div>
        <div class="kpi-body">
          <span class="kpi-label">当月总收入</span>
          <span class="kpi-value">{{ totalIncome }}<span class="kpi-unit">元</span></span>
        </div>
        <div class="kpi-trend">
          较上月 <span :class="trendClass">{{ trendText }}</span>
        </div>
      </div>

      <div class="kpi-card card-paid">
        <div class="kpi-icon-wrap">
          <svg viewBox="0 0 24 24" width="28" height="28"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41L9 16.17z" fill="currentColor"/></svg>
        </div>
        <div class="kpi-body">
          <span class="kpi-label">已缴费</span>
          <span class="kpi-value">{{ paidCount }}<span class="kpi-unit">/{{ totalCount }}间</span></span>
        </div>
        <div class="kpi-progress">
          <div class="progress-track">
            <div class="progress-fill progress-green" :style="{ width: paidPercent + '%' }" />
          </div>
          <span class="progress-text">{{ paidPercent }}%</span>
        </div>
      </div>

      <div class="kpi-card card-unpaid" style="cursor: pointer" @click="showUnpaidDialog">
        <div class="kpi-icon-wrap">
          <svg viewBox="0 0 24 24" width="28" height="28"><path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12 19 6.41z" fill="currentColor"/></svg>
        </div>
        <div class="kpi-body">
          <span class="kpi-label">未缴费 — 点击查看</span>
          <span class="kpi-value">{{ unpaidCount }}<span class="kpi-unit">间</span></span>
        </div>
        <div class="kpi-hint">待催收</div>
      </div>

      <div class="kpi-card card-occupancy">
        <div class="kpi-icon-wrap">
          <svg viewBox="0 0 24 24" width="28" height="28"><path d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z" fill="currentColor"/></svg>
        </div>
        <div class="kpi-body">
          <span class="kpi-label">出租率</span>
          <span class="kpi-value">{{ rentedCount }}<span class="kpi-unit">/{{ totalRooms }}间</span></span>
        </div>
        <div class="kpi-progress">
          <div class="progress-track">
            <div class="progress-fill progress-blue" :style="{ width: occupancyPercent + '%' }" />
          </div>
          <span class="progress-text">{{ occupancyPercent }}%</span>
        </div>
      </div>
    </div>

    <!-- 收入分解行 -->
    <div class="detail-row">
      <div class="detail-card card-elec">
        <div class="detail-card-header">
          <svg viewBox="0 0 24 24" width="22" height="22"><path d="M13.5 2L4.2 13.5h5.3L10 22l9.3-11.5H14L13.5 2z" fill="currentColor"/></svg>
          <span>电费</span>
        </div>
        <div class="detail-card-body">
          <div class="detail-metric">
            <span class="dm-label">用量</span>
            <span class="dm-value">{{ totalElecUsage }}<span class="dm-unit">度</span></span>
          </div>
          <div class="detail-divider" />
          <div class="detail-metric">
            <span class="dm-label">收入</span>
            <span class="dm-value dm-money">{{ totalElecCost }}<span class="dm-unit">元</span></span>
          </div>
        </div>
      </div>

      <div class="detail-card card-water">
        <div class="detail-card-header">
          <svg viewBox="0 0 24 24" width="22" height="22"><path d="M12 2c-5.33 4.55-8 8.48-8 11.8 0 4.98 3.8 8.2 8 8.2s8-3.22 8-8.2c0-3.32-2.67-7.25-8-11.8zM7.83 14c.37 0 .67.26.74.62.41 2.22 2.28 2.98 3.64 2.87.43-.02.79.32.79.75 0 .4-.32.73-.72.75-2.13.13-4.62-1.09-5.19-4.12-.07-.36.22-.69.58-.69.16 0 .31.06.4.17.09.11.14.25.16.4-.03.01-.06.02-.08.05-.1.14-.21.3-.32.5z" fill="currentColor"/></svg>
          <span>水费</span>
        </div>
        <div class="detail-card-body">
          <div class="detail-metric">
            <span class="dm-label">用量</span>
            <span class="dm-value">{{ totalWaterUsage }}<span class="dm-unit">吨</span></span>
          </div>
          <div class="detail-divider" />
          <div class="detail-metric">
            <span class="dm-label">收入</span>
            <span class="dm-value dm-money">{{ totalWaterCost }}<span class="dm-unit">元</span></span>
          </div>
        </div>
      </div>

      <div class="detail-card card-rent">
        <div class="detail-card-header">
          <svg viewBox="0 0 24 24" width="22" height="22"><path d="M20 4H4c-1.11 0-1.99.89-1.99 2L2 18c0 1.11.89 2 2 2h16c1.11 0 2-.89 2-2V6c0-1.11-.89-2-2-2zm0 14H4v-6h16v6zm0-10H4V6h16v2z" fill="currentColor"/></svg>
          <span>房租</span>
        </div>
        <div class="detail-card-body">
          <div class="detail-metric">
            <span class="dm-label">房间数</span>
            <span class="dm-value">{{ totalCount }}<span class="dm-unit">间</span></span>
          </div>
          <div class="detail-divider" />
          <div class="detail-metric">
            <span class="dm-label">收入</span>
            <span class="dm-value dm-money">{{ totalRentCost }}<span class="dm-unit">元</span></span>
          </div>
        </div>
      </div>
    </div>

    <!-- 未缴费列表弹窗 -->
    <el-dialog v-model="unpaidDialogVisible" title="待缴费房间" width="500px" top="10vh">
      <div class="unpaid-list" v-if="unpaidList.length > 0">
        <div class="unpaid-header">
          <span class="uh-room">房间号</span>
          <span class="uh-cost">总费用</span>
          <span class="uh-action">操作</span>
        </div>
        <div class="unpaid-row" v-for="item in unpaidList" :key="item.id">
          <span class="ur-room">{{ item.room_no }}</span>
          <span class="ur-cost">{{ item.total_cost }}元</span>
          <span class="ur-action">
            <el-button type="success" size="small" @click="markAsPaid(item)">标记已缴费</el-button>
          </span>
        </div>
      </div>
      <div class="unpaid-empty" v-else>🎉 本月所有房间已完成缴费</div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getBillList, getRooms, payBill } from '../api'
import { ElMessage } from 'element-plus'

const bills = ref([])
const allRooms = ref([])
const now = new Date()
const currentYear = now.getFullYear()
const currentMonthNum = now.getMonth() + 1
const currentMonth = `${currentYear}年${String(currentMonthNum).padStart(2, '0')}月`
const nowStr = ref(getNowStr())
const unpaidDialogVisible = ref(false)

function getNowStr() {
  const d = new Date()
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日 ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

// 当月账单统计
const totalCount = computed(() => bills.value.length)
const paidCount = computed(() => bills.value.filter(b => b.is_paid).length)
const unpaidCount = computed(() => bills.value.filter(b => !b.is_paid).length)
const paidPercent = computed(() => totalCount.value > 0 ? Math.round((paidCount.value / totalCount.value) * 100) : 0)

const totalIncome = computed(() => {
  let sum = 0
  bills.value.forEach(b => { const v = Number(b.total_cost); if (!isNaN(v)) sum += v })
  return sum.toFixed(2)
})

const totalElecCost = computed(() => {
  let sum = 0
  bills.value.forEach(b => { const v = Number(b.elec_cost); if (!isNaN(v)) sum += v })
  return sum.toFixed(2)
})

const totalWaterCost = computed(() => {
  let sum = 0
  bills.value.forEach(b => { const v = Number(b.water_cost); if (!isNaN(v)) sum += v })
  return sum.toFixed(2)
})

const totalRentCost = computed(() => {
  let sum = 0
  bills.value.forEach(b => { const v = Number(b.rent_cost); if (!isNaN(v)) sum += v })
  return sum.toFixed(2)
})

const totalElecUsage = computed(() => {
  let sum = 0
  bills.value.forEach(b => { const v = Number(b.elec_usage); if (!isNaN(v)) sum += v })
  return sum.toFixed(1)
})

const totalWaterUsage = computed(() => {
  let sum = 0
  bills.value.forEach(b => { const v = Number(b.water_usage); if (!isNaN(v)) sum += v })
  return sum.toFixed(1)
})

// 出租率
const rentedCount = computed(() => allRooms.value.filter(r => r.is_rented === 1).length)
const totalRooms = computed(() => allRooms.value.length)
const occupancyPercent = computed(() => totalRooms.value > 0 ? Math.round((rentedCount.value / totalRooms.value) * 100) : 0)

// 待缴费列表
const unpaidList = computed(() =>
  [...bills.value]
    .filter(b => !b.is_paid)
)

const trendClass = computed(() => 'trend-up')
const trendText = computed(() => '—')

function showUnpaidDialog() {
  unpaidDialogVisible.value = true
}

async function markAsPaid(item) {
  try {
    await payBill(item.id)
    ElMessage.success(`${item.room_no} 已标记为已缴费`)
    unpaidDialogVisible.value = false
    // 刷新数据
    await fetchData()
  } catch {
    ElMessage.error('操作失败，请重试')
  }
} 

async function fetchData() {
  try {
    const [billRes, roomRes] = await Promise.all([
      getBillList({ year: currentYear, month: currentMonthNum }),
      getRooms(undefined, undefined)
    ])
    bills.value = billRes.data || []
    allRooms.value = roomRes.data || []
  } catch { /* ignored */ }
}

onMounted(fetchData)
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 16px;
  height: calc(100vh - 24px);
  overflow-y: auto;
  padding: 2px 0;
}

/* ── 顶部标题 ── */
.dash-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
}

.dash-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 20px;
  font-weight: 700;
  color: #303133;
}

.dash-title svg {
  color: #409eff;
}

.dash-time {
  font-size: 13px;
  color: #909399;
  font-weight: 500;
}

/* ── KPI 卡片行 ── */
.kpi-row {
  display: grid;
  grid-template-columns: 1.3fr 1fr 1fr 1fr;
  gap: 14px;
  flex-shrink: 0;
}

.kpi-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 20px;
  border-radius: 14px;
  position: relative;
  overflow: hidden;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.kpi-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
}

.card-income {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}

.card-paid {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
  color: #fff;
}

.card-unpaid {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: #fff;
}

.card-occupancy {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: #fff;
}

.kpi-icon-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.2);
  flex-shrink: 0;
}

.kpi-body {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.kpi-label {
  font-size: 13px;
  font-weight: 500;
  opacity: 0.85;
  letter-spacing: 0.3px;
}

.kpi-value {
  font-size: 26px;
  font-weight: 800;
  line-height: 1.2;
}

.kpi-unit {
  font-size: 14px;
  font-weight: 500;
  opacity: 0.7;
  margin-left: 3px;
}

.kpi-trend {
  position: absolute;
  right: 16px;
  bottom: 12px;
  font-size: 12px;
  opacity: 0.8;
}

.kpi-hint {
  position: absolute;
  right: 16px;
  bottom: 12px;
  font-size: 12px;
  font-weight: 600;
  opacity: 0.7;
  background: rgba(0, 0, 0, 0.15);
  padding: 2px 10px;
  border-radius: 10px;
}

.trend-up { color: #67c23a; }

/* 进度条 */
.kpi-progress {
  position: absolute;
  right: 16px;
  bottom: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.progress-track {
  width: 60px;
  height: 6px;
  background: rgba(255, 255, 255, 0.25);
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.6s ease;
}

.progress-green { background: #fff; }
.progress-blue { background: #fff; }

.progress-text {
  font-size: 12px;
  font-weight: 700;
  opacity: 0.9;
}

/* ── 收入分解行 ── */
.detail-row {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 14px;
  flex-shrink: 0;
}

.detail-card {
  background: #fff;
  border-radius: 14px;
  padding: 18px 20px;
  border: 1px solid #ebeef5;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.detail-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.detail-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 700;
  padding-bottom: 14px;
  border-bottom: 1px solid #f0f2f5;
  margin-bottom: 14px;
}

.card-elec .detail-card-header { color: #e6a23c; }
.card-elec { background: linear-gradient(135deg, #fffbf0, #fff8e1); }

.card-water .detail-card-header { color: #409eff; }
.card-water { background: linear-gradient(135deg, #f0f9ff, #e6f4ff); }

.card-rent .detail-card-header { color: #67c23a; }
.card-rent { background: linear-gradient(135deg, #f0faf0, #e8f5e0); }

.detail-card-body {
  display: flex;
  align-items: center;
  gap: 0;
}

.detail-metric {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.dm-label {
  font-size: 12px;
  color: #909399;
  font-weight: 500;
}

.dm-value {
  font-size: 24px;
  font-weight: 800;
  color: #303133;
}

.dm-unit {
  font-size: 13px;
  font-weight: 500;
  color: #909399;
  margin-left: 2px;
}

.dm-money {
  color: #e6a23c;
}

.detail-divider {
  width: 1px;
  height: 40px;
  background: #ebeef5;
  flex-shrink: 0;
}

/* ── 未缴费弹窗列表 ── */
.unpaid-list {
  display: flex;
  flex-direction: column;
  max-height: 50vh;
  overflow-y: auto;
}

.unpaid-header {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  padding: 10px 12px;
  background: #f5f7fa;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  color: #909399;
  text-align: center;
  margin-bottom: 4px;
  flex-shrink: 0;
}

.unpaid-row {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  align-items: center;
  padding: 10px 12px;
  border-bottom: 1px solid #f0f2f5;
  text-align: center;
}

.unpaid-row:last-child {
  border-bottom: none;
}

.ur-room {
  font-size: 15px;
  font-weight: 700;
  color: #303133;
}

.ur-cost {
  font-size: 15px;
  font-weight: 700;
  color: #e6a23c;
}

.ur-action {
  display: flex;
  justify-content: center;
}

.unpaid-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 80px;
  color: #67c23a;
  font-size: 15px;
  font-weight: 500;
}
</style>
