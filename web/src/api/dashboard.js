import request from './request'

export const getDashboardOverview = () => request.get('/admin/dashboard/overview')
export const getDashboardTrend = (days = 7) => request.get('/admin/dashboard/trend', { params: { days } })
