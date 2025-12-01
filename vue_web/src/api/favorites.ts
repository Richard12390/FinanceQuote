import request from "@/utils/request"
import type { ApiResult, FavoriteItem } from "@/types/api"

export const fetchFavorites = () =>
  request.get<ApiResult<FavoriteItem[]>, ApiResult<FavoriteItem[]>>("/favorites")

export const addFavorite = (symbol: string) =>
  request.put<ApiResult<void>, ApiResult<void>>(`/favorites/${encodeURIComponent(symbol)}`)

export const removeFavorite = (symbol: string) =>
  request.delete<ApiResult<void>, ApiResult<void>>(`/favorites/${encodeURIComponent(symbol)}`)