import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/LoginView.vue'),
    meta: { public: true }
  },
  {
    path: '/',
    component: () => import('../layout/AdminLayout.vue'),
    children: [
      { path: '', name: 'dashboard', component: () => import('../views/DashboardView.vue'), meta: { title: '仪表盘' } },
      { path: 'article', name: 'article-list', component: () => import('../views/ArticleListView.vue'), meta: { title: '文章管理' } },
      { path: 'article/add', name: 'article-add', component: () => import('../views/ArticleEditView.vue'), meta: { title: '新建文章' } },
      { path: 'article/edit/:id', name: 'article-edit', component: () => import('../views/ArticleEditView.vue'), meta: { title: '编辑文章' } },
      { path: 'category', name: 'category', component: () => import('../views/CategoryView.vue'), meta: { title: '分类管理' } },
      { path: 'tag', name: 'tag', component: () => import('../views/TagView.vue'), meta: { title: '标签管理' } },
      { path: 'comment', name: 'comment', component: () => import('../views/CommentView.vue'), meta: { title: '评论管理' } },
      { path: 'link', name: 'link', component: () => import('../views/LinkView.vue'), meta: { title: '友链管理' } },
      { path: 'option', name: 'option', component: () => import('../views/OptionView.vue'), meta: { title: '站点配置' } },
      { path: 'server', name: 'server', component: () => import('../views/ServerView.vue'), meta: { title: '服务器管理' } },
      { path: 'agent', name: 'agent-chat', component: () => import('../views/AgentChatView.vue'), meta: { title: 'AI 助手' } }
    ]
  },
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to) => {
  const token = localStorage.getItem('token')
  if (!to.meta.public && !token) {
    return { path: '/login', query: to.fullPath === '/' ? {} : { redirect: to.fullPath } }
  }
  if (to.path === '/login' && token) {
    return '/'
  }
  document.title = to.meta.title ? `${to.meta.title} - gin-web 管理后台` : 'gin-web 管理后台'
  return true
})

export default router
