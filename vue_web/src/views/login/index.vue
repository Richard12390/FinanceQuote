<script setup lang="ts">
import { reactive, ref } from "vue"
import { useRouter, useRoute } from "vue-router"
import { toast } from "vue-sonner"
import { login, register } from "@/api/auth"

const router = useRouter()
const route = useRoute()

const form = reactive({
  account: "",
  userId: "",
  password: "",
  confirmPassword: "",
  displayName: "",
})
const isSubmitting = ref(false)
const isRegisterMode = ref(false)

const resetForm = () => {
  form.account = ""
  form.userId = ""
  form.password = ""
  form.confirmPassword = ""
  form.displayName = ""
}

const handleSubmit = async () => {
  if (isSubmitting.value) return
  isSubmitting.value = true
  try {
    if (isRegisterMode.value) {
      if (form.password !== form.confirmPassword) {
        toast.error("Passwords do not match")
        return
      }
      const { data: response } = await register({
        account: form.account,
        password: form.password,
        displayName: form.displayName || undefined,
      })
      toast.success(`Registration success! Your user ID is ${response.userId}`)
      form.userId = String(response.userId)
      form.password = ""
      form.confirmPassword = ""
      isRegisterMode.value = false
      return
    }

     const { data: response } = await login({
      account: form.account,
      userId: form.userId ? Number(form.userId) : undefined,
      password: form.password,
    })
    const payload = response
localStorage.setItem("loginUser", JSON.stringify(payload))
    toast.success(`Welcome back, ${response.displayName}`)
    const redirect = (route.query.redirect as string) ?? "/"
    await router.push(redirect)
  } catch {

  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="flex min-screen items-center justify-center bg-slate-850 py-16">
    <div class="w-full max-w-md rounded-2xl border border-slate-800 bg-slate-900/70 p-8 shadow-xl backdrop-blur">
      <header class="mb-4 text-center">
        <h1 class="text-3xl font-semibold text-white">
          {{ isRegisterMode ? "Create account" : "Sign in" }}
        </h1>
        <p class="mt-4 text-sm text-slate-300">
          {{ isRegisterMode ?  "Fill in your details to get started" : "Access your quote dashboard" }}
        </p>
      </header>

      <form class="space-y-6" @submit.prevent="handleSubmit">
        <label class="block text-sm font-medium text-slate-300">
          Account
          <input
            v-model.trim="form.account"
            type="text"
            required
            autocomplete="username"
            class="mt-2 w-full rounded-lg border border-slate-700 bg-slate-900 px-4 py-2 text-slate-100 placeholder-slate-500 focus:border-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-500"
          />
        </label>

        <label v-if="!isRegisterMode" class="block text-sm font-medium text-slate-300">
          User ID (optional)
          <input
            v-model.trim="form.userId"
            type="number"
            min="1"
            autocomplete="one-time-code"
            placeholder="Enter the ID assigned after registration"
            class="mt-2 w-full rounded-lg border border-slate-700 bg-slate-900 px-4 py-2 text-slate-100 placeholder-slate-500 focus:border-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-500"
          />
        </label>

        <label v-if="isRegisterMode" class="block text-sm font-medium text-slate-300">
          Display Name (optional)
          <input
            v-model.trim="form.displayName"
            type="text"
            maxlength="128"
            class="mt-2 w-full rounded-lg border border-slate-700 bg-slate-900 px-4 py-2 text-slate-100 placeholder-slate-500 focus:border-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-500"
          />
        </label>

        <label class="block text-sm font-medium text-slate-300">
          Password
          <input
            v-model="form.password"
            type="password"
            required
            autocomplete="current-password"
            class="mt-2 w-full rounded-lg border border-slate-700 bg-slate-900 px-4 py-2 text-slate-100 placeholder-slate-500 focus:border-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-500"
          />
        </label>

        <label v-if="isRegisterMode" class="block text-sm font-medium text-slate-300">
          Confirm Password
          <input
            v-model="form.confirmPassword"
            type="password"
            required
            autocomplete="new-password"
            class="mt-2 w-full rounded-lg border border-slate-700 bg-slate-900 px-4 py-2 text-slate-100 placeholder-slate-500 focus:border-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-500"
          />
        </label>

        <button
          type="submit"
          class="w-full rounded-lg bg-slate-100 px-4 py-2 text-sm font-semibold text-slate-900 transition-colors hover:bg-white/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-300 disabled:cursor-not-allowed disabled:bg-slate-600 disabled:text-slate-300"
          :disabled="isSubmitting"
        >
          {{ isSubmitting ? "Processing..." : (isRegisterMode ? "Register" : "Sign in") }}
        </button>
      </form>

      <div class="mt-6 text-center text-sm text-slate-300">
        <button
          type="button"
          class="font-medium text-slate-100 underline-offset-4 hover:underline"
          @click="() => { isRegisterMode = !isRegisterMode; resetForm() }"
        >
           {{ isRegisterMode ? "Sign in" : "Create an account" }}
        </button>
      </div>

      <footer class="mt-8 text-center text-xs text-slate-300">
        © {{ new Date().getFullYear() }} All rights reserved.
      </footer>
    </div>
  </div>
</template>