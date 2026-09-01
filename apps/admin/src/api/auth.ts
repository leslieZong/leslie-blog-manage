import request from '@/utils/request'
import type { LoginParams, LoginResult, UserInfo } from '@/types'

export function login(data: LoginParams) {
  return request.post<unknown, LoginResult>('/auth/login', data)
}

export function logout() {
  return request.post('/auth/logout')
}

export function getUserInfo() {
  return request.get<unknown, UserInfo>('/auth/user')
}
