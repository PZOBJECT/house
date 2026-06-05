<template>
  <div class="bill-manage">
    <!-- 头部 -->
    <div class="page-header">
      <h2>月度账单</h2>
      <div class="header-actions">
        <el-date-picker
          v-model="selectedMonth"
          type="month"
          placeholder="选择月份"
          format="YYYY年MM月"
          value-format="YYYY-MM"
          :clearable="false"
          style="width: 180px"
          @change="onMonthChange"
        />
        <el-select v-model="filterFloor" placeholder="楼层筛选" style="width: 120px" @change="onFilterChange">
          <el-option label="全部" :value="undefined" />
          <el-option v-for="f in [2,3,4,5]" :key="f" :label="`${f}楼`" :value="f" />
        </el-select>
        <el-button type="primary" @click="handleGenerate">生成当月账单</el-button>
      </div>
    </div>

    <!-- 表格区域 -->
    <div class="table-wrapper">
      <el-table
        ref="tableRef"
        :data="pagedBills"
        border
        stripe
        v-loading="loading"
        empty-text="暂无账单数据，请先生成账单"
        :row-class-name="tableRowClassName"
        show-summary
        :summary-method="getSummaries"
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
const filterFloor = ref(undefined)
const currentPage = ref(1)

const now = new Date()
const selectedMonth = ref(`${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`)

const queryParams = computed(() => {
  const [y, m] = selectedMonth.value.split('-')
  return { year: parseInt(y), month: parseInt(m) }
})

const pagedBills = computed(() => {
  const start = (currentPage.value - 1) * 9
  return bills.value.slice(start, start + 9)
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

function onMonthChange() {
  currentPage.value = 1
  fetchBills()
}

function onFilterChange() {
  currentPage.value = 1
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

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 12px;
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

.table-wrapper {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.table-wrapper :deep(.el-table__body-wrapper) {
  /* ensure remaining space is used */
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
