import request from './request'

export const getOptions = () => request.get('/common/option')
export const saveOptions = (options) =>
  request.put('/admin/option', { options })
