<template>
  <div class="page">
    <div class="page-head">
      <div>
        <h2 class="page-title">服务器管理</h2>
        <p class="page-desc">后端服务所在主机的系统信息与资源占用，每 {{ refreshSeconds }} 秒自动刷新。</p>
      </div>
      <div class="page-actions">
        <span v-if="updatedAt" class="sub">更新于 {{ updatedAt }}</span>
        <el-button :loading="loading" @click="load">
          <el-icon><Refresh /></el-icon>
          <span>刷新</span>
        </el-button>
      </div>
    </div>

    <div v-loading="loading && !info" class="card server-card">
      <template v-if="info">
        <div class="server-head">
          <div class="server-ident">
            <span class="server-mark"><el-icon :size="20"><Monitor /></el-icon></span>
            <div class="server-ident-copy">
              <div class="server-name">{{ info.hostname || '未知主机' }}</div>
              <div class="cell-slug">{{ info.ip || '未获取到 IP' }}</div>
            </div>
          </div>
          <span class="badge badge-success">运行中</span>
        </div>

        <div class="server-meta">
          <div v-for="item in metaList" :key="item.label" class="meta-item">
            <span class="meta-label">{{ item.label }}</span>
            <span class="meta-value" :title="item.value">{{ item.value }}</span>
          </div>
        </div>

        <div class="server-usage">
          <div v-for="item in usageList" :key="item.label" class="usage-item">
            <div class="usage-head">
              <span class="usage-label">{{ item.label }}</span>
              <span class="usage-detail">{{ item.detail }}</span>
            </div>
            <el-progress
              :percentage="item.percent"
              :status="item.status"
              :text-inside="true"
              :stroke-width="17"
            />
          </div>
        </div>
      </template>

      <el-empty v-else-if="!loading" description="未获取到服务器信息" />
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { Monitor, Refresh } from '@element-plus/icons-vue'
import { getServerInfo } from '../api/server'

const refreshSeconds = 5
const info = ref(null)
const loading = ref(false)
const updatedAt = ref('')
let timer = null

async function load() {
  loading.value = true
  try {
    info.value = await getServerInfo()
    updatedAt.value = new Date().toLocaleTimeString('zh-CN', { hour12: false })
  } catch (_) {
    // 请求层已统一提示，这里停止轮询避免错误信息刷屏
    stopPolling()
  } finally {
    loading.value = false
  }
}

function startPolling() {
  if (timer) return
  timer = window.setInterval(load, refreshSeconds * 1000)
}

function stopPolling() {
  if (!timer) return
  window.clearInterval(timer)
  timer = null
}

const metaList = computed(() => {
  const data = info.value
  if (!data) return []
  return [
    { label: '操作系统', value: data.platform || data.os || '-' },
    { label: '内核版本', value: data.kernel || '-' },
    { label: '系统架构', value: data.arch || '-' },
    { label: 'CPU 型号', value: data.cpu_model || '-' },
    { label: 'CPU 核心', value: data.cpu_cores ? `${data.cpu_cores} 核` : '-' },
    { label: '运行时长', value: formatUptime(data.uptime) },
    { label: '平均负载', value: data.load_avg ? String(data.load_avg) : '-' },
    { label: '更新时间', value: formatTime(data.collected_at) }
  ]
})

const usageList = computed(() => {
  const data = info.value || {}
  return [
    { label: 'CPU', percent: data.cpu_percent || 0, status: '', detail: data.cpu_cores ? `${data.cpu_cores} 核` : '-' },
    { label: '内存', percent: data.mem_percent || 0, status: 'success', detail: capacity(data.mem_used, data.mem_total) },
    { label: 'SWAP', percent: data.swap_percent || 0, status: 'warning', detail: data.swap_total ? capacity(data.swap_used, data.swap_total) : '未启用' },
    { label: '磁盘', percent: data.disk_percent || 0, status: 'success', detail: capacity(data.disk_used, data.disk_total) }
  ]
})

function capacity(used, total) {
  if (!total) return '-'
  return `${formatBytes(used)} / ${formatBytes(total)}`
}

function formatBytes(bytes) {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let value = bytes
  let index = 0
  while (value >= 1024 && index < units.length - 1) {
    value /= 1024
    index += 1
  }
  return `${index === 0 ? value : value.toFixed(1)} ${units[index]}`
}

function formatUptime(seconds) {
  if (!seconds) return '-'
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  if (days) return `${days} 天 ${hours} 小时`
  if (hours) return `${hours} 小时 ${minutes} 分`
  return `${minutes} 分`
}

function formatTime(timestamp) {
  if (!timestamp) return '-'
  return new Date(timestamp * 1000).toLocaleTimeString('zh-CN', { hour12: false })
}

onMounted(() => {
  load()
  startPolling()
})

onUnmounted(stopPolling)
</script>

<style scoped>
.server-card {
  padding: 22px 24px 26px;
}

.server-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding-bottom: 18px;
  border-bottom: 1px solid var(--art-card-border);
}

.server-ident {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.server-mark {
  flex: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  border-radius: 10px;
  background: var(--el-color-primary-light-9);
  color: var(--theme-color);
}

.server-ident-copy {
  min-width: 0;
}

.server-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.server-meta {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px 20px;
  padding: 18px 0 20px;
  border-bottom: 1px solid var(--art-card-border);
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: 5px;
  min-width: 0;
}

.meta-label {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
}

.meta-value {
  overflow: hidden;
  font-size: 13px;
  color: var(--el-text-color-primary);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.server-usage {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px 28px;
  padding-top: 20px;
}

.usage-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 7px;
}

.usage-label {
  font-size: 13px;
  color: var(--el-text-color-primary);
}

.usage-detail {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
  font-variant-numeric: tabular-nums;
}

@media (max-width: 1200px) {
  .server-meta {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 900px) {
  .server-meta {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .server-usage {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
