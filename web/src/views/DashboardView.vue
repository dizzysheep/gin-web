<template>
  <div class="page dashboard-page">
    <section class="welcome-card">
      <div class="welcome-copy">
        <span class="eyebrow">{{ greeting }}</span>
        <h1>{{ user.userInfo?.nickname || user.userInfo?.username || '管理员' }}，欢迎回来</h1>
        <p>今天也来保持内容的新鲜感吧，站点运行平稳。</p>
      </div>
      <div class="welcome-meta">
        <div class="welcome-date">{{ today }}</div>
        <div class="welcome-role"><span class="online-dot" />管理端在线</div>
      </div>
    </section>

    <div class="section-heading">
      <div>
        <h2>数据概览</h2>
        <span>站点内容的实时统计</span>
      </div>
      <el-button text class="refresh-button" :loading="loading" @click="loadStats">
        <el-icon><Refresh /></el-icon>刷新数据
      </el-button>
    </div>

    <el-row :gutter="14" class="stat-row">
      <el-col v-for="card in cards" :key="card.label" :xs="12" :sm="12" :md="6">
        <div class="stat-card card">
          <div class="stat-topline">
            <div class="stat-icon" :class="`tone-${card.tone}`"><el-icon><component :is="card.icon" /></el-icon></div>
            <span class="stat-trend" :class="card.trendTone">{{ card.trend }}</span>
          </div>
          <div class="stat-num">{{ card.value }}</div>
          <div class="stat-label">{{ card.label }}</div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="14" class="content-row">
      <el-col :xs="24" :lg="15">
        <div class="card activity-card">
          <div class="card-heading">
            <div><h3>内容趋势</h3><span>近 7 天文章发布情况</span></div>
            <span class="legend"><i />发布文章</span>
          </div>
          <div class="chart-wrap">
            <div class="chart-y-axis"><span v-for="value in chartAxis" :key="value">{{ value }}</span></div>
            <div class="bar-chart">
              <div class="chart-gridline grid-1" /><div class="chart-gridline grid-2" /><div class="chart-gridline grid-3" /><div class="chart-gridline grid-4" />
              <div v-for="point in weekly" :key="point.date" class="bar-column">
                <div class="bar-value">{{ point.total }}</div>
                <div class="bar" :style="{ height: `${point.total ? Math.max(point.total / chartMax * 80, 5) : 0}%` }" />
                <span>{{ point.label }}</span>
              </div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col :xs="24" :lg="9">
        <div class="card quick-card">
          <div class="card-heading"><div><h3>快捷操作</h3><span>常用功能入口</span></div></div>
          <div class="quick-grid">
            <button v-for="item in quickActions" :key="item.label" class="quick-item" @click="router.push(item.path)">
              <span class="quick-icon" :class="`tone-${item.tone}`"><el-icon><component :is="item.icon" /></el-icon></span>
              <span>{{ item.label }}</span>
              <el-icon class="quick-arrow"><ArrowRight /></el-icon>
            </button>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="14" class="content-row lower-row">
      <el-col :xs="24" :lg="15">
        <div class="card info-card">
          <div class="card-heading"><div><h3>最近动态</h3><span>账户与站点的最新状态</span></div><el-button text type="primary" @click="router.push('/article')">查看文章</el-button></div>
          <div class="activity-list">
            <div v-for="item in activities" :key="item.title" class="activity-item">
              <span class="activity-icon" :class="`tone-${item.tone}`"><el-icon><component :is="item.icon" /></el-icon></span>
              <div class="activity-copy"><strong>{{ item.title }}</strong><span>{{ item.desc }}</span></div>
              <time>{{ item.time }}</time>
            </div>
          </div>
        </div>
      </el-col>
      <el-col :xs="24" :lg="9">
        <div class="card account-card">
          <div class="card-heading"><div><h3>账户信息</h3><span>当前登录账号</span></div><el-icon class="heading-icon"><UserFilled /></el-icon></div>
          <div class="account-profile">
            <el-avatar :size="48" :src="user.userInfo?.avatar">{{ user.userInfo?.nickname?.charAt(0) || 'A' }}</el-avatar>
            <div><strong>{{ user.userInfo?.nickname || user.userInfo?.username || 'admin' }}</strong><span>{{ user.userInfo?.email || '未设置邮箱' }}</span></div>
          </div>
          <div class="account-meta"><span>用户名</span><strong>{{ user.userInfo?.username || '-' }}</strong></div>
          <div class="account-meta"><span>上次登录</span><strong>{{ formatTime(user.userInfo?.last_login) }}</strong></div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { formatTime } from '../utils/format'
import { getDashboardOverview, getDashboardTrend } from '../api/dashboard'
import {
  Refresh, ArrowRight, Document, Files, PriceTag, ChatDotRound,
  EditPen, Link, UserFilled, DataAnalysis
} from '@element-plus/icons-vue'

const router = useRouter()
const user = useUserStore()
const loading = ref(false)
const today = new Intl.DateTimeFormat('zh-CN', { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' }).format(new Date())
const hour = new Date().getHours()
const greeting = computed(() => (hour < 6 ? '夜深了' : hour < 12 ? '早上好' : hour < 18 ? '下午好' : '晚上好'))

const cards = ref([
  { label: '文章总数', value: '-', trend: '持续积累', trendTone: 'neutral', tone: 'blue', icon: Document },
  { label: '分类总数', value: '-', trend: '内容结构', trendTone: 'neutral', tone: 'violet', icon: Files },
  { label: '标签总数', value: '-', trend: '便于发现', trendTone: 'neutral', tone: 'amber', icon: PriceTag },
  { label: '待审评论', value: '-', trend: '需要处理', trendTone: 'warning', tone: 'green', icon: ChatDotRound }
])
const weekly = ref([
  { date: 'fallback-1', label: '周一', total: 2 }, { date: 'fallback-2', label: '周二', total: 5 },
  { date: 'fallback-3', label: '周三', total: 3 }, { date: 'fallback-4', label: '周四', total: 7 },
  { date: 'fallback-5', label: '周五', total: 4 }, { date: 'fallback-6', label: '周六', total: 6 },
  { date: 'fallback-7', label: '今天', total: 5 }
])
const quickActions = [
  { label: '写一篇文章', path: '/article/add', tone: 'blue', icon: EditPen },
  { label: '管理分类', path: '/category', tone: 'violet', icon: Files },
  { label: '查看评论', path: '/comment', tone: 'amber', icon: ChatDotRound },
  { label: '站点设置', path: '/option', tone: 'green', icon: DataAnalysis }
]
const activities = [
  { title: '内容中心已准备就绪', desc: '可以开始管理你的文章和分类', time: '现在', tone: 'blue', icon: Document },
  { title: '站点配置支持动态保存', desc: '修改后的配置会立即应用到前台', time: '提示', tone: 'green', icon: Link },
  { title: '保持内容持续更新', desc: '规律发布有助于读者发现新内容', time: '建议', tone: 'amber', icon: DataAnalysis }
]

const chartMax = computed(() => {
  const max = Math.max(...weekly.value.map((item) => item.total || 0), 0)
  return Math.max(8, Math.ceil(max / 2) * 2)
})
const chartAxis = computed(() => [
  chartMax.value,
  Math.ceil(chartMax.value * .75),
  Math.ceil(chartMax.value * .5),
  Math.ceil(chartMax.value * .25),
  0
])

async function loadStats() {
  loading.value = true
  try {
    const [overview, trend] = await Promise.allSettled([
      getDashboardOverview(),
      getDashboardTrend(7)
    ])
    if (overview.status === 'fulfilled') {
      const data = overview.value || {}
      cards.value[0].value = data.article_count ?? 0
      cards.value[1].value = data.category_count ?? 0
      cards.value[2].value = data.tag_count ?? 0
      cards.value[3].value = data.pending_comment_count ?? 0
    }
    if (trend.status === 'fulfilled' && Array.isArray(trend.value?.items)) {
      weekly.value = trend.value.items
    }
  } finally { loading.value = false }
}

onMounted(() => {
  if (!user.userInfo) user.fetchUserInfo().catch(() => {})
  loadStats()
})
</script>

<style scoped>
.dashboard-page { max-width: 1480px; margin: 0 auto; }
.welcome-card { display: flex; align-items: center; justify-content: space-between; min-height: 156px; padding: 28px 34px; border-radius: 12px; background: #5d87ff; color: #fff; box-shadow: 0 10px 24px rgba(93,135,255,.15); }
.eyebrow { display: block; margin-bottom: 8px; color: rgba(255,255,255,.72); font-size: 12px; letter-spacing: .06em; }
.welcome-copy h1 { margin: 0; font-size: 24px; font-weight: 600; letter-spacing: -.02em; }
.welcome-copy p { margin: 9px 0 0; color: rgba(255,255,255,.78); font-size: 13px; }
.welcome-meta { min-width: 148px; padding-left: 25px; border-left: 1px solid rgba(255,255,255,.25); text-align: right; }
.welcome-date { font-size: 13px; font-weight: 500; }
.welcome-role { display: flex; align-items: center; justify-content: flex-end; gap: 6px; margin-top: 13px; color: rgba(255,255,255,.75); font-size: 12px; }
.online-dot { width: 6px; height: 6px; border-radius: 50%; background: #a7f3d0; }
.section-heading { display: flex; align-items: flex-end; justify-content: space-between; margin: 27px 0 13px; }
.section-heading h2 { margin: 0; color: var(--art-gray-900); font-size: 17px; font-weight: 600; }
.section-heading span { display: block; margin-top: 4px; color: var(--art-gray-500); font-size: 12px; }
.refresh-button { padding: 5px 8px; color: var(--art-gray-600); font-size: 12px; }
.refresh-button .el-icon { margin-right: 5px; }
.stat-row { row-gap: 14px; }
.stat-card { min-height: 136px; padding: 18px 20px; border: 1px solid var(--art-card-border); transition: transform .18s ease, box-shadow .18s ease; }
.stat-card:hover { transform: translateY(-2px); box-shadow: 0 9px 20px rgba(31,35,41,.06); }
.stat-topline { display: flex; align-items: center; justify-content: space-between; }
.stat-icon, .quick-icon, .activity-icon { display: inline-flex; align-items: center; justify-content: center; border-radius: 9px; }
.stat-icon { width: 34px; height: 34px; font-size: 18px; }
.stat-trend { font-size: 11px; }
.stat-trend.neutral { color: var(--art-gray-500); }
.stat-trend.warning { color: #e6a23c; }
.stat-num { margin-top: 13px; color: var(--art-gray-900); font-size: 27px; font-weight: 650; font-variant-numeric: tabular-nums; line-height: 1; }
.stat-label { margin-top: 8px; color: var(--art-gray-600); font-size: 12px; }
.tone-blue { background: #edf3ff; color: #5d87ff; }
.tone-violet { background: #f1edff; color: #8b72e9; }
.tone-amber { background: #fff5e6; color: #e6a23c; }
.tone-green { background: #e9faf5; color: #13b994; }
.content-row { margin-top: 14px; row-gap: 14px; }
.lower-row { margin-top: 14px; }
.activity-card, .quick-card, .info-card, .account-card { min-height: 100%; padding: 20px; }
.card-heading { display: flex; align-items: flex-start; justify-content: space-between; }
.card-heading h3 { margin: 0; color: var(--art-gray-900); font-size: 15px; font-weight: 600; }
.card-heading span { display: block; margin-top: 4px; color: var(--art-gray-500); font-size: 12px; }
.heading-icon { color: var(--art-gray-400); font-size: 19px; }
.legend { display: flex !important; align-items: center; gap: 5px; color: var(--art-gray-500) !important; font-size: 11px !important; }
.legend i { width: 7px; height: 7px; border-radius: 50%; background: var(--theme-color); }
.chart-wrap { display: flex; height: 215px; margin-top: 19px; }
.chart-y-axis { display: flex; flex-direction: column; justify-content: space-between; width: 25px; padding: 2px 0 24px; color: var(--art-gray-500); font-size: 10px; }
.bar-chart { position: relative; display: flex; flex: 1; align-items: flex-end; justify-content: space-around; padding: 0 2% 0; }
.chart-gridline { position: absolute; left: 0; right: 0; border-top: 1px dashed var(--art-card-border); }
.grid-1 { bottom: 80%; }.grid-2 { bottom: 60%; }.grid-3 { bottom: 40%; }.grid-4 { bottom: 20%; }
.bar-column { position: relative; z-index: 1; display: flex; flex-direction: column; align-items: center; justify-content: flex-end; width: 10%; height: 100%; }
.bar { width: min(23px, 55%); min-height: 7px; border-radius: 5px 5px 2px 2px; background: #8eabff; transition: height .35s ease; }
.bar-column:last-child .bar { background: var(--theme-color); }
.bar-column span { margin-top: 10px; color: var(--art-gray-500); font-size: 11px; white-space: nowrap; }
.bar-value { height: 18px; color: var(--art-gray-600); font-size: 10px; }
.quick-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 9px; margin-top: 19px; }
.quick-item { display: flex; align-items: center; min-width: 0; height: 67px; padding: 0 11px; border: 1px solid var(--art-card-border); border-radius: 9px; background: #fff; color: var(--art-gray-800); font: inherit; font-size: 12px; text-align: left; cursor: pointer; transition: border-color .15s ease, background .15s ease; }
.quick-item:hover { border-color: var(--el-color-primary-light-5); background: var(--el-color-primary-light-9); }
.quick-icon { width: 29px; height: 29px; flex: none; margin-right: 8px; font-size: 15px; }
.quick-item > span:nth-child(2) { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.quick-arrow { margin-left: auto; color: var(--art-gray-400); font-size: 13px; }
.activity-list { margin-top: 11px; }
.activity-item { display: flex; align-items: center; gap: 11px; padding: 12px 0; border-bottom: 1px solid var(--art-card-border); }
.activity-item:last-child { border-bottom: 0; }
.activity-icon { width: 31px; height: 31px; flex: none; font-size: 15px; }
.activity-copy { display: flex; flex: 1; min-width: 0; flex-direction: column; }
.activity-copy strong { overflow: hidden; color: var(--art-gray-800); font-size: 12px; font-weight: 500; text-overflow: ellipsis; white-space: nowrap; }
.activity-copy span { overflow: hidden; margin-top: 3px; color: var(--art-gray-500); font-size: 11px; text-overflow: ellipsis; white-space: nowrap; }
.activity-item time { flex: none; color: var(--art-gray-400); font-size: 11px; }
.account-profile { display: flex; align-items: center; gap: 11px; margin: 20px 0 18px; }
.account-profile strong, .account-profile span { display: block; }
.account-profile strong { color: var(--art-gray-800); font-size: 13px; font-weight: 500; }
.account-profile span { margin-top: 4px; color: var(--art-gray-500); font-size: 11px; }
.account-meta { display: flex; align-items: center; justify-content: space-between; padding: 11px 0; border-top: 1px solid var(--art-card-border); color: var(--art-gray-500); font-size: 12px; }
.account-meta strong { color: var(--art-gray-700); font-size: 12px; font-weight: 500; }
@media (max-width: 640px) {
  .welcome-card { align-items: flex-start; flex-direction: column; min-height: 0; padding: 24px; }
  .welcome-copy h1 { font-size: 21px; }
  .welcome-meta { width: 100%; margin-top: 20px; padding: 14px 0 0; border-top: 1px solid rgba(255,255,255,.25); border-left: 0; text-align: left; }
  .welcome-role { justify-content: flex-start; }
  .section-heading { margin-top: 22px; }
  .chart-wrap { height: 190px; }
}
</style>
