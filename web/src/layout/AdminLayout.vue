<template>
  <el-container class="layout" :class="{ 'is-collapsed': sidebarCollapsed }">
    <el-aside
      :width="sidebarCollapsed ? '76px' : '236px'"
      class="aside"
      :class="{ 'is-collapsed': sidebarCollapsed }"
    >
      <div class="logo" @click="toggleSidebar">
        <span class="logo-mark">g</span>
        <div v-show="!sidebarCollapsed" class="logo-copy">
          <strong>gin-web</strong>
          <span>内容管理平台</span>
        </div>
      </div>

      <el-menu
        class="sidebar-menu"
        router
        :collapse="sidebarCollapsed"
        :collapse-transition="false"
        :default-active="activeMenu"
      >
        <el-menu-item index="/">
          <el-icon><Odometer /></el-icon>
          <template #title>仪表盘</template>
        </el-menu-item>

        <el-sub-menu index="content">
          <template #title>
            <el-icon><Reading /></el-icon>
            <span>内容管理</span>
          </template>
          <el-menu-item index="/article">
            <el-icon><Document /></el-icon>
            <template #title>文章管理</template>
          </el-menu-item>
          <el-menu-item index="/category">
            <el-icon><Files /></el-icon>
            <template #title>分类管理</template>
          </el-menu-item>
          <el-menu-item index="/tag">
            <el-icon><PriceTag /></el-icon>
            <template #title>标签管理</template>
          </el-menu-item>
          <el-menu-item index="/comment">
            <el-icon><ChatDotRound /></el-icon>
            <template #title>
              <span class="menu-title-with-badge">评论管理 <em v-if="pendingComments">{{ pendingComments }}</em></span>
            </template>
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="site">
          <template #title>
            <el-icon><Tools /></el-icon>
            <span>站点设置</span>
          </template>
          <el-menu-item index="/link">
            <el-icon><Link /></el-icon>
            <template #title>友链管理</template>
          </el-menu-item>
          <el-menu-item index="/option">
            <el-icon><Setting /></el-icon>
            <template #title>站点配置</template>
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="ops">
          <template #title>
            <el-icon><Operation /></el-icon>
            <span>运维管理</span>
          </template>
          <el-menu-item index="/server">
            <el-icon><Monitor /></el-icon>
            <template #title>服务器管理</template>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>

      <div class="aside-footer">
        <div class="aside-status"><span class="status-dot" />系统运行正常</div>
      </div>
    </el-aside>

    <el-container class="body" direction="vertical">
      <el-header class="header">
        <div class="header-left">
          <el-button class="icon-button sidebar-trigger" text @click="toggleSidebar" aria-label="折叠侧栏">
            <el-icon :size="18"><Expand v-if="sidebarCollapsed" /><Fold v-else /></el-icon>
          </el-button>
          <div class="crumb">
            <span class="crumb-home">管理后台</span>
            <el-icon class="crumb-arrow"><ArrowRight /></el-icon>
            <span class="crumb-current">{{ pageTitle }}</span>
          </div>
        </div>
        <div class="header-actions">
          <el-tooltip content="搜索" placement="bottom">
            <el-button class="icon-button" text aria-label="搜索" @click="focusSearch">
              <el-icon :size="18"><Search /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip content="全屏" placement="bottom">
            <el-button class="icon-button" text aria-label="全屏" @click="toggleFullscreen">
              <el-icon :size="18"><FullScreen /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip content="通知" placement="bottom">
            <el-button class="icon-button notification-button" text aria-label="通知" @click="showNotifications">
              <el-icon :size="18"><Bell /></el-icon>
              <span v-if="pendingComments" class="notification-dot" />
            </el-button>
          </el-tooltip>
          <span class="header-divider" />
          <el-dropdown @command="onCommand">
            <span class="user">
              <el-avatar :size="32" :src="user.userInfo?.avatar">
                {{ user.userInfo?.nickname?.charAt(0) || 'A' }}
              </el-avatar>
              <span class="user-copy">
                <strong>{{ user.userInfo?.nickname || user.userInfo?.username || 'admin' }}</strong>
                <small>管理员</small>
              </span>
              <el-icon class="caret"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="password">修改密码</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <div class="tabbar">
        <div class="tab-scroll">
          <div
            v-for="tab in tabs"
            :key="tab.path"
            class="tab"
            :class="{ 'is-active': tab.path === route.fullPath }"
            @click="router.push(tab.path)"
          >
            <span v-if="tab.path === route.fullPath" class="tab-dot" />
            <span class="tab-text">{{ tab.title }}</span>
            <el-icon v-if="!tab.affix" class="tab-close" @click.stop="closeTab(tab)"><Close /></el-icon>
          </div>
        </div>
        <el-dropdown class="tab-actions" @command="onTabCommand">
          <span class="tab-more" aria-label="标签页操作"><el-icon><MoreFilled /></el-icon></span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="others">关闭其他</el-dropdown-item>
              <el-dropdown-item command="all">关闭全部</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>

      <el-main class="main"><router-view /></el-main>
    </el-container>

    <el-dialog v-model="pwdVisible" title="修改密码" width="420px">
      <el-form label-width="80px">
        <el-form-item label="旧密码"><el-input v-model="pwdForm.old_password" type="password" show-password /></el-form-item>
        <el-form-item label="新密码"><el-input v-model="pwdForm.new_password" type="password" show-password placeholder="6-32 位" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pwdVisible = false">取消</el-button>
        <el-button type="primary" :loading="pwdSaving" @click="submitPassword">确定</el-button>
      </template>
    </el-dialog>
  </el-container>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Odometer, Document, Files, PriceTag, ChatDotRound, Link, Setting, Monitor,
  Reading, Tools, Operation,
  ArrowDown, ArrowRight, Close, MoreFilled, Expand, Fold, Search, Bell, FullScreen
} from '@element-plus/icons-vue'
import { useUserStore } from '../stores/user'
import { changePassword } from '../api/auth'

const route = useRoute()
const router = useRouter()
const user = useUserStore()
const sidebarCollapsed = ref(localStorage.getItem('sidebar-collapsed') === '1')
const pendingComments = ref(0)
const pageTitle = computed(() => route.meta.title || '仪表盘')
// 新建/编辑文章页不属于独立菜单项，统一高亮「文章管理」，el-menu 会据此自动展开父级子菜单
const activeMenu = computed(() => (route.path.startsWith('/article') ? '/article' : route.path))

const HOME_TAB = { path: '/', title: '仪表盘', affix: true }
const tabs = ref([{ ...HOME_TAB }])

watch(() => route.fullPath, () => {
  if (!route.meta?.title) return
  if (!tabs.value.some((tab) => tab.path === route.fullPath)) tabs.value.push({ path: route.fullPath, title: route.meta.title })
}, { immediate: true })

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem('sidebar-collapsed', sidebarCollapsed.value ? '1' : '0')
}

function closeTab(tab) {
  const index = tabs.value.findIndex((item) => item.path === tab.path)
  if (index === -1) return
  tabs.value.splice(index, 1)
  if (tab.path === route.fullPath) router.push((tabs.value[index - 1] || tabs.value[0] || HOME_TAB).path)
}

function onTabCommand(command) {
  if (command === 'others') {
    tabs.value = tabs.value.filter((tab) => tab.affix || tab.path === route.fullPath)
    return
  }
  tabs.value = [{ ...HOME_TAB }]
  router.push('/')
}

function focusSearch() { ElMessage.info('可在当前页面筛选内容') }

async function toggleFullscreen() {
  try {
    if (!document.fullscreenElement) await document.documentElement.requestFullscreen()
    else await document.exitFullscreen()
  } catch (_) {
    ElMessage.info('当前浏览器不支持全屏')
  }
}

function showNotifications() {
  if (pendingComments.value) router.push('/comment')
  else ElMessage.success('暂无新的通知')
}

const pwdVisible = ref(false)
const pwdSaving = ref(false)
const pwdForm = reactive({ old_password: '', new_password: '' })

onMounted(() => {
  if (!user.userInfo) user.fetchUserInfo().catch(() => {})
})

async function onCommand(command) {
  if (command === 'password') {
    pwdForm.old_password = ''
    pwdForm.new_password = ''
    pwdVisible.value = true
    return
  }
  if (command === 'logout') {
    await ElMessageBox.confirm('确定退出登录？', '提示', { type: 'warning' })
    await user.logout()
    router.push('/login')
  }
}

async function submitPassword() {
  if (!pwdForm.old_password || !pwdForm.new_password) return ElMessage.warning('请填写完整')
  if (pwdForm.new_password.length < 6 || pwdForm.new_password.length > 32) return ElMessage.warning('新密码长度需 6-32 位')
  pwdSaving.value = true
  try {
    await changePassword(pwdForm)
    ElMessage.success('密码已修改，请重新登录')
    pwdVisible.value = false
    await user.logout()
    router.push('/login')
  } finally {
    pwdSaving.value = false
  }
}
</script>

<style scoped>
.layout { height: 100%; }
.aside {
  flex: none; position: relative; z-index: 10; display: flex; flex-direction: column; overflow: hidden;
  background: #fff; border-right: 1px solid var(--art-card-border); transition: width .22s ease;
}
.logo { height: var(--gw-header-height); flex: none; display: flex; align-items: center; gap: 11px; padding: 0 20px; border-bottom: 1px solid var(--art-card-border); cursor: pointer; overflow: hidden; }
.logo-mark { width: 30px; height: 30px; flex: none; display: inline-flex; align-items: center; justify-content: center; border-radius: 9px; background: var(--theme-color); color: #fff; font-size: 17px; font-weight: 700; box-shadow: 0 5px 12px rgba(93, 135, 255, .24); }
.logo-copy { display: flex; flex-direction: column; line-height: 1.2; white-space: nowrap; }
.logo-copy strong { color: var(--art-gray-900); font-size: 15px; font-weight: 650; }
.logo-copy span { margin-top: 3px; color: var(--art-gray-500); font-size: 11px; }
.sidebar-menu { flex: 1; overflow-y: auto; padding: 14px 10px; border-right: 0; background: transparent; --el-menu-bg-color: transparent; --el-menu-text-color: var(--art-gray-700); --el-menu-hover-bg-color: var(--art-hover-color); --el-menu-active-color: var(--theme-color); }
.sidebar-menu :deep(.el-menu-item) { position: relative; display: flex; align-items: center; height: 44px; margin-bottom: 3px; padding: 0 12px !important; border-radius: 8px; color: var(--art-gray-700); font-size: 14px; transition: background .15s ease, color .15s ease; }
.sidebar-menu :deep(.el-menu-item .el-icon) { width: 19px; margin-right: 11px; font-size: 17px; color: var(--art-gray-500); }
.sidebar-menu :deep(.el-menu-item:hover) { color: var(--art-gray-900); background: var(--art-hover-color); }
.sidebar-menu :deep(.el-menu-item:hover .el-icon) { color: var(--art-gray-800); }
.sidebar-menu :deep(.el-menu-item.is-active) { color: var(--theme-color); background: var(--el-color-primary-light-9); font-weight: 500; }
.sidebar-menu :deep(.el-menu-item.is-active::before) { position: absolute; left: 0; width: 3px; height: 20px; border-radius: 0 3px 3px 0; background: var(--theme-color); content: ''; }
.sidebar-menu :deep(.el-menu-item.is-active .el-icon) { color: var(--theme-color); }
/* 可折叠分组标题，与一级菜单项同高同圆角 */
.sidebar-menu :deep(.el-sub-menu__title) { height: 44px; margin-bottom: 3px; padding: 0 12px !important; border-radius: 8px; color: var(--art-gray-700); font-size: 14px; transition: background .15s ease, color .15s ease; }
.sidebar-menu :deep(.el-sub-menu__title:hover) { color: var(--art-gray-900); background: var(--art-hover-color); }
.sidebar-menu :deep(.el-sub-menu__title .el-icon) { width: 19px; margin-right: 11px; font-size: 17px; color: var(--art-gray-500); }
.sidebar-menu :deep(.el-sub-menu__title:hover .el-icon) { color: var(--art-gray-800); }
.sidebar-menu :deep(.el-sub-menu__icon-arrow) { right: 12px; color: var(--art-gray-400); font-size: 12px; }
/* 内联展开的子项容器去掉默认底色，子项整体缩进 */
.sidebar-menu :deep(.el-menu--inline) { background: transparent; }
.sidebar-menu :deep(.el-menu--inline .el-menu-item) { padding-left: 44px !important; }
.sidebar-menu :deep(.el-menu--collapse) { padding: 14px 10px; }
.sidebar-menu :deep(.el-menu--collapse .el-menu-item),
.sidebar-menu :deep(.el-menu--collapse .el-sub-menu__title) { justify-content: center; padding: 0 !important; }
.sidebar-menu :deep(.el-menu--collapse .el-menu-item .el-icon),
.sidebar-menu :deep(.el-menu--collapse .el-sub-menu__title .el-icon) { margin: 0; }
.menu-title-with-badge { display: flex; align-items: center; gap: 7px; }
.menu-title-with-badge em { min-width: 18px; height: 18px; padding: 0 4px; border-radius: 9px; background: #fff0f0; color: #f56c6c; font-size: 11px; font-style: normal; line-height: 18px; text-align: center; }
.aside-footer { padding: 0 18px 18px; white-space: nowrap; }
.aside-status { display: flex; align-items: center; gap: 7px; color: var(--art-gray-500); font-size: 11px; }
.status-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--el-color-success); }
.is-collapsed .aside-footer { visibility: hidden; }
.body { min-width: 0; }
.header { height: var(--gw-header-height); flex: none; display: flex; align-items: center; justify-content: space-between; padding: 0 24px; background: #fff; border-bottom: 1px solid var(--art-card-border); }
.header-left, .header-actions, .crumb, .user { display: flex; align-items: center; }
.header-left { min-width: 0; gap: 12px; }
.header-actions { gap: 4px; }
.icon-button { width: 34px; height: 34px; padding: 0; border-radius: 8px; color: var(--art-gray-600); }
.icon-button:hover { color: var(--art-gray-900); background: var(--art-hover-color); }
.sidebar-trigger { color: var(--art-gray-700); }
.crumb { gap: 8px; min-width: 0; font-size: 13px; }
.crumb-home { color: var(--art-gray-500); }
.crumb-arrow { color: var(--art-gray-400); font-size: 13px; }
.crumb-current { overflow: hidden; color: var(--art-gray-900); font-weight: 550; text-overflow: ellipsis; white-space: nowrap; }
.header-divider { width: 1px; height: 20px; margin: 0 9px; background: var(--art-card-border); }
.notification-button { position: relative; }
.notification-dot { position: absolute; top: 7px; right: 7px; width: 5px; height: 5px; border: 1.5px solid #fff; border-radius: 50%; background: #f56c6c; }
.user { gap: 9px; padding: 3px 8px 3px 3px; border-radius: 9px; cursor: pointer; }
.user:hover { background: var(--art-hover-color); }
.user-copy { display: flex; flex-direction: column; line-height: 1.25; }
.user-copy strong { color: var(--art-gray-800); font-size: 13px; font-weight: 500; }
.user-copy small { margin-top: 2px; color: var(--art-gray-500); font-size: 10px; }
.caret { margin-left: 2px; color: var(--art-gray-500); font-size: 12px; }
.tabbar { flex: none; height: var(--gw-tabbar-height); display: flex; align-items: center; gap: 8px; padding: 0 16px; background: #fff; border-bottom: 1px solid var(--art-card-border); }
.tab-scroll { flex: 1; min-width: 0; display: flex; align-items: center; gap: 5px; overflow-x: auto; scrollbar-width: none; }
.tab-scroll::-webkit-scrollbar { display: none; }
.tab { flex: none; display: inline-flex; align-items: center; gap: 6px; height: 28px; padding: 0 10px; border-radius: 6px; color: var(--art-gray-600); font-size: 12px; cursor: pointer; white-space: nowrap; transition: background .15s ease, color .15s ease; }
.tab:hover { color: var(--art-gray-900); background: var(--art-hover-color); }
.tab.is-active { color: var(--theme-color); background: var(--el-color-primary-light-9); }
.tab-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--theme-color); }
.tab-close { padding: 2px; border-radius: 4px; font-size: 12px; }
.tab-close:hover { background: rgba(0,0,0,.08); }
.tab-more { display: inline-flex; align-items: center; justify-content: center; width: 28px; height: 28px; border-radius: 6px; color: var(--art-gray-600); cursor: pointer; }
.tab-more:hover { background: var(--art-hover-color); color: var(--art-gray-900); }
.main { overflow: auto; background: var(--gw-page-bg); padding: 20px 24px 28px; }
@media (max-width: 900px) {
  .aside { position: fixed; top: 0; bottom: 0; left: 0; height: 100%; box-shadow: 5px 0 22px rgba(31,35,41,.1); }
  .aside.is-collapsed { transform: translateX(-100%); }
  .header { padding: 0 16px; }
  .main { padding: 16px; }
  .user-copy, .header-divider { display: none; }
}
@media (max-width: 560px) {
  .crumb-home, .crumb-arrow, .header-actions .icon-button:nth-child(-n+2) { display: none; }
  .header { padding: 0 12px; }
  .tabbar { padding: 0 10px; }
}
</style>
