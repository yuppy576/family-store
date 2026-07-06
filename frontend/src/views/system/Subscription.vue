<template>
  <div class="subscription-page">
    <el-card class="subscription-card">
      <template #header>
        <div class="card-header">
          <h3>订阅管理</h3>
        </div>
      </template>

      <div v-if="loading" class="loading">
        <el-loading :text="'加载中...'" />
      </div>

      <div v-else class="subscription-content">
        <div class="subscription-status">
          <div class="status-badge" :class="statusClass">
            {{ statusText }}
          </div>
        </div>

        <div class="subscription-info">
          <div class="info-item">
            <label>订阅计划</label>
            <span class="value">{{ planText }}</span>
          </div>
          <div class="info-item">
            <label>开始日期</label>
            <span class="value">{{ subscription.start_date || '-' }}</span>
          </div>
          <div class="info-item">
            <label>到期日期</label>
            <span class="value" :class="{ 'expiring': isExpiring }">
              {{ subscription.end_date || '永久' }}
            </span>
          </div>
          <div v-if="daysRemaining !== null" class="info-item">
            <label>剩余天数</label>
            <span class="value" :class="{ 'expiring': daysRemaining <= 7 }">
              {{ daysRemaining }} 天
            </span>
          </div>
        </div>

        <div class="renew-section">
          <h4>续费/升级</h4>
          <el-form ref="renewForm" :model="renewFormData" label-width="80px">
            <el-form-item label="选择方案">
              <el-select v-model="renewFormData.plan" placeholder="请选择方案">
                <el-option label="个人版 (¥199/月)" value="PERSONAL" />
                <el-option label="专业版 (¥299/月)" value="PROFESSIONAL" />
                <el-option label="买断版 (¥2980)" value="LIFETIME" />
              </el-select>
            </el-form-item>
            <el-form-item label="购买时长">
              <el-select v-model="renewFormData.months" placeholder="请选择时长">
                <el-option label="1个月" :value="1" />
                <el-option label="3个月" :value="3" />
                <el-option label="6个月" :value="6" />
                <el-option label="12个月" :value="12" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <div class="price-display">
                合计：<span class="price">¥{{ totalPrice }}</span>
              </div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="renewLoading" @click="handleRenew">
                {{ renewLoading ? '处理中...' : '确认续费' }}
              </el-button>
              <el-button @click="handleActivate" v-if="subscription.status === 'FROZEN'">
                激活订阅
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <div class="plan-features">
          <h4>方案对比</h4>
          <el-table :data="planFeatures" border>
            <el-table-column prop="feature" label="功能" />
            <el-table-column prop="personal" label="个人版" />
            <el-table-column prop="professional" label="专业版" />
            <el-table-column prop="lifetime" label="买断版" />
          </el-table>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/api/request'

interface SubscriptionData {
  id: number
  store_id: number
  plan: string
  status: string
  start_date: string
  end_date: string | null
  created_at: string
  updated_at: string
}

const loading = ref(false)
const renewLoading = ref(false)
const subscription = ref<SubscriptionData>({} as SubscriptionData)

const renewForm = ref()
const renewFormData = reactive({
  plan: '',
  months: 1,
})

const planFeatures = ref([
  { feature: '寄卖人管理', personal: '✓', professional: '✓', lifetime: '✓' },
  { feature: '寄卖品管理', personal: '✓', professional: '✓', lifetime: '✓' },
  { feature: '结算管理', personal: '✓', professional: '✓', lifetime: '✓' },
  { feature: '过户进度', personal: '✓', professional: '✓', lifetime: '✓' },
  { feature: '商品管理', personal: '✓', professional: '✓', lifetime: '✓' },
  { feature: 'POS收银', personal: '✓', professional: '✓', lifetime: '✓' },
  { feature: '供应商管理', personal: '✓', professional: '✓', lifetime: '✓' },
  { feature: '进货管理', personal: '✓', professional: '✓', lifetime: '✓' },
  { feature: '审计日志', personal: '✓', professional: '✓', lifetime: '✓' },
  { feature: '库存预警', personal: '✓', professional: '✓', lifetime: '✓' },
  { feature: '销售报表', personal: '✓', professional: '✓', lifetime: '✓' },
  { feature: '数据导出', personal: '✓', professional: '✓', lifetime: '✓' },
  { feature: '多用户支持', personal: '3人', professional: '不限', lifetime: '不限' },
  { feature: '技术支持', personal: '邮件', professional: '优先', lifetime: '专属' },
])

const planText = computed(() => {
  const plans: Record<string, string> = {
    TRIAL: '免费试用',
    PERSONAL: '个人版',
    PROFESSIONAL: '专业版',
    LIFETIME: '买断版',
  }
  return plans[subscription.value.plan] || subscription.value.plan
})

const statusText = computed(() => {
  const statuses: Record<string, string> = {
    TRIAL: '试用中',
    ACTIVE: '已激活',
    FROZEN: '已冻结',
    EXPIRED: '已过期',
  }
  return statuses[subscription.value.status] || subscription.value.status
})

const statusClass = computed(() => {
  const classes: Record<string, string> = {
    TRIAL: 'status-trial',
    ACTIVE: 'status-active',
    FROZEN: 'status-frozen',
    EXPIRED: 'status-expired',
  }
  return classes[subscription.value.status] || 'status-active'
})

const daysRemaining = computed(() => {
  if (!subscription.value.end_date) return null
  const endDate = new Date(subscription.value.end_date)
  const now = new Date()
  const diff = endDate.getTime() - now.getTime()
  return Math.max(0, Math.ceil(diff / (1000 * 60 * 60 * 24)))
})

const isExpiring = computed(() => {
  return daysRemaining.value !== null && daysRemaining.value <= 7
})

const totalPrice = computed(() => {
  const prices: Record<string, number> = {
    PERSONAL: 199,
    PROFESSIONAL: 299,
    LIFETIME: 2980,
  }
  const basePrice = prices[renewFormData.plan] || 0
  if (renewFormData.plan === 'LIFETIME') {
    return basePrice
  }
  return basePrice * renewFormData.months
})

async function loadSubscription() {
  loading.value = true
  try {
    const res = await request.get('/subscription')
    const body = res.data
    if (body.success) {
      subscription.value = body.data
    }
  } catch (err: any) {
    ElMessage.error('获取订阅信息失败')
  } finally {
    loading.value = false
  }
}

async function handleRenew() {
  if (!renewFormData.plan) {
    ElMessage.warning('请选择订阅方案')
    return
  }
  if (!renewFormData.months) {
    ElMessage.warning('请选择购买时长')
    return
  }

  renewLoading.value = true
  try {
    const res = await request.post('/subscription/renew', {
      plan: renewFormData.plan,
      months: renewFormData.months,
    })
    const body = res.data
    if (body.success) {
      ElMessage.success('续费成功')
      subscription.value = body.data
      renewFormData.plan = ''
      renewFormData.months = 1
    }
  } catch (err: any) {
    ElMessage.error('续费失败')
  } finally {
    renewLoading.value = false
  }
}

async function handleActivate() {
  try {
    const res = await request.post('/subscription/activate')
    const body = res.data
    if (body.success) {
      ElMessage.success('激活成功')
      subscription.value = body.data
    }
  } catch (err: any) {
    ElMessage.error('激活失败')
  }
}

onMounted(loadSubscription)
</script>

<style scoped>
.subscription-page { padding: 0; }
.subscription-card { max-width: 700px; margin: 0 auto; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.card-header h3 { margin: 0; font-size: 18px; }
.loading { display: flex; justify-content: center; padding: 40px; }
.subscription-content { padding-top: 20px; }

.subscription-status { margin-bottom: 24px; }
.status-badge {
  display: inline-block;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
}
.status-trial { background: #fff7e6; color: #e6a23c; }
.status-active { background: #f0f9eb; color: #67c23a; }
.status-frozen { background: #fef0f0; color: #f56c6c; }
.status-expired { background: #f5f5f5; color: #909399; }

.subscription-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 32px;
  padding: 20px;
  background: #fafafa;
  border-radius: 8px;
}
.info-item { display: flex; flex-direction: column; }
.info-item label { font-size: 13px; color: #909399; margin-bottom: 4px; }
.info-item .value { font-size: 15px; font-weight: 500; color: #303133; }
.info-item .value.expiring { color: #f56c6c; }

.renew-section { margin-bottom: 32px; padding: 20px; background: #fafafa; border-radius: 8px; }
.renew-section h4 { margin: 0 0 16px; font-size: 16px; }
.price-display { padding: 10px 0; }
.price-display .price { font-size: 24px; font-weight: 600; color: #e6a23c; }

.plan-features { margin-top: 20px; }
.plan-features h4 { margin: 0 0 16px; font-size: 16px; }

@media (max-width: 768px) {
  .subscription-info { grid-template-columns: 1fr; }
}
</style>