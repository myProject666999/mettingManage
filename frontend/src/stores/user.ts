import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<any>(JSON.parse(localStorage.getItem('userInfo') || 'null'))
  const role = ref<string>(localStorage.getItem('role') || '')

  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const isUser = computed(() => role.value === 'user')
  const isAdmin = computed(() => role.value === 'admin' || role.value === 'super_admin')
  const isSuperAdmin = computed(() => role.value === 'super_admin')

  // 方法
  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUserInfo = (info: any) => {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  const setRole = (newRole: string) => {
    role.value = newRole
    localStorage.setItem('role', newRole)
  }

  const login = (newToken: string, info: any, newRole: string) => {
    setToken(newToken)
    setUserInfo(info)
    setRole(newRole)
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    role.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    localStorage.removeItem('role')
  }

  return {
    token,
    userInfo,
    role,
    isLoggedIn,
    isUser,
    isAdmin,
    isSuperAdmin,
    setToken,
    setUserInfo,
    setRole,
    login,
    logout
  }
})
