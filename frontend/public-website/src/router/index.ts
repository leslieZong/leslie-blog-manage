import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/home/index.vue'),
    meta: {
      hasHeader: true,
      zhTitle: '首页',
      enTitle: 'Home',
    },
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login/index.vue'),
    meta: {
      hasHeader: false,
      zhTitle: '登录',
      enTitle: 'Login',
    },
  },
  {
    path: '/article',
    name: 'article',
    component: () => import('@/views/article/index.vue'),
    meta: {
      hasHeader: true,
      zhTitle: '文章',
      enTitle: 'Article',
    },
  },
  {
    path: '/category',
    name: 'category',
    component: () => import('@/views/category/index.vue'),
    meta: {
      hasHeader: true,
      zhTitle: '分类',
      enTitle: 'Category',
    },
  },
  {
    path: '/project',
    name: 'project',
    component: () => import('@/views/project/index.vue'),
    meta: {
      hasHeader: true,
      zhTitle: '项目',
      enTitle: 'Project',
    },
  },
  {
    path: '/about',
    name: 'about',
    component: () => import('@/views/about/index.vue'),
    meta: {
      hasHeader: true,
      zhTitle: '关于',
      enTitle: 'About',
    },
  },
]
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})
router.beforeEach((to, from) => {
  console.log("🚀 ~ from:", from)
  // 设置文档标题
  const lang = localStorage.getItem('lang') || 'zh'
  let title = lang === 'zh' ? (to.meta.zhTitle as string) : (to.meta.enTitle as string)
  if (!title) {
    title = to.name as string
  }
  document.title = `${title} | Leslie Blog`
})

export default router
