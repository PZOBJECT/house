<template>
  <div class="room-manage">
    <!-- 头部 -->
    <div class="page-header">
      <h2>房间管理</h2>
      <div class="header-actions">
        <el-radio-group v-model="filterRented" @change="fetchRooms" class="status-filter">
          <el-radio-button :value="undefined">全部</el-radio-button>
          <el-radio-button :value="1">已出租</el-radio-button>
          <el-radio-button :value="0">未出租</el-radio-button>
        </el-radio-group>
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
              active-text="已出租"
              inactive-text="未出租"
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
  } catch {
    row.is_rented = row.is_rented === 1 ? 0 : 1
  }
}

onMounted(fetchRooms)
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

/* ── 楼层卡片筛选 ── */
.floor-filter {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 12px;
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
