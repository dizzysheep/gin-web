import request from './request'

export const getAgentTools = () => request.get('/admin/agent/tools')

export async function chatWithAgent(message, onEvent) {
  const token = localStorage.getItem('token')
  const response = await fetch('/api/v1/admin/agent/chat', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: 'Bearer ' + token } : {})
    },
    body: JSON.stringify({ message })
  })

  const contentType = response.headers.get('content-type') || ''
  if (!contentType.includes('text/event-stream')) {
    const body = await response.json().catch(() => ({}))
    if ([210001, 210002, 210003].includes(body.code)) {
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      window.location.href = '/login'
    }
    throw new Error(body.msg || 'Agent 请求失败')
  }

  const reader = response.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  let finalResult
  let finished = false
  while (!finished) {
    const { value, done } = await reader.read()
    buffer += decoder.decode(value || new Uint8Array(), { stream: !done })
    const frames = buffer.split(/\r?\n\r?\n/)
    buffer = frames.pop() || ''
    for (const frame of frames) {
      const eventName = frame.match(/^event:\s*(.+)$/m)?.[1]
      const data = frame.match(/^data:\s*(.*)$/m)?.[1]
      if (!eventName || data === undefined) continue
      const payload = JSON.parse(data)
      if (eventName === 'error') throw new Error(payload.message || 'Agent 请求失败')
      if (eventName === 'message') finalResult = payload
      onEvent?.(eventName, payload)
      if (eventName === 'done') finished = true
    }
    if (done) finished = true
  }
  if (!finalResult) throw new Error('Agent 流提前结束，未收到最终回答')
  return finalResult
}
