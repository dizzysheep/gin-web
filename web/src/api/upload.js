import request from './request'

// FormData 的 file 字段，返回 {url}
export const uploadImage = (formData) =>
  request.post('/upload/image', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
