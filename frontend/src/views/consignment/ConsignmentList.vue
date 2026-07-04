<template>
  <div class="consignment-list">
    <!-- Filter Bar -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filterForm" size="default">
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部状态" clearable style="width: 140px">
            <el-option label="在售" value="ON_SALE" />
            <el-option label="已售出" value="SOLD" />
            <el-option label="已过期" value="EXPIRED" />
            <el-option label="已取回" value="RETURNED" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键字">
          <el-input
            v-model="filterForm.keyword"
            placeholder="物品名称"
            clearable
            style="width: 200px"
            @keyup.enter="handleFilter"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleFilter">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Toolbar -->
    <div class="toolbar">
      <el-button type="primary" @click="openAddDialog">
        <el-icon><Plus /></el-icon>
        新增寄卖品
      </el-button>
    </div>

    <!-- Table -->
    <el-card shadow="never">
      <el-table
        :data="tableData"
        v-loading="loading"
        stripe
        border
        style="width: 100%"
        @sort-change="handleSortChange"
      >
        <el-table-column type="expand" width="50">
          <template #default="{ row }">
            <div v-if="row.isVehicle" class="vehicle-detail">
              <el-descriptions title="车辆详情" :column="3" border size="small">
                <el-descriptions-item label="品牌">{{ row.vehicle?.brand || '-' }}</el-descriptions-item>
                <el-descriptions-item label="型号">{{ row.vehicle?.model || '-' }}</el-descriptions-item>
                <el-descriptions-item label="车牌号">{{ row.vehicle?.plate_number || row.vehicle?.plateNumber || '-' }}</el-descriptions-item>
                <el-descriptions-item label="年份">{{ row.vehicle?.year || '-' }}</el-descriptions-item>
                <el-descriptions-item label="里程(km)">{{ row.vehicle?.mileage ? `${row.vehicle.mileage}km` : '-' }}</el-descriptions-item>
                <el-descriptions-item label="颜色">{{ row.vehicle?.color || '-' }}</el-descriptions-item>
                <el-descriptions-item label="VIN码">{{ row.vehicle?.vin || '-' }}</el-descriptions-item>
                <el-descriptions-item label="年检到期">{{ formatTime(row.vehicle?.inspection_expire || row.vehicle?.inspectionExpire) }}</el-descriptions-item>
                <el-descriptions-item label="保险到期">{{ formatTime(row.vehicle?.insurance_expire || row.vehicle?.insuranceExpire) }}</el-descriptions-item>
              </el-descriptions>
            </div>
            <div v-else class="no-vehicle-detail">
              <el-empty description="非车辆物品，无车辆详情" />
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="id" label="ID" width="70" sortable="custom" />
        <el-table-column label="物品名称" min-width="150">
          <template #default="{ row }">{{ row.name || row.itemName || '-' }}</template>
        </el-table-column>
        <el-table-column label="车辆标记" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.isVehicle" type="warning" size="small">车辆</el-tag>
            <el-tag v-else type="info" size="small">普通</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="寄卖人" min-width="120">
          <template #default="{ row }">
            {{ row.consignorName || row.consignor?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag
              :type="statusTagType(row.status)"
              size="small"
            >
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="期望价格" width="110" align="right">
          <template #default="{ row }">
            {{ renderPrice(row, 'expected_price') }}
          </template>
        </el-table-column>
        <el-table-column label="成交价格" width="110" align="right">
          <template #default="{ row }">
            {{ renderPrice(row, 'final_price') }}
          </template>
        </el-table-column>
        <el-table-column label="佣金率" width="80" align="right">
          <template #default="{ row }">{{ row.commission_rate || row.commission || '-' }}{{ row.commission_rate ? '%' : '' }}</template>
        </el-table-column>
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">{{ formatTime(row.created_at || row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="openEditDialog(row)">
              编辑
            </el-button>
            <el-button size="small" type="success" link @click="goSettlement(row)">
              结算
            </el-button>
            <el-button v-if="row.is_vehicle || row.isVehicle" size="small" type="warning" link @click="goTransfer(row)">
              过户
            </el-button>
            <el-popconfirm
              title="确定要删除此寄卖品吗？"
              @confirm="handleDelete(row.id)"
            >
              <template #reference>
                <el-button size="small" type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>

    <!-- Add / Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEditing ? '编辑寄卖品' : '新增寄卖品'"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
        size="default"
      >
        <el-form-item label="物品名称" prop="itemName">
          <el-input v-model="form.itemName" placeholder="请输入物品名称" />
        </el-form-item>
        <el-form-item label="物品描述">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="2"
            placeholder="物品描述（选填）"
          />
        </el-form-item>
        <el-form-item label="寄卖人" prop="consignorId">
          <el-select
            v-model="form.consignorId"
            placeholder="请选择寄卖人"
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="c in consignorOptions"
              :key="c.id"
              :label="`${c.name} (${c.phone})`"
              :value="c.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="是否车辆">
          <el-switch v-model="form.isVehicle" />
        </el-form-item>
        <template v-if="form.isVehicle">
          <el-divider content-position="left">车辆信息</el-divider>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="品牌" prop="vehicleBrand">
                <el-input v-model="form.vehicleBrand" placeholder="丰田/本田/宝马..." />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="型号" prop="vehicleModel">
                <el-input v-model="form.vehicleModel" placeholder="凯美瑞/雅阁/X5..." />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="车牌号">
                <el-input v-model="form.vehiclePlate" placeholder="粤B12345" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="车架号(VIN)">
                <el-input v-model="form.vehicleVin" placeholder="17位车架号" maxlength="17" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="16">
            <el-col :span="8">
              <el-form-item label="年份">
                <el-input-number v-model="form.vehicleYear" :min="2000" :max="2030" style="width:100%" />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="里程(km)">
                <el-input-number v-model="form.vehicleMileage" :min="0" :step="1000" style="width:100%" />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="颜色">
                <el-input v-model="form.vehicleColor" placeholder="黑色/白色..." />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="年检到期">
                <el-date-picker v-model="form.vehicleInspection" type="date" placeholder="选择日期" style="width:100%" value-format="YYYY-MM-DD" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="保险到期">
                <el-date-picker v-model="form.vehicleInsurance" type="date" placeholder="选择日期" style="width:100%" value-format="YYYY-MM-DD" />
              </el-form-item>
            </el-col>
          </el-row>
        </template>
        <el-form-item label="期望价格">
          <el-input-number
            v-model="form.expectedPrice"
            :min="0"
            :precision="2"
            style="width: 100%"
            placeholder="期望售价"
          />
        </el-form-item>
        <el-form-item label="佣金">
          <el-input-number
            v-model="form.commission"
            :min="0"
            :precision="2"
            style="width: 100%"
            placeholder="佣金金额"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="2"
            placeholder="备注（选填）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import {
  loadAllConsignments,
  createConsignment,
  updateConsignment,
  deleteConsignment,
  loadAllConsignors,
  createVehicle,
  updateVehicle,
  getVehicle,
} from '@/api/consignment'
import type { ConsignmentData } from '@/api/consignment'

interface ConsignorOption {
  id: number
  name: string
  phone: string
}

interface VehicleInfo {
  brand?: string
  model?: string
  plateNumber?: string
  year?: number
  mileage?: number
  color?: string
  vin?: string
  engine?: string
  registrationDate?: string
}

interface ConsignmentRecord {
  id: number
  itemName: string
  description?: string
  consignorId: number
  consignorName?: string
  consignor?: ConsignorOption
  expectedPrice?: number
  sellingPrice?: number
  commission?: number
  status: string
  isVehicle: boolean
  images?: string[]
  remark?: string
  createdAt: string
  vehicle?: VehicleInfo
}

const loading = ref(false)
const tableData = ref<ConsignmentRecord[]>([])
const dialogVisible = ref(false)
const isEditing = ref(false)
const saving = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()
const consignorOptions = ref<ConsignorOption[]>([])

const router = useRouter()
function goSettlement(row: any) {
  router.push({ path: '/consignment/settlements', query: { id: row.id } })
}
function goTransfer(row: any) {
  router.push({ path: '/consignment/transfer', query: { id: row.id } })
}

const filterForm = reactive({
  status: '',
  keyword: '',
})

const defaultForm = (): any => ({
  itemName: '', description: '', consignorId: 0,
  expectedPrice: undefined, commission: undefined,
  status: 'ON_SALE', isVehicle: false, remark: '',
  vehicleBrand: '', vehicleModel: '', vehiclePlate: '', vehicleVin: '',
  vehicleYear: 0, vehicleMileage: 0, vehicleColor: '',
  vehicleInspection: '', vehicleInsurance: '',
})

const form = reactive<ConsignmentData>(defaultForm())

const formRules: FormRules = {
  itemName: [{ required: true, message: '请输入物品名称', trigger: 'blur' }],
  consignorId: [{ required: true, message: '请选择寄卖人', trigger: 'change' }],
}

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0,
})

const sort = reactive({
  field: '',
  order: '',
})

onMounted(async () => {
  await loadConsignorOptions()
  loadData()
})

async function loadConsignorOptions() {
  try {
    const res = await loadAllConsignors({ skip: 0, limit: 999 })
    const body = res.data
    const d = body.success ? body.data : body
    const list = d.consignors || d.list || d.rows || d.data || d
    consignorOptions.value = Array.isArray(list) ? list : (Array.isArray(d) ? d : [])
  } catch {
    consignorOptions.value = []
  }
}

async function loadData() {
  loading.value = true
  try {
    const params: any = {
      skip: (pagination.page - 1) * pagination.pageSize,
      limit: pagination.pageSize,
    }
    if (filterForm.status) {
      params.status = filterForm.status
    }
    if (filterForm.keyword) {
      params.keyword = filterForm.keyword
    }
    if (sort.field) {
      params.sortField = sort.field
      params.sortOrder = sort.order === 'ascending' ? 'asc' : 'desc'
    }
    const res = await loadAllConsignments(params)
    const body = res.data
    const d = body.success ? body.data : body
    if (d.consignments) {
      tableData.value = d.consignments
      pagination.total = d.meta?.total || d.consignments.length
    } else if (Array.isArray(d)) {
      tableData.value = d
      pagination.total = d.length
    } else if (d.list) {
      tableData.value = d.list
      pagination.total = d.total || d.list.length
    } else if (d.rows) {
      tableData.value = d.rows
      pagination.total = d.total || d.rows.length
    } else {
      tableData.value = d.data || []
      pagination.total = d.total || tableData.value.length
    }
  } catch {
    // handled by interceptor
  } finally {
    loading.value = false
  }
}

function handleFilter() {
  pagination.page = 1
  loadData()
}

function renderPrice(row: any, field: string) {
  const val = row[field] ?? row[field.replace(/_([a-z])/g, (_, l) => l.toUpperCase())] ?? null
  return val ? `¥${Number(val).toLocaleString()}` : '-'
}
function formatTime(t: string) {
  return t ? t.replace('T', ' ').substring(0, 19) : '-'
}
function resetFilter() {
  filterForm.status = ''
  filterForm.keyword = ''
  pagination.page = 1
  loadData()
}

function handleSortChange(sortInfo: any) {
  sort.field = sortInfo.prop || ''
  sort.order = sortInfo.order || ''
  loadData()
}

function resetForm() {
  const d = defaultForm()
  Object.assign(form, d)
}

function openAddDialog() {
  isEditing.value = false
  editingId.value = null
  resetForm()
  dialogVisible.value = true
}

function openEditDialog(row: ConsignmentRecord) {
  isEditing.value = true
  editingId.value = row.id
  form.itemName = row.name || row.itemName || ''
  form.description = row.description || ''
  form.consignorId = row.consignor_id ?? row.consignorId ?? 0
  form.expectedPrice = row.expected_price ?? row.expectedPrice
  form.commission = row.commission_rate ?? row.commission
  form.status = row.status || 'ON_SALE'
  form.isVehicle = !!(row.is_vehicle ?? row.isVehicle)
  form.remark = row.memo || row.remark || ''
  form.vehicleBrand = ''; form.vehicleModel = ''; form.vehiclePlate = ''; form.vehicleVin = ''
  form.vehicleYear = 0; form.vehicleMileage = 0; form.vehicleColor = ''
  form.vehicleInspection = ''; form.vehicleInsurance = ''
  // 如果是车辆，异步加载车辆信息
  if (form.isVehicle && row.id) {
    getVehicle(row.id).then((res: any) => {
      const v = res.data?.data || res.data
      if (v && v.id) {
        form.vehicleBrand = v.brand || ''
        form.vehicleModel = v.model || ''
        form.vehiclePlate = v.plate_number || ''
        form.vehicleVin = v.vin || ''
        form.vehicleYear = v.year || 0
        form.vehicleMileage = v.mileage || 0
        form.vehicleColor = v.color || ''
        form.vehicleInspection = v.inspection_expire || ''
        form.vehicleInsurance = v.insurance_expire || ''
      }
    }).catch(() => {})
  }
  dialogVisible.value = true
}

async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    const payload: any = {
      name: form.itemName,
      description: form.description,
      consignor_id: form.consignorId,
      expected_price: form.expectedPrice,
      commission_rate: form.commission,
      memo: form.remark,
      is_vehicle: form.isVehicle,
    }
    let itemId = editingId.value
    if (isEditing.value && itemId) {
      await updateConsignment(itemId, payload)
      ElMessage.success('更新成功')
    } else {
      const res = await createConsignment(payload)
      const body = res.data
      itemId = body.data?.id || (body.success && body.data?.id)
      ElMessage.success('创建成功')
    }
    // 如果是车辆，同步创建/更新车辆信息
    if (form.isVehicle && itemId) {
      const vPayload: any = {
        brand: form.vehicleBrand, model: form.vehicleModel,
        plate_number: form.vehiclePlate, vin: form.vehicleVin,
        year: form.vehicleYear || undefined,
        mileage: form.vehicleMileage || undefined,
        color: form.vehicleColor,
        inspection_expire: form.vehicleInspection || undefined,
        insurance_expire: form.vehicleInsurance || undefined,
      }
      try {
        const existing = await getVehicle(itemId)
        if (existing.data?.data?.id) {
          await updateVehicle(itemId, vPayload)
        } else {
          await createVehicle(itemId, vPayload)
        }
      } catch {
        await createVehicle(itemId, vPayload)
      }
    }
    dialogVisible.value = false
    loadData()
  } catch {
    // handled by interceptor
  } finally {
    saving.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await deleteConsignment(id)
    ElMessage.success('删除成功')
    loadData()
  } catch {
    // handled by interceptor
  }
}

function statusTagType(status: string): string {
  const map: Record<string, string> = {
    ON_SALE: 'success',
    SOLD: 'danger',
    EXPIRED: 'warning',
    RETURNED: 'info',
    CANCELLED: 'info',
  }
  return map[status] || 'info'
}

function statusLabel(status: string): string {
  const map: Record<string, string> = {
    ON_SALE: '在售',
    SOLD: '已售出',
    EXPIRED: '已过期',
    RETURNED: '已取回',
    CANCELLED: '已取消',
  }
  return map[status] || status
}
</script>

<style scoped>
.consignment-list {
  max-width: 1400px;
}

.filter-card {
  margin-bottom: 16px;
}

.toolbar {
  margin-bottom: 12px;
}

.pagination-wrapper {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.vehicle-detail {
  padding: 16px 32px;
}

.no-vehicle-detail {
  padding: 16px;
}
</style>
