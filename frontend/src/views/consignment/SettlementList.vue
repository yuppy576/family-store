<template>
  <div>
    <!-- 寄卖品选择 -->
    <el-card shadow="never" class="mb-4">
      <el-form :inline="true" size="default">
        <el-form-item label="寄卖品">
          <el-select v-model="consignmentId" filterable placeholder="请选择寄卖品" style="width:300px" @change="loadData">
            <el-option v-for="c in consignmentOptions" :key="c.id" :label="`#${c.id} ${c.name || c.itemName}`" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="consignmentId">
          <el-button type="primary" @click="showCreateDialog('SOLD_SETTLEMENT')">卖出结算</el-button>
          <el-button @click="showCreateDialog('RENEWAL')">续费</el-button>
          <el-button @click="showCreateDialog('RETURN_SETTLEMENT')">到期取回</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 结算列表 -->
    <el-card shadow="never" v-if="consignmentId">
      <template #header>
        <span>结算记录 - 寄卖品 #{{ consignmentId }}</span>
      </template>
      <el-table :data="tableData" v-loading="loading" stripe border style="width:100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="settlementTagType(row.type)" size="small">{{ settlementLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="成交价" width="120" align="right">
          <template #default="{ row }">{{ row.sale_price ? '¥' + Number(row.sale_price).toLocaleString() : '-' }}</template>
        </el-table-column>
        <el-table-column label="佣金" width="120" align="right">
          <template #default="{ row }">{{ row.commission_amount ? '¥' + Number(row.commission_amount).toLocaleString() : '-' }}</template>
        </el-table-column>
        <el-table-column label="结算金额" width="140" align="right">
          <template #default="{ row }">{{ row.settlement_amount ? '¥' + Number(row.settlement_amount).toLocaleString() : '-' }}</template>
        </el-table-column>
        <el-table-column label="续费金额" width="120" align="right">
          <template #default="{ row }">{{ row.renewal_fee ? '¥' + Number(row.renewal_fee).toLocaleString() : '-' }}</template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="180" show-overflow-tooltip />
        <el-table-column label="时间" width="170">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!loading && tableData.length === 0" description="暂无结算记录" />
    </el-card>

    <!-- 创建结算对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="结算类型">
          <el-tag>{{ settlementLabel(form.type) }}</el-tag>
        </el-form-item>
        <el-form-item v-if="form.type === 'SOLD_SETTLEMENT'" label="成交价" prop="sale_price">
          <el-input-number v-model="form.sale_price" :min="0" :precision="2" style="width:100%" />
        </el-form-item>
        <el-form-item v-if="form.type === 'SOLD_SETTLEMENT'" label="佣金" prop="commission_amount">
          <el-input-number v-model="form.commission_amount" :min="0" :precision="2" style="width:100%" />
        </el-form-item>
        <el-form-item v-if="form.type === 'SOLD_SETTLEMENT'" label="实付金额">
          <el-input-number v-model="form.settlement_amount" :min="0" :precision="2" style="width:100%" />
        </el-form-item>
        <el-form-item v-if="form.type === 'RENEWAL'" label="续费金额" prop="renewal_fee">
          <el-input-number v-model="form.renewal_fee" :min="0" :precision="2" style="width:100%" />
        </el-form-item>
        <el-form-item v-if="form.type === 'RENEWAL'" label="续费月数">
          <el-input-number v-model="form.renewal_months" :min="1" :max="12" style="width:100%" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible=false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { loadAllConsignments, createSettlement, listSettlements } from '@/api/consignment'

const loading = ref(false)
const submitting = ref(false)
const consignmentId = ref<number>(0)
const tableData = ref<any[]>([])
const dialogVisible = ref(false)
const formRef = ref()
const consignmentOptions = ref<any[]>([])

const form = ref<any>({ type: 'SOLD_SETTLEMENT', sale_price: 0, commission_amount: 0, settlement_amount: 0, renewal_fee: 0, renewal_months: 1, remark: '' })
const rules: any = { sale_price: [{ required: true, message: '请输入成交价' }] }

const dialogTitle = ref('')

function settlementTagType(type: string) {
  const m: Record<string,string> = { SOLD_SETTLEMENT: 'success', RETURN_SETTLEMENT: 'warning', RENEWAL: 'info' }
  return m[type] || ''
}
function settlementLabel(type: string) {
  const m: Record<string,string> = { SOLD_SETTLEMENT: '卖出结算', RETURN_SETTLEMENT: '到期取回', RENEWAL: '续费' }
  return m[type] || type
}
function formatTime(t: string) { return t ? t.replace('T', ' ').substring(0, 19) : '-' }

async function loadConsignmentOptions() {
  try {
    const res = await loadAllConsignments({ skip: 0, limit: 999 })
    const body = res.data
    const d = body.success ? body.data : body
    consignmentOptions.value = d.consignments || d.list || []
  } catch { consignmentOptions.value = [] }
}

async function loadData() {
  if (!consignmentId.value) return
  loading.value = true
  try {
    const res = await listSettlements(consignmentId.value)
    const body = res.data
    const d = body.success ? body.data : body
    tableData.value = Array.isArray(d) ? d : d.list || d.settlements || []
  } finally { loading.value = false }
}

function showCreateDialog(type: string) {
  form.value = { type, sale_price: 0, commission_amount: 0, settlement_amount: 0, renewal_fee: 0, renewal_months: 1, remark: '' }
  dialogTitle.value = settlementLabel(type)
  dialogVisible.value = true
}

async function submitForm() {
  submitting.value = true
  try {
    await createSettlement(consignmentId.value, form.value)
    ElMessage.success('结算成功')
    dialogVisible.value = false
    loadData()
  } catch {} finally { submitting.value = false }
}

const route = useRoute()
onMounted(async () => {
  await loadConsignmentOptions()
  if (route.query.id) {
    consignmentId.value = Number(route.query.id)
    loadData()
  }
})
</script>
