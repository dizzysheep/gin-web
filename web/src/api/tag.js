import request from './request'

// 分页列表，参数 page/page_size/name/state
export const listTags = (params) => request.get('/admin/tag', { params })

export const addTag = (data) => request.post('/tag', data)
export const editTag = (id, data) => request.patch(`/tag/${id}`, data)
export const deleteTag = (id) => request.delete(`/tag/${id}`)
