<script setup lang="ts">
import { onBeforeUnmount, onMounted } from "vue"
import { Client } from "@stomp/stompjs"
import SockJS from "sockjs-client"
import { wsEndpoint } from "@/config/env"
import { fetchFavorites, removeFavorite } from "@/api/favorites"

import QuoteTable from "@/components/quotes/QuoteTable.vue"
import { useQuoteTable } from "@/composables/useQuoteTable"
import type { FavoriteItem } from "@/types/api"



const {
  segmentOptions,
  selectedSegment,
  getSegmentLabel,
  tableGroups,
  quotes,
  favorites,
  isFavorite,
} = useQuoteTable("all")

const client = new Client({
  webSocketFactory: () => new SockJS(wsEndpoint),
  reconnectDelay: 5000,
})

client.onConnect = () => {
  client.subscribe("/topic/quotes", message => {
    try {
      const payload = JSON.parse(message.body)
      if (!favorites.value.has(payload.symbol)) {
        return
      }

      const priceValue =
        payload.lastPrice != null
          ? Number(payload.lastPrice)
          : payload.price != null
            ? Number(payload.price)
            : null

      const changePercent =
        payload.changePercent != null ? Number(payload.changePercent) : null
      const volume =
        payload.volume != null ? Number(payload.volume) : null

      quotes.set(payload.symbol, {
        assetType: payload.assetType
          ? String(payload.assetType).toLowerCase()
          : "unknown",
        symbol: payload.symbol,
        price: priceValue,
        changePercent,
        volume,
      })
    } catch (error) {
      console.warn("Failed to parse quote payload", error)
    }
  })
}

const loadFavorites = async () => {
  if (!localStorage.getItem("loginUser")) {
    return
  }  
  try {
    const response = await fetchFavorites()
    const { success, data, message } = response
    if (success) {
      const items = (data ?? []) as FavoriteItem[]
      const symbols = new Set<string>(
        items.map((item: FavoriteItem) => item.symbol),
      )
      favorites.value = symbols
      items.forEach((item: FavoriteItem) => {
        if (!quotes.has(item.symbol)) {
          quotes.set(item.symbol, {
            assetType: "unknown",
            symbol: item.symbol,
            price: null,
            changePercent: null,
            volume: null,
          })
        }
      })
    } else {
      favorites.value.clear()
      quotes.clear()
      console.warn(
        "Unable to load favorites.",
        message ?? "Please try again later.",
      )
    }
  } catch (error) {
    favorites.value.clear()
    quotes.clear()
    console.warn("Unable to load favorites. Please try again later.", error)
  }
}

const removeSymbol = async (symbol: string) => {
  if (!favorites.value.has(symbol)) {
    return
  }

  const previousFavorites = new Set(favorites.value)
  const previousQuote = quotes.get(symbol)

  favorites.value.delete(symbol)
  quotes.delete(symbol)

  try {
    const response = await removeFavorite(symbol)
    const { success, message } = response  
    
    if (!success) {
      favorites.value = previousFavorites
      if (previousQuote) {
        quotes.set(symbol, previousQuote)
      }
      console.warn(
        "Failed to remove favorite.",
        message ?? "Please try again later."
      )
    }
  } catch (error) {
    favorites.value = previousFavorites
    if (previousQuote) {
      quotes.set(symbol, previousQuote)
    }
    console.warn("Failed to remove favorite. Please try again later.", error)
  }
}

onMounted(async () => {
  await loadFavorites()
  client.activate()
})

onBeforeUnmount(() => client.deactivate())
</script>

<template>
  <div class="w-full px-6 pb-8">
    <QuoteTable
      v-model:selected-segment="selectedSegment"
      title="Favorites"
      :segment-options="segmentOptions"
      :get-segment-label="getSegmentLabel"
      :table-groups="tableGroups"
      :on-toggle-favorite="removeSymbol"
      :is-favorite="isFavorite"
      :show-favorite-column="true"
    />
  </div>
</template>
