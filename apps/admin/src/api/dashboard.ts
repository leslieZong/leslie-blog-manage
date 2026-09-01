import request from '@/utils/request'
import type { DashboardStats, RecentItem } from '@/types'

export function getStats() {
  return request.get<unknown, DashboardStats>('/dashboard/stats')
}

export function getRecentPosts(limit = 5) {
  return request.get<unknown, RecentItem[]>('/dashboard/recent-posts', {
    params: { limit },
  })
}

export function getRecentComments(limit = 5) {
  return request.get<unknown, RecentItem[]>('/dashboard/recent-comments', {
    params: { limit },
  })
}
