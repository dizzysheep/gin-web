import request from './request'

export const login = (data) => request.post('/user/login', data)
export const logout = () => request.post('/user/logout')
export const getUserInfo = () => request.get('/user/info')
export const changePassword = (data) => request.patch('/user/password', data)
