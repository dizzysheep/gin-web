import axios from 'axios'
import { ElMessage } from 'element-plus'

// 后端 token 相关错误码，见 internal/errcode/errcode.go
const TOKEN_CODES = [210001, 210002, 210003]

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 15000
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

function clearAuth() {
  localStorage.removeItem('token')
  localStorage.removeItem('userInfo')
  if (!window.location.hash.includes('/login')) {
    window.location.href = '/login'
  }
}

request.interceptors.response.use(
  (res) => {
    const { code, msg, data } = res.data
    if (code === 0) {
      return data
    }
    if (TOKEN_CODES.includes(code)) {
      ElMessage.error('登录已失效，请重新登录')
      clearAuth()
    } else {
      ElMessage.error(msg || '请求失败')
    }
    return Promise.reject(new Error(msg || `code: ${code}`))
  },
  (err) => {
    // 参数校验失败等场景后端返回 400 + {code,msg}
    const body = err.response?.data
    if (body && typeof body.code !== 'undefined') {
      if (TOKEN_CODES.includes(body.code)) {
        ElMessage.error('登录已失效，请重新登录')
        clearAuth()
      } else {
        ElMessage.error(body.msg || '请求失败')
      }
    } else {
      ElMessage.error(err.message || '网络异常')
    }
    return Promise.reject(err)
  }
)

export default request
