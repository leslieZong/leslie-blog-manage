export interface MenuItem {
  path: string
  title: string
  icon: string
  // 隐藏的编辑页不显示
  hidden?: boolean
  parent?: string
}

// 由路由 meta 自动派生（取 router children 中非 hidden 的项）
export const routerMenu: MenuItem[] = [
  { path: '/dashboard', title: '仪表盘', icon: 'Odometer' },
  { path: '/posts', title: '文章管理', icon: 'Document' },
  { path: '/categories', title: '分类管理', icon: 'Files' },
  { path: '/projects', title: '项目管理', icon: 'Folder' },
  { path: '/tech-stack', title: '技术栈管理', icon: 'Cpu' },
  { path: '/media', title: '媒体库', icon: 'Picture' },
  { path: '/comments', title: '评论管理', icon: 'ChatDotRound' },
  { path: '/settings', title: '系统设置', icon: 'Setting' },
]
