<template>
  <div>
    <el-card shadow="never" class="mb-4">
      <el-form :inline="true" :model="searchForm" size="default">
        <el-form-item label="关键字">
          <el-input v-model="searchForm.q" placeholder="姓名/手机号" clearable @keyup.enter="loadData" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
          <el-button @click="searchForm.q='';loadData()">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <div class="mb-3">
      <el-button type="primary" @click="dialogVisible=true;isEdit=false;form={}">
        <el-icon><Plus /></el-icon> 新增寄卖人
      </el-button>
      <el-button type="success" @click="handleExport">
        <el-icon><Download /></el-icon> 导出数据
      </el-button>
    </div>

    <el-card shadow="never">
      <el-table :data="tableData" v-loading="loading" stripe border style="width:100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="姓名" min-width="120" />
        <el-table-column prop="phone" label="手机号" width="140" />
        <el-table-column prop="id_card" label="身份证号" width="180" />
        <el-table-column prop="address" label="地址" min-width="200" show-overflow-tooltip />
        <el-table-column prop="memo" label="备注" min-width="160" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" link @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-if="total > 0"
        v-model:current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadData"
        class="mt-3"
      />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit?'编辑寄卖人':'新增寄卖人'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="身份证号">
          <el-input v-model="form.id_card" />
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="form.address" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.memo" type="textarea" :rows="2" />
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
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Download } from '@element-plus/icons-vue'
import { loadAllConsignors, createConsignor, updateConsignor, deleteConsignor } from '@/api/consignment'
import { exportConsignors } from '@/utils/export'

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()
const searchForm = reactive({ q: '' })
const form = reactive<any>({ name: '', phone: '', id_card: '', address: '', memo: '' })
const rules = { name: [{ required: true, message: '请输入姓名', trigger: 'blur' }], phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }] }

async function loadData() {
  loading.value = true
  try {
    const res = await loadAllConsignors({ q: searchForm.q || undefined, skip: (page.value - 1) * pageSize.value, limit: pageSize.value })
    const body = res.data
    if (body.success && body.data) {
      const d = body.data
      if (d.consignors) { tableData.value = d.consignors; total.value = d.meta?.total || d.consignors.length }
      else if (Array.isArray(d)) { tableData.value = d; total.value = d.length }
      else if (Array.isArray(d.list)) { tableData.value = d.list; total.value = d.total || d.list.length }
    }
  } finally { loading.value = false }
}

function handleEdit(row: any) {
  isEdit.value = true
  Object.assign(form, { name: row.name, phone: row.phone, id_card: row.id_card || '', address: row.address || '', memo: row.memo || '' })
  ;(form as any).id = row.id
  dialogVisible.value = true
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确认删除这条记录？', '提示')
    await deleteConsignor(id)
    ElMessage.success('删除成功')
    loadData()
  } catch {}
}

async function submitForm() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    const data = { name: form.name, phone: form.phone, id_card: form.id_card || undefined, address: form.address || undefined, memo: form.memo || undefined }
    if (isEdit.value) { await updateConsignor((form as any).id, data); ElMessage.success('修改成功') }
    else { await createConsignor(data); ElMessage.success('新增成功') }
    dialogVisible.value = false
    loadData()
  } finally { submitting.value = false }
}

function handleExport(){exportConsignors(tableData.value)}
onMounted(loadData)
</script>
<style scoped>
.mb-4 { margin-bottom: 16px; }
.mb-3 { margin-bottom: 12px; }
.mt-3 { margin-top: 12px; }
</style>
