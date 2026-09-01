import request from '@/utils/request'
import type { SiteSettings } from '@/types'

export function getSettings() {
  return request.get<unknown, SiteSettings>('/settings')
}

export function updateSettings(data: SiteSettings) {
  return request.put('/settings', data)
}
