<template>
  <div class="login-wrap">
    <section class="login-brand">
      <div class="brand-top">
        <span class="brand-mark">g</span>
        <span class="brand-name">gin-web</span>
      </div>
      <div class="brand-content">
        <span class="brand-kicker">CONTENT WORKSPACE</span>
        <h1>让每一次创作<br />都值得被看见</h1>
        <p>一个简洁、专注的内容管理工作台，帮助你把想法稳定地整理成作品。</p>
        <div class="brand-points">
          <span><el-icon><CircleCheck /></el-icon>清晰的内容结构</span>
          <span><el-icon><CircleCheck /></el-icon>顺手的编辑体验</span>
          <span><el-icon><CircleCheck /></el-icon>可靠的站点配置</span>
        </div>
      </div>
      <div class="brand-footer">© gin-web admin workspace</div>
    </section>

    <section class="login-panel">
      <div class="login-card">
        <div class="mobile-brand"><span class="brand-mark">g</span><span>gin-web</span></div>
        <div class="login-heading">
          <h2>欢迎回来</h2>
          <p>登录管理后台，继续你的内容工作</p>
        </div>
        <el-form :model="form" @keyup.enter="submit">
          <el-form-item label="用户名">
            <el-input v-model="form.username" placeholder="请输入用户名" size="large">
              <template #prefix><el-icon><User /></el-icon></template>
            </el-input>
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password size="large">
              <template #prefix><el-icon><Lock /></el-icon></template>
            </el-input>
          </el-form-item>
          <div class="login-options"><el-checkbox>记住登录状态</el-checkbox><span>安全登录</span></div>
          <el-button class="login-btn" type="primary" size="large" :loading="loading" @click="submit">登录</el-button>
        </el-form>
      </div>
    </section>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, CircleCheck } from '@element-plus/icons-vue'
import { useUserStore } from '../stores/user'

const route = useRoute()
const router = useRouter()
const user = useUserStore()
const form = reactive({ username: '', password: '' })
const loading = ref(false)

async function submit() {
  if (!form.username || !form.password) return ElMessage.warning('请输入用户名和密码')
  loading.value = true
  try {
    await user.login({ ...form })
    user.fetchUserInfo().catch(() => {})
    ElMessage.success('登录成功')
    router.push(route.query.redirect || '/')
  } finally { loading.value = false }
}
</script>

<style scoped>
.login-wrap { display: flex; width: 100%; height: 100%; min-height: 560px; background: #fff; }
.login-brand { position: relative; display: flex; width: 42%; min-width: 440px; flex-direction: column; padding: 38px 58px; overflow: hidden; background: #5d87ff; color: #fff; }
.login-brand::after { position: absolute; right: -100px; bottom: -110px; width: 340px; height: 340px; border: 1px solid rgba(255,255,255,.19); border-radius: 50%; box-shadow: 0 0 0 34px rgba(255,255,255,.05), 0 0 0 69px rgba(255,255,255,.04), 0 0 0 104px rgba(255,255,255,.03); content: ''; }
.brand-top, .mobile-brand { display: flex; align-items: center; gap: 10px; }
.brand-mark { display: inline-flex; align-items: center; justify-content: center; width: 32px; height: 32px; border-radius: 9px; background: #fff; color: var(--theme-color); font-size: 18px; font-weight: 700; }
.brand-name { font-size: 16px; font-weight: 600; letter-spacing: .01em; }
.brand-content { position: relative; z-index: 1; max-width: 380px; margin: auto 0; }
.brand-kicker { color: rgba(255,255,255,.62); font-size: 11px; letter-spacing: .15em; }
.brand-content h1 { margin: 18px 0 16px; font-size: 35px; font-weight: 600; line-height: 1.3; letter-spacing: -.03em; }
.brand-content p { max-width: 330px; margin: 0; color: rgba(255,255,255,.77); font-size: 14px; line-height: 1.8; }
.brand-points { display: flex; flex-direction: column; gap: 13px; margin-top: 30px; color: rgba(255,255,255,.88); font-size: 13px; }
.brand-points span { display: flex; align-items: center; gap: 9px; }
.brand-points .el-icon { color: #b8f4dd; font-size: 16px; }
.brand-footer { position: relative; z-index: 1; color: rgba(255,255,255,.53); font-size: 11px; }
.login-panel { display: flex; flex: 1; align-items: center; justify-content: center; padding: 40px 28px; background: #f7f8fa; }
.login-card { width: min(100%, 380px); padding: 12px 16px; }
.mobile-brand { display: none; color: var(--art-gray-900); font-size: 16px; font-weight: 600; }
.mobile-brand .brand-mark { background: var(--theme-color); color: #fff; }
.login-heading { margin-bottom: 28px; }
.login-heading h2 { margin: 0; color: var(--art-gray-900); font-size: 26px; font-weight: 600; letter-spacing: -.02em; }
.login-heading p { margin: 8px 0 0; color: var(--art-gray-500); font-size: 13px; }
.login-card :deep(.el-form-item) { margin-bottom: 21px; }
.login-card :deep(.el-form-item__label) { height: auto; margin-bottom: 7px; color: var(--art-gray-700); font-size: 13px; line-height: 1.3; }
.login-card :deep(.el-input__wrapper) { background: #fff; }
.login-options { display: flex; align-items: center; justify-content: space-between; margin: 0 0 21px; color: var(--art-gray-500); font-size: 12px; }
.login-options span { color: var(--theme-color); }
.login-btn { width: 100%; height: 42px !important; border-radius: 8px; font-size: 14px; }
@media (max-width: 800px) {
  .login-brand { display: none; }
  .login-panel { padding: 24px; }
  .login-card { padding: 0; }
  .mobile-brand { display: flex; margin-bottom: 42px; }
}
@media (max-width: 420px) {
  .login-panel { padding: 20px; }
  .mobile-brand { margin-bottom: 30px; }
  .login-heading h2 { font-size: 23px; }
}
</style>
