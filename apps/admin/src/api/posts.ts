import request from '@/utils/request'
import type { PageQuery, PageResult, Post, PostForm } from '@/types'

export interface PostQuery extends PageQuery {
  categoryId?: number
  tag?: string
}

export function getPosts(params: PostQuery) {
  return request.get<unknown, PageResult<Post>>('/posts', { params })
}

export function getPost(id: number) {
  return request.get<unknown, Post>(`/posts/${id}`)
}

export function createPost(data: PostForm) {
  return request.post('/posts', data)
}

export function updatePost(id: number, data: PostForm) {
  return request.put(`/posts/${id}`, data)
}

export function deletePost(id: number) {
  return request.delete(`/posts/${id}`)
}

export function publishPost(id: number) {
  return request.patch(`/posts/${id}/publish`)
}

export function unpublishPost(id: number) {
  return request.patch(`/posts/${id}/unpublish`)
}

export function toggleTop(id: number, isTop: boolean) {
  return request.patch(`/posts/${id}/top`, { isTop })
}
