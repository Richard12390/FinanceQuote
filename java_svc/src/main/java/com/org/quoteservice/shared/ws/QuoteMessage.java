package com.org.quoteservice.shared.ws;

public record QuoteMessage(
                String assetType,
                String symbol,
                Double price,
                Double change,
                Double changePercent,
                Double volume,
                Long ts,
                Long tradeTimestamp,
                String tradeTime) {
}
