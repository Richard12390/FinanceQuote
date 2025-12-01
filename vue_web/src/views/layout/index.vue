<script setup lang="ts">
import AppSidebar from '@/components/AppSidebar.vue'
import { SidebarProvider, SidebarTrigger, SidebarInset } from "@/components/ui/sidebar"
import { Toaster } from "@/components/ui/sonner"
</script>

<template>
  <SidebarProvider>
    <div class="flex min-h-screen">
      <AppSidebar />
      <SidebarInset class="flex-1 overflow-auto">
        <main class="min-h-screen w-full px-6 py-6">
          <SidebarTrigger class="mb-4" />
          <router-view />
        </main>
      </SidebarInset>
    </div>
    <Toaster />
  </SidebarProvider>
</template>
