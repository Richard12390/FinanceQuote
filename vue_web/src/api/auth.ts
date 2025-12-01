import request from "@/utils/request"

export type LoginRequest = {
  account: string
  userId?: number
  password: string
}

export type LoginResponse = {
  token: string
  userId: number
  displayName: string
}

export type RegisterRequest = {
  account: string
  password: string
  displayName?: string
}

export type RegisterResponse = {
  userId: number
  account: string
  displayName?: string
}

export const login = (payload: LoginRequest) =>
  request.post<LoginResponse>("/auth/login", payload)

export const logout = () =>
  request.post('/auth/logout'); 

export const register = (payload: RegisterRequest) =>
  request.post<RegisterResponse>("/auth/register", payload)