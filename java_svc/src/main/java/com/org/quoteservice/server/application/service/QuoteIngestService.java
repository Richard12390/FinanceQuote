package com.org.quoteservice.server.application.service;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.org.quoteservice.shared.ws.QuoteMessage;
import com.org.quoteservice.shared.ws.QuoteStreamProducer;
import java.time.Instant;
import java.time.ZoneOffset;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import com.org.quoteservice.server.application.port.in.QuoteIngestUseCase;
import com.org.quoteservice.server.application.port.out.QuoteRepositoryPort;
import com.org.quoteservice.server.application.port.out.SymbolRepositoryPort;
import com.org.quoteservice.server.domain.model.QuoteNorm;
import com.org.quoteservice.server.domain.model.Symbol;
import com.org.quoteservice.server.domain.service.SymbolDomainService;
import org.springframework.stereotype.Service;

@Slf4j
@Service
@RequiredArgsConstructor
public class QuoteIngestService implements QuoteIngestUseCase {
    private final ObjectMapper objectMapper;
    private final SymbolRepositoryPort symbolRepository;
    private final QuoteRepositoryPort quoteRepository;
    private final SymbolDomainService symbolDomainService;
    private final QuoteStreamProducer quoteStreamProducer;
    private final Map<String, Double> lastPriceCache = new ConcurrentHashMap<>();
    private final Map<String, Double> snapshotPriceCache = new ConcurrentHashMap<>();
    private final Map<String, Double> lastChangePercentCache = new ConcurrentHashMap<>();

    @Override
    public void handleMessage(String subject, String jsonPayload) {
        try {
            QuoteNorm quote = objectMapper.readValue(jsonPayload, QuoteNorm.class);
            int categoryId = symbolDomainService.resolveCategoryId(quote.getAssetType());
            Integer exchangeId = symbolRepository.findExchangeIdByCode(quote.getExchange())
                    .orElse(null);

            Symbol symbol = Symbol.builder()
                    .source(quote.getSource())
                    .categoryId(categoryId)
                    .symbol(quote.getSymbol())
                    .symbolOrigin(quote.getSymbolOrigin())
                    .exchangeId(exchangeId)
                    .currency(quote.getCurrency())
                    .isActive(1)
                    .build();
            long symbolId = symbolRepository.upsertSymbol(symbol);
            if (symbolId == 0L) {
                throw new IllegalStateException("symbol_id not resolved for " + quote.getSymbol());
            }

            String symbolCode = quote.getSymbol();
            Double newPrice = quote.getPrice();
            Double prevClose = quote.getPrevClose();
            String marketState = quote.getMarketState();
            String assetType = quote.getAssetType();

            boolean hasNewPrice = newPrice != null && Double.compare(newPrice, 0.0d) != 0;
            Double price = newPrice;
            if (hasNewPrice) {
                lastPriceCache.put(symbolCode, newPrice);
            } else {
                if (prevClose != null && Double.compare(prevClose, 0.0d) != 0) {
                    price = prevClose;
                } else {
                    price = lastPriceCache.get(symbolCode);
                }
            }

            Double change = (price != null && prevClose != null)
                    ? price - prevClose
                    : null;
            Double changePercent = null;
            Double snapshotPrice = snapshotPriceCache.get(symbolCode);
            boolean isStockOrEtf = "stock".equalsIgnoreCase(assetType) || "etf".equalsIgnoreCase(assetType);

            if (hasNewPrice) {
                if (prevClose != null && Double.compare(prevClose, 0.0d) != 0) {
                    changePercent = (newPrice - prevClose) / prevClose * 100;
                } else if (snapshotPrice != null && Double.compare(snapshotPrice, 0.0d) != 0) {
                    changePercent = (newPrice - snapshotPrice) / snapshotPrice * 100;
                }
            }

            boolean isTrading = "TRADING".equalsIgnoreCase(marketState);
            Double cachedChangePercent = lastChangePercentCache.get(symbolCode);
            if (isStockOrEtf && !isTrading) {
                if (changePercent == null) {
                    changePercent = cachedChangePercent;
                }
            } else {
                if (changePercent == null) {
                    changePercent = cachedChangePercent;
                }
            }
            if (changePercent != null && hasNewPrice) {
                lastChangePercentCache.put(symbolCode, changePercent);
            } else if (changePercent != null && cachedChangePercent == null) {
                lastChangePercentCache.put(symbolCode, changePercent);
            }

            if (hasNewPrice && newPrice != null) {
                snapshotPriceCache.put(symbolCode, newPrice);
            }

            Double volume = quote.getVolume();
            Long timestamp = quote.getTsUnixMs() != null ? quote.getTsUnixMs() : System.currentTimeMillis();
            Long tradeTimestamp = quote.getTsUnixMs() != null ? quote.getTsUnixMs() : timestamp;
            String tradeTime = quote.getTsLocal();

            boolean isDaily = subject != null && subject.endsWith(".daily");
            if (isDaily) {
                java.sql.Date tradeDate = toUtcDate(quote.getTsUnixMs(), quote.getTsLocal());
                quoteRepository.upsertDaily(symbolId, quote, jsonPayload, tradeDate);
            } else {
                QuoteMessage quoteMessage = new QuoteMessage(
                        quote.getAssetType(),
                        quote.getSymbol(),
                        price,
                        change,
                        changePercent,
                        volume,
                        System.currentTimeMillis(),
                        tradeTimestamp,
                        tradeTime);
                quoteStreamProducer.publish(quoteMessage);
            }
        } catch (Exception e) {
            log.error("ingest failed: {}", e.getMessage(), e);
        }
    }

    private static java.sql.Date toUtcDate(Long tsUnixMs, String tsIso) {
        try {
            if (tsUnixMs != null) {
                return java.sql.Date.valueOf(Instant.ofEpochMilli(tsUnixMs).atZone(ZoneOffset.UTC).toLocalDate());
            }
            if (tsIso != null) {
                return java.sql.Date.valueOf(Instant.parse(tsIso).atZone(ZoneOffset.UTC).toLocalDate());
            }
        } catch (Exception ignored) {
        }
        return java.sql.Date.valueOf(Instant.now().atZone(ZoneOffset.UTC).toLocalDate());
    }
}
