import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/home/index.vue'),
    meta: {
      hasHeader: true,
    },
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login/index.vue'),
  },
  {
    path: '/article',
    name: 'article',
    component: () => import('@/views/article/index.vue'),
    meta: {
      hasHeader: true,
    },
  },
  {
    path: '/category',
    name: 'category',
    component: () => import('@/views/category/index.vue'),
    meta: {
      hasHeader: true,
    },
  },
  {
    path: '/project',
    name: 'project',
    component: () => import('@/views/project/index.vue'),
    meta: {
      hasHeader: true,
    },
  },
  {
    path: '/about',
    name: 'about',
    component: () => import('@/views/about/index.vue'),
    meta: {
      hasHeader: true,
    },
  },
]
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})
// router.beforeEach((to, from) => {

// })

export default router
