import request from './request'

// 全量列表，参数 name/state 可选
export const listCategories = (params) => request.get('/admin/category', { params })

export const addCategory = (data) => request.post('/category', data)
export const editCategory = (id, data) => request.patch(`/category/${id}`, data)
export const deleteCategory = (id) => request.delete(`/category/${id}`)
