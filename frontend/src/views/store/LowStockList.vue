<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>库存预警</span>
          <el-input v-model="search" placeholder="搜索商品名称" style="width:200px" @input="handleSearch" />
        </div>
      </template>
      <el-table :data="tableData" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="商品名称" />
        <el-table-column prop="category.name" label="分类" />
        <el-table-column prop="stock" label="当前库存" width="120">
          <template #default="{ row }">
            <el-tag :type="row.stock < 10 ? 'danger' : 'warning'">{{ row.stock }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="price" label="单价" width="120">
          <template #default="{ row }">¥{{ row.price.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">{{ row.created_at?.substring(0, 19) }}</template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!loading && tableData.length === 0" description="暂无低库存商品" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/api/request'

const loading = ref(false)
const search = ref('')
const tableData = ref<any[]>([])

async function loadData() {
  loading.value = true
  try {
    const res = await request.get('/products/low-stock')
    const body = res.data
    if (body.success) {
      tableData.value = body.data || []
    }
  } catch {
    tableData.value = []
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  loadData()
}

onMounted(loadData)
</script>

<style scoped>
.page-container { padding: 0; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
