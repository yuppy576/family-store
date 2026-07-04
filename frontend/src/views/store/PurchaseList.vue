<template>
  <div>
    <div class="mb-3"><el-button type="primary" @click="showAdd"><el-icon><Plus /></el-icon> 新增进货</el-button></div>
    <el-card shadow="never">
      <el-table :data="tableData" v-loading="loading" stripe border style="width:100%">
        <el-table-column prop="id" label="单号" width="80" />
        <el-table-column label="供应商" min-width="150"><template #default="{row}">{{ row.supplier?.name||'-' }}</template></el-table-column>
        <el-table-column prop="operator" label="经办人" width="100" />
        <el-table-column label="金额" width="120" align="right"><template #default="{row}">¥{{ row.total_amount?Number(row.total_amount).toLocaleString():'0' }}</template></el-table-column>
        <el-table-column label="状态" width="80" align="center">
          <template #default="{row}"><el-tag :type="row.status==='COMPLETED'?'success':'info'" size="small">{{ row.status==='COMPLETED'?'已完成':row.status }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column label="时间" width="160"><template #default="{row}">{{ formatTime(row.created_at) }}</template></el-table-column>
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{row}"><el-button size="small" type="primary" link @click="showDetail(row)">详情</el-button></template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!loading&&tableData.length===0" description="暂无进货记录" />
    </el-card>

    <!-- 新增进货对话框 -->
    <el-dialog v-model="addVisible" title="新增进货" width="650px" :close-on-click-modal="false">
      <el-form :model="purchaseForm" ref="purchaseFormRef" label-width="80px">
        <el-form-item label="供应商" prop="supplier_id" :rules="[{required:true,message:'请选择供应商'}]">
          <el-select v-model="purchaseForm.supplier_id" filterable placeholder="选择供应商" style="width:100%">
            <el-option v-for="s in suppliers" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="purchaseForm.remark" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <el-divider content-position="left">进货商品</el-divider>
      <div v-for="(item,idx) in items" :key="idx" style="margin-bottom:8px;padding:8px;background:#f9f9fb;border-radius:6px">
        <el-row :gutter="8">
          <el-col :span="9"><el-select v-model="item.product_id" filterable placeholder="选择商品" style="width:100%" @change="(id:number)=>{const p=products.find(x=>x.id===id);if(p)item.unit_price=p.price||0}">
            <el-option v-for="p in products" :key="p.id" :label="`${p.name} (¥${p.price})`" :value="p.id" />
          </el-select></el-col>
          <el-col :span="5"><el-input-number v-model="item.quantity" :min="1" placeholder="数量" style="width:100%" /></el-col>
          <el-col :span="5"><el-input-number v-model="item.unit_price" :min="0" :precision="2" placeholder="单价" style="width:100%" /></el-col>
          <el-col :span="3" style="text-align:right;line-height:32px">¥{{ (item.quantity*item.unit_price).toFixed(2) }}</el-col>
          <el-col :span="2"><el-button type="danger" link @click="items.splice(idx,1)">✕</el-button></el-col>
        </el-row>
      </div>
      <el-button type="primary" link @click="items.push({product_id:null,quantity:1,unit_price:0})"><el-icon><Plus /></el-icon> 添加商品</el-button>
      <template #footer>
        <el-button @click="addVisible=false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitPurchase" :disabled="items.length===0">保存进货</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="进货详情" width="550px">
      <el-descriptions :column="2" border size="small">
        <el-descriptions-item label="单号">{{ detail?.id }}</el-descriptions-item>
        <el-descriptions-item label="供应商">{{ detail?.supplier?.name||'-' }}</el-descriptions-item>
        <el-descriptions-item label="经办人">{{ detail?.operator||'-' }}</el-descriptions-item>
        <el-descriptions-item label="总金额">¥{{ detail?.total_amount?Number(detail.total_amount).toLocaleString():'0' }}</el-descriptions-item>
        <el-descriptions-item label="状态"><el-tag size="small">{{ detail?.status }}</el-tag></el-descriptions-item>
        <el-descriptions-item label="时间">{{ formatTime(detail?.created_at) }}</el-descriptions-item>
      </el-descriptions>
      <el-divider content-position="left">商品明细</el-divider>
      <el-table :data="detailItems" border size="small">
        <el-table-column label="商品" min-width="150"><template #default="{row}">{{ row.product?.name||'已删除' }}</template></el-table-column>
        <el-table-column prop="quantity" label="数量" width="80" />
        <el-table-column label="单价" width="100"><template #default="{row}">¥{{ Number(row.unit_price).toFixed(2) }}</template></el-table-column>
        <el-table-column label="小计" width="100"><template #default="{row}">¥{{ Number(row.total_price).toFixed(2) }}</template></el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/api/request'
const loading=ref(false),submitting=ref(false),addVisible=ref(false),detailVisible=ref(false),purchaseFormRef=ref()
const tableData=ref<any[]>([]),suppliers=ref<any[]>([]),products=ref<any[]>([]),detailItems=ref<any[]>([])
const detail=ref<any>(null)
const purchaseForm=reactive<any>({supplier_id:null,remark:''})
const items=ref<any[]>([])

function formatTime(t:string){return t?t.replace('T',' ').substring(0,19):'-'}

async function loadData(){
  loading.value=true
  try{const res=await request.get('/purchases',{params:{skip:0,limit:99}});const d=res.data;const b=d.success?d.data:d;tableData.value=b.purchases||b.list||[]}
  finally{loading.value=false}
}
async function loadOptions(){
  try{
    const[sr,pr]=await Promise.all([request.get('/suppliers?skip=0&limit=99'),request.get('/products?skip=1&limit=99')])
    const sd=sr.data;const pd=pr.data
    suppliers.value=(sd.success?sd.data:sd).suppliers||[]
    const pp=pd.success?pd.data:pd;products.value=pp.products||[]
  }catch{}
}
function showAdd(){purchaseForm.supplier_id=null;purchaseForm.remark='';items.value=[{product_id:null,quantity:1,unit_price:0}];addVisible.value=true}
async function submitPurchase(){
  if(!purchaseForm.supplier_id){ElMessage.warning('请选择供应商');return}
  if(items.value.length===0){ElMessage.warning('请添加至少一个商品');return}
  submitting.value=true
  try{
    await request.post('/purchases',{supplier_id:purchaseForm.supplier_id,remark:purchaseForm.remark||undefined,items:items.value.map(i=>({product_id:i.product_id,quantity:i.quantity,unit_price:i.unit_price}))})
    ElMessage.success('进货成功');addVisible.value=false;loadData()
  }catch{}finally{submitting.value=false}
}
async function showDetail(row:any){
  detail.value=row
  try{
    const res=await request.get(`/purchases/${row.id}/items`).catch(()=>null)
    if(res){const d=res.data;detailItems.value=(d.success?d.data:d).list||[]}else{detailItems.value=[]}
  }catch{detailItems.value=[]}
  detailVisible.value=true
}
onMounted(()=>{loadData();loadOptions()})
</script>
<style scoped>
.mb-3{margin-bottom:12px}
</style>
