<template>
  <div class="login-container">
    <el-card class="login-card" shadow="always">
      <template #header>
        <div class="login-header">
          <h2>家族门店管理系统</h2>
          <p class="login-subtitle">寄卖行 · 零售管理</p>
          <p class="login-domain" v-if="currentDomain">{{ currentDomain }}</p>
        </div>
      </template>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="0"
        size="large"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="email">
          <el-input
            v-model="form.email"
            placeholder="请输入邮箱"
            :prefix-icon="Message"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            show-password
            :prefix-icon="Lock"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            style="width: 100%"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登 录' }}
          </el-button>
        </el-form-item>
        <el-form-item class="register-link">
          <a href="https://store.yuppy576.top" target="_blank" class="register-btn">还没有账号？立即注册免费试用</a>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import { Message, Lock } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const authStore = useAuthStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  email: '',
  password: '',
})

const rules: FormRules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
  ],
}

const currentDomain = computed(() => {
  const host = window.location.hostname
  if (host.endsWith('.store.yuppy576.top')) {
    return host
  }
  return ''
})

async function handleLogin() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await authStore.login(form.email, form.password)
  } catch (err: any) {
    const msg = err?.response?.data?.messages?.[0] || '登录失败，请检查账号密码'
    ElMessage.error(msg)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  border-radius: 12px;
}

.login-header {
  text-align: center;
}

.login-header h2 {
  margin: 0 0 4px;
  font-size: 22px;
  color: #303133;
}

.login-subtitle {
  margin: 0;
  font-size: 13px;
  color: #909399;
}

.login-domain {
  margin: 8px 0 0;
  font-size: 12px;
  color: #667eea;
  font-weight: 500;
}

.register-link {
  margin-bottom: 0;
  text-align: center;
}

.register-btn {
  font-size: 13px;
  color: #667eea;
  text-decoration: none;
}

.register-btn:hover {
  text-decoration: underline;
}
</style>
