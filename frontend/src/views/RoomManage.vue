<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: space-between; align-items: center">
      <h2 style="margin: 0">房间管理</h2>
      <el-radio-group v-model="filterRented" @change="fetchRooms">
        <el-radio-button :value="undefined">全部</el-radio-button>
        <el-radio-button :value="1">已出租</el-radio-button>
        <el-radio-button :value="0">未出租</el-radio-button>
      </el-radio-group>
    </div>

    <el-table :data="rooms" border stripe v-loading="loading" empty-text="暂无房间数据">
      <el-table-column prop="room_no" label="房间号" width="100" align="center" />
      <el-table-column prop="rent_price" label="月租金（元）" width="140" align="center">
        <template #default="{ row }">
          <template v-if="editingRoom === row.room_no">
            <el-input v-model="editForm.rent_price" size="small" style="width: 100px" />
          </template>
          <template v-else>{{ row.rent_price }}</template>
        </template>
      </el-table-column>
      <el-table-column prop="elec_price" label="电费单价（元/度）" width="160" align="center">
        <template #default="{ row }">
          <template v-if="editingRoom === row.room_no">
            <el-input v-model="editForm.elec_price" size="small" style="width: 100px" />
          </template>
          <template v-else>{{ row.elec_price }}</template>
        </template>
      </el-table-column>
      <el-table-column prop="water_price" label="水费单价（元/吨）" width="160" align="center">
        <template #default="{ row }">
          <template v-if="editingRoom === row.room_no">
            <el-input v-model="editForm.water_price" size="small" style="width: 100px" />
          </template>
          <template v-else>{{ row.water_price }}</template>
        </template>
      </el-table-column>
      <el-table-column prop="is_rented" label="出租状态" width="120" align="center">
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
      <el-table-column label="操作" width="140" align="center">
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
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getRooms, updateRoom } from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'

const rooms = ref([])
const loading = ref(false)
const filterRented = ref(undefined)
const editingRoom = ref(null)
const editForm = reactive({ rent_price: '', elec_price: '', water_price: '' })

async function fetchRooms() {
  loading.value = true
  try {
    const res = await getRooms(filterRented.value)
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
