<script setup lang="ts">
import type {
  AssetSegment,
  QuoteGroup,
} from "@/composables/useQuoteTable"

const FILLED_STAR = "\u2605"
const EMPTY_STAR = "\u2606"

const props = defineProps<{
  title: string
  description?: string
  segmentOptions: AssetSegment[]
  selectedSegment: AssetSegment
  getSegmentLabel: (value: AssetSegment) => string
  tableGroups: QuoteGroup[]
  onToggleFavorite?: (symbol: string) => void
  isFavorite?: (symbol: string) => boolean
  showFavoriteColumn?: boolean
}>()

const emit = defineEmits<{
  (e: "update:selectedSegment", value: AssetSegment): void
}>()

const toggleFavorite = (symbol: string) => {
  props.onToggleFavorite?.(symbol)
}
</script>

<template>
  <div class="w-full space-y-8 rounded-2xl border border-slate-250 bg-white shadow">
    <header class="flex flex-wrap items-center justify-between gap-4 px-6 py-4">
      <div class="space-y-1">
        <h1 class="text-xl font-semibold text-slate-900">{{ title }}</h1>
        <p v-if="description" class="text-sm text-slate-500">
          {{ description }}
        </p>
      </div>
      <div class="flex items-center gap-2 text-sm text-slate-500">
        <button
          v-for="option in segmentOptions"
          :key="option"
          class="rounded-lg px-4 py-2 font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-400"
          :class="option === selectedSegment
            ? 'bg-slate-900 text-white shadow'
            : 'bg-slate-100 text-slate-600 hover:bg-slate-200'"
          @click="emit('update:selectedSegment', option)"
        >
          {{ getSegmentLabel(option) }}
        </button>
      </div>
    </header>

    <section
      v-for="group in tableGroups"
      :key="group.key"
      class="px-6 pb-6"
    >
      <header class="flex items-center justify-between px-6 py-3">
        <h2 class="text-lg font-semibold text-slate-900">{{ group.label }}</h2>
      </header>

      <div class="rounded-xl border border-slate-250 px-6 py-6">
        <table class="w-full border-collapse table-fixed border border-slate-250 text-sm text-slate-700">
          <thead class="bg-slate-50 text-xs tracking-wide text-slate-500">
            <tr>
              <th
                v-if="showFavoriteColumn"
                class="w-16 border border-slate-250 px-4 py-3 text-[1.1rem] font-bold"
              >
                {{ FILLED_STAR }}
              </th>
              <th class="border border-slate-250 px-4 py-3 text-[1.1rem] font-bold">
                Symbol
              </th>
              <th class="border border-slate-250 px-4 py-3 text-center text-[1.1rem] font-bold">
                Price
              </th>
              <th class="border border-slate-250 px-4 py-3 text-center text-[1.1rem] font-bold">
                Change (%)
              </th>
              <th class="border border-slate-250 px-4 py-3 text-center text-[1.1rem] font-bold">
                Volume
              </th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in group.rows"
              :key="row.symbol"
              class="hover:bg-slate-50"
            >
              <td
                v-if="showFavoriteColumn"
                class="border border-slate-250 px-4 py-3 text-center"
              >
                <button
                  type="button"
                  class="text-base transition-colors"
                  :class="isFavorite?.(row.symbol)
                    ? 'text-amber-400'
                    : 'text-slate-300 hover:text-slate-500'"
                  @click="toggleFavorite(row.symbol)"
                >
                  {{
                    isFavorite?.(row.symbol)
                      ? FILLED_STAR
                      : EMPTY_STAR
                  }}
                </button>
              </td>
              <td class="border border-slate-250 px-4 py-3 text-center font-semibold text-blue-600">
                {{ row.symbol }}
              </td>
              <td class="border border-slate-250 px-4 py-3 text-center font-semibold text-slate-600">
                {{
                  row.price == null
                    ? "---"
                    : row.price.toLocaleString()
                }}
              </td>
              <td class="border border-slate-250 px-4 py-3 text-center font-semibold text-slate-600">
                <span
                  :class="row.changePercent == null
                    ? 'text-slate-500'
                    : row.changePercent >= 0
                      ? 'text-emerald-500'
                      : 'text-rose-500'"
                >
                  {{
                    row.changePercent == null
                      ? "---"
                      : row.changePercent.toFixed(2) + " %"
                  }}
                </span>
              </td>
              <td class="border border-slate-250 px-4 py-3 text-center font-semibold text-slate-600">
                {{ row.volume?.toLocaleString() ?? "--" }}
              </td>
            </tr>
            <tr v-if="group.rows.length === 0">
              <td
                :colspan="showFavoriteColumn ? 5 : 4"
                class="px-4 py-6 text-center text-slate-400"
              >
                No data yet
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <section
      v-if="tableGroups.length === 0"
      class="px-6 pb-8"
    >
      <div class="flex min-h-[18rem] items-center justify-center rounded-xl border border-dashed border-slate-300 bg-slate-50 text-slate-500">
        No market data available yet.
      </div>
    </section>
  </div>
</template>
