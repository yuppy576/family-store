<template>
  <div>
    <el-card shadow="never" class="mb-4">
      <el-form :inline="true" :model="searchForm" size="default">
        <el-form-item label="关键字">
          <el-input v-model="searchForm.q" placeholder="名称/手机号" clearable @keyup.enter="loadData" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
          <el-button @click="searchForm.q='';loadData()">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    <div class="mb-3"><el-button type="primary" @click="dialogVisible=true;isEdit=false;form={}"><el-icon><Plus /></el-icon> 新增供应商</el-button></div>
    <el-card shadow="never">
      <el-table :data="tableData" v-loading="loading" stripe border style="width:100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="名称" min-width="150" />
        <el-table-column prop="contact_person" label="联系人" min-width="120" />
        <el-table-column prop="phone" label="电话" width="140" />
        <el-table-column prop="address" label="地址" min-width="200" show-overflow-tooltip />
        <el-table-column prop="memo" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" link @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination v-if="total>0" v-model:current-page="page" :page-size="pageSize" :total="total" layout="total,prev,pager,next" @current-change="loadData" class="mt-3" />
    </el-card>
    <el-dialog v-model="dialogVisible" :title="isEdit?'编辑供应商':'新增供应商'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="名称" prop="name"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="联系人"><el-input v-model="form.contact_person" /></el-form-item>
        <el-form-item label="电话"><el-input v-model="form.phone" /></el-form-item>
        <el-form-item label="地址"><el-input v-model="form.address" type="textarea" :rows="2" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.memo" type="textarea" :rows="2" /></el-form-item>
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
import request from '@/api/request'

const loading=ref(false),submitting=ref(false),dialogVisible=ref(false),isEdit=ref(false),formRef=ref()
const tableData=ref<any[]>([]),total=ref(0),page=ref(1),pageSize=ref(20)
const searchForm=reactive({q:''})
const form=reactive<any>({name:'',contact_person:'',phone:'',address:'',memo:''})
const rules={name:[{required:true,message:'请输入名称',trigger:'blur'}]}

async function loadData(){
  loading.value=true
  try{
    const res=await request.get('/suppliers',{params:{q:searchForm.q||undefined,skip:(page.value-1)*pageSize.value,limit:pageSize.value}})
    const d=res.data;const b=d.success?d.data:d
    tableData.value=b.suppliers||b.list||b;total.value=b.meta?.total||tableData.value.length
  }finally{loading.value=false}
}
function handleEdit(row:any){
  isEdit.value=true;Object.assign(form,{name:row.name,contact_person:row.contact_person||'',phone:row.phone||'',address:row.address||'',memo:row.memo||''});(form as any).id=row.id;dialogVisible.value=true
}
async function handleDelete(id:number){
  try{await ElMessageBox.confirm('确认删除？','提示');await request.delete(`/suppliers/${id}`);ElMessage.success('已删除');loadData()}catch{}
}
async function submitForm(){
  if(!await formRef.value?.validate().catch(()=>false))return
  submitting.value=true
  try{
    const data={name:form.name,contact_person:form.contact_person||undefined,phone:form.phone||undefined,address:form.address||undefined,memo:form.memo||undefined}
    if(isEdit.value){await request.put(`/suppliers/${(form as any).id}`,data);ElMessage.success('已更新')}
    else{await request.post('/suppliers',data);ElMessage.success('已创建')}
    dialogVisible.value=false;loadData()
  }finally{submitting.value=false}
}
onMounted(loadData)
</script>
<style scoped>
.mb-4{margin-bottom:16px}.mb-3{margin-bottom:12px}.mt-3{margin-top:12px}
</style>
