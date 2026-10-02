import request from './request'

export const getServerInfo = () => request.get('/admin/server')
