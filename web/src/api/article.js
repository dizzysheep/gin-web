import request from './request'

// 管理端列表（含草稿）
export const adminListArticle = (params) => request.get('/admin/article', { params })
export const adminGetArticle = (id) => request.get(`/admin/article/${id}`)

export const addArticle = (data) => request.post('/article', data)
export const editArticle = (id, data) => request.patch(`/article/${id}`, data)
export const deleteArticle = (id) => request.delete(`/article/${id}`)

// publish=true 发布（is_draft=0），false 转草稿
export const publishArticle = (id, publish) =>
  request.patch(`/article/${id}/publish`, { publish })

// state 1启用 0禁用
export const setArticleState = (id, state) =>
  request.patch(`/article/${id}/state`, { state })

export const setArticleTop = (id, isTop) =>
  request.patch(`/article/${id}/top`, { is_top: isTop })
