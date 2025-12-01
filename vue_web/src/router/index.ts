import { createRouter, createWebHistory } from "vue-router"
import IndexView from "@/views/index/index.vue"
import LayoutView from "@/views/layout/index.vue"
import FavoritesView from "@/views/favorites/index.vue"
import LoginView from "@/views/login/index.vue"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "",
      component: LayoutView,
      redirect: "/index",
      meta: { requiresAuth: true },
      children: [
        { path: "/index", name: "index", component: IndexView },
        { path: "/favorites", name: "favorites", component: FavoritesView },
      ],
    },
    { path: "/login", name: "login", component: LoginView },
  ],
})

type StoredLoginUser = { token: string; expiresAt?: number; }

const getStoredLoginUser = (): StoredLoginUser | null => {
  const raw = localStorage.getItem("loginUser")
  if (!raw) return null
  try {
    return JSON.parse(raw) as StoredLoginUser
  } catch {
    localStorage.removeItem("loginUser")
    return null
  }
}

const isTokenValid = (user: StoredLoginUser | null) => {
  if (!user?.token) return false
  if (user.expiresAt && Date.now() >= user.expiresAt * 1000) return false
  try {
    const [, payload = ""] = user.token.split(".")
    const decoded = JSON.parse(atob(payload))
    return decoded.exp ? Date.now() < decoded.exp * 1000 : true
  } catch {
    return false
  }
}

router.beforeEach((to, _from, next) => {
  const user = getStoredLoginUser()
  const isAuthenticated = isTokenValid(user)

  const rawRedirect = to.query.redirect
  const redirect: string =
    typeof rawRedirect === "string"
      ? rawRedirect
      : Array.isArray(rawRedirect)
        ? rawRedirect[0] ?? "/"
        : "/"

  if (to.matched.some(record => record.meta.requiresAuth) && !isAuthenticated) {
    next({ name: "login", query: { redirect } })
  } else if (to.name === "login" && isAuthenticated) {
    next({ name: "index" })
  } else {
    next()
  }
})



export default router
