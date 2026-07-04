<template>
  <div>
    <!-- 选择车辆 -->
    <el-card shadow="never" class="mb-4">
      <el-form :inline="true" size="default">
        <el-form-item label="寄卖品">
          <el-select v-model="consignmentId" filterable placeholder="选择车辆寄卖品" style="width:350px" @change="onConsignmentChange">
            <el-option v-for="c in vehicleOptions" :key="c.id" :label="`#${c.id} ${c.name} - ${c.vehicle?.plate_number || c.vehicle?.plateNumber || '未上牌'}`" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="vehicleId">
          <el-button type="primary" @click="showAddDialog">新增进度</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 进度状态机 -->
    <el-card shadow="never" v-if="vehicleId">
      <template #header>
        <span>过户进度 - 寄卖品 #{{ consignmentId }}</span>
      </template>

      <!-- 状态步骤条 -->
      <el-steps :active="currentStepIndex" align-center style="margin-bottom:30px">
        <el-step title="待验车" description="等待验车" />
        <el-step title="验车完成" description="车辆已检验" />
        <el-step title="过户办理中" description="正在办理过户" />
        <el-step title="过户完成" description="过户手续已完成" />
        <el-step title="已结算" description="款项已结清" />
      </el-steps>

      <!-- 进度明细 -->
      <el-timeline>
        <el-timeline-item
          v-for="(item, idx) in progressList"
          :key="idx"
          :timestamp="formatTime(item.created_at)"
          :color="item.status === currentStatus ? '#409eff' : '#e0e0e0'"
        >
          <div class="progress-item">
            <el-tag :type="statusTagType(item.status)" size="small">{{ statusLabel(item.status) }}</el-tag>
            <span class="progress-remark" v-if="item.remark">{{ item.remark }}</span>
            <span class="progress-operator" v-if="item.operator">- {{ item.operator }}</span>
          </div>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-if="progressList.length === 0" description="暂无过户进度" />
    </el-card>

    <!-- 新增进度对话框 -->
    <el-dialog v-model="dialogVisible" title="新增过户进度" width="450px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="当前状态" prop="status">
          <el-select v-model="form.status" style="width:100%">
            <el-option label="待验车" value="PENDING_INSPECTION" />
            <el-option label="验车完成" value="INSPECTED" />
            <el-option label="过户办理中" value="TRANSFERRING" />
            <el-option label="过户完成" value="TRANSFERRED" />
            <el-option label="已结算" value="SETTLED" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="备注说明" />
        </el-form-item>
        <el-form-item label="经办人">
          <el-input v-model="form.operator" placeholder="经办人姓名" />
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
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { loadAllConsignments, getVehicle, listTransferProgress, createTransferProgress } from '@/api/consignment'

const route = useRoute()
const loading = ref(false)
const submitting = ref(false)
const consignmentId = ref(0)
const vehicleId = ref(0)
const progressList = ref<any[]>([])
const vehicleOptions = ref<any[]>([])
const dialogVisible = ref(false)
const formRef = ref()

const form = ref<any>({ status: 'PENDING_INSPECTION', remark: '', operator: '' })
const rules = { status: [{ required: true, message: '请选择状态' }] }

const stepStatuses = ['PENDING_INSPECTION', 'INSPECTED', 'TRANSFERRING', 'TRANSFERRED', 'SETTLED']
const currentStatus = computed(() => progressList.value[0]?.status || '')
const currentStepIndex = computed(() => {
  const idx = stepStatuses.indexOf(currentStatus.value)
  return idx >= 0 ? idx : 0
})

function statusTagType(s: string) {
  const m: Record<string,string> = { PENDING_INSPECTION:'info', INSPECTED:'primary', TRANSFERRING:'warning', TRANSFERRED:'success', SETTLED:'success' }
  return m[s] || 'info'
}
function statusLabel(s: string) {
  const m: Record<string,string> = { PENDING_INSPECTION:'待验车', INSPECTED:'验车完成', TRANSFERRING:'过户办理中', TRANSFERRED:'过户完成', SETTLED:'已结算' }
  return m[s] || s
}
function formatTime(t: string) { return t ? t.replace('T', ' ').substring(0, 19) : '-' }

async function loadVehicleOptions() {
  try {
    const res = await loadAllConsignments({ skip: 0, limit: 999 })
    const body = res.data
    const d = body.success ? body.data : body
    const items = d.consignments || d.list || []
    vehicleOptions.value = items.filter((i: any) => i.is_vehicle || i.isVehicle)
    // 顺便获取车辆详情
    for (const item of vehicleOptions.value) {
      try {
        const vr = await getVehicle(item.id)
        const vb = vr.data
        item.vehicle = vb.success ? vb.data : vb
      } catch {}
    }
  } catch { vehicleOptions.value = [] }
}

async function loadProgress() {
  if (!vehicleId.value) return
  loading.value = true
  try {
    const res = await listTransferProgress(vehicleId.value)
    const body = res.data
    const d = body.success ? body.data : body
    progressList.value = Array.isArray(d) ? d : d.list || d.progress || []
  } finally { loading.value = false }
}

async function onConsignmentChange() {
  if (!consignmentId.value) { vehicleId.value = 0; progressList.value = []; return }
  try {
    const vr = await getVehicle(consignmentId.value)
    const vb = vr.data
    vehicleId.value = (vb.success ? vb.data : vb)?.id || 0
    loadProgress()
  } catch { vehicleId.value = 0 }
}

function showAddDialog() {
  form.value = { status: 'PENDING_INSPECTION', remark: '', operator: '' }
  dialogVisible.value = true
}

async function submitForm() {
  submitting.value = true
  try {
    await createTransferProgress(vehicleId.value, form.value)
    ElMessage.success('进度已更新')
    dialogVisible.value = false
    loadProgress()
  } catch {} finally { submitting.value = false }
}

onMounted(async () => {
  await loadVehicleOptions()
  if (route.query.id) {
    consignmentId.value = Number(route.query.id)
    onConsignmentChange()
  }
})
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
.progress-item { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.progress-remark { color: #606266; font-size: 13px; }
.progress-operator { color: #909399; font-size: 12px; }
</style>
