<template>
  <div class="product-list">
    <el-card shadow="never">
      <el-empty description="商店功能开发中" :image-size="200">
        <template #image>
          <el-icon :size="120" color="#c0c4cc">
            <Tools />
          </el-icon>
        </template>
        <template #description>
          <p class="empty-text">商店管理模块正在开发中，敬请期待</p>
          <p class="empty-subtext">功能包括：商品上架、库存管理、订单处理、收银结算</p>
        </template>
      </el-empty>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { Tools } from '@element-plus/icons-vue'
</script>

<style scoped>
.product-list {
  max-width: 1400px;
}

.empty-text {
  font-size: 16px;
  color: #606266;
  margin: 0 0 8px;
}

.empty-subtext {
  font-size: 13px;
  color: #909399;
  margin: 0;
}
</style>
