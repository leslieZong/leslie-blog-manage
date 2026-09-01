import request from '@/utils/request'
import type { PageQuery, PageResult, Project, ProjectForm } from '@/types'

export function getProjects(params: PageQuery) {
  return request.get<unknown, PageResult<Project>>('/projects', { params })
}

export function getProject(id: number) {
  return request.get<unknown, Project>(`/projects/${id}`)
}

export function createProject(data: ProjectForm) {
  return request.post('/projects', data)
}

export function updateProject(id: number, data: ProjectForm) {
  return request.put(`/projects/${id}`, data)
}

export function deleteProject(id: number) {
  return request.delete(`/projects/${id}`)
}
