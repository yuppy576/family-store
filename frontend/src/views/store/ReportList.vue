<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>报表中心</span>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            @change="loadData"
          />
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="6">
          <el-statistic title="订单总数" :value="salesStats?.order_count || 0" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="总销售额" :value="salesStats?.total_amount || 0" prefix="¥" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="平均客单价" :value="avgOrderAmount" prefix="¥" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="商品种类" :value="salesStats?.product_count || 0" />
        </el-col>
      </el-row>

      <div style="margin-top:20px">
        <h4>每日销售趋势</h4>
        <div style="height:300px">
          <canvas ref="chartCanvas"></canvas>
        </div>
      </div>

      <div style="margin-top:20px">
        <h4>每日销售明细</h4>
        <el-table :data="dailySales" v-loading="loading">
          <el-table-column prop="date" label="日期" width="120" />
          <el-table-column prop="order_count" label="订单数" width="100" />
          <el-table-column prop="total_amount" label="销售额" width="120">
            <template #default="{ row }">¥{{ row.total_amount.toFixed(2) }}</template>
          </el-table-column>
        </el-table>
        <el-empty v-if="!loading && dailySales.length === 0" description="暂无数据" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import request from '@/api/request'

const loading = ref(false)
const dateRange = ref<Date[]>([])
const salesStats = ref<any>(null)
const dailySales = ref<any[]>([])
const chartCanvas = ref<HTMLCanvasElement | null>(null)

const avgOrderAmount = computed(() => {
  if (salesStats.value && salesStats.value.order_count > 0) {
    return (salesStats.value.total_amount / salesStats.value.order_count).toFixed(2)
  }
  return 0
})

async function loadData() {
  loading.value = true
  try {
    let startDate = ''
    let endDate = ''
    if (dateRange.value && dateRange.value.length === 2) {
      startDate = dateRange.value[0].toISOString().split('T')[0]
      endDate = dateRange.value[1].toISOString().split('T')[0]
    }

    const [statsRes, dailyRes] = await Promise.all([
      request.get('/reports/sales/stats', { params: { start_date: startDate, end_date: endDate } }),
      request.get('/reports/sales/daily', { params: { start_date: startDate, end_date: endDate } }),
    ])

    const statsBody = statsRes.data
    if (statsBody.success) {
      salesStats.value = statsBody.data
    }

    const dailyBody = dailyRes.data
    if (dailyBody.success) {
      dailySales.value = dailyBody.data || []
    }

    await nextTick()
    renderChart()
  } catch {
    salesStats.value = null
    dailySales.value = []
  } finally {
    loading.value = false
  }
}

function renderChart() {
  const canvas = chartCanvas.value
  if (!canvas || dailySales.value.length === 0) return

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  canvas.width = canvas.offsetWidth * 2
  canvas.height = canvas.offsetHeight * 2
  ctx.scale(2, 2)

  const width = canvas.offsetWidth
  const height = canvas.offsetHeight
  const padding = 40

  const labels = dailySales.value.map((d: any) => d.date?.substring(5) || '')
  const values = dailySales.value.map((d: any) => d.total_amount || 0)
  const maxValue = Math.max(...values, 1)

  ctx.clearRect(0, 0, width, height)

  ctx.strokeStyle = '#e4e7ed'
  ctx.lineWidth = 1
  for (let i = 0; i <= 5; i++) {
    const y = padding + (height - padding * 2) * (i / 5)
    ctx.beginPath()
    ctx.moveTo(padding, y)
    ctx.lineTo(width - padding, y)
    ctx.stroke()

    ctx.fillStyle = '#909399'
    ctx.font = '12px sans-serif'
    ctx.textAlign = 'right'
    ctx.fillText(((maxValue * (5 - i)) / 5).toFixed(0), padding - 10, y + 4)
  }

  ctx.fillStyle = '#909399'
  ctx.font = '12px sans-serif'
  ctx.textAlign = 'center'
  labels.forEach((label: string, i: number) => {
    const x = padding + (width - padding * 2) * (i / (labels.length - 1))
    ctx.fillText(label, x, height - 10)
  })

  const gradient = ctx.createLinearGradient(0, padding, 0, height - padding)
  gradient.addColorStop(0, 'rgba(102, 126, 234, 0.3)')
  gradient.addColorStop(1, 'rgba(102, 126, 234, 0)')

  ctx.beginPath()
  ctx.moveTo(padding, height - padding)
  values.forEach((value: number, i: number) => {
    const x = padding + (width - padding * 2) * (i / (values.length - 1))
    const y = height - padding - ((value / maxValue) * (height - padding * 2))
    ctx.lineTo(x, y)
  })
  ctx.lineTo(width - padding, height - padding)
  ctx.closePath()
  ctx.fillStyle = gradient
  ctx.fill()

  ctx.beginPath()
  values.forEach((value: number, i: number) => {
    const x = padding + (width - padding * 2) * (i / (values.length - 1))
    const y = height - padding - ((value / maxValue) * (height - padding * 2))
    if (i === 0) ctx.moveTo(x, y)
    else ctx.lineTo(x, y)
  })
  ctx.strokeStyle = '#667eea'
  ctx.lineWidth = 2
  ctx.stroke()

  values.forEach((value: number, i: number) => {
    const x = padding + (width - padding * 2) * (i / (values.length - 1))
    const y = height - padding - ((value / maxValue) * (height - padding * 2))
    ctx.beginPath()
    ctx.arc(x, y, 4, 0, Math.PI * 2)
    ctx.fillStyle = '#667eea'
    ctx.fill()
    ctx.strokeStyle = '#fff'
    ctx.lineWidth = 2
    ctx.stroke()
  })
}

onMounted(() => {
  const now = new Date()
  const twoWeeksAgo = new Date()
  twoWeeksAgo.setDate(twoWeeksAgo.getDate() - 14)
  dateRange.value = [twoWeeksAgo, now]
  loadData()
})
</script>

<style scoped>
.page-container { padding: 0; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
