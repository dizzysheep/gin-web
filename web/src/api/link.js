import request from './request'

// 全量列表，参数 name/state 可选
export const listLinks = (params) => request.get('/admin/link', { params })

export const addLink = (data) => request.post('/link', data)
export const editLink = (id, data) => request.patch(`/link/${id}`, data)
export const deleteLink = (id) => request.delete(`/link/${id}`)
