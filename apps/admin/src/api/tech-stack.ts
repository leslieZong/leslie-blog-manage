import request from '@/utils/request'
import type { TechStack, TechStackForm } from '@/types'

export function getTechStack(params?: { keyword?: string }) {
  return request.get<unknown, TechStack[]>('/tech-stack', { params })
}

export function createTechStack(data: TechStackForm) {
  return request.post('/tech-stack', data)
}

export function updateTechStack(id: number, data: TechStackForm) {
  return request.put(`/tech-stack/${id}`, data)
}

export function deleteTechStack(id: number) {
  return request.delete(`/tech-stack/${id}`)
}
