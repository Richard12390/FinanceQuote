<!-- <script setup lang="ts">
import { onMounted, onBeforeUnmount, ref, watch, nextTick } from 'vue'
import type { IChartApi, ISeriesApi, LineData } from 'lightweight-charts'

const props = defineProps<{
  data: LineData<'Line'>[]
  height?: number
}>()
const container = ref<HTMLDivElement | null>(null)
let chart: IChartApi | undefined
let series: ISeriesApi<'Line'> | undefined
let resizeObserver: ResizeObserver | undefined

const applyData = () => series?.setData(props.data as unknown as LineData<'Line'>[])

onMounted(async () => {
  await nextTick()
  if (!container.value) return

  const { createChart} = await import('lightweight-charts')

  chart = createChart(container.value, {
    width: container.value.clientWidth,
    height: props.height ?? 360,
    layout: { background: { color: '#0b0b0d' }, textColor: '#e0e0e0' },
    grid: { vertLines: { color: '#1f2937' }, horzLines: { color: '#1f2937' } },
    crosshair: { mode: 0 },
  })
  series = chart.addLineSeries({ color: '#10b981', lineWidth: 2 })

  applyData()

  resizeObserver = new ResizeObserver(entries => {
    const entry = entries[0]
    if (!entry?.contentRect || !chart) return
    chart.applyOptions({
      width: Math.floor(entry.contentRect.width),
      height: Math.floor(entry.contentRect.height),
    })
    chart.timeScale().fitContent()
  })
  resizeObserver.observe(container.value)
})

onBeforeUnmount(() => {
  resizeObserver?.disconnect()
  chart?.remove()
  chart = undefined
  series = undefined
})

watch(() => props.data, applyData, { deep: true })
</script>

<template>
  <div ref="container" class="w-full h-full rounded-xl overflow-hidden"></div>
</template> -->
