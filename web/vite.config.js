import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// 开发时代理到 Go 后端（本地默认 :8080，如与 config/app.toml 的 appAddr 不一致以实际端口为准）
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://127.0.0.1:8080',
      '/uploads': 'http://127.0.0.1:8080',
      '/rss.xml': 'http://127.0.0.1:8080',
      '/sitemap.xml': 'http://127.0.0.1:8080'
    }
  }
})
