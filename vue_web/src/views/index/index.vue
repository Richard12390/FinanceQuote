<script setup lang="ts">
import { onBeforeUnmount, onMounted } from "vue"
import { Client } from "@stomp/stompjs"
import SockJS from "sockjs-client"
import { wsEndpoint } from "@/config/env"
import { fetchFavorites, addFavorite, removeFavorite } from "@/api/favorites"

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
      favorites.value = new Set(items.map((item) => item.symbol))      
    } else {
      console.warn(
        "Unable to load favorites.",
        message ?? "Please try again later."
      )
    }
  } catch (error) {
    console.warn("Unable to load favorites. Please try again later.", error)
  }
}

const toggleFavorite = async (symbol: string) => {
  const favoriteList = new Set(favorites.value)
  const addingFavorite = !favoriteList.has(symbol)

  if (addingFavorite) {
    favoriteList.add(symbol)
    favorites.value = new Set(favoriteList)
    try {
      const response = await addFavorite(symbol)
      const { success, message } = response

      if (!success) {
        favoriteList.delete(symbol)
        favorites.value = new Set(favoriteList)
        console.warn(
          "Failed to add favorite.",
          message ?? "Please try again later."
        )
      }
    } catch (error) {
      favoriteList.delete(symbol)
      favorites.value = new Set(favoriteList)
      console.warn("Failed to add favorite. Please try again later.", error)
    }
  } else {
    favoriteList.delete(symbol)
    favorites.value = new Set(favoriteList)
    try {
      const response = await removeFavorite(symbol)
      const { success, message } = response
      if (!success) {
        favoriteList.add(symbol)
        favorites.value = new Set(favoriteList)
        console.warn(
          "Failed to remove favorite.",
          message ?? "Please try again later."
        )
      }
    } catch (error) {
      favoriteList.add(symbol)
      favorites.value = new Set(favoriteList)
      console.warn("Failed to remove favorite. Please try again later.", error)
    }
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
      title="Market Overview"
      :segment-options="segmentOptions"
      :get-segment-label="getSegmentLabel"
      :table-groups="tableGroups"
      :on-toggle-favorite="toggleFavorite"
      :is-favorite="isFavorite"
      :show-favorite-column="true"
    />
  </div>
</template>
