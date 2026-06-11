<template>
  <div class="room-manage">
    <!-- 头部 -->
    <div class="page-header">
      <h2>房间管理</h2>
      <div class="header-actions">
        <el-radio-group v-model="filterRented" @change="fetchRooms" class="status-filter">
          <el-radio-button :value="undefined">全部</el-radio-button>
          <el-radio-button :value="1">已租</el-radio-button>
          <el-radio-button :value="0">未租</el-radio-button>
        </el-radio-group>
      </div>
    </div>

    <!-- 出租统计栏 -->
    <div class="stats-bar">
      <div class="stat-item stat-total">
        <svg class="stat-icon" viewBox="0 0 24 24" width="18" height="18">
          <path d="M19 9.3V4h-3v2.6L12 3 2 12h3v8h5v-6h4v6h5v-8h3l-3-2.7z" fill="currentColor"/>
        </svg>
        <div class="stat-body">
          <span class="stat-label">总房间</span>
          <span class="stat-value">{{ statsTotal }}<span class="stat-unit">间</span></span>
        </div>
      </div>
      <div class="stat-divider" />
      <div class="stat-item stat-rented">
        <svg class="stat-icon" viewBox="0 0 24 24" width="18" height="18">
          <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41L9 16.17z" fill="currentColor"/>
        </svg>
        <div class="stat-body">
          <span class="stat-label">已租</span>
          <span class="stat-value">{{ statsRented }}<span class="stat-unit">间</span></span>
        </div>
      </div>
      <div class="stat-divider" />
      <div class="stat-item stat-unrented">
        <svg class="stat-icon" viewBox="0 0 24 24" width="18" height="18">
          <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12 19 6.41z" fill="currentColor"/>
        </svg>
        <div class="stat-body">
          <span class="stat-label">未租</span>
          <span class="stat-value">{{ statsUnrented }}<span class="stat-unit">间</span></span>
        </div>
      </div>
      <div class="stat-divider" />
      <div class="stat-item stat-rate">
        <svg class="stat-icon" viewBox="0 0 24 24" width="18" height="18">
          <path d="M11.99 2C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8zm.5-13H11v6l5.25 3.15.75-1.23-4.5-2.67z" fill="currentColor"/>
        </svg>
        <div class="stat-body">
          <span class="stat-label">出租率</span>
          <span class="stat-value">{{ statsRate }}<span class="stat-unit">%</span></span>
        </div>
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
          <span class="floor-desc">{{ f.desc }}</span>
          <span v-if="filterFloor === f.value" class="floor-check">✓</span>
        </div>
      </div>
    </div>

    <!-- 表格区域 -->
    <div class="table-wrapper">
      <el-table
        :data="rooms"
        border
        stripe
        v-loading="loading"
        empty-text="暂无房间数据"
        height="100%"
      >
        <el-table-column prop="room_no" label="房间号" width="120" align="center">
          <template #default="{ row }">
            <span class="cell-room">{{ row.room_no }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="rent_price" label="月租金（元）" width="170" align="center">
          <template #default="{ row }">
            <template v-if="editingRoom === row.room_no">
              <el-input v-model="editForm.rent_price" size="small" style="width: 100px" />
            </template>
            <template v-else><span class="cell-value">{{ row.rent_price }}</span></template>
          </template>
        </el-table-column>
        <el-table-column prop="elec_price" label="电费单价（元/度）" width="190" align="center">
          <template #default="{ row }">
            <template v-if="editingRoom === row.room_no">
              <el-input v-model="editForm.elec_price" size="small" style="width: 100px" />
            </template>
            <template v-else><span class="cell-value">{{ row.elec_price }}</span></template>
          </template>
        </el-table-column>
        <el-table-column prop="water_price" label="水费单价（元/吨）" width="190" align="center">
          <template #default="{ row }">
            <template v-if="editingRoom === row.room_no">
              <el-input v-model="editForm.water_price" size="small" style="width: 100px" />
            </template>
            <template v-else><span class="cell-value">{{ row.water_price }}</span></template>
          </template>
        </el-table-column>
        <el-table-column prop="is_rented" label="出租状态" width="140" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_rented"
              :active-value="1"
              :inactive-value="0"
              active-text="已租"
              inactive-text="未租"
              @change="handleRentedChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" align="center">
          <template #default="{ row }">
            <template v-if="editingRoom === row.room_no">
              <el-button type="primary" size="small" @click="saveRoom(row)">保存</el-button>
              <el-button size="small" @click="cancelEdit">取消</el-button>
            </template>
            <template v-else>
              <el-button type="primary" size="small" @click="startEdit(row)">编辑</el-button>
            </template>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getRooms, updateRoom } from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'

const rooms = ref([])
const loading = ref(false)
const filterRented = ref(undefined)
const filterFloor = ref(2)
const editingRoom = ref(null)
const editForm = reactive({ rent_price: '', elec_price: '', water_price: '' })

// 全局统计数据（不受筛选影响）
const allRooms = ref([])
const statsTotal = ref(0)
const statsRented = ref(0)
const statsUnrented = ref(0)
const statsRate = ref(0)

const floorOptions = [
  { value: undefined, label: '全部', desc: '2-5楼' },
  { value: 2, label: '2楼', desc: '201-209' },
  { value: 3, label: '3楼', desc: '301-309' },
  { value: 4, label: '4楼', desc: '401-409' },
  { value: 5, label: '5楼', desc: '501-509' }
]

function selectFloor(val) {
  filterFloor.value = val
  fetchRooms()
}

function updateStats(data) {
  allRooms.value = data
  statsTotal.value = data.length
  const rented = data.filter(r => r.is_rented === 1)
  statsRented.value = rented.length
  statsUnrented.value = data.length - rented.length
  statsRate.value = data.length > 0 ? Math.round((rented.length / data.length) * 100) : 0
}

async function fetchStats() {
  try {
    const res = await getRooms(undefined, undefined)
    updateStats(res.data || [])
  } catch { /* ignore */ }
}

async function fetchRooms() {
  loading.value = true
  try {
    const res = await getRooms(filterRented.value, filterFloor.value)
    rooms.value = res.data || []
  } finally {
    loading.value = false
  }
}

function startEdit(row) {
  editingRoom.value = row.room_no
  editForm.rent_price = row.rent_price
  editForm.elec_price = row.elec_price
  editForm.water_price = row.water_price
}

function cancelEdit() {
  editingRoom.value = null
}

async function saveRoom(row) {
  const rent = parseFloat(editForm.rent_price)
  const elec = parseFloat(editForm.elec_price)
  const water = parseFloat(editForm.water_price)
  if (isNaN(rent) || isNaN(elec) || isNaN(water)) {
    ElMessage.warning('请输入有效的价格')
    return
  }
  if (rent <= 0 || elec <= 0 || water <= 0) {
    ElMessage.warning('价格必须大于0')
    return
  }
  await updateRoom(row.room_no, {
    rent_price: rent,
    elec_price: elec,
    water_price: water
  })
  ElMessage.success('保存成功')
  editingRoom.value = null
  await fetchRooms()
}

async function handleRentedChange(row) {
  try {
    await updateRoom(row.room_no, { is_rented: row.is_rented })
    ElMessage.success('状态已更新')
    // 同步刷新统计数据
    await fetchStats()
  } catch {
    row.is_rented = row.is_rented === 1 ? 0 : 1
  }
}

onMounted(async () => {
  await fetchStats()
  await fetchRooms()
})
</script>

<style scoped>
.room-manage {
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
  gap: 12px;
  margin-bottom: 10px;
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
}

/* 状态筛选 */
.status-filter :deep(.el-radio-button__inner) {
  font-size: 14px;
  padding: 8px 20px;
}

/* ── 出租统计栏 ── */
.stats-bar {
  display: flex;
  align-items: center;
  gap: 0;
  margin-bottom: 10px;
  flex-shrink: 0;
  background: linear-gradient(135deg, #f8f6f0 0%, #f5f2ed 100%);
  border-radius: 12px;
  padding: 12px 18px;
  border: 1px solid #e4e9f2;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
  padding: 0 4px;
}

.stat-icon {
  flex-shrink: 0;
  width: 20px;
  height: 20px;
}

.stat-body {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.stat-label {
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.3px;
  line-height: 1.2;
}

.stat-value {
  font-size: 20px;
  font-weight: 800;
  line-height: 1.2;
}

.stat-unit {
  font-size: 12px;
  font-weight: 500;
  margin-left: 2px;
  opacity: 0.7;
}

/* 各项目不同的颜色 */
.stat-total .stat-icon { color: #606266; }
.stat-total .stat-label { color: #909399; }
.stat-total .stat-value { color: #303133; }

.stat-rented .stat-icon { color: #67c23a; }
.stat-rented .stat-label { color: #67c23a; }
.stat-rented .stat-value { color: #67c23a; }

.stat-unrented .stat-icon { color: #e6a23c; }
.stat-unrented .stat-label { color: #e6a23c; }
.stat-unrented .stat-value { color: #e6a23c; }

.stat-rate .stat-icon { color: #409eff; }
.stat-rate .stat-label { color: #409eff; }
.stat-rate .stat-value { color: #409eff; }

.stat-divider {
  width: 1px;
  height: 40px;
  background: #e0e4ea;
  margin: 0 8px;
  flex-shrink: 0;
}

/* ── 楼层卡片筛选 ── */
.floor-filter {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
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
  height: 52px;
}

.table-wrapper :deep(.el-table__row td) {
  padding: 6px 0;
}

/* 单元格字体加粗加大 */
.cell-room {
  font-size: 18px;
  font-weight: 700;
  color: #303133;
  letter-spacing: 1px;
}

.cell-value {
  font-size: 17px;
  font-weight: 600;
  color: #409eff;
}

/* 编辑模式下输入框 */
.table-wrapper :deep(.el-input__inner) {
  font-size: 15px;
  font-weight: 600;
}

/* Switch 字体 */
.table-wrapper :deep(.el-switch__label) {
  font-size: 13px;
  font-weight: 500;
}

/* 操作按钮 */
.table-wrapper :deep(.el-button--small) {
  font-size: 13px;
  padding: 6px 14px;
}
</style>
