import { defineStore } from 'pinia'
import type { UserInfo } from '@/types'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo:
      (JSON.parse(localStorage.getItem('userInfo') || 'null') as UserInfo | null) ?? null,
  }),
  getters: {
    isLogin: (state) => !!state.token,
  },
  actions: {
    setAuth(token: string) {
      this.token = token
      localStorage.setItem('token', token)
    },
    setUserInfo(userInfo: UserInfo) {
      this.userInfo = userInfo
      localStorage.setItem('userInfo', JSON.stringify(userInfo))
    },
    clearAuth() {
      this.token = ''
      this.userInfo = null
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
    },
  },
})
