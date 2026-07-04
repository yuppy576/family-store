<template>
  <div class="main-layout">
    <!-- Header -->
    <header class="layout-header">
      <div class="header-left">
        <el-button class="menu-toggle" @click="menuOpen=!menuOpen">
          <el-icon><Menu /></el-icon>
        </el-button>
        <span class="app-title">家族门店</span>
        <el-button-group class="mode-switcher">
          <el-button :type="currentMode==='store'?'primary':'default'" size="small" @click="switchMode('store')">商店</el-button>
          <el-button :type="currentMode==='consignment'?'primary':'default'" size="small" @click="switchMode('consignment')">寄卖</el-button>
        </el-button-group>
      </div>
      <div class="header-right">
        <span class="user-name">{{ authStore.user?.name||authStore.user?.email }}</span>
        <el-button type="danger" size="small" plain @click="handleLogout">退出</el-button>
      </div>
    </header>

    <!-- 到期提醒 -->
    <div v-if="expiringItems.length>0" class="expiring-banner">
      <el-alert :title="`${expiringItems.length} 件寄卖品即将到期`" type="warning" show-icon :closable="false">
        <template #default>
          <div v-for="item in expiringItems.slice(0,5)" :key="item.id" style="font-size:13px;margin:2px 0">
            · {{ item.name }}（到期：{{ formatDate(item.contract_end) }}）
          </div>
          <div v-if="expiringItems.length>5" style="font-size:12px;color:#909399">还有 {{ expiringItems.length-5 }} 件...</div>
        </template>
      </el-alert>
    </div>

    <div class="layout-body">
      <!-- 遮罩层（手机端菜单打开时） -->
      <div v-if="menuOpen" class="menu-overlay" @click="menuOpen=false" />

      <!-- Sidebar -->
      <aside class="layout-sidebar" :class="{ 'sidebar-open': menuOpen }">
        <el-menu :default-active="activeMenu" router @select="onMenuSelect">
          <template v-if="currentMode==='consignment'">
            <el-sub-menu index="consignment">
              <template #title><el-icon><List /></el-icon><span>寄卖管理</span></template>
              <el-menu-item index="/consignment/consignors"><el-icon><User /></el-icon><span>寄卖人管理</span></el-menu-item>
              <el-menu-item index="/consignment/items"><el-icon><Goods /></el-icon><span>寄卖品管理</span></el-menu-item>
              <el-menu-item index="/consignment/settlements"><el-icon><Money /></el-icon><span>结算管理</span></el-menu-item>
              <el-menu-item index="/consignment/transfer"><el-icon><Guide /></el-icon><span>过户进度</span></el-menu-item>
            </el-sub-menu>
          </template>
          <template v-if="currentMode==='store'">
            <el-sub-menu index="store">
              <template #title><el-icon><Shop /></el-icon><span>商店管理</span></template>
              <el-menu-item index="/store/products"><el-icon><Box /></el-icon><span>商品管理</span></el-menu-item>
              <el-menu-item index="/store/suppliers"><el-icon><User /></el-icon><span>供应商管理</span></el-menu-item>
              <el-menu-item index="/store/purchases"><el-icon><Download /></el-icon><span>进货管理</span></el-menu-item>
              <el-menu-item index="/store/pos"><el-icon><ShoppingCart /></el-icon><span>POS收银</span></el-menu-item>
            </el-sub-menu>
          </template>
        </el-menu>
      </aside>

      <!-- Main Content -->
      <main class="layout-main" @click="menuOpen=false">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import request from '@/api/request'
import { Menu, List, User, Goods, Shop, Box, ShoppingCart, Money, Guide, Download } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const menuOpen = ref(false)
const expiringItems = ref<any[]>([])
function formatDate(d:string) { return d?d.substring(0,10):'-' }
async function loadExpiring() {
  try {
    const res = await request.get('/consignment/expiring')
    const body = res.data
    expiringItems.value = body.success ? (body.data||[]) : []
  } catch { expiringItems.value = [] }
}
function onMenuSelect() { menuOpen.value = false }
onMounted(loadExpiring)

const activeMenu = computed(() => route.path)

const currentMode = computed(() => {
  if (route.path.startsWith('/store')) return 'store'
  return 'consignment'
})

function switchMode(mode: string) {
  if (mode === 'store') {
    router.push('/store/products')
  } else {
    router.push('/consignment/consignors')
  }
}

function handleLogout() {
  authStore.logout()
}
</script>

<style scoped>
.main-layout { display:flex; flex-direction:column; height:100vh; }
.layout-header { display:flex; justify-content:space-between; align-items:center; height:56px; padding:0 16px; background:#fff; border-bottom:1px solid #e4e7ed; z-index:100; }
.header-left { display:flex; align-items:center; gap:12px; }
.app-title { font-size:18px; font-weight:600; color:#303133; white-space:nowrap; }
.mode-switcher { flex-shrink:0; }
.header-right { display:flex; align-items:center; gap:8px; }
.user-name { font-size:14px; color:#606266; display:none; }
.layout-body { display:flex; flex:1; overflow:hidden; position:relative; }

/* 手机菜单按钮 */
.menu-toggle { display:flex !important; }

/* 遮罩层 */
.menu-overlay { display:none; }
.expiring-banner { padding:0 16px; background:#fdf6ec; border-bottom:1px solid #faecd8; }

/* 侧边栏 */
.layout-sidebar { width:220px; background:#fff; border-right:1px solid #e4e7ed; overflow-y:auto; flex-shrink:0; transition:transform .25s; }
.layout-sidebar .el-menu { border-right:none; }
.layout-main { flex:1; padding:16px; background:#f5f7fa; overflow-y:auto; }

/* 手机端适配 */
@media (max-width:768px) {
  .user-name { display:none; }
  .app-title { font-size:15px; }
  .layout-sidebar { position:fixed; top:56px; left:0; bottom:0; z-index:200; transform:translateX(-100%); }
  .layout-sidebar.sidebar-open { transform:translateX(0); }
  .menu-overlay { display:block; position:fixed; top:56px; left:0; right:0; bottom:0; background:rgba(0,0,0,.3); z-index:199; }
  .layout-main { padding:10px; }
}
@media (min-width:769px) {
  .menu-toggle { display:none !important; }
  .user-name { display:inline; }
}
</style>
