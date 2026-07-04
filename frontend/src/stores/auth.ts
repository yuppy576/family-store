import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi } from '@/api/auth'
import router from '@/router'

export interface UserInfo {
  id: number
  email: string
  name: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<UserInfo | null>(() => {
    const raw = localStorage.getItem('user')
    if (raw) {
      try {
        return JSON.parse(raw)
      } catch {
        return null
      }
    }
    return null
  })

  const isLoggedIn = computed(() => !!token.value)

  async function login(email: string, password: string): Promise<void> {
    const res = await loginApi({ email, password })
    const body = res.data
    const tok = body.data.token
    token.value = tok
    user.value = { id: 0, email, name: email.split('@')[0] }
    localStorage.setItem('token', tok)
    localStorage.setItem('user', JSON.stringify(user.value))
    router.push('/')
  }

  function logout(): void {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    window.location.href = '/store/login'
  }

  return { token, user, isLoggedIn, login, logout }
})
