<template>
  <div>
    <el-card shadow="never" class="mb-4">
      <el-form :inline="true" size="default">
        <el-form-item label="搜索"><el-input v-model="search.q" placeholder="商品名称" clearable @keyup.enter="loadData" style="width:200px" /></el-form-item>
        <el-form-item><el-button type="primary" @click="loadData">查询</el-button><el-button @click="search.q='';loadData()">重置</el-button></el-form-item>
      </el-form>
    </el-card>
    <div class="mb-3">
      <el-button type="primary" @click="showAdd"><el-icon><Plus /></el-icon> 新增商品</el-button>
      <el-tag v-if="lowStockCount>0" type="danger" class="ml-2">库存预警：{{ lowStockCount }} 种商品库存不足</el-tag>
      <span class="ml-2" style="font-size:12px;color:#909399">预警阈值：</span>
      <el-input-number v-model="lowStockThreshold" :min="1" :max="99" size="small" style="width:80px" />
    </div>
    <el-card shadow="never">
      <el-table :data="tableData" v-loading="loading" stripe border style="width:100%" :row-class-name="rowClassName">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="图片" width="70"><template #default="{row}"><el-image v-if="row.image" :src="row.image" style="width:40px;height:40px" fit="cover" /><span v-else>-</span></template></el-table-column>
        <el-table-column prop="name" label="名称" min-width="160" />
        <el-table-column label="分类" width="120"><template #default="{row}">{{ row.category?.name||'-' }}</template></el-table-column>
        <el-table-column label="单位" width="80"><template #default="{row}">{{ row.unit||'-' }}{{ row.base_unit?`(${row.conversion_rate||1}${row.base_unit}/${row.unit})`:'' }}</template></el-table-column>
        <el-table-column prop="price" label="售价" width="100" align="right"><template #default="{row}">¥{{ Number(row.price).toLocaleString() }}</template></el-table-column>
        <el-table-column prop="stock" label="库存" width="80" align="center">
          <template #default="{row}">
            <span :style="{color: row.stock<lowStockThreshold?'#f56c6c':'inherit',fontWeight:row.stock<lowStockThreshold?'700':'inherit'}">{{ row.stock }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{row}"><el-button size="small" type="primary" link @click="showEdit(row)">编辑</el-button><el-button size="small" type="danger" link @click="handleDelete(row.id)">删除</el-button></template>
        </el-table-column>
      </el-table>
      <el-pagination v-if="meta.total>0" v-model:current-page="page" :page-size="pageSize" :total="meta.total" layout="total,prev,pager,next" @current-change="loadData" class="mt-3" />
    </el-card>
    <el-dialog v-model="dialogVisible" :title="isEdit?'编辑商品':'新增商品'" width="550px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="名称" prop="name"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="分类"><el-select v-model="form.category_id" filterable placeholder="选择分类" style="width:100%"><el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" /></el-select></el-form-item>
        <el-row :gutter="16"><el-col :span="12"><el-form-item label="售价" prop="price"><el-input-number v-model="form.price" :min="0" :precision="2" style="width:100%" /></el-form-item></el-col><el-col :span="12"><el-form-item label="库存" prop="stock"><el-input-number v-model="form.stock" :min="0" style="width:100%" /></el-form-item></el-col></el-row>
        <el-row :gutter="16">
          <el-col :span="8"><el-form-item label="单位"><el-select v-model="form.unit" placeholder="单位" style="width:100%" allow-create filterable><el-option label="条" value="条" /><el-option label="箱" value="箱" /><el-option label="瓶" value="瓶" /><el-option label="包" value="包" /><el-option label="斤" value="斤" /><el-option label="盒" value="盒" /></el-select></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="拆零单位"><el-select v-model="form.base_unit" placeholder="拆零" style="width:100%" allow-create filterable><el-option label="包" value="包" /><el-option label="瓶" value="瓶" /><el-option label="盒" value="盒" /><el-option label="支" value="支" /><el-option label="两" value="两" /></el-select></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="1单位="><el-input-number v-model="form.conversion_rate" :min="1" :max="999" style="width:100%" /><span style="font-size:12px;color:#909399;margin-left:4px">{{ form.base_unit||'?' }}</span></el-form-item></el-col>
        </el-row>
        <el-form-item label="图片"><el-input v-model="form.image" placeholder="图片URL" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogVisible=false">取消</el-button><el-button type="primary" @click="submitForm" :loading="submitting">保存</el-button></template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/api/request'
const loading=ref(false),submitting=ref(false),dialogVisible=ref(false),isEdit=ref(false),formRef=ref()
const tableData=ref<any[]>([]),categories=ref<any[]>([])
const lowStockThreshold=ref(10)
const lowStockCount=computed(()=>tableData.value.filter(p=>p.stock<lowStockThreshold.value).length)
function rowClassName({row}:any){return row.stock<lowStockThreshold.value?'low-stock-row':''}
const meta=ref<any>({total:0,limit:10,skip:1}),page=ref(1),pageSize=ref(10)
const search=reactive({q:''})
const form=reactive<any>({name:'',category_id:null,price:0,stock:0,image:'',unit:'',base_unit:'',conversion_rate:1})
const rules:any={name:[{required:true,message:'请输入名称'}],price:[{required:true,message:'请输入售价'}],stock:[{required:true,message:'请输入库存'}]}
async function loadData(){
  loading.value=true
  try{
    const res=await request.get('/products',{params:{q:search.q||undefined,skip:page.value*pageSize.value-pageSize.value+1,limit:pageSize.value}})
    const d=res.data;const b=d.success?d.data:d
    tableData.value=b.products||b.list||[];meta.value=b.meta||{total:0}
  }finally{loading.value=false}
}
async function loadCategories(){
  try{const res=await request.get('/categories',{params:{skip:1,limit:99}});const d=res.data;categories.value=(d.success?d.data:d).categories||[]}catch{}
}
function showAdd(){isEdit.value=false;Object.assign(form,{name:'',category_id:null,price:0,stock:0,image:'',unit:'',base_unit:'',conversion_rate:1});dialogVisible.value=true}
function showEdit(row:any){
  isEdit.value=true;Object.assign(form,{name:row.name,category_id:row.category_id,price:row.price,stock:row.stock,image:row.image||'',unit:row.unit||'',base_unit:row.base_unit||'',conversion_rate:row.conversion_rate||1});(form as any).id=row.id;dialogVisible.value=true
}
async function handleDelete(id:number){
  try{await ElMessageBox.confirm('确认删除？','提示');await request.delete(`/products/${id}`);ElMessage.success('已删除');loadData()}catch{}
}
async function submitForm(){
  if(!await formRef.value?.validate().catch(()=>false))return
  submitting.value=true
  try{
    const data={name:form.name,category_id:form.category_id,price:form.price,stock:form.stock,image:form.image||undefined,unit:form.unit||undefined,base_unit:form.base_unit||undefined,conversion_rate:form.conversion_rate||1}
    if(isEdit.value){await request.put(`/products/${(form as any).id}`,data);ElMessage.success('已更新')}
    else{await request.post('/products',data);ElMessage.success('已创建')}
    dialogVisible.value=false;loadData()
  }finally{submitting.value=false}
}
onMounted(()=>{loadData();loadCategories()})
</script>
<style scoped>
.mb-4{margin-bottom:16px}.mb-3{margin-bottom:12px}.mt-3{margin-top:12px}.ml-2{margin-left:8px}:deep(.low-stock-row){background-color:#fef0f0!important}
</style>
