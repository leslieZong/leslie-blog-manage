// 通用分页 / 响应结构
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PageQuery {
  page?: number
  pageSize?: number
  keyword?: string
  status?: string | number
}

export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

// 用户 / 鉴权
export interface UserInfo {
  id: number
  username: string
  nickname: string
  avatar?: string
  email?: string
  role?: string
}

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  token: string
  userInfo: UserInfo
}

// Dashboard 统计
export interface DashboardStats {
  postCount: number
  publishedCount: number
  categoryCount: number
  commentCount: number
  projectCount: number
  techStackCount: number
  mediaCount: number
  viewCount: number
}

export interface RecentItem {
  id: number
  title: string
  createdAt: string
  status?: number
}

// 文章
export interface Post {
  id: number
  title: string
  slug?: string
  summary?: string
  content: string
  cover?: string
  categoryId?: number
  categoryName?: string
  tags?: string[]
  status: number // 0 草稿 1 已发布 2 已下线
  viewCount: number
  commentCount: number
  isTop?: boolean
  publishedAt?: string
  createdAt: string
  updatedAt: string
}

export interface PostForm {
  id?: number
  title: string
  slug?: string
  summary?: string
  content: string
  cover?: string
  categoryId?: number
  tags?: string[]
  status: number
  isTop?: boolean
}

// 分类
export interface Category {
  id: number
  name: string
  slug?: string
  description?: string
  parentId?: number
  postCount?: number
  sortOrder?: number
  createdAt: string
}

export interface CategoryForm {
  id?: number
  name: string
  slug?: string
  description?: string
  parentId?: number
  sortOrder?: number
}

// 项目
export interface Project {
  id: number
  name: string
  description?: string
  cover?: string
  demoUrl?: string
  repoUrl?: string
  techStack?: string[]
  status: number
  sortOrder?: number
  createdAt: string
}

export interface ProjectForm {
  id?: number
  name: string
  description?: string
  cover?: string
  demoUrl?: string
  repoUrl?: string
  techStack?: string[]
  status: number
  sortOrder?: number
}

// 技术栈
export interface TechStack {
  id: number
  name: string
  icon?: string
  category?: string
  level?: number
  description?: string
  sortOrder?: number
  createdAt: string
}

export interface TechStackForm {
  id?: number
  name: string
  icon?: string
  category?: string
  level?: number
  description?: string
  sortOrder?: number
}

// 媒体
export interface MediaItem {
  id: number
  name: string
  url: string
  type: string
  size: number
  mimeType?: string
  createdAt: string
}

// 评论
export interface Comment {
  id: number
  postId: number
  postTitle?: string
  content: string
  author: string
  email?: string
  avatar?: string
  parentId?: number
  status: number // 0 待审 1 通过 2 拒绝
  ip?: string
  createdAt: string
}

// 设置
export interface SiteSettings {
  siteName: string
  logo?: string
  description?: string
  keywords?: string
  author?: string
  icp?: string
  social?: {
    github?: string
    email?: string
    twitter?: string
  }
}
