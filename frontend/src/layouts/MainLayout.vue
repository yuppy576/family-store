<template>
  <div class="main-layout">
    <!-- Header -->
    <header class="layout-header">
      <div class="header-left">
        <span class="app-title">家族门店管理系统</span>
        <el-button-group class="mode-switcher">
          <el-button
            :type="currentMode === 'store' ? 'primary' : 'default'"
            size="small"
            @click="switchMode('store')"
          >
            商店
          </el-button>
          <el-button
            :type="currentMode === 'consignment' ? 'primary' : 'default'"
            size="small"
            @click="switchMode('consignment')"
          >
            寄卖
          </el-button>
        </el-button-group>
      </div>
      <div class="header-right">
        <span class="user-name">{{ authStore.user?.name || authStore.user?.email }}</span>
        <el-button type="danger" size="small" plain @click="handleLogout">退出</el-button>
      </div>
    </header>

    <div class="layout-body">
      <!-- Sidebar -->
      <aside class="layout-sidebar">
        <el-menu
          :default-active="activeMenu"
          router
          :collapse="false"
        >
          <template v-if="currentMode === 'consignment'">
            <el-sub-menu index="consignment">
              <template #title>
                <el-icon><List /></el-icon>
                <span>寄卖管理</span>
              </template>
              <el-menu-item index="/consignment/consignors">
                <el-icon><User /></el-icon>
                <span>寄卖人管理</span>
              </el-menu-item>
              <el-menu-item index="/consignment/items">
                <el-icon><Goods /></el-icon>
                <span>寄卖品管理</span>
              </el-menu-item>
              <el-menu-item index="/consignment/settlements">
                <el-icon><Money /></el-icon>
                <span>结算管理</span>
              </el-menu-item>
              <el-menu-item index="/consignment/transfer">
                <el-icon><Guide /></el-icon>
                <span>过户进度</span>
              </el-menu-item>
            </el-sub-menu>
          </template>
          <template v-if="currentMode === 'store'">
            <el-sub-menu index="store">
              <template #title>
                <el-icon><Shop /></el-icon>
                <span>商店管理</span>
              </template>
              <el-menu-item index="/store/products">
                <el-icon><Box /></el-icon>
                <span>商品管理</span>
              </el-menu-item>
            </el-sub-menu>
          </template>
        </el-menu>
      </aside>

      <!-- Main Content -->
      <main class="layout-main">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { List, User, Goods, Shop, Box } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

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
.main-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.layout-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 56px;
  padding: 0 20px;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  z-index: 100;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.app-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  white-space: nowrap;
}

.mode-switcher {
  flex-shrink: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-name {
  font-size: 14px;
  color: #606266;
}

.layout-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.layout-sidebar {
  width: 220px;
  background: #fff;
  border-right: 1px solid #e4e7ed;
  overflow-y: auto;
  flex-shrink: 0;
}

.layout-sidebar .el-menu {
  border-right: none;
}

.layout-main {
  flex: 1;
  padding: 20px;
  background: #f5f7fa;
  overflow-y: auto;
}
</style>
