import request from './request'

// 参数 page/page_size/article_id/state（state: 0待审 1通过 2拒绝）
export const adminListComment = (params) => request.get('/admin/comment', { params })

// 审核：state 1通过 2拒绝
export const auditComment = (id, state) =>
  request.patch(`/admin/comment/${id}/state`, { state })

export const deleteComment = (id) => request.delete(`/admin/comment/${id}`)
