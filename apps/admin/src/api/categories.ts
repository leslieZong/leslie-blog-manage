import request from '@/utils/request'
import type { Category, CategoryForm } from '@/types'

export function getCategories(params?: { keyword?: string }) {
  return request.get<unknown, Category[]>('/categories', { params })
}

export function getCategory(id: number) {
  return request.get<unknown, Category>(`/categories/${id}`)
}

export function createCategory(data: CategoryForm) {
  return request.post('/categories', data)
}

export function updateCategory(id: number, data: CategoryForm) {
  return request.put(`/categories/${id}`, data)
}

export function deleteCategory(id: number) {
  return request.delete(`/categories/${id}`)
}
