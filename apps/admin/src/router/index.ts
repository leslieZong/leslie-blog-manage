import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login/index.vue'),
    meta: { public: true, title: '登录' },
  },
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '仪表盘', icon: 'Odometer' },
      },
      {
        path: 'posts',
        name: 'posts',
        component: () => import('@/views/posts/index.vue'),
        meta: { title: '文章管理', icon: 'Document' },
      },
      {
        path: 'posts/editor',
        name: 'post-create',
        component: () => import('@/views/posts/editor.vue'),
        meta: { title: '新建文章', parent: 'posts', hidden: true },
      },
      {
        path: 'posts/editor/:id',
        name: 'post-edit',
        component: () => import('@/views/posts/editor.vue'),
        meta: { title: '编辑文章', parent: 'posts', hidden: true },
      },
      {
        path: 'categories',
        name: 'categories',
        component: () => import('@/views/categories/index.vue'),
        meta: { title: '分类管理', icon: 'Files' },
      },
      {
        path: 'projects',
        name: 'projects',
        component: () => import('@/views/projects/index.vue'),
        meta: { title: '项目管理', icon: 'Folder' },
      },
      {
        path: 'tech-stack',
        name: 'tech-stack',
        component: () => import('@/views/tech-stack/index.vue'),
        meta: { title: '技术栈管理', icon: 'Cpu' },
      },
      {
        path: 'media',
        name: 'media',
        component: () => import('@/views/media/index.vue'),
        meta: { title: '媒体库', icon: 'Picture' },
      },
      {
        path: 'comments',
        name: 'comments',
        component: () => import('@/views/comments/index.vue'),
        meta: { title: '评论管理', icon: 'ChatDotRound' },
      },
      {
        path: 'settings',
        name: 'settings',
        component: () => import('@/views/settings/index.vue'),
        meta: { title: '系统设置', icon: 'Setting' },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/views/error/404.vue'),
    meta: { public: true, title: '页面不存在' },
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  document.title = `${(to.meta.title as string) || 'Admin'} | Leslie Blog Admin`
  if (!to.meta.public && !auth.isLogin) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.name === 'login' && auth.isLogin) {
    return { name: 'dashboard' }
  }
  return true
})

export default router
