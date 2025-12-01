import { computed, reactive, ref } from "vue"

export type AssetSegment = "all" | "crypto" | "stock" | "etf"

export type QuoteRow = {
  assetType: string
  symbol: string
  price: number | null
  changePercent: number | null
  volume: number | null
}

export type QuoteGroup = {
  key: string
  label: string
  rows: QuoteRow[]
}

export function useQuoteTable(initialSegment: AssetSegment = "all") {
  const segmentLabelMap: Record<Exclude<AssetSegment, "all">, string> = {
    crypto: "Crypto",
    stock: "Stock",
    etf: "ETF",
  }
  const otherLabel = "Others"

  const segmentOptions: AssetSegment[] = ["all", "crypto", "stock", "etf"]
  const selectedSegment = ref<AssetSegment>(initialSegment)

  const quotes = reactive(new Map<string, QuoteRow>())
  const favorites = ref(new Set<string>())

  const getSegmentLabel = (value: AssetSegment) =>
    value === "all" ? "All" : segmentLabelMap[value]

  const allRows = computed(() => Array.from(quotes.values()))

  const cryptoRows = computed(() =>
    allRows.value.filter(row => row.assetType === "crypto"),
  )
  const stockRows = computed(() =>
    allRows.value.filter(row => row.assetType === "stock"),
  )
  const etfRows = computed(() =>
    allRows.value.filter(row => row.assetType === "etf"),
  )
  const otherRows = computed(() =>
    allRows.value.filter(
      row =>
        row.assetType !== "stock" &&
        row.assetType !== "etf" &&
        row.assetType !== "crypto",
    ),
  )

  const groupedRows = computed(() => ({
    stock: stockRows.value,
    etf: etfRows.value,
    crypto: cryptoRows.value,
    others: otherRows.value,
  }))

  const tableGroups = computed<QuoteGroup[]>(() => {
    if (selectedSegment.value === "all") {
      const order: Array<Exclude<AssetSegment, "all"> | "others"> = [
        "crypto",
        "stock",
        "etf",
        "others",
      ]
      const groups = order
        .map(key => {
          if (key === "others") {
            return {
              key,
              label: otherLabel,
              rows: groupedRows.value.others,
            }
          }
          return {
            key,
            label: segmentLabelMap[key],
            rows: groupedRows.value[key],
          }
        })
        .filter(group => group.rows.length > 0)

      if (groups.length > 0) {
        return groups
      }
      return [
        {
          key: "all",
          label: "All",
          rows: allRows.value,
        },
      ]
    }

    const key = selectedSegment.value as Exclude<AssetSegment, "all">
    return [
      {
        key,
        label: segmentLabelMap[key],
        rows: groupedRows.value[key],
      },
    ]
  })

  const isFavorite = (symbol: string) => favorites.value.has(symbol)

  return {
    segmentOptions,
    selectedSegment,
    getSegmentLabel,
    tableGroups,
    quotes,
    favorites,
    allRows,
    cryptoRows,
    stockRows,
    etfRows,
    otherRows,
    groupedRows,
    isFavorite,
  }
}
