export type ApiResult<T> = {
  success: boolean
  message?: string
  data?: T
}

export type FavoriteItem = {
  userId: number
  symbol: string
}
