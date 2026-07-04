<template>
  <div class="pos-container">
    <el-row :gutter="16">
      <!-- 左侧：商品列表 -->
      <el-col :span="14">
        <el-card shadow="never" class="mb-2">
          <el-row :gutter="8">
            <el-col :span="8">
              <el-select v-model="categoryFilter" filterable clearable placeholder="全部分类" style="width:100%" @change="loadProducts">
                <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
              </el-select>
            </el-col>
            <el-col :span="10">
              <el-input v-model="searchText" placeholder="搜索商品名称/条码" clearable @keyup.enter="loadProducts" />
            </el-col>
            <el-col :span="6">
              <el-button type="primary" @click="loadProducts" style="width:100%">搜索</el-button>
            </el-col>
          </el-row>
        </el-card>
        <el-card shadow="never" style="height:calc(100vh - 280px);overflow-y:auto">
          <el-row :gutter="12">
            <el-col :span="8" v-for="p in products" :key="p.id" class="mb-2">
              <el-card shadow="hover" class="product-card" @click="addToCart(p)">
                <div style="font-size:13px;font-weight:600;margin-bottom:4px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">{{ p.name }}</div>
                <div style="color:#e6a23c;font-size:16px;font-weight:700">¥{{ Number(p.price).toFixed(2) }}<span style="font-size:11px;color:#909399">/{{ p.base_unit||p.unit||'个' }}</span></div>
                <div style="color:#909399;font-size:11px">库存: {{ p.stock }}{{ p.unit||'' }}{{ p.base_unit?`(${p.stock*p.conversion_rate}${p.base_unit})`:'' }}</div>
              </el-card>
            </el-col>
          </el-row>
          <el-empty v-if="products.length===0" description="暂无商品" />
        </el-card>
      </el-col>

      <!-- 右侧：购物车 -->
      <el-col :span="10">
        <el-card shadow="never" style="height:calc(100vh - 240px);display:flex;flex-direction:column">
          <template #header><strong>购物车 ({{ cart.length }})</strong></template>
          <div style="flex:1;overflow-y:auto">
            <div v-for="(item,idx) in cart" :key="idx" class="cart-item">
              <div class="cart-item-info">
                <div class="cart-item-name">{{ item.name }}</div>
                <div class="cart-item-price">¥{{ (item.price * item.qty).toFixed(2) }}</div>
              </div>
              <div style="font-size:11px;color:#909399;margin-bottom:4px">{{ item.qty }}{{ item.base_unit||item.unit||'个' }}</div>
              <div class="cart-item-actions">
                <el-button type="danger" link size="small" @click="item.qty>1?item.qty--:removeItem(idx)">−</el-button>
                <span class="cart-qty">{{ item.qty }}</span>
                <el-button type="primary" link size="small" @click="item.qty++">+</el-button>
                <el-button type="danger" link size="small" @click="removeItem(idx)">✕</el-button>
              </div>
            </div>
            <el-empty v-if="cart.length===0" description="请选择商品" :image-size="60" />
          </div>
          <el-divider style="margin:8px 0" />
          <div class="cart-total">
            <span>合计:</span>
            <span class="total-amount">¥{{ totalAmount.toFixed(2) }}</span>
          </div>
          <div class="mb-2">
            <el-select v-model="paymentId" placeholder="支付方式" style="width:100%">
              <el-option v-for="pm in payments" :key="pm.id" :label="pm.name" :value="pm.id" />
            </el-select>
          </div>
          <el-button type="success" :disabled="cart.length===0||!paymentId" :loading="submitting" size="large" style="width:100%" @click="submitOrder">
            结算收款 ¥{{ totalAmount.toFixed(2) }}
          </el-button>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/api/request'

const products=ref<any[]>([]),categories=ref<any[]>([]),payments=ref<any[]>([])
const searchText=ref(''),categoryFilter=ref(null),paymentId=ref(null),submitting=ref(false)
const cart=ref<any[]>([])

const totalAmount=computed(()=>cart.value.reduce((s,i)=>s+i.price*i.qty,0))

async function loadProducts(){
  try{
    const params:any={skip:1,limit:50}
    if(searchText.value)params.q=searchText.value
    if(categoryFilter.value)params.category_id=categoryFilter.value
    const res=await request.get('/products',{params})
    const d=res.data;const b=d.success?d.data:d
    products.value=b.products||b.list||[]
  }catch{}
}
async function loadCategories(){
  try{const res=await request.get('/categories',{params:{skip:1,limit:99}});const d=res.data;categories.value=(d.success?d.data:d).categories||[]}catch{}
}
async function loadPayments(){
  try{const res=await request.get('/payments',{params:{skip:1,limit:99}});const d=res.data;const b=d.success?d.data:d;payments.value=b.payments||b.list||[];if(payments.value.length)paymentId.value=payments.value[0].id}catch{}
}
function addToCart(p:any){
  const exist=cart.value.find(i=>i.id===p.id)
  if(exist){exist.qty++}else{cart.value.push({id:p.id,name:p.name,price:p.price,qty:1,unit:p.unit,base_unit:p.base_unit})}
}
function removeItem(idx:number){cart.value.splice(idx,1)}
async function submitOrder(){
  if(!paymentId.value){ElMessage.warning('请选择支付方式');return}
  submitting.value=true
  try{
    const res=await request.post('/orders',{
      payment_id:paymentId.value,
      total_paid:totalAmount.value,
      customer_name:'散客',
      products:cart.value.map(i=>({product_id:i.id,qty:i.qty}))
    })
    ElMessage.success(`收款成功！¥${totalAmount.value.toFixed(2)}`)
    cart.value=[]
    loadProducts()
  }catch{}finally{submitting.value=false}
}
onMounted(()=>{loadProducts();loadCategories();loadPayments()})
</script>
<style scoped>
.pos-container{padding:0;background:#f5f5f7;height:calc(100vh - 120px)}
.mb-2{margin-bottom:8px}
.product-card{cursor:pointer;transition:all .2s;user-select:none}
.product-card:hover{transform:translateY(-2px);box-shadow:0 4px 12px rgba(0,0,0,.1)}
.cart-item{padding:8px 0;border-bottom:1px solid #f0f0f0}
.cart-item-info{display:flex;justify-content:space-between;margin-bottom:4px}
.cart-item-name{font-size:13px}
.cart-item-price{color:#e6a23c;font-weight:600}
.cart-item-actions{display:flex;align-items:center;gap:4px;justify-content:flex-end}
.cart-qty{width:24px;text-align:center;font-weight:600}
.cart-total{display:flex;justify-content:space-between;align-items:center;font-size:18px;font-weight:700;margin-bottom:8px}
.total-amount{color:#e6a23c;font-size:24px}
</style>
