<template>
  <div class="page agent-page">
    <section class="agent-heading">
      <div>
        <div class="agent-kicker"><span class="agent-status" />博客内容助手</div>
        <h1>AI 助手</h1>
        <p>用自然语言查找已发布文章。需要时，助手会调用文章查询工具并展示执行结果。</p>
      </div>
      <el-button plain :icon="Refresh" :disabled="sending" @click="clearChat">新对话</el-button>
    </section>

    <section class="chat-card card">
      <div class="chat-toolbar">
        <div class="assistant-avatar"><el-icon><ChatLineRound /></el-icon></div>
        <div class="toolbar-copy"><strong>内容管理助手</strong><span>工具调用已连接</span></div>
        <el-tag v-if="tools.length" round effect="plain" type="success">已注册 {{ tools.length }} 个工具</el-tag>
        <el-tag v-else round effect="plain" type="info">文章只读查询</el-tag>
      </div>

      <div ref="threadEl" class="chat-thread" aria-live="polite">
        <div v-if="messages.length === 0" class="welcome-state">
          <div class="welcome-icon"><el-icon><ChatLineRound /></el-icon></div>
          <h2>你好，需要查找什么内容？</h2>
          <p>我可以帮你搜索已发布文章，并说明实际调用了哪些工具。</p>
          <div class="suggestions">
            <button v-for="prompt in suggestions" :key="prompt" :disabled="sending" @click="send(prompt)">
              <span>{{ prompt }}</span><el-icon><ArrowRight /></el-icon>
            </button>
          </div>
        </div>

        <article v-for="(message, index) in messages" :key="index" class="message-row" :class="'from-' + message.role">
          <div v-if="message.role === 'assistant'" class="message-avatar"><el-icon><ChatLineRound /></el-icon></div>
          <div class="message-content">
            <div class="message-author">{{ message.role === 'user' ? '你' : '内容管理助手' }}</div>
            <div v-if="message.text" class="message-bubble" :class="{ 'user-bubble': message.role === 'user' }">
              <p>{{ message.text }}</p>
              <div v-if="message.error" class="message-error">本次请求未完成，请检查 Agent 配置或稍后重试。</div>
            </div>
            <div v-if="message.toolCalls?.length" class="tool-list">
              <div v-for="(call, callIndex) in message.toolCalls" :key="callIndex" class="tool-card">
                <div class="tool-title">
                  <span class="tool-check"><el-icon><Check /></el-icon></span>
                  <div><strong>{{ toolLabel(call.name) }}</strong><small>工具已执行</small></div>
                  <el-tag size="small" effect="plain">{{ call.name }}</el-tag>
                </div>
                <el-collapse>
                  <el-collapse-item title="查看参数与工具结果">
                    <div class="tool-section"><span>调用参数</span><pre>{{ pretty(call.arguments) }}</pre></div>
                    <div class="tool-section"><span>返回结果</span><pre>{{ pretty(call.output) }}</pre></div>
                  </el-collapse-item>
                </el-collapse>
              </div>
            </div>
          </div>
          <div v-if="message.role === 'user'" class="user-avatar"><el-icon><UserFilled /></el-icon></div>
        </article>

        <div v-if="sending" class="message-row from-assistant">
          <div class="message-avatar"><el-icon><ChatLineRound /></el-icon></div>
          <div class="message-content">
            <div class="message-author">内容管理助手</div>
            <div class="message-bubble typing"><i /><i /><i /><span>正在查询并整理回答</span></div>
          </div>
        </div>
      </div>

      <form class="composer" @submit.prevent="send()">
        <el-input
          v-model="draft"
          type="textarea"
          :autosize="{ minRows: 1, maxRows: 4 }"
          maxlength="4000"
          resize="none"
          placeholder="例如：找一下最近写的 Go 相关文章"
          :disabled="sending"
          @keydown.enter.exact.prevent="send()"
        />
        <div class="composer-bottom">
          <span>Enter 发送 · Shift + Enter 换行</span>
          <el-button type="primary" :icon="Position" :loading="sending" :disabled="!draft.trim()" native-type="submit">发送</el-button>
        </div>
      </form>
    </section>
    <div class="agent-footnote"><el-icon><Lock /></el-icon>助手仅查询已发布文章；对话需要通过后台登录验证。</div>
  </div>
</template>

<script setup>
import { nextTick, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowRight, ChatLineRound, Check, Lock, Position, Refresh, UserFilled } from '@element-plus/icons-vue'
import { chatWithAgent, getAgentTools } from '../api/agent'

const messages = ref([])
const tools = ref([])
const draft = ref('')
const sending = ref(false)
const threadEl = ref(null)
const suggestions = [
  '列出最近发布的文章',
  '搜索标题或内容里包含 Go 的文章',
  '有哪些文章介绍了工具调用？'
]

onMounted(async () => {
  try {
    tools.value = await getAgentTools() || []
  } catch (_) {
    // Shared request interceptor displays authentication and network errors.
  }
})

function toolLabel(name) {
  return name === 'list_articles' ? '查询已发布文章' : name
}

function pretty(value) {
  if (typeof value === 'string') {
    try { return JSON.stringify(JSON.parse(value), null, 2) } catch (_) { return value }
  }
  return JSON.stringify(value ?? null, null, 2)
}

async function scrollToBottom() {
  await nextTick()
  if (threadEl.value) threadEl.value.scrollTop = threadEl.value.scrollHeight
}

async function send(text = draft.value) {
  const message = String(text || '').trim()
  if (!message || sending.value) return
  draft.value = ''
  messages.value.push({ role: 'user', text: message })
  const pendingToolCalls = []
  let assistantMessage = null
  const ensureAssistant = () => {
    if (!assistantMessage) {
      assistantMessage = { role: 'assistant', text: '', toolCalls: pendingToolCalls }
      messages.value.push(assistantMessage)
    }
    return assistantMessage
  }
  const answerText = (answer) => {
    const value = String(answer || '').trim()
    return value || '请求已完成，但模型没有返回文本回答。'
  }
  sending.value = true
  await scrollToBottom()
  try {
    const result = await chatWithAgent(message, (eventName, payload) => {
      if (eventName === 'tool_call') {
        pendingToolCalls.push(payload)
        if (assistantMessage) assistantMessage.toolCalls = pendingToolCalls
      }
      if (eventName === 'message') {
        const row = ensureAssistant()
        row.text = answerText(payload?.answer)
        row.toolCalls = payload?.tool_calls?.length ? payload.tool_calls : pendingToolCalls
      }
      scrollToBottom()
    })
    const row = ensureAssistant()
    row.text = answerText(result?.answer)
    row.toolCalls = result?.tool_calls?.length ? result.tool_calls : pendingToolCalls
  } catch (_) {
    const row = ensureAssistant()
    row.text = '连接 Agent 服务失败。请确认服务端已配置中转站地址、模型和 API Key。'
    row.toolCalls = pendingToolCalls
    row.error = true
    ElMessage.error('AI 助手请求失败')
  } finally {
    sending.value = false
    await scrollToBottom()
  }
}

function clearChat() {
  if (sending.value) return
  messages.value = []
}
</script>

<style scoped>
.agent-page { max-width: 1120px; margin: 0 auto; }
.agent-heading { display:flex; align-items:center; justify-content:space-between; gap:20px; margin:2px 0 18px; }
.agent-kicker { display:flex; align-items:center; gap:7px; color:#5d87ff; font-size:12px; font-weight:600; }
.agent-status { width:7px; height:7px; border-radius:50%; background:#13b994; box-shadow:0 0 0 3px #e6faf5; }
.agent-heading h1 { margin:5px 0 3px; color:#1d2129; font-size:23px; font-weight:650; }
.agent-heading p { margin:0; color:#7c7f85; font-size:13px; }
.chat-card { display:flex; flex-direction:column; height:min(690px, calc(100vh - 210px)); min-height:460px; overflow:hidden; }
.chat-toolbar { display:flex; align-items:center; gap:11px; padding:15px 20px; border-bottom:1px solid #f0f2f5; }
.assistant-avatar,.message-avatar,.user-avatar { display:flex; align-items:center; justify-content:center; flex:none; }
.assistant-avatar { width:36px; height:36px; border-radius:11px; background:#eef3ff; color:#5d87ff; font-size:18px; }
.toolbar-copy { display:flex; flex:1; flex-direction:column; gap:2px; }
.toolbar-copy strong { color:#29343d; font-size:13px; font-weight:600; }
.toolbar-copy span { color:#13a883; font-size:11px; }
.chat-thread { flex:1; overflow:auto; padding:23px clamp(16px, 5vw, 58px); scroll-behavior:smooth; }
.welcome-state { max-width:500px; margin:9vh auto 0; text-align:center; }
.welcome-icon { display:grid; width:48px; height:48px; place-items:center; margin:auto; border-radius:15px; background:#eef3ff; color:#5d87ff; font-size:23px; }
.welcome-state h2 { margin:15px 0 5px; color:#29343d; font-size:18px; font-weight:600; }
.welcome-state p { margin:0; color:#909399; font-size:12px; }
.suggestions { display:grid; grid-template-columns:1fr; gap:8px; margin-top:22px; text-align:left; }
.suggestions button { display:flex; align-items:center; justify-content:space-between; padding:11px 13px; border:1px solid #ebedf0; border-radius:9px; background:white; color:#5c5f66; font:inherit; font-size:12px; cursor:pointer; transition:.15s; }
.suggestions button:hover { border-color:#aec3ff; background:#f8faff; color:#4a6ccc; }
.message-row { display:flex; align-items:flex-start; gap:10px; margin:0 0 22px; }
.from-user { justify-content:flex-end; }
.message-avatar,.user-avatar { width:28px; height:28px; margin-top:19px; border-radius:9px; }
.message-avatar { background:#eef3ff; color:#5d87ff; }
.user-avatar { background:#f1edff; color:#8b72e9; }
.message-content { max-width:min(76%, 660px); min-width:0; }
.from-user .message-content { display:flex; align-items:flex-end; flex-direction:column; }
.message-author { margin:0 0 6px 2px; color:#909399; font-size:11px; }
.message-bubble { padding:12px 15px; border:1px solid #ebedf0; border-radius:3px 12px 12px 12px; background:#fff; color:#29343d; font-size:13px; line-height:1.75; white-space:pre-wrap; overflow-wrap:anywhere; }
.message-bubble p { margin:0; }
.user-bubble { border-color:#5d87ff; border-radius:12px 3px 12px 12px; background:#5d87ff; color:white; }
.message-error { margin-top:5px; color:#d46b08; font-size:11px; }
.typing { display:flex; align-items:center; gap:4px; color:#909399; font-size:11px; }
.typing i { width:5px; height:5px; border-radius:50%; background:#8eabff; animation:pulse 1s infinite ease-in-out; }
.typing i:nth-child(2) { animation-delay:.15s; }.typing i:nth-child(3) { animation-delay:.3s; }
.typing span { margin-left:5px; }
@keyframes pulse { 0%,80%,100% { opacity:.35; transform:translateY(0); } 40% { opacity:1; transform:translateY(-3px); } }
.tool-list { display:grid; gap:8px; margin-top:9px; }
.tool-card { overflow:hidden; border:1px solid #e8edf5; border-radius:9px; background:#fbfcfe; }
.tool-title { display:flex; align-items:center; gap:8px; padding:10px 11px 7px; }
.tool-check { display:grid; width:22px; height:22px; place-items:center; border-radius:50%; background:#e7fcf8; color:#0fa78a; font-size:12px; }
.tool-title > div { display:flex; flex:1; flex-direction:column; gap:1px; }
.tool-title strong { color:#3d4652; font-size:11px; font-weight:600; }
.tool-title small { color:#a3a6ad; font-size:10px; }
.tool-title :deep(.el-tag) { font-size:10px; }
.tool-card :deep(.el-collapse) { border:0; background:transparent; }
.tool-card :deep(.el-collapse-item__header) { height:30px; padding:0 11px; border:0; background:transparent; color:#7c7f85; font-size:11px; }
.tool-card :deep(.el-collapse-item__wrap) { border:0; background:transparent; }
.tool-card :deep(.el-collapse-item__content) { padding:0 11px 10px; }
.tool-section { margin-top:7px; }
.tool-section > span { color:#909399; font-size:10px; }
.tool-section pre { max-height:160px; overflow:auto; margin:4px 0 0; padding:8px; border-radius:6px; background:#f1f3f7; color:#596273; font:11px/1.5 ui-monospace, SFMono-Regular, Menlo, monospace; white-space:pre-wrap; overflow-wrap:anywhere; }
.composer { margin:0 20px 16px; padding:10px 12px 9px; border:1px solid #dfe4ed; border-radius:11px; transition:border-color .15s, box-shadow .15s; }
.composer:focus-within { border-color:#8eabff; box-shadow:0 0 0 3px rgba(93,135,255,.08); }
.composer :deep(.el-textarea__inner) { padding:3px 2px; border:0; box-shadow:none; color:#29343d; font-size:13px; line-height:1.65; }
.composer :deep(.el-textarea__inner:focus) { box-shadow:none; }
.composer-bottom { display:flex; align-items:center; justify-content:space-between; gap:12px; margin-top:7px; }
.composer-bottom > span { color:#a3a6ad; font-size:10px; }
.composer-bottom :deep(.el-button) { min-height:30px; padding:6px 12px; }
.agent-footnote { display:flex; align-items:center; justify-content:center; gap:5px; margin-top:10px; color:#a3a6ad; font-size:10px; }
@media (max-width:700px) {
  .agent-heading { align-items:flex-start; }
  .chat-card { height:calc(100vh - 205px); min-height:400px; }
  .chat-thread { padding:18px 13px; }
  .message-content { max-width:calc(100% - 35px); }
  .composer { margin:0 10px 10px; }
}
</style>
