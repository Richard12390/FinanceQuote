<script lang="ts" setup>
import type { ToasterProps } from "vue-sonner"
import { Toaster as Sonner } from "vue-sonner"

const props = withDefaults(defineProps<ToasterProps>(), {
  position: "top-center",
  toastOptions: () => ({
    class: "bg-slate-900/80 text-white border border-white/20 shadow-2xl backdrop-blur",
    style: {
      backgroundColor: "rgba(15,23,42,0.85)",
    },
  }),
})
</script>

<template>
  <Sonner
    class="toaster group"
    v-bind="props"
    :style="{
      '--normal-bg': 'var(--popover)',
      '--normal-text': 'var(--popover-foreground)',
      '--normal-border': 'var(--border)',
    }"
  />
</template>
:global([data-sonner-toaster][data-x-position='center']) {
  top: 50% !important;
  left: 50% !important;
  transform: translate(-50%, -50%) !important;
}
<style scoped>
</style>
