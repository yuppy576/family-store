<template>
  <div>
    <el-card shadow="never" class="mb-4">
      <el-form :inline="true" size="default">
        <el-form-item label="操作类型"><el-select v-model="search.action" placeholder="全部" clearable style="width:120px"><el-option label="创建" value="CREATE" /><el-option label="更新" value="UPDATE" /><el-option label="删除" value="DELETE" /></el-select></el-form-item>
        <el-form-item label="资源类型"><el-select v-model="search.resource_type" placeholder="全部" clearable style="width:140px"><el-option label="用户" value="users" /><el-option label="支付方式" value="payments" /><el-option label="分类" value="categories" /><el-option label="商品" value="products" /><el-option label="订单" value="orders" /><el-option label="寄卖人" value="consignors" /><el-option label="寄卖品" value="consignments" /><el-option label="结算" value="settlements" /><el-option label="供应商" value="suppliers" /><el-option label="进货" value="purchases" /></el-select></el-form-item>
        <el-form-item label="用户ID"><el-input v-model="search.user_id" type="number" placeholder="用户ID" clearable style="width:100px" /></el-form-item>
        <el-form-item label="开始时间"><el-date-picker v-model="search.start_time" type="date" placeholder="开始日期" style="width:140px" /></el-form-item>
        <el-form-item label="结束时间"><el-date-picker v-model="search.end_time" type="date" placeholder="结束日期" style="width:140px" /></el-form-item>
        <el-form-item><el-button type="primary" @click="loadData">查询</el-button><el-button @click="resetSearch">重置</el-button></el-form-item>
      </el-form>
    </el-card>
    <el-card shadow="never">
      <el-table :data="tableData" v-loading="loading" stripe border style="width:100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="user_id" label="用户ID" width="90" />
        <el-table-column prop="user_name" label="用户名" width="120" />
        <el-table-column label="操作" width="100">
          <template #default="{row}">
            <el-tag :type="getActionType(row.action)">{{ getActionLabel(row.action) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" label="资源类型" width="120"><template #default="{row}">{{ getResourceLabel(row.resource_type) }}</template></el-table-column>
        <el-table-column prop="resource_id" label="资源ID" width="90" />
        <el-table-column prop="ip_address" label="IP地址" width="140" />
        <el-table-column prop="created_at" label="操作时间" width="170" />
        <el-table-column label="详情" width="120" fixed="right">
          <template #default="{row}"><el-button size="small" type="primary" link @click="showDetail(row)">详情</el-button></template>
        </el-table-column>
      </el-table>
      <el-pagination v-if="meta.total>0" v-model:current-page="page" :page-size="pageSize" :total="meta.total" layout="total,prev,pager,next" @current-change="loadData" class="mt-3" />
    </el-card>
    <el-dialog v-model="detailVisible" title="操作详情" width="600px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="操作类型">{{ getActionLabel(detailRow.action) }}</el-descriptions-item>
        <el-descriptions-item label="资源类型">{{ getResourceLabel(detailRow.resource_type) }}</el-descriptions-item>
        <el-descriptions-item label="资源ID">{{ detailRow.resource_id }}</el-descriptions-item>
        <el-descriptions-item label="操作人">{{ detailRow.user_name }} (ID:{{ detailRow.user_id }})</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ detailRow.ip_address }}</el-descriptions-item>
        <el-descriptions-item label="操作时间">{{ detailRow.created_at }}</el-descriptions-item>
        <el-descriptions-item label="原始数据" v-if="detailRow.old_data"><template #content><pre style="max-height:200px;overflow:auto">{{ detailRow.old_data }}</pre></template></el-descriptions-item>
        <el-descriptions-item label="新数据" v-if="detailRow.new_data"><template #content><pre style="max-height:200px;overflow:auto">{{ detailRow.new_data }}</pre></template></el-descriptions-item>
      </el-descriptions>
      <template #footer><el-button @click="detailVisible=false">关闭</el-button></template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import request from '@/api/request'
const loading=ref(false),detailVisible=ref(false)
const tableData=ref<any[]>([]),detailRow=ref<any>({})
const meta=ref<any>({total:0,limit:10,skip:1}),page=ref(1),pageSize=ref(10)
const search=reactive({action:'',resource_type:'',user_id:'',start_time:'',end_time:''})
function getActionType(action:string){
  switch(action){case 'CREATE':return 'success';case 'UPDATE':return 'warning';case 'DELETE':return 'danger';default:return 'info'}
}
function getActionLabel(action:string){
  switch(action){case 'CREATE':return '创建';case 'UPDATE':return '更新';case 'DELETE':return '删除';default:return action}
}
function getResourceLabel(type:string){
  const map:any={users:'用户',payments:'支付方式',categories:'分类',products:'商品',orders:'订单',consignors:'寄卖人',consignments:'寄卖品',settlements:'结算',suppliers:'供应商',purchases:'进货'}
  return map[type]||type
}
function resetSearch(){Object.assign(search,{action:'',resource_type:'',user_id:'',start_time:'',end_time:''});page.value=1;loadData()}
async function loadData(){
  loading.value=true
  try{
    const params:any={skip:page.value*pageSize.value-pageSize.value+1,limit:pageSize.value}
    if(search.action)params.action=search.action
    if(search.resource_type)params.resource_type=search.resource_type
    if(search.user_id)params.user_id=search.user_id
    if(search.start_time)params.start_time=search.start_time
    if(search.end_time)params.end_time=search.end_time
    const res=await request.get('/audit-logs',{params})
    const d=res.data;const b=d.success?d.data:d
    tableData.value=b.audit_logs||b.list||[];meta.value=b.meta||{total:0}
  }finally{loading.value=false}
}
function showDetail(row:any){detailRow.value=row;detailVisible.value=true}
onMounted(loadData)
</script>
<style scoped>
.mb-4{margin-bottom:16px}.mt-3{margin-top:12px}
</style>
