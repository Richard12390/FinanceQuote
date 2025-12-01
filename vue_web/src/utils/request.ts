import axios from "axios"
import router from "@/router"

type LoginUser = {
  token: string
  userId: number
  name: string
}

const request = axios.create({
  baseURL: "/api",
  timeout: 600000,
})

request.interceptors.request.use(
  config => {

    const stored = localStorage.getItem("loginUser")
    if (stored) {
      try {
        const loginUser: LoginUser = JSON.parse(stored)
        if (loginUser?.token) {
          config.headers = config.headers ?? {}
          config.headers.Authorization = `Bearer ${loginUser.token}`
        }
      } catch {
        localStorage.removeItem("loginUser")
      }
    }
    return config

  },
  error => Promise.reject(error),
)

request.interceptors.response.use(
  response => response.data,
  error => {
    const status = error.response?.status
    if (status === 401 || status === 403) {
      localStorage.removeItem("loginUser")
      const currentRoute = router.currentRoute.value
      const redirect = currentRoute?.fullPath ?? "/"
      if (currentRoute?.name !== "login") {
        router.push({ name: "login", query: { redirect } })
      }
    } 
    return Promise.reject(error)
  },
)

export default request
