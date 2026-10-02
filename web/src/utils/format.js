// 后端时间字段为 uint32 Unix 秒
export function formatTime(ts) {
  if (!ts) return '-'
  const d = new Date(ts * 1000)
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`
}

export const stateText = { 0: '禁用', 1: '启用' }
export const commentStateText = { 0: '待审核', 1: '已通过', 2: '已拒绝' }
export const commentStateTag = { 0: 'warning', 1: 'success', 2: 'danger' }
