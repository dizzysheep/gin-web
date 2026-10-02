import { defineStore } from 'pinia'
import { login as loginAPI, logout as logoutAPI, getUserInfo } from '../api/auth'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: JSON.parse(localStorage.getItem('userInfo') || 'null')
  }),
  actions: {
    async login(form) {
      const data = await loginAPI(form)
      this.token = data.jwt
      localStorage.setItem('token', data.jwt)
    },
    async fetchUserInfo() {
      this.userInfo = await getUserInfo()
      localStorage.setItem('userInfo', JSON.stringify(this.userInfo))
    },
    async logout() {
      try {
        await logoutAPI()
      } catch (e) {
        // 登出接口失败也继续清理本地状态
      }
      this.token = ''
      this.userInfo = null
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
    }
  }
})
