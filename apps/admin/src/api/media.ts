import request from '@/utils/request'
import type { PageQuery, PageResult, MediaItem } from '@/types'

export function getMediaList(params: PageQuery & { type?: string }) {
  return request.get<unknown, PageResult<MediaItem>>('/media', { params })
}

export function uploadMedia(file: File, onProgress?: (percent: number) => void) {
  const form = new FormData()
  form.append('file', file)
  return request.post<unknown, MediaItem>('/media/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress(e) {
      if (e.total && onProgress) {
        onProgress(Math.round((e.loaded / e.total) * 100))
      }
    },
  })
}

export function deleteMedia(id: number) {
  return request.delete(`/media/${id}`)
}
