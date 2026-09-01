import request from '@/utils/request'
import type { PageQuery, PageResult, Comment } from '@/types'

export function getComments(params: PageQuery) {
  return request.get<unknown, PageResult<Comment>>('/comments', { params })
}

export function approveComment(id: number) {
  return request.patch(`/comments/${id}/approve`)
}

export function rejectComment(id: number) {
  return request.patch(`/comments/${id}/reject`)
}

export function deleteComment(id: number) {
  return request.delete(`/comments/${id}`)
}
