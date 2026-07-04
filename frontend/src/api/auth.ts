import request from './request'

export interface LoginParams {
  email: string
  password: string
}

export interface LoginResult {
  token: string
  user: {
    id: number
    email: string
    name: string
  }
}

export function login(data: LoginParams) {
  return request.post<LoginResult>('/users/login', data)
}
