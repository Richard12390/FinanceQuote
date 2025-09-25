package com.org.quoteservice.server.application.service;

import com.fasterxml.jackson.databind.ObjectMapper;
import java.time.Instant;
import java.time.ZoneOffset;
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

            boolean isDaily = subject != null && subject.endsWith(".daily");
            if (isDaily) {
                java.sql.Date tradeDate = toUtcDate(quote.getTsUnixMs(), quote.getTsLocal());
                quoteRepository.upsertDaily(symbolId, quote, jsonPayload, tradeDate);
            }
            // else: future realtime tick handling (e.g. WebSocket push)
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